package blob

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &StoreBlobEvent{}

type StoreBlobEvent struct {
	BlobHash   hashing.HashValue
	FieldSizes []uint32
}

func (e *StoreBlobEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncStoreBlob.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncStoreBlob.Hname(): %w", err))
	}
	return w.Bytes()
}

func (e *StoreBlobEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes32(&w, e.BlobHash.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.BlobHash: %w", err))
	}
	for _, v := range e.FieldSizes {
		if err := util.WriteUint32(&w, v); err != nil {
			panic(fmt.Errorf("failed to write event.BlobHash: %w", err))
		}
	}
	return w.Bytes()
}

func (e *StoreBlobEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
}
