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

	"github.com/samber/lo"

	"github.com/iotaledger/hive.go/ds/shrinkingmap"
	"github.com/iotaledger/hive.go/logger"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/isc"
)

// On the input, the consumed* field are always non-nil; Then the consensus can have the following cases:
//   - The produced* fields are nil -- new TX has to be built on top of consumed*.
//   - The produced* fields non-nil -- existing TX has to be re-issued.
type ConsensusInputEntry struct {
	producedChainOutputs    *isc.ChainOutputs   // The produced chain outputs.
	producedTransaction     *iotago.Transaction // The transaction publishing the chain outputs.
	consumedAnchorOutputID  iotago.OutputID     // Consumed in the TX.
	consumedAccountOutputID iotago.OutputID     // Consumed in the TX.
}

// On the output from the consensus, all the fields have to be non-nil.
// The output is 0 or more such entries. It will be for a particular log index.
type ConsensusOutputEntry struct {
	producedChainOutputs    *isc.ChainOutputs   // The produced chain outputs.
	producedTransaction     *iotago.Transaction // The transaction publishing the chain outputs.
	producedBlock           *iotago.Block       // Block produced to publish the TX.
	consumedAnchorOutputID  iotago.OutputID     // Consumed in the TX.
	consumedAccountOutputID iotago.OutputID     // Consumed in the TX.
}

// Block might be nil, so check it before calling this.
func (coe *ConsensusOutputEntry) MustBlockID() iotago.BlockID {
	blockID, err := coe.producedBlock.ID()
	if err != nil {
		panic(fmt.Errorf("failed to get BlockID: %v", err))
	}
	return blockID
}

// Transaction will always be set, so it should be safe to call this.
func (coe *ConsensusOutputEntry) MustTransactionID() iotago.TransactionID {
	txID, err := coe.producedTransaction.ID()
	if err != nil {
		panic(fmt.Errorf("failed to get TX ID: %v", err))
	}
	return txID
}

type VarLocalView interface {
	//
	// Corresponds to the `tx_posted` event in the specification.
	// Returns true, if the proposed BaseAnchorOutput has changed.
	ConsensusOutputDone(
		logIndex LogIndex,
		consOutEntries []*ConsensusOutputEntry,
	)
	//
	// Called upon receiving confirmation from the L1.
	// In a normal scenario both (Anchor/Account) outputs will be confirmed together,
	// because they are in the same TX. If someone moved one of those outputs externally,
	// they can be moved independently. In such a case, one of the anc/acc parameters will be nil.
	//   - It is important to get these events in the correct order, otherwise the out-of-order
	//     events will be considered as an reorg.
	ChainOutputsConfirmed(
		confirmedAnchor *iotago.AnchorOutput,
		confirmedAnchorID *iotago.OutputID,
		confirmedAccount *iotago.AccountOutput,
		confirmedAccountID *iotago.OutputID,
	)
	//
	// Called when the TX containing the specified outputs was rejected.
	// The outputs cannot be rejected independently because they are in the same TX.
	ChainOutputsRejected(rejected *isc.ChainOutputs)
	//
	//
	BlockExpired(blockID *iotago.BlockID)
	//
	// Support functions.
	StatusString() string
}

// This is the result of the chain tip tracking.
// Here we decide the latest block to build on,
// optionally a block to use as a tip and
// a list of transactions that should be resubmitted
// (by producing and signing new blocks).
type VarLocalViewOutput interface {
	BaseCO() *isc.ChainOutputs           // Will always be non-nill.
	BaseBlock() *iotago.Block            // Can be nil.
	ReattachTXes() []*iotago.Transaction // Can be empty.
}

type varLocalViewOutput struct {
	baseCO       *isc.ChainOutputs
	baseBlock    *iotago.Block
	reattachTXes []*iotago.Transaction
}

func newVarLocalViewOutput(
	baseCO *isc.ChainOutputs,
	baseBlock *iotago.Block,
	reattachTXes []*iotago.Transaction,
) *varLocalViewOutput {
	return &varLocalViewOutput{
		baseCO:       baseCO,
		baseBlock:    baseBlock,
		reattachTXes: reattachTXes,
	}
}

func (o *varLocalViewOutput) BaseCO() *isc.ChainOutputs           { return o.baseCO }
func (o *varLocalViewOutput) BaseBlock() *iotago.Block            { return o.baseBlock }
func (o *varLocalViewOutput) ReattachTXes() []*iotago.Transaction { return o.reattachTXes }
func (o *varLocalViewOutput) Equals(other *varLocalViewOutput) bool {
	if !o.baseCO.Equals(o.baseCO) {
		return false
	}
	if (o.baseBlock == nil && other.baseBlock != nil) || (o.baseBlock != nil && other.baseBlock == nil) {
		return false
	}
	if o.baseBlock != nil && other.baseBlock != nil && o.baseBlock.MustID() != other.baseBlock.MustID() {
		return false
	}
	if len(o.reattachTXes) != len(other.reattachTXes) {
		return false
	}
	for i := range o.reattachTXes {
		id1, err1 := o.reattachTXes[i].ID()
		id2, err2 := other.reattachTXes[i].ID()
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
	consOutEntry *ConsensusOutputEntry // The data from the consensus.
	blockExpired bool                  // True, if the block has been expired.
	rejected     bool                  // True, if the AO as rejected. We keep them to detect the other rejected AOs.
	logIndex     LogIndex              // LogIndex of the consensus produced the output, if any.
}

type varLocalViewImpl struct {
	//
	// The latest confirmed CO, as received from L1.
	// All the pending entries are built on top of this one.
	// It can be nil, if the latest AO is unclear (either not received yet).
	//
	// We don't use the isc.ChainOutputs structure here, because we can
	// receive the anchor/account outputs separately.
	confirmedAnchor    *iotago.AnchorOutput
	confirmedAnchorID  *iotago.OutputID
	confirmedAccount   *iotago.AccountOutput
	confirmedAccountID *iotago.OutputID
	confirmed          *isc.ChainOutputs // Derived from the above, when all of them are received.
	//
	// AOs and blocks produced by this committee, but not confirmed yet.
	// It is possible to have several AOs for a StateIndex in the case of
	// Recovery/Timeout notices. Then the next consensus is started o build a TX.
	// Both of them can still produce a TX, but only one of them will be confirmed.
	pending *shrinkingmap.ShrinkingMap[uint32, []*varLocalViewEntry]
	//
	// Limit pipelining (the maximal number of unconfirmed TXes to build).
	// -1 -- infinite, 0 -- disabled, L > 0 -- up to L TXes ahead.
	pipeliningLimit int
	//
	// Callback for the TIP changes.
	outputCB func(out *varLocalViewOutput)
	output   *varLocalViewOutput
	//
	// Just a logger.
	log *logger.Logger
}

func NewVarLocalView(pipeliningLimit int, outputCB func(out *varLocalViewOutput), log *logger.Logger) VarLocalView {
	log.Debugf("NewVarLocalView, pipeliningLimit=%v", pipeliningLimit)
	return &varLocalViewImpl{
		confirmedAnchor:    nil,
		confirmedAnchorID:  nil,
		confirmedAccount:   nil,
		confirmedAccountID: nil,
		confirmed:          nil,
		pending:            shrinkingmap.New[uint32, []*varLocalViewEntry](),
		pipeliningLimit:    pipeliningLimit,
		outputCB:           outputCB,
		output:             nil,
		log:                log,
	}
}

func (lvi *varLocalViewImpl) ConsensusOutputDone(
	logIndex LogIndex,
	consOutEntries []*ConsensusOutputEntry,
) {
	lvi.log.Debugf("ConsensusOutputDone: logIndex=%v, |consOutEntries|=", logIndex, len(consOutEntries))
	for i := range consOutEntries {
		consOutEntry := consOutEntries[i]
		stateIndex := consOutEntry.producedChainOutputs.GetStateIndex()
		if lvi.confirmed != nil && lvi.confirmed.GetStateIndex() >= stateIndex {
			// We already know it is outdated, so don't add it.
			continue
		}

		var pendingForSI []*varLocalViewEntry
		pendingForSI, ok := lvi.pending.Get(stateIndex)
		if !ok {
			pendingForSI = []*varLocalViewEntry{}
		}
		var replaced = false
		pendingForSI = lo.FilterMap(pendingForSI, func(e *varLocalViewEntry, _ int) (*varLocalViewEntry, bool) {
			if e.blockExpired && e.consOutEntry.MustTransactionID() == consOutEntry.MustTransactionID() {
				if replaced {
					// Already replaced the existing block. Just remove remaining expired blocks for that TX.
					return nil, false
				}
				// Existing block as expired, but now we publish the same TX again.
				replaced = true
				e.blockExpired = false
				e.consOutEntry = consOutEntry
				return e, true
			}
			if e.consOutEntry.MustBlockID() == consOutEntry.MustBlockID() {
				// This is a duplicate, just ignore it.
				replaced = true
				return e, true
			}
			// Keep other elements as it.
			return e, true
		})
		if !replaced {
			// This is new entry, append it.
			pendingForSI = append(pendingForSI, &varLocalViewEntry{
				consOutEntry: consOutEntry,
				blockExpired: false,
				rejected:     false,
				logIndex:     logIndex,
			})
		}
		lvi.pending.Set(stateIndex, pendingForSI)
	}
	lvi.outputIfChanged()
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
	confirmedAnchor *iotago.AnchorOutput,
	confirmedAnchorID *iotago.OutputID,
	confirmedAccount *iotago.AccountOutput,
	confirmedAccountID *iotago.OutputID,
) {
	if confirmedAnchor != nil {
		lvi.confirmedAnchor = confirmedAnchor
		lvi.confirmedAnchorID = confirmedAnchorID
	}
	if confirmedAccount != nil {
		lvi.confirmedAccount = confirmedAccount
		lvi.confirmedAccountID = confirmedAccountID
	}
	if lvi.confirmedAnchor == nil || lvi.confirmedAccount == nil {
		// Have no both outputs confirmed yet, wait longer.
		lvi.log.Debugf("AnchorOutputConfirmed: confirmed=%v", lvi.confirmed)
		return
	}
	lvi.confirmed = isc.NewChainOutputs(
		lvi.confirmedAnchor,
		*lvi.confirmedAnchorID,
		lvi.confirmedAccount,
		*lvi.confirmedAccountID,
	)
	lvi.log.Debugf("AnchorOutputConfirmed: confirmed=%v", lvi.confirmed)

	if lvi.isConfirmedPending() {
		confirmedStateIndex := lvi.confirmed.GetStateIndex()
		lvi.pending.ForEachKey(func(si uint32) bool {
			if si <= confirmedStateIndex {
				lvi.pending.Delete(si)
			}
			return true
		})
	} else {
		lvi.pending.Clear()
	}
	lvi.outputIfChanged()
}

// Mark the specified AO as rejected.
// Trim the suffix of rejected AOs.
func (lvi *varLocalViewImpl) ChainOutputsRejected(rejected *isc.ChainOutputs) {
	lvi.log.Debugf("AnchorOutputRejected: rejected=%v", rejected)
	stateIndex := rejected.GetStateIndex()
	//
	// Mark the output as rejected, as well as all the outputs depending on it.
	if entries, ok := lvi.pending.Get(stateIndex); ok {
		for _, entry := range entries {
			if entry.consOutEntry.producedChainOutputs.Equals(rejected) {
				lvi.log.Debugf("⊳ Entry marked as rejected.")
				entry.rejected = true
				lvi.markDependentAsRejected(rejected)
			}
		}
	}
	lvi.outputIfChanged()
}

func (lvi *varLocalViewImpl) BlockExpired(blockID *iotago.BlockID) {
	lvi.pending.ForEach(func(si uint32, es []*varLocalViewEntry) bool {
		found := false
		for i := range es {
			// Mark it expired, and if the re are other entries for the same
			// transaction, then remove it altogether.
			if es[i].consOutEntry.producedBlock.MustID() == *blockID {
				found = true
				es[i].blockExpired = true
				for j := range es {
					if i != j && es[j].consOutEntry.MustTransactionID() == es[i].consOutEntry.MustTransactionID() {
						es = append(es[:i], es[i+1:]...)
					}
				}
			}
		}
		if found {
			lvi.pending.Set(si, es)
		}
		return !found
	})
}

func (lvi *varLocalViewImpl) markDependentAsRejected(co *isc.ChainOutputs) {
	accRejected := map[iotago.OutputID]struct{}{co.AnchorOutputID: {}}
	for si := co.GetStateIndex() + 1; ; si++ {
		es, esFound := lvi.pending.Get(si)
		if !esFound {
			break
		}
		for _, e := range es {
			if _, ok := accRejected[e.consOutEntry.consumedAnchorOutputID]; ok && !e.rejected {
				lvi.log.Debugf("⊳ Also marking %v as rejected.", e.consOutEntry.producedChainOutputs)
				e.rejected = true
				accRejected[e.consOutEntry.producedChainOutputs.AnchorOutputID] = struct{}{}
			}
		}
	}
}

func (lvi *varLocalViewImpl) normalizePending() {
	if !lvi.allRejectedOrExpired() || lvi.pending.IsEmpty() {
		return
	}
	if lvi.confirmed == nil {
		return
	}
	lvi.log.Debugf("⊳ All entries are rejected or expired, clearing them.")
	//
	// Only keep a prefix of entries forming a continuous chain
	// with no forks nor rejections.
	latestCO := lvi.confirmed
	confirmedSI := lvi.confirmedAnchor.StateIndex
	pendingSICount := uint32(lvi.pending.Size())
	remainingPendingEntries := []*varLocalViewEntry{}
	for i := uint32(0); i < pendingSICount; i++ {
		siEntries, ok := lvi.pending.Get(confirmedSI + i + 1)
		if !ok {
			// Should have chain without gaps.
			break
		}
		if len(siEntries) != 1 {
			// Should have no pending forks.
			break
		}
		siEntry := siEntries[0]
		if siEntry.rejected {
			// Should have no unresolved rejections.
			break
		}
		if latestCO.AnchorOutputID != siEntry.consOutEntry.consumedAnchorOutputID {
			// Should have chain without gaps.
			break
		}
		latestCO = siEntry.consOutEntry.producedChainOutputs
		if siEntry.blockExpired {
			remainingPendingEntries = append(remainingPendingEntries, siEntry)
		} else {
			break
		}
	}
	lvi.pending.Clear()
	for _, e := range remainingPendingEntries {
		lvi.pending.Set(
			e.consOutEntry.producedChainOutputs.GetStateIndex(),
			[]*varLocalViewEntry{e},
		)
	}
}

func (lvi *varLocalViewImpl) allRejectedOrExpired() bool {
	all := true
	lvi.pending.ForEach(func(si uint32, es []*varLocalViewEntry) bool {
		containsPending := lo.ContainsBy(es, func(e *varLocalViewEntry) bool {
			return !e.rejected && !e.blockExpired
		})
		all = !containsPending
		return !containsPending
	})
	return all
}

func (lvi *varLocalViewImpl) outputIfChanged() {
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
}

func (lvi *varLocalViewImpl) StatusString() string {
	var tip *isc.ChainOutputs = nil
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
	latestCO := lvi.confirmed
	var latestBlock *iotago.Block = nil
	confirmedSI := lvi.confirmedAnchor.StateIndex
	pendingSICount := uint32(lvi.pending.Size())
	if lvi.pipeliningLimit >= 0 && pendingSICount > uint32(lvi.pipeliningLimit) {
		// pipeliningLimit < 0 ==> no limit on the pipelining.
		// pipeliningLimit = 0 ==> there is no pipelining, we wait each TX to be confirmed first.
		// pipeliningLimit > 0 ==> up to pipeliningLimit TXes can be build unconfirmed.
		return nil
	}
	withoutBlocks := []*iotago.Transaction{}
	for i := uint32(0); i < pendingSICount; i++ {
		siEntries, ok := lvi.pending.Get(confirmedSI + i + 1)
		if !ok {
			// Should have chain without gaps.
			return nil
		}
		if len(siEntries) != 1 {
			// Should have no pending forks.
			return nil
		}
		siEntry := siEntries[0]
		if siEntry.rejected {
			// Should have no unresolved rejections.
			return nil
		}
		if latestCO.AnchorOutputID != siEntry.consOutEntry.consumedAnchorOutputID {
			// Should have chain without gaps.
			return nil
		}
		latestCO = siEntry.consOutEntry.producedChainOutputs
		if siEntry.blockExpired {
			latestBlock = nil
			withoutBlocks = append(withoutBlocks, siEntry.consOutEntry.producedTransaction)
		} else {
			latestBlock = siEntry.consOutEntry.producedBlock
		}
	}
	withoutBlocksCount := len(withoutBlocks)
	if 0 < withoutBlocksCount && withoutBlocksCount < int(pendingSICount) {
		// All or none of the blocks have to have blocks expired.
		return nil
	}
	return newVarLocalViewOutput(
		latestCO,      // Cannot be nil.
		latestBlock,   // Can be nil.
		withoutBlocks, // Will contain all the TXes, or none of them. They will form a chain.
	)
}

func (lvi *varLocalViewImpl) isConfirmedPending() bool {
	found := false
	lvi.pending.ForEach(func(si uint32, es []*varLocalViewEntry) bool {
		found = lo.ContainsBy(es, func(e *varLocalViewEntry) bool {
			return e.consOutEntry.producedChainOutputs.Equals(lvi.confirmed)
		})
		return !found
	})
	return found
}
