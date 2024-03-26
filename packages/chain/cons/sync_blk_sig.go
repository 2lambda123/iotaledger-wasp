package cons

import (
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/gpa"
)

type SyncBlkSig interface {
	HaveBlock(bl *iotago.Block) gpa.OutMessages
	HaveSig(sig []byte) gpa.OutMessages
}

type syncBlkSig struct {
	readyCB func(bl *iotago.Block, sig []byte) gpa.OutMessages
	bl      *iotago.Block
	sig     []byte
}

func NewSyncBlkSig(
	readyCB func(bl *iotago.Block, sig []byte) gpa.OutMessages,
) SyncBlkSig {
	return &syncBlkSig{readyCB: readyCB}
}

func (s *syncBlkSig) HaveBlock(bl *iotago.Block) gpa.OutMessages {
	if s.bl != nil {
		return nil
	}
	s.bl = bl
	return s.tryOutput()
}

func (s *syncBlkSig) HaveSig(sig []byte) gpa.OutMessages {
	if s.sig != nil {
		return nil
	}
	s.sig = sig
	return s.tryOutput()
}

func (s *syncBlkSig) tryOutput() gpa.OutMessages {
	if s.bl == nil || s.sig == nil {
		return nil
	}
	if s.readyCB == nil {
		return nil
	}
	cb := s.readyCB
	s.readyCB = nil
	return cb(s.bl, s.sig)
}
