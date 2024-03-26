package cons

import (
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/gpa"
)

// Interaction with the NodeConnection.
// To get the tip proposals.
type SyncNC interface {
	// TXCreateInputReceived(baseCO *isc.ChainOutputs, blockToRefer *iotago.Block) gpa.OutMessages
	// BlockOnlyInputReceived(txToPublish *iotago.SignedTransaction, blockToRefer *iotago.Block) gpa.OutMessages

	// StateMgrProposalReceived() gpa.OutMessages
	// MempoolProposalReceived() gpa.OutMessages
	// DSStIndexProposalReceived() gpa.OutMessages
	// DSSbIndexProposalReceived() gpa.OutMessages
	// TimeUpdateReceived() gpa.OutMessages

	BlockTipSetNeeded() gpa.OutMessages
	BlockTipSetReceived(strongParents iotago.BlockIDs) gpa.OutMessages
}

type syncNCImpl struct {
	blockTipSetNeededCB   func() gpa.OutMessages
	blockTipSetReceivedCB func(strongParents iotago.BlockIDs) gpa.OutMessages
	// blockToRefer           *iotago.Block // Optional, can be left nil. Set along with tx or block inputs.
	// txCreateInputReceived  *isc.ChainOutputs
	// blockOnlyInputReceived *iotago.SignedTransaction

	// stateMgrProposalReceived bool
	// mempoolProposalReceived  []*isc.RequestRef

	// dssTIndexProposal []int // Index proposals from the DSS for signing the TX.
	// dssBIndexProposal []int // Index proposals from the DSS for signing the Block.
	// timeData          time.Time

	// inputsReady   bool
	// inputsReadyCB func() gpa.OutMessages

	// // Output is all the inputs plus the tip proposal.
	// strongParents iotago.BlockIDs
	// outputReadyCB func(
	// 	blockToRefer *iotago.Block,
	// 	txCreateInputReceived *isc.ChainOutputs,
	// 	blockOnlyInputReceived *iotago.SignedTransaction,
	// 	mempoolProposalReceived []*isc.RequestRef,
	// 	dssTIndexProposal []int,
	// 	dssBIndexProposal []int,
	// 	timeData time.Time,
	// 	strongParents iotago.BlockIDs,
	// ) gpa.OutMessages
}

func NewSyncNC(
	blockTipSetNeededCB func() gpa.OutMessages,
	blockTipSetReceivedCB func(strongParents iotago.BlockIDs) gpa.OutMessages,
) SyncNC {
	return &syncNCImpl{
		blockTipSetNeededCB:   blockTipSetNeededCB,
		blockTipSetReceivedCB: blockTipSetReceivedCB,
	}
}

func (s *syncNCImpl) BlockTipSetNeeded() gpa.OutMessages {
	if s.blockTipSetNeededCB == nil {
		return nil // Already done.
	}
	cb := s.blockTipSetNeededCB
	s.blockTipSetNeededCB = nil
	return cb()
}

func (s *syncNCImpl) BlockTipSetReceived(strongParents iotago.BlockIDs) gpa.OutMessages {
	if s.blockTipSetReceivedCB == nil {
		return nil // Already done.
	}
	cb := s.blockTipSetReceivedCB
	s.blockTipSetReceivedCB = nil
	return cb(strongParents)
}

// func NewSyncNC(
// 	inputsReadyCB func() gpa.OutMessages,
// 	outputReadyCB func(
// 		blockToRefer *iotago.Block,
// 		txCreateInputReceived *isc.ChainOutputs,
// 		blockOnlyInputReceived *iotago.SignedTransaction,
// 		mempoolProposalReceived []*isc.RequestRef,
// 		dssTIndexProposal []int,
// 		dssBIndexProposal []int,
// 		timeData time.Time,
// 		strongParents iotago.BlockIDs,
// 	) gpa.OutMessages,
// ) SyncNC {
// 	return &syncNCImpl{
// 		inputsReadyCB: inputsReadyCB,
// 		outputReadyCB: outputReadyCB,
// 	}
// }

// func (s *syncNCImpl) TXCreateInputReceived(baseCO *isc.ChainOutputs, blockToRefer *iotago.Block) gpa.OutMessages {
// 	if s.txCreateInputReceived != nil || s.blockOnlyInputReceived != nil {
// 		return nil
// 	}
// 	s.txCreateInputReceived = baseCO
// 	s.blockToRefer = blockToRefer
// 	return s.tryInputsReady()
// }

// func (s *syncNCImpl) BlockOnlyInputReceived(txToPublish *iotago.SignedTransaction, blockToRefer *iotago.Block) gpa.OutMessages {
// 	if s.txCreateInputReceived != nil || s.blockOnlyInputReceived != nil {
// 		return nil
// 	}
// 	s.blockOnlyInputReceived = txToPublish
// 	s.blockToRefer = blockToRefer
// 	return s.tryInputsReady()
// }

// func (s *syncNCImpl) StateMgrProposalReceived() gpa.OutMessages {
// 	if s.stateMgrProposalReceived {
// 		return nil
// 	}
// 	s.stateMgrProposalReceived = true
// 	return s.tryInputsReady()
// }

// func (s *syncNCImpl) MempoolProposalReceived(requestRefs []*isc.RequestRef) gpa.OutMessages {
// 	if s.mempoolProposalReceived != nil {
// 		return nil
// 	}
// 	s.mempoolProposalReceived = requestRefs
// 	return s.tryInputsReady()
// }

// func (s *syncNCImpl) DSStIndexProposalReceived(indexProposal []int) gpa.OutMessages {
// 	if s.dssTIndexProposal != nil {
// 		return nil
// 	}
// 	s.dssTIndexProposal = indexProposal
// 	return s.tryInputsReady()
// }

// func (s *syncNCImpl) DSSbIndexProposalReceived(indexProposal []int) gpa.OutMessages {
// 	if s.dssBIndexProposal != nil {
// 		return nil
// 	}
// 	s.dssBIndexProposal = indexProposal
// 	return s.tryInputsReady()
// }

// func (s *syncNCImpl) TimeUpdateReceived(timeData time.Time) gpa.OutMessages {
// 	if timeData.Before(s.timeData) || timeData.Equal(s.timeData) {
// 		return nil
// 	}
// 	s.timeData = timeData
// 	return s.tryInputsReady()
// }

// func (s *syncNCImpl) BlockTipSetReceived(strongParents iotago.BlockIDs) gpa.OutMessages {
// 	if strongParents == nil {
// 		panic(fmt.Errorf("nil as strongParents in cons.sync_nc."))
// 	}
// 	if s.strongParents != nil {
// 		return nil // Received already.
// 	}
// 	if !s.inputsReady {
// 		return nil // Too early.
// 	}
// 	s.strongParents = strongParents
// 	return s.outputReadyCB(
// 		s.blockToRefer,
// 		s.txCreateInputReceived,
// 		s.blockOnlyInputReceived,
// 		s.mempoolProposalReceived,
// 		s.dssTIndexProposal,
// 		s.dssBIndexProposal,
// 		s.timeData,
// 		s.strongParents,
// 	)
// }

// func (s *syncNCImpl) tryInputsReady() gpa.OutMessages {
// 	if s.inputsReady {
// 		return nil // Done already.
// 	}
// 	if s.txCreateInputReceived == nil && s.blockOnlyInputReceived == nil {
// 		return nil // At least one of these is required.
// 	}
// 	if s.txCreateInputReceived != nil {
// 		if !s.stateMgrProposalReceived || s.mempoolProposalReceived == nil {
// 			return nil // Mempool and StateMgr are needed if we are going to build a TX.
// 		}
// 	}
// 	if s.dssTIndexProposal == nil || s.dssBIndexProposal == nil || s.timeData.IsZero() {
// 		return nil // These are required in any case.
// 	}
// 	s.inputsReady = true
// 	return s.inputsReadyCB()
// }

// // Try to provide useful human-readable compact status.
// func (s *syncNCImpl) String() string {
// 	str := "NC"
// 	// if sub.indexProposalReady && sub.outputReady { // TODO: ...
// 	// 	return str + statusStrOK
// 	// }
// 	// if sub.indexProposalReady {
// 	// 	str += "/idx=OK"
// 	// } else {
// 	// 	str += fmt.Sprintf("/idx[initialInputsReady=%v,indexProposalReady=%v]", sub.initialInputsReady, sub.indexProposalReady)
// 	// }
// 	// if sub.outputReady {
// 	// 	str += "/sig=OK"
// 	// } else if sub.signingInputsReady {
// 	// 	str += "/sig[WaitingForDSS]"
// 	// } else {
// 	// 	wait := []string{}
// 	// 	if sub.MessageToSign == nil {
// 	// 		wait = append(wait, "MessageToSign")
// 	// 	}
// 	// 	if sub.DecidedIndexProposals == nil {
// 	// 		wait = append(wait, "DecidedIndexProposals")
// 	// 	}
// 	// 	str += fmt.Sprintf("/sig=WAIT[%v]", strings.Join(wait, ","))
// 	// }
// 	return str
// }
