// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// Here we implement the local view of a tangle, maintained by a committee to decide which
// anchor and account outputs to propose to the consensus. The anchor and account outputs
// decided by the consensus will be used as an input for TX we build.
//
// The LocalView maintains a list of Anchor/Account Outputs (AOs). They are chained based on consumed/produced
// AOs in a transactions we publish. The goal here is to track the unconfirmed anchor/account outputs, update
// the list based on confirmations/rejections from the L1.
//
// In overall, the LocalView acts as a filter between the L1 and LogIndex assignment in varLogIndex.
// It has to distinguish between AOs that are confirming a prefix of the posted transaction (pipelining),
// from other changes in L1 (rotations, rollbacks, rejections, reorgs, etc.).
//
// We have several inputs in this algorithm.
//
// Events from the L1:
//
//   - **Anchor OR Account Output Confirmed**.
//     They can posted by this committee,
//     as well as by another committee (e.g. chain was rotated to other committee and then back)
//     or a user (e.g. external rotation TX).
//     The anchor/account output confirmation events will be received independently of each other.
//
//   - **Anchor/Account Output Rejected**.
//     These events are always for TXes posted by this committee.
//     We assume for each TX we will get either Confirmation or Rejection.
//
//   - **Block outdated**.
//     The block was posted to the L1 network, but was not confirmed not rejected for long enough.
//     This event means we assume the block will not be confirmed anymore because of the max
//     block depth (L1 parameter) or is too old to be included into the current time slot.
//
// Events from the consensus:
//
//   - **Consensus Done**.
//     Consensus produced a TX, and will post it to the L1.
//     Only the successful decisions (non-skip) are interesting here, as they
//     provide new information for the local view. E.g. consensus output with the decision
//     SKIP will only cause increase of log index, which is irrelevant here.
//
// On the pipelining:
//
//   - During the normal operation, if consensus produces a TX, the chain can use it
//     to build next TX on it. That's pipelining. It allows to produce multiple TXes without
//     waiting L1 confirmations. This component tracks the AOs build but not confirmed yet.
//
//   - If an AO produced by the consensus is rejected, then all the AOs build on top of it will
//     be rejected eventually as well, because they use the rejected AO as an input. On the
//     other hand, it is unclear if unconfirmed AOs chained before the rejected output will be
//     confirmed or rejected, we will wait until L1 decides on all of them.
//
//   - If we get a confirmed AO that is not one of the AOs we have posted (and still waiting
//     for a decision), then someone from the outside of a committee transitioned the chain.
//     In this case all our produced/pending transactions are not meaningful anymore and we
//     have to start building all the chain from the newly received AO.
//     It is possible that the received unseen AO is posted by other nodes in the current
//     committee (e.g. the current node is lagging).
//
// Note on the AO as an input for a consensus. The provided AO is just a proposal. After ACS
// is completed, the participants will select the actual AO, which can differ from the one
// proposed by this node.
//
// On the rejections. When we get a rejection of an AO, we cannot mark all the subsequent
// StateIndexes as rejected, because it is possible that the rejected AO was started to publish
// before a reorg/reject. Thus, only that single AO has to be marked as rejected. Nevertheless,
// the AOs explicitly (via consumed AO) depending on the rejected AO can be cleaned up.
package cmt_log

import (
	"fmt"

	"github.com/iotaledger/hive.go/ds/shrinkingmap"
	"github.com/iotaledger/hive.go/logger"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/chain/cons"
	"github.com/iotaledger/wasp/packages/isc"
)

type VarLocalView interface {
	//
	// Corresponds to the `tx_posted` event in the specification.
	// Returns true, if the proposed BaseAnchorOutput has changed.
	ConsensusOutputDone(
		logIndex LogIndex,
		consResult *cons.Result,
		eventOutputCB VarLocalViewOutputCB,
	)
	//
	// Called upon receiving confirmation from the L1.
	// In a normal scenario both (Anchor/Account) outputs will be confirmed together,
	// because they are in the same TX. If someone moved one of those outputs externally,
	// they can be moved independently. In such a case, one of the anc/acc parameters will be nil.
	//   - It is important to get these events in the correct order, otherwise the out-of-order
	//     events will be considered as an reorg.
	ChainOutputsConfirmed(
		confirmedOutputs *isc.ChainOutputs,
		eventOutputCB VarLocalViewOutputCB,
	) LogIndex
	//
	// Called when the TX containing the specified outputs was rejected.
	// The outputs cannot be rejected independently because they are in the same TX.
	ChainOutputsRejected(
		rejected *isc.ChainOutputs,
		eventOutputCB VarLocalViewOutputCB,
	)
	//
	//
	BlockExpired(
		blockID iotago.BlockID,
		eventOutputCB VarLocalViewOutputCB,
	)
	//
	// Support functions.
	StatusString() string
}

type varLocalViewOutput struct { // Implements cons.Input
	baseBlock  *iotago.Block
	baseCO     *isc.ChainOutputs
	reattachTX *iotago.SignedTransaction
}

type VarLocalViewOutputCB = func(consInput cons.Input)

func newVarLocalViewOutput(
	baseBlock *iotago.Block,
	baseCO *isc.ChainOutputs,
	reattachTX *iotago.SignedTransaction,
) *varLocalViewOutput {
	return &varLocalViewOutput{
		baseBlock:  baseBlock,
		baseCO:     baseCO,
		reattachTX: reattachTX,
	}
}

func (o *varLocalViewOutput) BaseBlock() *iotago.Block              { return o.baseBlock }
func (o *varLocalViewOutput) BaseCO() *isc.ChainOutputs             { return o.baseCO }
func (o *varLocalViewOutput) ReattachTX() *iotago.SignedTransaction { return o.reattachTX }
func (o *varLocalViewOutput) Equals(other *varLocalViewOutput) bool {
	//
	// Compare the BaseBlock
	if (o.baseBlock == nil) != (other.baseBlock == nil) {
		return false
	}
	if o.baseBlock != nil && other.baseBlock != nil && o.baseBlock.MustID() != other.baseBlock.MustID() {
		return false
	}
	//
	// Compare the BaseCO
	if (o.baseCO == nil) != (other.baseCO == nil) {
		return false
	}
	if o.baseCO != nil && other.baseCO != nil && !o.baseCO.Equals(other.baseCO) {
		return false
	}
	//
	// Compare the ReattachTX
	if (o.reattachTX == nil) != (other.reattachTX == nil) {
		return false
	}
	if o.reattachTX != nil && other.reattachTX != nil {
		id1, err1 := o.reattachTX.ID()
		id2, err2 := other.reattachTX.ID()
		if err1 != nil {
			panic(fmt.Errorf("cannot extract TX ID: %v", err1))
		}
		if err2 != nil {
			panic(fmt.Errorf("cannot extract TX ID: %v", err2))
		}
		if id1 != id2 {
			return false
		}
	}
	return true
}

type varLocalViewEntry struct {
	producedChainOutputs    *isc.ChainOutputs                // The produced chain outputs.
	producedTransaction     *iotago.SignedTransaction        // The transaction publishing the chain outputs.
	consumedAnchorOutputID  iotago.OutputID                  // Consumed in the TX.
	consumedAccountOutputID iotago.OutputID                  // Consumed in the TX.
	blocks                  map[iotago.BlockID]*iotago.Block // All the non-expired blocks for this TX.
	reuse                   bool                             // True, if the TX should be reused.
	rejected                bool                             // True, if the TX as rejected. We keep them to detect the other rejected TXes.
	logIndex                LogIndex                         // LogIndex of the consensus produced the output, if any.
}

func (e *varLocalViewEntry) isTxReusable() bool {
	return len(e.blocks) == 0 && e.reuse && !e.rejected
}

func (e *varLocalViewEntry) txID() iotago.SignedTransactionID {
	id, err := e.producedTransaction.ID()
	if err != nil {
		panic(fmt.Errorf("cannot extract TX ID: %v", err))
	}
	return id
}

type varLocalViewImpl struct {
	//
	// The latest confirmed CO, as received from L1.
	// All the pending entries are built on top of this one.
	// It can be nil, if the latest AO is unclear (either not received yet).
	//
	// We don't use the isc.ChainOutputs structure here, because we can
	// receive the anchor/account outputs separately.
	confirmedAnchor  *isc.AnchorOutputWithID
	confirmedAccount *isc.AccountOutputWithID
	confirmedCO      *isc.ChainOutputs // Derived from the above, when all of them are received.
	//
	// AOs and blocks produced by this committee, but not confirmed yet.
	// It is possible to have several AOs for a StateIndex in the case of
	// Recovery/Timeout notices. Then the next consensus is started o build a TX.
	// Both of them can still produce a TX, but only one of them will be confirmed.
	pending *shrinkingmap.ShrinkingMap[uint32, map[iotago.SignedTransactionID]*varLocalViewEntry]
	//
	// Limit pipelining (the maximal number of unconfirmed TXes to build).
	// -1 -- infinite, 0 -- disabled, L > 0 -- up to L TXes ahead.
	pipeliningLimit int
	//
	// Callback for the TIP changes.
	outputCB VarLocalViewOutputCB
	output   *varLocalViewOutput
	//
	// Just a logger.
	log *logger.Logger
}

func NewVarLocalView(pipeliningLimit int, outputCB VarLocalViewOutputCB, log *logger.Logger) VarLocalView {
	log.Debugf("NewVarLocalView, pipeliningLimit=%v", pipeliningLimit)
	return &varLocalViewImpl{
		confirmedAnchor:  nil,
		confirmedAccount: nil,
		confirmedCO:      nil,
		pending:          shrinkingmap.New[uint32, map[iotago.SignedTransactionID]*varLocalViewEntry](),
		pipeliningLimit:  pipeliningLimit,
		outputCB:         outputCB,
		output:           nil,
		log:              log,
	}
}

func (lvi *varLocalViewImpl) ConsensusOutputDone(
	logIndex LogIndex,
	consResult *cons.Result,
	eventOutputCB VarLocalViewOutputCB,
) {
	lvi.log.Debugf("ConsensusOutputDone: logIndex=%v, consResult=", logIndex, consResult)
	stateIndex := consResult.ProducedChainOutputs().GetStateIndex()
	if lvi.confirmedCO != nil && lvi.confirmedCO.GetStateIndex() >= stateIndex {
		// We already know it is outdated, so don't add it.
		return
	}

	var pendingForSI map[iotago.SignedTransactionID]*varLocalViewEntry
	pendingForSI, ok := lvi.pending.Get(stateIndex)
	if !ok {
		pendingForSI = map[iotago.SignedTransactionID]*varLocalViewEntry{}
		lvi.pending.Set(stateIndex, pendingForSI)
	}
	txID := consResult.MustSignedTransactionID()
	blID := consResult.MustIotaBlockID()
	entry, ok := pendingForSI[txID]
	if !ok {
		entry = &varLocalViewEntry{
			producedChainOutputs:    consResult.ProducedChainOutputs(),
			producedTransaction:     consResult.ProducedTransaction(),
			consumedAnchorOutputID:  consResult.ConsumedAnchorOutputID(),
			consumedAccountOutputID: consResult.ConsumedAccountOutputID(),
			reuse:                   false, // TODO: Reconsider this field.
			rejected:                false,
			logIndex:                logIndex,
		}
		pendingForSI[txID] = entry
	}
	entry.blocks[blID] = consResult.ProducedIotaBlock()
	lvi.outputIfChanged(eventOutputCB)
}

// A confirmed Anchor/Account output is received from L1. Based on that, we either
// truncate our local history until the received CO (if we know it was posted before),
// or we replace the entire history with an unseen CO (probably produced not by this chain×cmt).
//
// The input here can contain either both - account and anchor outputs, of one of them.
// This is needed to keep the case of both outputs atomic, while supporting out-of-pair
// updates of the outputs.
//
// In the TLA+ spec this function corresponds to:
//   - BothOutputsConfirmed,
//   - AnchorOutputConfirmed,
//   - AccountOutputConfirmed.
func (lvi *varLocalViewImpl) ChainOutputsConfirmed(
	confirmedOutputs *isc.ChainOutputs,
	eventOutputCB VarLocalViewOutputCB,
) LogIndex {
	lvi.confirmedCO = confirmedOutputs
	lvi.log.Debugf("AnchorOutputConfirmed: confirmed=%v", lvi.confirmedCO)

	confirmedLogIndex := NilLogIndex()
	if pending, cnfLI := lvi.isConfirmedPending(); pending {
		confirmedLogIndex = cnfLI
		confirmedStateIndex := lvi.confirmedCO.GetStateIndex()
		lvi.pending.ForEachKey(func(si uint32) bool {
			if si <= confirmedStateIndex {
				lvi.pending.Delete(si)
			}
			return true
		})
	} else {
		lvi.pending.Clear()
	}
	lvi.outputIfChanged(eventOutputCB)
	return confirmedLogIndex
}

// Mark the specified AO as rejected.
// Trim the suffix of rejected AOs.
func (lvi *varLocalViewImpl) ChainOutputsRejected(rejected *isc.ChainOutputs, eventOutputCB VarLocalViewOutputCB) {
	lvi.log.Debugf("AnchorOutputRejected: rejected=%v", rejected)
	stateIndex := rejected.GetStateIndex()
	//
	// Mark the output as rejected, as well as all the outputs depending on it.
	if entries, ok := lvi.pending.Get(stateIndex); ok {
		for _, entry := range entries {
			if entry.producedChainOutputs.Equals(rejected) {
				lvi.log.Debugf("⊳ Entry marked as rejected.")
				entry.rejected = true
				lvi.markDependentAsRejected(rejected)
			}
		}
	}
	lvi.outputIfChanged(eventOutputCB)
}

func (lvi *varLocalViewImpl) BlockExpired(blockID iotago.BlockID, eventOutputCB VarLocalViewOutputCB) {
	found := false
	lvi.pending.ForEach(func(si uint32, es map[iotago.SignedTransactionID]*varLocalViewEntry) bool {
		for _, e := range es {
			if _, ok := e.blocks[blockID]; ok {
				delete(e.blocks, blockID)
				found = true
				break
			}
		}
		return !found
	})
	if found {
		lvi.outputIfChanged(eventOutputCB)
	}
}

func (lvi *varLocalViewImpl) markDependentAsRejected(co *isc.ChainOutputs) {
	accRejected := map[iotago.OutputID]struct{}{co.AnchorOutputID: {}}
	for si := co.GetStateIndex() + 1; ; si++ {
		es, esFound := lvi.pending.Get(si)
		if !esFound {
			break
		}
		for _, e := range es {
			if _, ok := accRejected[e.consumedAnchorOutputID]; ok && !e.rejected {
				lvi.log.Debugf("⊳ Also marking %v as rejected.", e.producedChainOutputs)
				e.rejected = true
				accRejected[e.producedChainOutputs.AnchorOutputID] = struct{}{}
			}
		}
	}
}

func (lvi *varLocalViewImpl) normalizePending() {
	if !lvi.allRejectedOrExpired() || lvi.pending.IsEmpty() {
		return
	}
	if lvi.confirmedCO == nil {
		return
	}
	lvi.log.Debugf("⊳ All entries are rejected or expired, clearing them.")
	//
	// Only keep a prefix of entries forming a continuous chain
	// with no forks nor rejections.
	latestCO := lvi.confirmedCO
	pendingSICount := uint32(lvi.pending.Size())
	remainingPendingEntries := map[iotago.SignedTransactionID]*varLocalViewEntry{}
	for i := uint32(0); i < pendingSICount; i++ {
		nextSIEntry := lvi.nextSinglePendingEntry(latestCO)
		if nextSIEntry == nil {
			// The pending entries don't form a continuous non-forked non-rejected chain.
			break
		}
		if len(nextSIEntry.blocks) == 0 {
			remainingPendingEntries[nextSIEntry.txID()] = nextSIEntry
		} else {
			break
		}
		latestCO = nextSIEntry.producedChainOutputs
	}
	lvi.pending.Clear()
	for txID, e := range remainingPendingEntries {
		e.reuse = true
		lvi.pending.Set(
			e.producedChainOutputs.GetStateIndex(),
			map[iotago.SignedTransactionID]*varLocalViewEntry{txID: e},
		)
	}
}

func (lvi *varLocalViewImpl) allRejectedOrExpired() bool {
	all := true
	lvi.pending.ForEach(func(si uint32, es map[iotago.SignedTransactionID]*varLocalViewEntry) bool {
		for _, e := range es {
			if !e.rejected || len(e.blocks) != 0 {
				all = false
			}
		}
		return all
	})
	return all
}

func (lvi *varLocalViewImpl) outputIfChanged(eventOutputCB VarLocalViewOutputCB) {
	lvi.normalizePending()
	newOutput := lvi.deriveOutput()
	if newOutput == nil && lvi.output == nil {
		return
	}
	if newOutput != nil && lvi.output != nil {
		if newOutput.Equals(lvi.output) {
			return
		}
	}
	lvi.output = newOutput
	lvi.outputCB(newOutput)
	if eventOutputCB != nil {
		eventOutputCB(newOutput)
	}
}

func (lvi *varLocalViewImpl) StatusString() string {
	var tip *isc.ChainOutputs
	if lvi.output != nil {
		tip = lvi.output.baseCO
	}
	return fmt.Sprintf("{varLocalView: tip=%v, |pendingSIs|=%v}", tip, lvi.pending.Size())
}

// This implements TLA+ spec operators: HaveOutput and Output.
// Additionally, the pipelining limit is considered here.
func (lvi *varLocalViewImpl) deriveOutput() *varLocalViewOutput {
	if lvi.confirmedAnchor == nil || lvi.confirmedAccount == nil {
		// Should have a confirmed base.
		return nil
	}
	pendingSICount := uint32(lvi.pending.Size())
	if lvi.pipeliningLimit >= 0 && pendingSICount > uint32(lvi.pipeliningLimit) {
		// pipeliningLimit < 0 ==> no limit on the pipelining.
		// pipeliningLimit = 0 ==> there is no pipelining, we wait each TX to be confirmed first.
		// pipeliningLimit > 0 ==> up to pipeliningLimit TXes can be build unconfirmed.
		return nil
	}
	var reusableEntry *varLocalViewEntry // First reusable TX found.
	var reusableParent *isc.ChainOutputs // Parent outputs of the reusableEntry.
	var latestBlock *iotago.Block        // A block before the proposed TX or CO.
	latestCO := lvi.confirmedCO
	for i := uint32(0); i < pendingSICount; i++ {
		nextSIEntry := lvi.nextSinglePendingEntry(latestCO)
		if nextSIEntry == nil {
			// The pending entries don't form a continuous non-forked non-rejected chain.
			return nil
		}
		if nextSIEntry.isTxReusable() {
			// If this is the first entry that contains a reusable TX, record it.
			if reusableEntry == nil {
				reusableEntry = nextSIEntry
				reusableParent = latestCO
			}
			latestBlock = nil
		} else {
			// If we saw a reusable entry before, but the current is not reusable,
			// we cannot reuse it yet and the chain is not clear. Thus, nothing to propose.
			if reusableEntry != nil {
				return nil
			}
			if len(nextSIEntry.blocks) != 1 {
				return nil
			}
			for _, latestBlock = range nextSIEntry.blocks {
				break // Just take first/single element
			}
		}
		latestCO = nextSIEntry.producedChainOutputs
	}
	if reusableEntry != nil {
		return newVarLocalViewOutput(
			nil,                               // If we are reusing a TX, the parent block will be too old to be a tip.
			reusableParent,                    // Cannot be nil.
			reusableEntry.producedTransaction, // Will contain all the TXes, or none of them. They will form a chain.
		)
	}
	return newVarLocalViewOutput(
		latestBlock, // Can be nil.
		latestCO,    // Cannot be nil.
		nil,         // Will contain all the TXes, or none of them. They will form a chain.
	)
}

func (lvi *varLocalViewImpl) nextSinglePendingEntry(prevCO *isc.ChainOutputs) *varLocalViewEntry {
	prevSI := prevCO.GetStateIndex()
	nextSIEntries, ok := lvi.pending.Get(prevSI + 1)
	if !ok {
		// Should have chain without gaps.
		return nil
	}
	if len(nextSIEntries) != 1 {
		// Should have no pending forks.
		return nil
	}
	var nextSIEntry *varLocalViewEntry
	for _, nextSIEntry = range nextSIEntries {
		break // Just take the first (a single) element
	}
	if nextSIEntry.rejected {
		// Should have no unresolved rejections.
		return nil
	}
	if prevCO.AnchorOutputID != nextSIEntry.consumedAnchorOutputID {
		// Should have chain without gaps.
		return nil
	}
	return nextSIEntry
}

func (lvi *varLocalViewImpl) isConfirmedPending() (bool, LogIndex) {
	found := false
	logIndex := NilLogIndex()
	lvi.pending.ForEach(func(si uint32, es map[iotago.SignedTransactionID]*varLocalViewEntry) bool {
		for _, e := range es {
			if e.producedChainOutputs.Equals(lvi.confirmedCO) {
				found = true
				logIndex = e.logIndex
				break
			}
		}
		return !found
	})
	return found, logIndex
}
