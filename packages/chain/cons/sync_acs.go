// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"fmt"
	"strings"
	"time"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/gpa/acs"
	"github.com/iotaledger/wasp/packages/isc"
)

type SyncACS interface {
	TXCreateInputReceived(baseCO *isc.ChainOutputs, blockToRefer *iotago.Block) gpa.OutMessages
	BlockOnlyInputReceived(txToPublish *iotago.SignedTransaction, blockToRefer *iotago.Block) gpa.OutMessages

	MempoolRequestsReceived(requestRefs []*isc.RequestRef) gpa.OutMessages
	StateMgrProposalReceived(proposedBaseAnchorOutput *isc.ChainOutputs) gpa.OutMessages
	DSStIndexProposalReceived(dssIndexProposal []int) gpa.OutMessages
	DSSbIndexProposalReceived(dssIndexProposal []int) gpa.OutMessages
	TimeUpdateReceived(timeData time.Time) gpa.OutMessages
	BlockTipSetProposalReceived(strongParents iotago.BlockIDs) gpa.OutMessages
	// ACSInputsReceived(
	// 	blockToRefer *iotago.Block,
	// 	txCreateInputReceived *isc.ChainOutputs,
	// 	blockOnlyInputReceived *iotago.SignedTransaction,
	// 	mempoolProposalReceived []*isc.RequestRef,
	// 	dssTIndexProposal []int,
	// 	dssBIndexProposal []int,
	// 	timeData time.Time,
	// 	strongParents iotago.BlockIDs,
	// ) gpa.OutMessages
	ACSOutputReceived(output gpa.Output) gpa.OutMessages
	String() string
}

type SyncACSBlockTipsNeededCB = func() gpa.OutMessages

type SyncACSInputsReadyCB = func(
	blockToRefer *iotago.Block,
	txCreateInputReceived *isc.ChainOutputs,
	blockOnlyInputReceived *iotago.SignedTransaction,
	mempoolProposalReceived []*isc.RequestRef,
	dssTIndexProposal []int,
	dssBIndexProposal []int,
	timeData time.Time,
	strongParents iotago.BlockIDs,
) gpa.OutMessages

type SyncACSOutputsReadyCB = func(
	output map[gpa.NodeID][]byte,
) gpa.OutMessages

// > UPON Reception of responses from Mempool, StateMgr and DSS NonceIndexes:
// >     Produce a batch proposal.
// >     Start the ACS.
type syncACSImpl struct {
	blockToRefer           *iotago.Block
	txCreateInputReceived  *isc.ChainOutputs
	blockOnlyInputReceived *iotago.SignedTransaction

	stateMgrProposalReceived *isc.ChainOutputs // Should be same as txCreateInputReceived
	mempoolProposalReceived  []*isc.RequestRef
	dssTIndexProposal        []int // Index proposals from the DSS for signing the TX.
	dssBIndexProposal        []int // Index proposals from the DSS for signing the Block.
	timeData                 time.Time
	strongParents            iotago.BlockIDs

	blockTipsNeededCB SyncACSBlockTipsNeededCB
	// inputsReady       bool
	inputsReadyCB SyncACSInputsReadyCB
	outputReady   bool
	outputReadyCB SyncACSOutputsReadyCB
	terminated    bool
	terminatedCB  func()
}

func NewSyncACS(
	blockTipsNeededCB SyncACSBlockTipsNeededCB,
	inputsReadyCB SyncACSInputsReadyCB,
	outputReadyCB SyncACSOutputsReadyCB,
	terminatedCB func(),
) SyncACS {
	return &syncACSImpl{
		blockTipsNeededCB: blockTipsNeededCB,
		inputsReadyCB:     inputsReadyCB,
		outputReadyCB:     outputReadyCB,
		terminatedCB:      terminatedCB,
	}
}

func (sub *syncACSImpl) TXCreateInputReceived(baseCO *isc.ChainOutputs, blockToRefer *iotago.Block) gpa.OutMessages {
	if sub.txCreateInputReceived != nil || sub.blockOnlyInputReceived != nil {
		return nil
	}
	sub.txCreateInputReceived = baseCO
	sub.blockToRefer = blockToRefer
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) BlockOnlyInputReceived(txToPublish *iotago.SignedTransaction, blockToRefer *iotago.Block) gpa.OutMessages {
	if sub.txCreateInputReceived != nil || sub.blockOnlyInputReceived != nil {
		return nil
	}
	sub.blockOnlyInputReceived = txToPublish
	sub.blockToRefer = blockToRefer
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) StateMgrProposalReceived(baseCO *isc.ChainOutputs) gpa.OutMessages {
	if sub.stateMgrProposalReceived != nil {
		return nil
	}
	sub.stateMgrProposalReceived = baseCO
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) MempoolRequestsReceived(requestRefs []*isc.RequestRef) gpa.OutMessages {
	if sub.mempoolProposalReceived != nil {
		return nil
	}
	sub.mempoolProposalReceived = requestRefs
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) DSStIndexProposalReceived(indexProposal []int) gpa.OutMessages {
	if sub.dssTIndexProposal != nil {
		return nil
	}
	sub.dssTIndexProposal = indexProposal
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) DSSbIndexProposalReceived(indexProposal []int) gpa.OutMessages {
	if sub.dssBIndexProposal != nil {
		return nil
	}
	sub.dssBIndexProposal = indexProposal
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) TimeUpdateReceived(timeData time.Time) gpa.OutMessages {
	if timeData.After(sub.timeData) {
		sub.timeData = timeData
		return sub.tryCompleteInput()
	}
	return nil
}

func (sub *syncACSImpl) BlockTipSetProposalReceived(strongParents iotago.BlockIDs) gpa.OutMessages {
	if sub.strongParents != nil {
		return nil // Already.
	}
	sub.strongParents = strongParents
	return sub.tryCompleteInput()
}

func (sub *syncACSImpl) tryCompleteInput() gpa.OutMessages {
	if sub.inputsReadyCB == nil {
		return nil // Done already.
	}
	if sub.txCreateInputReceived == nil && sub.blockOnlyInputReceived == nil {
		return nil // At least one of these is required.
	}
	if sub.txCreateInputReceived != nil {
		if sub.stateMgrProposalReceived == nil || sub.mempoolProposalReceived == nil {
			return nil // Mempool and StateMgr are needed if we are going to build a TX.
		}
	}
	if sub.dssTIndexProposal == nil || sub.dssBIndexProposal == nil || sub.timeData.IsZero() {
		return nil // These are required in any case.
	}

	msgs := gpa.NoMessages()
	if sub.blockTipsNeededCB != nil {
		msgs.AddAll(sub.blockTipsNeededCB())
		sub.blockTipsNeededCB = nil
	}

	if sub.strongParents == nil {
		return msgs
	}

	cb := sub.inputsReadyCB
	sub.inputsReadyCB = nil
	return msgs.AddAll(cb(
		sub.blockToRefer,
		sub.txCreateInputReceived,
		sub.blockOnlyInputReceived,
		sub.mempoolProposalReceived,
		sub.dssTIndexProposal,
		sub.dssBIndexProposal,
		sub.timeData,
		sub.strongParents,
	))
}

func (sub *syncACSImpl) ACSOutputReceived(output gpa.Output) gpa.OutMessages {
	if output == nil {
		return nil
	}
	acsOutput, ok := output.(*acs.Output)
	if !ok {
		panic(fmt.Errorf("acs returned unexpected output: %v", output))
	}
	if !sub.terminated && acsOutput.Terminated {
		sub.terminated = true
		sub.terminatedCB()
	}
	if sub.outputReady {
		return nil
	}
	sub.outputReady = true
	return sub.outputReadyCB(acsOutput.Values)
}

// Try to provide useful human-readable compact status.
func (sub *syncACSImpl) String() string {
	str := "ACS"
	if sub.outputReady {
		str += statusStrOK
	} else if sub.inputsReadyCB == nil {
		str += "/WAIT[ACS to complete]"
	} else {
		wait := []string{}
		if sub.stateMgrProposalReceived == nil {
			wait = append(wait, "BaseAnchorOutput")
		}
		if sub.mempoolProposalReceived == nil {
			wait = append(wait, "RequestRefs")
		}
		if sub.dssTIndexProposal == nil {
			wait = append(wait, "DSStIndexProposal")
		}
		if sub.dssBIndexProposal == nil {
			wait = append(wait, "DSSbIndexProposal")
		}
		if sub.timeData.IsZero() {
			wait = append(wait, "TimeData")
		}
		str += fmt.Sprintf("/WAIT[%v]", strings.Join(wait, ","))
	}
	return str
}
