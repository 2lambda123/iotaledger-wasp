package cons

import (
	"time"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
)

type SyncBlkData interface {
	HaveTipsProposal(tipsFn func(randomness hashing.HashValue) iotago.BlockIDs) gpa.OutMessages
	HaveRandomness(randomness hashing.HashValue) gpa.OutMessages
	HaveTimestamp(timestamp time.Time) gpa.OutMessages
	HaveSignedTX(tx *iotago.SignedTransaction) gpa.OutMessages
}

type SyncBlkDataCB = func(
	tipsFn func(randomness hashing.HashValue) iotago.BlockIDs,
	randomness hashing.HashValue,
	timestamp time.Time,
	tx *iotago.SignedTransaction,
) gpa.OutMessages

type syncBlkData struct {
	readyCB    SyncBlkDataCB
	tipsFn     func(randomness hashing.HashValue) iotago.BlockIDs
	randomness *hashing.HashValue
	timestamp  *time.Time
	tx         *iotago.SignedTransaction
}

func NewSyncBlkData(readyCB SyncBlkDataCB) SyncBlkData {
	return &syncBlkData{readyCB: readyCB}
}

func (s *syncBlkData) HaveTipsProposal(tipsFn func(randomness hashing.HashValue) iotago.BlockIDs) gpa.OutMessages {
	if s.tipsFn != nil {
		return nil
	}
	s.tipsFn = tipsFn
	return s.tryOutput()
}

func (s *syncBlkData) HaveRandomness(randomness hashing.HashValue) gpa.OutMessages {
	if s.randomness != nil {
		return nil
	}
	s.randomness = &randomness
	return s.tryOutput()
}

func (s *syncBlkData) HaveTimestamp(timestamp time.Time) gpa.OutMessages {
	if s.timestamp != nil {
		return nil
	}
	s.timestamp = &timestamp
	return s.tryOutput()
}

func (s *syncBlkData) HaveSignedTX(tx *iotago.SignedTransaction) gpa.OutMessages {
	if s.tx != nil {
		return nil
	}
	s.tx = tx
	return s.tryOutput()
}

func (s *syncBlkData) tryOutput() gpa.OutMessages {
	if s.tipsFn == nil || s.randomness == nil || s.timestamp == nil || s.tx == nil {
		return nil // Not yet.
	}
	if s.readyCB == nil {
		return nil // Already called.
	}
	cb := s.readyCB
	s.readyCB = nil
	return cb(s.tipsFn, *s.randomness, *s.timestamp, s.tx)
}
