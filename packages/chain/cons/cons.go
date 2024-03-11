// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// Consensus. Single instance of it.
//
// Used sub-protocols (on the same thread):
//   - DSS -- Distributed Schnorr Signature
//   - ACS -- Asynchronous Common Subset
//
// Used components (running on other threads):
//   - Mempool
//   - StateMgr
//   - VM
//
// > INPUT: baseAnchorOutputID
// > ON Startup:
// >     Start a DSS.
// >     Ask Mempool for backlog (based on baseAnchorOutputID).
// >     Ask StateMgr for a virtual state (based on baseAnchorOutputID).
// > UPON Reception of responses from Mempool, StateMgr and DSS NonceIndexes:
// >     Produce a batch proposal.
// >     Start the ACS.
// > UPON Reception of ACS output:
// >     IF result is possible THEN
// >         Submit agreed NonceIndexes to DSS.
// >         Send the BLS partial signature.
// >     ELSE
// >         OUTPUT SKIP
// > UPON Reception of N-2F BLS partial signatures:
// >     Start VM.
// > UPON Reception of VM Result:
// >     IF result is non-empty THEN
// >         Save the produced block to SM.
// >         Submit the result hash to the DSS.
// >     ELSE
// >         OUTPUT SKIP
// > UPON Reception of VM Result and a signature from the DSS
// >     IF rotation THEN
// >        OUTPUT Signed Governance TX.
// >     ELSE
// >        Save the block to the StateMgr.
// >        OUTPUT Signed State Transition TX
//
// We move all the synchronization logic to separate objects (upon_...). They are
// responsible for waiting specific data and then triggering the next state action
// once. This way we hope to solve a lot of race conditions gracefully. The `upon`
// predicates and the corresponding done functions should not depend on each other.
// If some data is needed at several places, it should be passed to several predicates.
//
// TODO: Handle the requests gracefully in the VM before getting the initTX.
// TODO: Reconsider the termination. Do we need to wait for DSS, RND?
package cons

import (
	"crypto/ed25519"
	"encoding/binary"
	"fmt"
	"time"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"

	"github.com/iotaledger/hive.go/log"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/iota.go/v4/api"
	"github.com/iotaledger/iota.go/v4/builder"
	"github.com/iotaledger/wasp/packages/chain/cons/bp"
	"github.com/iotaledger/wasp/packages/chain/dss"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/gpa/acs"
	"github.com/iotaledger/wasp/packages/gpa/cc/blssig"
	"github.com/iotaledger/wasp/packages/gpa/cc/semi"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/rotate"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/tcrypto"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm"
	"github.com/iotaledger/wasp/packages/vm/processors"
)

type Cons interface {
	AsGPA() gpa.GPA
}

// This is the result of the chain tip tracking.
// Here we decide the latest block to build on,
// optionally a block to use as a tip and
// a list of transactions that should be resubmitted
// (by producing and signing new blocks).

type Input interface {
	BaseBlock() *iotago.Block              // Can be nil or present in all cases.
	BaseCO() *isc.ChainOutputs             // Either BaseCO
	ReattachTX() *iotago.SignedTransaction // or reattachTX will be present.
}

type OutputStatus byte

func (os OutputStatus) String() string {
	switch os {
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Skipped:
		return "Skipped"
	default:
		return fmt.Sprintf("Unexpected-%v", byte(os))
	}
}

const (
	Running   OutputStatus = iota // Instance is still running.
	Completed                     // Consensus reached, TX is prepared for publication.
	Skipped                       // Consensus reached, no TX should be posted for this LogIndex.
)

type Output struct {
	Status     OutputStatus
	Terminated bool
	//
	// Requests for other components.
	NeedMempoolProposal       *isc.ChainOutputs // Requests for the mempool are needed for this Base Alias Output.
	NeedMempoolRequests       []*isc.RequestRef // Request payloads are needed from mempool for this IDs/Hash.
	NeedStateMgrStateProposal *isc.ChainOutputs // Query for a proposal for Virtual State (it will go to the batch proposal).
	NeedStateMgrDecidedState  *isc.ChainOutputs // Query for a decided Virtual State to be used by VM.
	NeedStateMgrSaveBlock     state.StateDraft  // Ask StateMgr to save the produced block.
	NeedNodeConnBlockTipSet   bool              // We need a tip set for a block now. // TODO: Handle it.
	NeedVMResult              *vm.VMTask        // VM Result is needed for this (agreed) batch.
	//
	// Following is the final result.
	// All the fields are filled, if State == Completed.
	Result *Result
}

type Result struct {
	producedChainOutputs    *isc.ChainOutputs         // The produced chain outputs.
	producedTransaction     *iotago.SignedTransaction // The TX for committing the block.
	producedIotaBlock       *iotago.Block             // Block produced to publish the TX.
	producedStateBlock      state.Block               // The state diff produced.
	consumedAnchorOutputID  iotago.OutputID           // Consumed in the TX.
	consumedAccountOutputID iotago.OutputID           // Consumed in the TX.
	// TODO: Cleanup the following.
	// Transaction      *iotago.SignedTransaction // The TX for committing the block.
	// BaseAnchorOutput iotago.OutputID           // AO consumed in the TX.
	// NextAnchorOutput *isc.ChainOutputs         // AO produced in the TX.
	// Block            state.Block               // The state diff produced.
}

func (r *Result) String() string {
	txID, err := r.producedTransaction.ID()
	if err != nil {
		txID = iotago.SignedTransactionID{}
	}
	return fmt.Sprintf("{cons.Result, txID=%v, baseAO=%v, nextAO=%v}", txID, r.consumedAnchorOutputID.ToHex(), r.producedChainOutputs)
}

func (r *Result) ProducedChainOutputs() *isc.ChainOutputs        { return r.producedChainOutputs }
func (r *Result) ProducedTransaction() *iotago.SignedTransaction { return r.producedTransaction }
func (r *Result) ProducedIotaBlock() *iotago.Block               { return r.producedIotaBlock }
func (r *Result) ProducedStateBlock() state.Block                { return r.producedStateBlock }
func (r *Result) ConsumedAnchorOutputID() iotago.OutputID        { return r.consumedAnchorOutputID }
func (r *Result) ConsumedAccountOutputID() iotago.OutputID       { return r.consumedAccountOutputID }

// Block might be nil, so check it before calling this.
func (r *Result) MustIotaBlockID() iotago.BlockID {
	blockID, err := r.producedIotaBlock.ID()
	if err != nil {
		panic(fmt.Errorf("failed to get BlockID: %v", err))
	}
	return blockID
}

// Transaction will always be set, so it should be safe to call this.
func (r *Result) MustSignedTransactionID() iotago.SignedTransactionID {
	txID, err := r.producedTransaction.ID()
	if err != nil {
		panic(fmt.Errorf("failed to get TX ID: %v", err))
	}
	return txID
}

type consImpl struct {
	instID           []byte // Consensus Instance ID.
	chainID          isc.ChainID
	chainStore       state.Store
	edSuite          suites.Suite // For signatures.
	blsSuite         suites.Suite // For randomness only.
	dkShare          tcrypto.DKShare
	l1APIProvider    iotago.APIProvider
	tokenInfo        *api.InfoResBaseToken
	processorCache   *processors.Cache
	nodeIDs          []gpa.NodeID
	me               gpa.NodeID
	f                int
	asGPA            gpa.GPA
	dssT             dss.DSS
	dssB             dss.DSS
	acs              acs.ACS
	subMP            SyncMP         // Mempool.
	subSM            SyncSM         // StateMgr.
	subNC            SyncNC         // NodeConn.
	subDSSt          SyncDSS        // Distributed Schnorr Signature to sign the TX.
	subDSSb          SyncDSS        // Distributed Schnorr Signature to sign the block.
	subACS           SyncACS        // Asynchronous Common Subset.
	subRND           SyncRND        // Randomness.
	subVM            SyncVM         // Virtual Machine.
	subTXS           SyncTXSig      // Building final TX.
	subBlkD          SyncBlkData    // Builds the block, not signed yet.
	subBlkS          SyncBlkSig     // Builds the signed block.
	subRes           SyncRes        // Collects the consensus result.
	term             *termCondition // To detect, when this instance can be terminated.
	msgWrapper       *gpa.MsgWrapper
	output           *Output
	validatorAgentID isc.AgentID
	log              log.Logger
}

const (
	subsystemTypeDSS byte = iota
	subsystemTypeACS
)

const (
	subsystemTypeDSSIndexT int = iota
	subsystemTypeDSSIndexB
)

var (
	_ gpa.GPA = &consImpl{}
	_ Cons    = &consImpl{}
)

func New(
	l1APIProvider iotago.APIProvider,
	tokenInfo *api.InfoResBaseToken,
	chainID isc.ChainID,
	chainStore state.Store,
	me gpa.NodeID,
	mySK *cryptolib.PrivateKey,
	dkShare tcrypto.DKShare,
	processorCache *processors.Cache,
	instID []byte,
	nodeIDFromPubKey func(pubKey *cryptolib.PublicKey) gpa.NodeID,
	validatorAgentID isc.AgentID,
	log log.Logger,
) Cons {
	edSuite := tcrypto.DefaultEd25519Suite()
	blsSuite := tcrypto.DefaultBLSSuite()

	dkShareNodePubKeys := dkShare.GetNodePubKeys()
	nodeIDs := make([]gpa.NodeID, len(dkShareNodePubKeys))
	nodePKs := map[gpa.NodeID]kyber.Point{}
	for i := range dkShareNodePubKeys {
		var err error
		nodeIDs[i] = nodeIDFromPubKey(dkShareNodePubKeys[i])
		nodePKs[nodeIDs[i]], err = dkShareNodePubKeys[i].AsKyberPoint()
		if err != nil {
			panic(fmt.Errorf("cannot convert nodePK[%v] to kyber.Point: %w", i, err))
		}
	}

	f := len(dkShareNodePubKeys) - int(dkShare.GetT())
	myKyberKeys, err := mySK.AsKyberKeyPair()
	if err != nil {
		panic(fmt.Errorf("cannot convert node's SK to kyber.Scalar: %w", err))
	}
	longTermDKS := dkShare.DSS()
	acsLog := log.NewChildLogger("ACS")
	acsCCInstFunc := func(nodeID gpa.NodeID, round int) gpa.GPA {
		var roundBin [4]byte
		binary.BigEndian.PutUint32(roundBin[:], uint32(round))
		sid := hashing.HashDataBlake2b(instID, nodeID[:], roundBin[:]).Bytes()
		realCC := blssig.New(blsSuite, nodeIDs, dkShare.BLSCommits(), dkShare.BLSPriShare(), int(dkShare.BLSThreshold()), me, sid, acsLog)
		return semi.New(round, realCC)
	}
	c := &consImpl{
		instID:           instID,
		chainID:          chainID,
		chainStore:       chainStore,
		edSuite:          edSuite,
		blsSuite:         blsSuite,
		dkShare:          dkShare,
		processorCache:   processorCache,
		nodeIDs:          nodeIDs,
		l1APIProvider:    l1APIProvider,
		tokenInfo:        tokenInfo,
		me:               me,
		f:                f,
		dssT:             dss.New(edSuite, nodeIDs, nodePKs, f, me, myKyberKeys.Private, longTermDKS, log.NewChildLogger("DSSt")),
		dssB:             dss.New(edSuite, nodeIDs, nodePKs, f, me, myKyberKeys.Private, longTermDKS, log.NewChildLogger("DSSb")),
		acs:              acs.New(nodeIDs, me, f, acsCCInstFunc, acsLog),
		output:           &Output{Status: Running},
		log:              log,
		validatorAgentID: validatorAgentID,
	}
	c.asGPA = gpa.NewOwnHandler(me, c)
	c.msgWrapper = gpa.NewMsgWrapper(msgTypeWrapped, c.msgWrapperFunc)
	c.subMP = NewSyncMP(
		c.uponMPProposalInputsReady,
		c.uponMPProposalReceived,
		c.uponMPRequestsNeeded,
		c.uponMPRequestsReceived,
	)
	c.subSM = NewSyncSM(
		c.uponSMStateProposalQueryInputsReady,
		c.uponSMStateProposalReceived,
		c.uponSMDecidedStateQueryInputsReady,
		c.uponSMDecidedStateReceived,
		c.uponSMSaveProducedBlockInputsReady,
		c.uponSMSaveProducedBlockDone,
	)
	c.subNC = NewSyncNC(
		c.uponNCBlockTipSetNeeded,
		c.uponNCBlockTipSetReceived,
	)
	c.subDSSt = NewSyncDSS(
		c.uponDSStInitialInputsReady,
		c.uponDSStIndexProposalReady,
		c.uponDSStSigningInputsReceived,
		c.uponDSStOutputReady,
	)
	c.subDSSb = NewSyncDSS(
		c.uponDSSbInitialInputsReady,
		c.uponDSSbIndexProposalReady,
		c.uponDSSbSigningInputsReceived,
		c.uponDSSbOutputReady,
	)
	c.subACS = NewSyncACS(
		c.uponACSTipsRequired,
		c.uponACSInputsReceived,
		c.uponACSOutputReceived,
		c.uponACSTerminated,
	)
	c.subRND = NewSyncRND(
		int(dkShare.BLSThreshold()),
		c.uponRNDInputsReady,
		c.uponRNDSigSharesReady,
	)
	c.subVM = NewSyncVM(
		c.uponVMInputsReceived,
		c.uponVMOutputReceived,
	)
	c.subTXS = NewSyncTX(
		c.uponTXInputsReady,
	)
	c.subBlkD = NewSyncBlkData(
		c.uponBlkDataInputsReady,
	)
	c.subBlkS = NewSyncBlkSig(
		c.uponBlkSigInputsReady,
	)
	c.subRes = NewSyncRes(
		c.uponResInputsReady,
	)
	c.term = newTermCondition(
		c.uponTerminationCondition,
	)
	return c
}

// Used to select a target subsystem for a wrapped message received.
func (c *consImpl) msgWrapperFunc(subsystem byte, index int) (gpa.GPA, error) {
	if subsystem == subsystemTypeDSS {
		switch index {
		case subsystemTypeDSSIndexT:
			return c.dssT.AsGPA(), nil
		case subsystemTypeDSSIndexB:
			return c.dssB.AsGPA(), nil
		}
		return nil, fmt.Errorf("unexpected DSS index: %v", index)
	}
	if subsystem == subsystemTypeACS {
		if index != 0 {
			return nil, fmt.Errorf("unexpected ACS index: %v", index)
		}
		return c.acs.AsGPA(), nil
	}
	return nil, fmt.Errorf("unexpected subsystem: %v", subsystem)
}

func (c *consImpl) AsGPA() gpa.GPA {
	return c.asGPA
}

func (c *consImpl) Input(input gpa.Input) gpa.OutMessages {
	switch input := input.(type) {
	case *inputTimeData:
		// ignore this to filter out ridiculously excessive logging
	default:
		c.log.LogDebugf("Input %T: %+v", input, input)
	}

	switch input := input.(type) {
	case *inputProposal:
		c.log.LogInfof("Consensus started, received %v", input.String())
		msgs := gpa.NoMessages()
		msgs = msgs.
			AddAll(c.subDSSt.InitialInputReceived()).
			AddAll(c.subDSSb.InitialInputReceived())
		if input.params.BaseCO() != nil {
			return msgs.
				AddAll(c.subACS.TXCreateInputReceived(input.params.BaseCO(), input.params.BaseBlock())).
				AddAll(c.subMP.BaseAnchorOutputReceived(input.params.BaseCO())).
				AddAll(c.subSM.ProposedBaseAnchorOutputReceived(input.params.BaseCO()))
		}
		return msgs.
			AddAll(c.subACS.BlockOnlyInputReceived(input.params.ReattachTX(), input.params.BaseBlock()))
	case *inputMempoolProposal:
		return c.subMP.ProposalReceived(input.requestRefs)
	case *inputMempoolRequests:
		return c.subMP.RequestsReceived(input.requests)
	case *inputStateMgrProposalConfirmed:
		return c.subSM.StateProposalConfirmedByStateMgr()
	case *inputStateMgrDecidedVirtualState:
		return c.subSM.DecidedVirtualStateReceived(input.chainState)
	case *inputStateMgrBlockSaved:
		return c.subSM.BlockSaved(input.block)
	case *inputNodeConnBlockTipSet:
		return c.subNC.BlockTipSetReceived(input.strongParents)
	case *inputTimeData:
		return c.subACS.TimeUpdateReceived(input.timeData)
	case *inputVMResult:
		return c.subVM.VMResultReceived(input.task)
	}
	panic(fmt.Errorf("unexpected input: %v", input))
}

// Implements the gpa.GPA interface.
// Here we route all the messages.
func (c *consImpl) Message(msg gpa.Message) gpa.OutMessages {
	switch msgT := msg.(type) {
	case *msgBLSPartialSig:
		return c.subRND.BLSPartialSigReceived(msgT.Sender(), msgT.partialSig)
	case *gpa.WrappingMsg:
		sub, subMsgs, err := c.msgWrapper.DelegateMessage(msgT)
		if err != nil {
			c.log.LogWarnf("unexpected wrapped message: %w", err)
			return nil
		}
		msgs := gpa.NoMessages().AddAll(subMsgs)
		switch msgT.Subsystem() {
		case subsystemTypeACS:
			return msgs.AddAll(c.subACS.ACSOutputReceived(sub.Output()))
		case subsystemTypeDSS:
			switch msgT.Index() {
			case subsystemTypeDSSIndexT:
				return msgs.AddAll(c.subDSSt.DSSOutputReceived(sub.Output()))
			case subsystemTypeDSSIndexB:
				return msgs.AddAll(c.subDSSb.DSSOutputReceived(sub.Output()))
			default:
				c.log.LogWarnf("unexpected DSS index after check: %+v", msg)
				return nil
			}
		default:
			c.log.LogWarnf("unexpected subsystem after check: %+v", msg)
			return nil
		}
	}
	panic(fmt.Errorf("unexpected message: %v", msg))
}

func (c *consImpl) Output() gpa.Output {
	return c.output // Always non-nil.
}

func (c *consImpl) StatusString() string {
	// We con't include RND here, maybe that's less important, and visible from the VM status.
	return fmt.Sprintf("{consImpl⟨%v⟩,%v,%v,%v,%v,%v,%v,%v}",
		c.output.Status,
		c.subSM.String(),
		c.subMP.String(),
		c.subDSSt.String(),
		c.subDSSb.String(),
		c.subACS.String(),
		c.subVM.String(),
		c.subTXS.String(),
	)
}

////////////////////////////////////////////////////////////////////////////////
// MP -- MemPool

func (c *consImpl) uponMPProposalInputsReady(baseAnchorOutput *isc.ChainOutputs) gpa.OutMessages {
	c.output.NeedMempoolProposal = baseAnchorOutput
	return nil
}

func (c *consImpl) uponMPProposalReceived(requestRefs []*isc.RequestRef) gpa.OutMessages {
	c.output.NeedMempoolProposal = nil
	return gpa.NoMessages().
		// AddAll(c.subNC.MempoolProposalReceived()).
		AddAll(c.subACS.MempoolRequestsReceived(requestRefs))
}

func (c *consImpl) uponMPRequestsNeeded(requestRefs []*isc.RequestRef) gpa.OutMessages {
	c.output.NeedMempoolRequests = requestRefs
	return nil
}

func (c *consImpl) uponMPRequestsReceived(requests []isc.Request) gpa.OutMessages {
	c.output.NeedMempoolRequests = nil
	return c.subVM.RequestsReceived(requests)
}

////////////////////////////////////////////////////////////////////////////////
// SM -- StateManager

func (c *consImpl) uponSMStateProposalQueryInputsReady(baseAnchorOutput *isc.ChainOutputs) gpa.OutMessages {
	c.output.NeedStateMgrStateProposal = baseAnchorOutput
	return nil
}

func (c *consImpl) uponSMStateProposalReceived(proposedAnchorOutput *isc.ChainOutputs) gpa.OutMessages {
	c.output.NeedStateMgrStateProposal = nil
	return gpa.NoMessages().
		// AddAll(c.subNC.StateMgrProposalReceived()).
		AddAll(c.subACS.StateMgrProposalReceived(proposedAnchorOutput))
}

func (c *consImpl) uponSMDecidedStateQueryInputsReady(decidedBaseAnchorOutput *isc.ChainOutputs) gpa.OutMessages {
	c.output.NeedStateMgrDecidedState = decidedBaseAnchorOutput
	return nil
}

func (c *consImpl) uponSMDecidedStateReceived(chainState state.State) gpa.OutMessages {
	c.output.NeedStateMgrDecidedState = nil
	return c.subVM.DecidedStateReceived(chainState)
}

func (c *consImpl) uponSMSaveProducedBlockInputsReady(producedBlock state.StateDraft) gpa.OutMessages {
	if producedBlock == nil {
		// Don't have a block to save in the case of self-governed rotation.
		// So mark it as saved immediately.
		return c.subSM.BlockSaved(nil)
	}
	c.output.NeedStateMgrSaveBlock = producedBlock
	return nil
}

func (c *consImpl) uponSMSaveProducedBlockDone(block state.Block) gpa.OutMessages {
	c.output.NeedStateMgrSaveBlock = nil
	return gpa.NoMessages().
		AddAll(c.subTXS.BlockSaved(block)).
		AddAll(c.subRes.HaveStateBlock(block))
}

// //////////////////////////////////////////////////////////////////////////////
// NC

func (c *consImpl) uponNCBlockTipSetNeeded() gpa.OutMessages {
	c.output.NeedNodeConnBlockTipSet = true
	return nil
}

func (c *consImpl) uponNCBlockTipSetReceived(strongParents iotago.BlockIDs) gpa.OutMessages {
	c.output.NeedNodeConnBlockTipSet = false
	return c.subACS.BlockTipSetProposalReceived(strongParents)
}

// func (c *consImpl) uponNCInputsReady() gpa.OutMessages {
// 	c.output.NeedNodeConnBlockTipSet = true
// 	return nil
// }

// func (c *consImpl) uponNCOutputReady(
// 	blockToRefer *iotago.Block,
// 	txCreateInputReceived *isc.ChainOutputs,
// 	blockOnlyInputReceived *iotago.SignedTransaction,
// 	mempoolProposalReceived []*isc.RequestRef,
// 	dssTIndexProposal []int,
// 	dssBIndexProposal []int,
// 	timeData time.Time,
// 	strongParents iotago.BlockIDs,
// ) gpa.OutMessages {
// 	return c.subACS.ACSInputsReceived(
// 		blockToRefer,
// 		txCreateInputReceived,
// 		blockOnlyInputReceived,
// 		mempoolProposalReceived,
// 		dssTIndexProposal,
// 		dssBIndexProposal,
// 		timeData,
// 		strongParents,
// 	)
// }

// //////////////////////////////////////////////////////////////////////////////
// DSS_t

func (c *consImpl) uponDSStInitialInputsReady() gpa.OutMessages {
	c.log.LogDebugf("uponDSStInitialInputsReady")
	sub, subMsgs, err := c.msgWrapper.DelegateInput(subsystemTypeDSS, subsystemTypeDSSIndexT, dss.NewInputStart())
	if err != nil {
		panic(fmt.Errorf("cannot provide input to DSSt: %w", err))
	}
	return gpa.NoMessages().
		AddAll(subMsgs).
		AddAll(c.subDSSt.DSSOutputReceived(sub.Output()))
}

func (c *consImpl) uponDSStIndexProposalReady(indexProposal []int) gpa.OutMessages {
	c.log.LogDebugf("uponDSStIndexProposalReady")
	return gpa.NoMessages().
		// AddAll(c.subNC.DSStIndexProposalReceived()).
		AddAll(c.subACS.DSStIndexProposalReceived(indexProposal))
}

func (c *consImpl) uponDSStSigningInputsReceived(decidedIndexProposals map[gpa.NodeID][]int, messageToSign []byte) gpa.OutMessages {
	c.log.LogDebugf("uponDSStSigningInputsReceived(decidedIndexProposals=%+v, H(messageToSign)=%v)", decidedIndexProposals, hashing.HashDataBlake2b(messageToSign))
	dssDecidedInput := dss.NewInputDecided(decidedIndexProposals, messageToSign)
	subDSSt, subMsgs, err := c.msgWrapper.DelegateInput(subsystemTypeDSS, subsystemTypeDSSIndexT, dssDecidedInput)
	if err != nil {
		panic(fmt.Errorf("cannot provide inputs for signing: %w", err))
	}
	return gpa.NoMessages().
		AddAll(subMsgs).
		AddAll(c.subDSSt.DSSOutputReceived(subDSSt.Output()))
}

func (c *consImpl) uponDSStOutputReady(signature []byte) gpa.OutMessages {
	c.log.LogDebugf("uponDSStOutputReady")
	return c.subTXS.SignatureReceived(signature)
}

// //////////////////////////////////////////////////////////////////////////////
// DSS_b

func (c *consImpl) uponDSSbInitialInputsReady() gpa.OutMessages {
	c.log.LogDebugf("uponDSSbInitialInputsReady")
	sub, subMsgs, err := c.msgWrapper.DelegateInput(subsystemTypeDSS, subsystemTypeDSSIndexB, dss.NewInputStart())
	if err != nil {
		panic(fmt.Errorf("cannot provide input to DSSb: %w", err))
	}
	return gpa.NoMessages().
		AddAll(subMsgs).
		AddAll(c.subDSSb.DSSOutputReceived(sub.Output()))
}

func (c *consImpl) uponDSSbIndexProposalReady(indexProposal []int) gpa.OutMessages {
	c.log.LogDebugf("uponDSSbIndexProposalReady")
	return gpa.NoMessages().
		// AddAll(c.subNC.DSSbIndexProposalReceived()).
		AddAll(c.subACS.DSSbIndexProposalReceived(indexProposal))
}

func (c *consImpl) uponDSSbSigningInputsReceived(decidedIndexProposals map[gpa.NodeID][]int, messageToSign []byte) gpa.OutMessages {
	c.log.LogDebugf("uponDSSbSigningInputsReceived(decidedIndexProposals=%+v, H(messageToSign)=%v)", decidedIndexProposals, hashing.HashDataBlake2b(messageToSign))
	dssDecidedInput := dss.NewInputDecided(decidedIndexProposals, messageToSign)
	subDSSb, subMsgs, err := c.msgWrapper.DelegateInput(subsystemTypeDSS, subsystemTypeDSSIndexB, dssDecidedInput)
	if err != nil {
		panic(fmt.Errorf("cannot provide inputs for signing: %w", err))
	}
	return gpa.NoMessages().
		AddAll(subMsgs).
		AddAll(c.subDSSb.DSSOutputReceived(subDSSb.Output()))
}

func (c *consImpl) uponDSSbOutputReady(signature []byte) gpa.OutMessages {
	c.log.LogDebugf("uponDSSbOutputReady")
	return c.subBlkS.HaveSig(signature)
}

////////////////////////////////////////////////////////////////////////////////
// ACS

func (c *consImpl) uponACSTipsRequired() gpa.OutMessages {
	return c.subNC.BlockTipSetNeeded()
}

func (c *consImpl) uponACSInputsReceived(
	blockToRefer *iotago.Block,
	baseCO *isc.ChainOutputs,
	resignTX *iotago.SignedTransaction,
	requestRefs []*isc.RequestRef,
	dssTIndexProposal []int,
	dssBIndexProposal []int,
	timeData time.Time,
	strongParents iotago.BlockIDs,
) gpa.OutMessages {
	batchProposal := bp.NewBatchProposal(
		c.l1APIProvider.LatestAPI(),
		*c.dkShare.GetIndex(),
		blockToRefer,
		strongParents,
		baseCO,
		resignTX,
		util.NewFixedSizeBitVector(c.dkShare.GetN()).SetBits(dssTIndexProposal),
		util.NewFixedSizeBitVector(c.dkShare.GetN()).SetBits(dssBIndexProposal),
		timeData,
		c.validatorAgentID,
		requestRefs,
	)
	subACS, subMsgs, err := c.msgWrapper.DelegateInput(subsystemTypeACS, 0, batchProposal.Bytes())
	if err != nil {
		panic(fmt.Errorf("cannot provide input to the ACS: %w", err))
	}
	return gpa.NoMessages().
		AddAll(subMsgs).
		AddAll(c.subACS.ACSOutputReceived(subACS.Output()))
}

func (c *consImpl) uponACSOutputReceived(outputValues map[gpa.NodeID][]byte) gpa.OutMessages {
	aggr := bp.AggregateBatchProposals(outputValues, c.nodeIDs, c.f, c.l1APIProvider.LatestAPI(), c.log)
	if aggr.ShouldBeSkipped() {
		// Cannot proceed with such proposals.
		// Have to retry the consensus after some time with the next log index.
		c.log.LogInfof("Terminating consensus with status=Skipped, there is no way to aggregate batch proposal.")
		c.output.Status = Skipped
		c.term.haveOutputProduced()
		return nil
	}

	msgs := gpa.NoMessages().
		AddAll(c.subRND.CanProceed(c.instID)).
		AddAll(c.subDSSb.DecidedIndexProposalsReceived(aggr.DecidedDSSbIndexProposals())).
		AddAll(c.subBlkD.HaveTimestamp(aggr.AggregatedTime())).
		AddAll(c.subBlkD.HaveTipsProposal(func(randomness hashing.HashValue) iotago.BlockIDs { return aggr.DecidedStrongParents(randomness) }))

	//
	// Either we are going to build a fresh TX
	if aggr.ShouldBuildNewTX() {
		bao := aggr.DecidedBaseCO()
		reqs := aggr.DecidedRequestRefs()
		c.log.LogDebugf("ACS decision: baseAO=%v, requests=%v", bao, reqs)
		return msgs.
			AddAll(c.subMP.RequestsNeeded(reqs)).
			AddAll(c.subSM.DecidedVirtualStateNeeded(bao)).
			AddAll(c.subVM.DecidedBatchProposalsReceived(aggr)).
			AddAll(c.subDSSt.DecidedIndexProposalsReceived(aggr.DecidedDSStIndexProposals()))
	}
	// Or we are going to reuse the existing TX.
	return msgs.
		AddAll(c.subRes.ReuseTX(aggr.DecidedReattachTX())).
		AddAll(c.subBlkD.HaveSignedTX(aggr.DecidedReattachTX()))
}

func (c *consImpl) uponACSTerminated() {
	c.term.haveAcsTerminated()
}

////////////////////////////////////////////////////////////////////////////////
// RND

func (c *consImpl) uponRNDInputsReady(dataToSign []byte) gpa.OutMessages {
	sigShare, err := c.dkShare.BLSSignShare(dataToSign)
	if err != nil {
		panic(fmt.Errorf("cannot sign share for randomness: %w", err))
	}
	msgs := gpa.NoMessages()
	for _, nid := range c.nodeIDs {
		msgs.Add(newMsgBLSPartialSig(c.blsSuite, nid, sigShare))
	}
	return msgs
}

func (c *consImpl) uponRNDSigSharesReady(dataToSign []byte, partialSigs map[gpa.NodeID][]byte) (bool, gpa.OutMessages) {
	partialSigArray := make([][]byte, 0, len(partialSigs))
	for nid := range partialSigs {
		partialSigArray = append(partialSigArray, partialSigs[nid])
	}
	sig, err := c.dkShare.BLSRecoverMasterSignature(partialSigArray, dataToSign)
	if err != nil {
		c.log.LogWarnf("Cannot reconstruct BLS signature from %v/%v sigShares: %v", len(partialSigs), c.dkShare.GetN(), err)
		return false, nil // Continue to wait for other sig shares.
	}
	randomness := hashing.HashDataBlake2b(sig.Signature.Bytes())
	return true, gpa.NoMessages().
		AddAll(c.subVM.RandomnessReceived(randomness)).
		AddAll(c.subBlkD.HaveRandomness(randomness))
}

////////////////////////////////////////////////////////////////////////////////
// VM

func (c *consImpl) uponVMInputsReceived(aggregatedProposals *bp.AggregatedBatchProposals, chainState state.State, randomness *hashing.HashValue, requests []isc.Request) gpa.OutMessages {
	// TODO: chainState state.State is not used for now. That's because VM takes it form the store by itself.
	// The decided base anchor output can be different from that we have proposed!
	decidedBaseCO := aggregatedProposals.DecidedBaseCO()
	c.output.NeedVMResult = &vm.VMTask{
		Processors:           c.processorCache,
		Inputs:               decidedBaseCO,
		Store:                c.chainStore,
		Requests:             aggregatedProposals.OrderedRequests(requests, *randomness),
		Timestamp:            aggregatedProposals.AggregatedTime(),
		Entropy:              *randomness,
		ValidatorFeeTarget:   aggregatedProposals.ValidatorFeeTarget(*randomness),
		EstimateGasMode:      false,
		EnableGasBurnLogging: false,
		Log:                  c.log.NewChildLogger("VM"),
		L1APIProvider:        c.l1APIProvider,
		TokenInfo:            c.tokenInfo,
	}
	return nil
}

func (c *consImpl) uponVMOutputReceived(vmResult *vm.VMTaskResult) gpa.OutMessages {
	c.output.NeedVMResult = nil
	if len(vmResult.RequestResults) == 0 {
		// No requests were processed, don't have what to do.
		// Will need to retry the consensus with the next log index some time later.
		c.log.LogInfof("Terminating consensus with status=Skipped, 0 requests processed.")
		c.output.Status = Skipped
		c.term.haveOutputProduced()
		return nil
	}

	if vmResult.RotationAddress != nil {
		// Rotation by the Self-Governed Committee.
		l1API := c.l1APIProvider.APIForTime(vmResult.Task.Timestamp)
		tx, _, err := rotate.MakeRotationTransactionForSelfManagedChain(
			vmResult.RotationAddress,
			vmResult.Task.Inputs,
			l1API.TimeProvider().SlotFromTime(vmResult.Task.Timestamp),
			l1API,
		)
		if err != nil {
			c.log.LogWarnf("Cannot create rotation TX, failed to make TX essence: %w", err)
			c.output.Status = Skipped
			c.term.haveOutputProduced()
			return nil
		}
		vmResult.Transaction = tx
		vmResult.StateDraft = nil
	}

	signingMsg, err := vmResult.Transaction.SigningMessage()
	if err != nil {
		panic(fmt.Errorf("uponVMOutputReceived: cannot obtain signing message: %w", err))
	}

	chained, err := isc.ChainOutputsFromTx(vmResult.Transaction, c.chainID.AsAddress())
	if err != nil {
		panic(fmt.Errorf("cannot get AnchorOutput from produced TX: %w", err))
	}

	consumedAnchorOutputID := vmResult.Task.Inputs.AnchorOutputID
	var consumedAccountOutputID iotago.OutputID
	if accountOutputID, _, hasAccountOutputID := vmResult.Task.Inputs.AccountOutput(); hasAccountOutputID {
		consumedAccountOutputID = accountOutputID
	}
	return gpa.NoMessages().
		AddAll(c.subSM.BlockProduced(vmResult.StateDraft)).
		AddAll(c.subTXS.VMResultReceived(vmResult)).
		AddAll(c.subDSSt.MessageToSignReceived(signingMsg)).
		AddAll(c.subRes.HaveTransition(chained, consumedAnchorOutputID, consumedAccountOutputID))
}

////////////////////////////////////////////////////////////////////////////////
// TX

// Everything is ready for the output TX, produce it.
func (c *consImpl) uponTXInputsReady(vmResult *vm.VMTaskResult, block state.Block, signature []byte) gpa.OutMessages {
	// resultTx := vmResult.Transaction
	// publicKey := c.dkShare.GetSharedPublic()
	// var signatureArray [ed25519.SignatureSize]byte
	// copy(signatureArray[:], signature)
	// signatureForUnlock := &iotago.Ed25519Signature{
	// 	PublicKey: publicKey.AsKey(),
	// 	Signature: signatureArray,
	// }

	// resultInputs, err := resultTx.Inputs()
	// if err != nil {
	// 	panic(fmt.Errorf("cannot get inputs from result TX: %w", err))
	// }

	tx := &iotago.SignedTransaction{
		API: c.l1APIProvider.LatestAPI(), // TODO: Use the decided timestamp?
		Transaction: &iotago.Transaction{
			TransactionEssence: vmResult.Transaction.TransactionEssence,
		},
		// TODO: Unlocks: vmResult.Transaction.MakeSignatureAndReferenceUnlocks(len(resultInputs), signatureForUnlock),
	}

	return gpa.NoMessages().
		AddAll(c.subBlkD.HaveSignedTX(tx)).
		AddAll(c.subRes.BuiltTX(tx))
}

////////////////////////////////////////////////////////////////////////////////
// BLK

// readyCB func(tipsFn func(randomness hashing.HashValue) iotago.BlockIDs, randomness hashing.HashValue, timestamp time.Time, tx *iotago.SignedTransaction) gpa.OutMessages
func (c *consImpl) uponBlkDataInputsReady(
	tipsFn func(randomness hashing.HashValue) iotago.BlockIDs,
	randomness hashing.HashValue,
	timestamp time.Time,
	tx *iotago.SignedTransaction,
) gpa.OutMessages {
	strongParents := tipsFn(randomness)
	blk, err := builder.
		NewBasicBlockBuilder(c.l1APIProvider.APIForTime(timestamp)).
		StrongParents(strongParents).
		IssuingTime(timestamp).
		Payload(tx).
		Build()
	if err != nil {
		panic(fmt.Errorf("cannot build iota block: %v", err))
	}

	blkSigMsg, err := blk.SigningMessage()
	if err != nil {
		panic(fmt.Errorf("cannot build iota block: %v", err))
	}

	return gpa.NoMessages().
		AddAll(c.subBlkS.HaveBlock(blk)).
		AddAll(c.subDSSb.MessageToSignReceived(blkSigMsg))
}

func (c *consImpl) uponBlkSigInputsReady(
	bl *iotago.Block,
	sig []byte,
) gpa.OutMessages {
	var signatureArray [ed25519.SignatureSize]byte
	copy(signatureArray[:], sig)
	bl.Signature = &iotago.Ed25519Signature{
		PublicKey: c.dkShare.GetSharedPublic().AsKey(),
		Signature: signatureArray,
	}
	return c.subRes.HaveIotaBlock(bl)
}

////////////////////////////////////////////////////////////////////////////////
// RES

func (c *consImpl) uponResInputsReady(
	transactionReused bool,
	transaction *iotago.SignedTransaction,
	producedIotaBlock *iotago.Block,
	producedChainOutputs *isc.ChainOutputs,
	producedStateBlock state.Block,
	consumedAnchorOutputID iotago.OutputID,
	consumedAccountOutputID iotago.OutputID,
) gpa.OutMessages {
	transactionID, err := transaction.ID()
	if err != nil {
		panic(fmt.Errorf("cannot get ID from the produced TX: %w", err))
	}

	c.output.Result = &Result{
		producedTransaction:     transaction,
		producedChainOutputs:    producedChainOutputs,
		producedIotaBlock:       producedIotaBlock,
		producedStateBlock:      producedStateBlock,
		consumedAnchorOutputID:  consumedAnchorOutputID,
		consumedAccountOutputID: consumedAccountOutputID,
	}
	c.output.Status = Completed
	c.log.LogInfof(
		"Terminating consensus with status=Completed, produced tx.ID=%v, nextAO=%v, baseAO.ID=%v",
		transactionID.ToHex(), producedChainOutputs, consumedAnchorOutputID.ToHex(),
	)
	c.term.haveOutputProduced()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// TERM

func (c *consImpl) uponTerminationCondition() {
	c.output.Terminated = true
}
