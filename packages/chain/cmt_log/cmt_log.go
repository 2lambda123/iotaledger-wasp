// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// package cmtLog is responsible for producing a log of chain's block decisions
// for a particular committee. The main tasks for this module are:
//   - Track the head of the chain log for a committee.
//   - Track which blocks are approved, pending or reverted.
//   - Handle startup and recovery scenarios.
//
// Main principles of the algorithm:
// >   - Propose to go to the next log index when
// >       - a) consensus for the latest known LogIndex has terminated (confirmed | rejected | skip | recover).
// >            All the previous can be left uncompleted.
// >       - b) confirmed anchor output is received from L1, that we haven't posted in this node.
// >       - and we have a tip AO (from the local view).
// >   - If there is a single clean chain, we can use pipelining (consider consDone event
// >     instead of waiting for confirmed | rejected).
// >   - We start a consensus instance whenever VLI reports a quorum and we have not started it yet.
// >     Even if we provided no proposal for that LI.
//
// The algorithm at a high level:
//
// > ON Startup:
// >     Let prevLI <- TRY restoring the last started LogIndex ELSE 0
// >     MinLI <- prevLI + 1
// >     LogIndex.Start(prevLI)
// > UPON AnchorOutput (AO) {Confirmed | Rejected} by L1:
// >     LocalView.Update(AO)
// >     IF LocalView changed THEN
// >         LogIndex.L1ReplacedBaseAnchorOutput()
// >         TryProposeConsensus()
//
// > ON Startup:
// >     Let prevLI <- TRY restoring the last started LogIndex ELSE 0
// >     MinLI <- prevLI + 1
// >     LogIndex.Start(prevLI)
// >     TryProposeConsensus()
// > UPON AnchorOutput (AO) {Confirmed | Rejected} by L1:
// >     LocalView.Update(AO)
// >     IF LocalView changed THEN
// >         LogIndex.L1ReplacedBaseAnchorOutput()
// >         TryProposeConsensus()
// > ON ConsensusOutput/DONE (CD)
// >     LocalView.Update(CD)
// >     IF LocalView changed THEN
// >         LogIndex.ConsensusOutput(CD.LogIndex)
// >         TryProposeConsensus()
// > ON ConsensusOutput/SKIP (CS)
// >     LogIndex.ConsensusOutput(CS.LogIndex)
// >     TryProposeConsensus()
// > ON ConsensusTimeout (CT)
// >     LogIndex.ConsensusTimeout(CT.LogIndex)
// >     TryProposeConsensus()
// > ON Suspend:
// >     Suspended <- TRUE
// >     TryProposeConsensus()
// > ON Reception of ⟨NextLI, •⟩ message:
// >     LogIndex.Receive(⟨NextLI, •⟩ message).
// >     TryProposeConsensus()
// > PROCEDURE TryProposeConsensus:
// >     IF ∧ LocalView.BaseAO ≠ NIL
// >        ∧ LogIndex > ConsensusLI
// >        ∧ LogIndex ≥ MinLI // ⇒ LogIndex ≠ NIL
// >        ∧ ¬ Suspended
// >     THEN
// >         Persist LogIndex
// >         ConsensusLI <- LogIndex
// >         Propose LocalView.BaseAO for LogIndex
// >     ELSE
// >         Don't propose any consensus.
//
// See `WaspChainRecovery.tla` for more precise specification.
//
// Notes and invariants:
//   - Here only a single consensus instance will be considered needed for this node at a time.
//     Other instances may continue running, but their results will be ignored. That's
//     because a consensus takes an input from the previous consensus output (the base
//     anchor ID and other parts that depend on it).
//   - A consensus is started when we have new log index greater than that we have
//     crashed with, and there is an anchor output received.
//
// ## Summary.
//
// Inputs expected:
//   - Consensus: Start -> Done | Timeout.
//   - AnchorOutput: Confirmed | Rejected -> {}.
//   - Suspend.
//
// Messages exchanged:
//   - NextLogIndex (private, between cmtLog instances).
package cmt_log

import (
	"errors"
	"fmt"

	"github.com/iotaledger/hive.go/log"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/chain/cons"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/metrics"

	"github.com/iotaledger/wasp/packages/tcrypto"
	"github.com/iotaledger/wasp/packages/util/byz_quorum"
)

// Public interface for this algorithm.
type CmtLog interface {
	AsGPA() gpa.GPA
}

type State struct {
	LogIndex LogIndex
}

// Interface used to store and recover the existing persistent state.
// To be implemented by the registry.
type ConsensusStateRegistry interface {
	Get(chainID isc.ChainID, committeeAddress iotago.Address) (*State, error) // Can return ErrCmtLogStateNotFound.
	Set(chainID isc.ChainID, committeeAddress iotago.Address, state *State) error
}

var ErrCmtLogStateNotFound = errors.New("errCmtLogStateNotFound")

// Output for this protocol indicates, what instance of a consensus
// is currently required to be run. The unique identifier here is the
// logIndex (there will be no different baseAnchorOutputs for the same logIndex).
type Output struct {
	logIndex       LogIndex
	consensusInput cons.Input
}

func makeOutput(logIndex LogIndex, consensusInput cons.Input) *Output {
	return &Output{logIndex: logIndex, consensusInput: consensusInput}
}

func (o *Output) GetLogIndex() LogIndex {
	return o.logIndex
}

func (o *Output) ConsensusInput() cons.Input {
	return o.consensusInput
}

func (o *Output) String() string {
	return fmt.Sprintf("{Output, logIndex=%v, consensusInput=%v}", o.logIndex, o.consensusInput)
}

// Protocol implementation.
type cmtLogImpl struct {
	chainID                isc.ChainID            // Chain, for which this log is maintained by this committee.
	cmtAddr                iotago.Address         // Address of the committee running this chain.
	consensusStateRegistry ConsensusStateRegistry // Persistent storage.
	varLogIndex            VarLogIndex            // Calculates the current log index.
	varLocalView           VarLocalView           // Tracks the pending anchor outputs.
	varOutput              VarOutput              // Calculate the output.
	asGPA                  gpa.GPA                // This object, just with all the needed wrappers.
	log                    log.Logger
}

var _ gpa.GPA = &cmtLogImpl{}

// Construct new node instance for this protocol.
//
// > ON Startup:
// >     Let prevLI <- TRY restoring the last started LogIndex ELSE 0
// >     MinLI <- prevLI + 1
// >     ...
func New(
	me gpa.NodeID,
	chainID isc.ChainID,
	dkShare tcrypto.DKShare,
	consensusStateRegistry ConsensusStateRegistry,
	nodeIDFromPubKey func(pubKey *cryptolib.PublicKey) gpa.NodeID,
	deriveAOByQuorum bool,
	pipeliningLimit int,
	cclMetrics *metrics.ChainCmtLogMetrics,
	log log.Logger,
) (CmtLog, error) {
	cmtAddr := dkShare.GetSharedPublic().AsEd25519Address()
	//
	// Load the last LogIndex we were working on.
	var prevLI LogIndex
	state, err := consensusStateRegistry.Get(chainID, cmtAddr)
	if err != nil {
		if !errors.Is(err, ErrCmtLogStateNotFound) {
			return nil, fmt.Errorf("cannot load cmtLogState for %v: %w", cmtAddr, err)
		}
		prevLI = NilLogIndex()
	} else {
		// Don't participate in the last stored LI, because maybe we have already sent some messages.
		prevLI = state.LogIndex
	}
	//
	// Make node IDs.
	nodePKs := dkShare.GetNodePubKeys()
	nodeIDs := make([]gpa.NodeID, len(nodePKs))
	for i := range nodeIDs {
		nodeIDs[i] = nodeIDFromPubKey(nodePKs[i])
	}
	//
	// Construct the object.
	n := len(nodeIDs)
	f := dkShare.DSS().MaxFaulty()
	if f > byz_quorum.MaxF(n) {
		log.LogPanicf("invalid f=%v for n=%v", f, n)
	}
	//
	// Log important info.
	log.LogInfof("Committee: N=%v, F=%v, address=%v", n, f, cmtAddr.String())
	for i := range nodePKs {
		log.LogInfof("Committee node[%v]=%v", i, nodePKs[i])
	}
	//
	// Create it.
	cl := &cmtLogImpl{
		chainID:                chainID,
		cmtAddr:                cmtAddr,
		consensusStateRegistry: consensusStateRegistry,
		varLogIndex:            nil, // Set bellow.
		varLocalView:           nil, // Set bellow.
		varOutput:              nil, // Set bellow.
		log:                    log,
	}
	cl.varOutput = NewVarOutput(func(li LogIndex) {
		if err := consensusStateRegistry.Set(chainID, cmtAddr, &State{LogIndex: li}); err != nil {
			// Nothing to do, if we cannot persist this.
			panic(fmt.Errorf("cannot persist the cmtLog state: %w", err))
		}
	}, log.NewChildLogger("VO"))
	cl.varLogIndex = NewVarLogIndex(nodeIDs, n, f, prevLI, cl.varOutput.LogIndexAgreed, cclMetrics, log.NewChildLogger("VLI"))
	cl.varLocalView = NewVarLocalView(pipeliningLimit, cl.varOutput.ConsInputChanged, log.NewChildLogger("VLV"))
	cl.asGPA = gpa.NewOwnHandler(me, cl)
	return cl, nil
}

// Implements the CmtLog interface.
func (cl *cmtLogImpl) AsGPA() gpa.GPA {
	return cl.asGPA
}

// Implements the gpa.GPA interface.
func (cl *cmtLogImpl) Input(input gpa.Input) gpa.OutMessages {
	cl.log.LogDebugf("Input %T: %+v", input, input)
	switch input := input.(type) {
	case *inputAnchorOutputConfirmed:
		return cl.handleInputChainOutputsConfirmed(input)
	case *inputConsensusOutputDone:
		return cl.handleInputConsensusOutputDone(input)
	case *inputConsensusOutputSkip:
		return cl.handleInputConsensusOutputSkip(input)
	case *inputConsensusOutputConfirmed:
		return cl.handleInputConsensusOutputConfirmed(input)
	case *inputConsensusOutputRejected:
		return cl.handleInputConsensusOutputRejected(input)
	case *inputConsensusTimeout:
		return cl.handleInputConsensusTimeout(input)
	case *inputCanPropose:
		cl.handleInputCanPropose()
		return nil
	case *inputSuspend:
		cl.handleInputSuspend()
		return nil
	}
	panic(fmt.Errorf("unexpected input %T: %+v", input, input))
}

// Implements the gpa.GPA interface.
func (cl *cmtLogImpl) Message(msg gpa.Message) gpa.OutMessages {
	msgNLI, ok := msg.(*MsgNextLogIndex)
	if !ok {
		cl.log.LogWarnf("dropping unexpected message %T: %+v", msg, msg)
		return nil
	}
	return cl.handleMsgNextLogIndex(msgNLI)
}

// > UPON AnchorOutput (AO) {Confirmed | Rejected} by L1:
// >   ...
func (cl *cmtLogImpl) handleInputChainOutputsConfirmed(input *inputAnchorOutputConfirmed) gpa.OutMessages {
	msgs := gpa.NoMessages()
	tipChanged := false
	cnfLogIndex := cl.varLocalView.ChainOutputsConfirmed(
		input.confirmedOutputs,
		func(consInput cons.Input) {
			cl.varOutput.Suspended(false)
			tipChanged = true
			msgs.AddAll(cl.varLogIndex.L1ReplacedBaseAnchorOutput())
		},
	)
	if !tipChanged && !cnfLogIndex.IsNil() {
		msgs.AddAll(cl.varLogIndex.L1ConfirmedAnchorOutput(cnfLogIndex))
	}
	return msgs
}

// >   ...
func (cl *cmtLogImpl) handleInputConsensusOutputConfirmed(input *inputConsensusOutputConfirmed) gpa.OutMessages {
	return cl.varLogIndex.ConsensusOutputReceived(input.logIndex) // This should be superfluous, always follows handleInputConsensusOutputDone.
}

// >   ...
func (cl *cmtLogImpl) handleInputConsensusOutputRejected(input *inputConsensusOutputRejected) gpa.OutMessages {
	msgs := gpa.NoMessages()
	msgs.AddAll(cl.varLogIndex.ConsensusOutputReceived(input.logIndex)) // This should be superfluous, always follows handleInputConsensusOutputDone.
	cl.varLocalView.ChainOutputsRejected((input.anchorOutput), func(consInput cons.Input) {
		msgs.AddAll(cl.varLogIndex.L1ReplacedBaseAnchorOutput())
	})
	return msgs
}

// > ON ConsensusOutput/DONE (CD)
// >   ...
func (cl *cmtLogImpl) handleInputConsensusOutputDone(input *inputConsensusOutputDone) gpa.OutMessages {
	cl.varLocalView.ConsensusOutputDone(input.logIndex, input.result, func(consInput cons.Input) {})
	return cl.varLogIndex.ConsensusOutputReceived(input.logIndex)
}

// > ON ConsensusOutput/SKIP (CS)
// >   ...
func (cl *cmtLogImpl) handleInputConsensusOutputSkip(input *inputConsensusOutputSkip) gpa.OutMessages {
	return cl.varLogIndex.ConsensusOutputReceived(input.logIndex)
}

// > ON ConsensusTimeout (CT)
// >   ...
func (cl *cmtLogImpl) handleInputConsensusTimeout(input *inputConsensusTimeout) gpa.OutMessages {
	return cl.varLogIndex.ConsensusRecoverReceived(input.logIndex)
}

func (cl *cmtLogImpl) handleInputCanPropose() {
	cl.varOutput.CanPropose()
}

// > ON Suspend:
// >   ...
func (cl *cmtLogImpl) handleInputSuspend() {
	cl.varOutput.Suspended(true)
}

// > ON Reception of ⟨NextLI, •⟩ message:
// >   ...
func (cl *cmtLogImpl) handleMsgNextLogIndex(msg *MsgNextLogIndex) gpa.OutMessages {
	return cl.varLogIndex.MsgNextLogIndexReceived(msg)
}

// Implements the gpa.GPA interface.
func (cl *cmtLogImpl) Output() gpa.Output {
	out := cl.varOutput.Value()
	if out == nil {
		return nil // Untyped nil.
	}
	return out
}

// Implements the gpa.GPA interface.
func (cl *cmtLogImpl) StatusString() string {
	return fmt.Sprintf(
		"{cmtLogImpl, %v, %v, %v}",
		cl.varOutput.StatusString(),
		cl.varLocalView.StatusString(),
		cl.varLogIndex.StatusString(),
	)
}
