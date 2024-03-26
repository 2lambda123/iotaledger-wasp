package cons

import (
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/state"
)

type SyncRes interface {
	ReuseTX(tx *iotago.SignedTransaction) gpa.OutMessages
	BuiltTX(tx *iotago.SignedTransaction) gpa.OutMessages
	HaveTransition(
		producedChainOutputs *isc.ChainOutputs,
		consumedAnchorOutputID iotago.OutputID,
		consumedAccountOutputID iotago.OutputID,
	) gpa.OutMessages
	HaveStateBlock(producedStateBlock state.Block) gpa.OutMessages
	HaveIotaBlock(producedIotaBlock *iotago.Block) gpa.OutMessages
}

type syncResCB = func(
	transactionReused bool,
	transaction *iotago.SignedTransaction,
	producedIotaBlock *iotago.Block,
	producedChainOutputs *isc.ChainOutputs,
	producedStateBlock state.Block,
	consumedAnchorOutputID iotago.OutputID,
	consumedAccountOutputID iotago.OutputID,
) gpa.OutMessages

type syncRes struct {
	readyCB syncResCB

	transactionReused   bool
	transaction         *iotago.SignedTransaction
	transactionReceived bool

	producedChainOutputs    *isc.ChainOutputs
	producedStateBlock      state.Block
	consumedAnchorOutputID  iotago.OutputID
	consumedAccountOutputID iotago.OutputID
	transitionReceived      bool

	producedIotaBlock *iotago.Block
}

func NewSyncRes(readyCB syncResCB) SyncRes {
	return &syncRes{readyCB: readyCB}
}

func (s *syncRes) ReuseTX(tx *iotago.SignedTransaction) gpa.OutMessages {
	if s.transactionReceived {
		return nil // Already received.
	}
	if s.transitionReceived {
		panic("transition received, but wer are going to reuse the TX.")
	}
	s.transactionReceived = true
	s.transactionReused = true
	s.transaction = tx
	return s.tryOutput()
}

func (s *syncRes) BuiltTX(tx *iotago.SignedTransaction) gpa.OutMessages {
	if s.transactionReceived {
		return nil // Already received.
	}
	s.transactionReceived = true
	s.transactionReused = false
	s.transaction = tx
	return s.tryOutput()
}

func (s *syncRes) HaveTransition(
	producedChainOutputs *isc.ChainOutputs,
	consumedAnchorOutputID iotago.OutputID,
	consumedAccountOutputID iotago.OutputID,
) gpa.OutMessages {
	if s.transactionReused {
		panic("transaction is reused but we received the transition")
	}
	if s.transitionReceived {
		return nil // Already received.
	}
	s.transitionReceived = true
	s.producedChainOutputs = producedChainOutputs
	s.consumedAnchorOutputID = consumedAnchorOutputID
	s.consumedAccountOutputID = consumedAccountOutputID
	return s.tryOutput()
}

func (s *syncRes) HaveStateBlock(producedStateBlock state.Block) gpa.OutMessages {
	if s.transactionReused {
		panic("transaction is reused but we received the transition")
	}
	if s.producedStateBlock != nil {
		panic("state block already received")
	}
	s.producedStateBlock = producedStateBlock
	return s.tryOutput()
}

func (s *syncRes) HaveIotaBlock(producedIotaBlock *iotago.Block) gpa.OutMessages {
	if s.producedIotaBlock != nil {
		return nil // Have already.
	}
	s.producedIotaBlock = producedIotaBlock
	return s.tryOutput()
}

func (s *syncRes) tryOutput() gpa.OutMessages {
	if !s.transactionReceived || !s.transitionReceived || s.producedIotaBlock == nil {
		return nil // Not yet.
	}
	if s.transactionReceived && !s.transactionReused && s.producedStateBlock == nil {
		return nil // Have to wait for the block.
	}

	if s.readyCB == nil {
		return nil // Already
	}
	cb := s.readyCB
	s.readyCB = nil
	return cb(
		s.transactionReused,
		s.transaction,
		s.producedIotaBlock,
		s.producedChainOutputs,
		s.producedStateBlock,
		s.consumedAnchorOutputID,
		s.consumedAccountOutputID,
	)
}
