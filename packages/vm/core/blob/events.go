package blob

import (
	"bytes"
	"fmt"
	"time"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &StoreBlobEvent{}

type StoreBlobEvent struct {
	Timestamp  uint64
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
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
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

func (e *StoreBlobEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}
	if err := util.ReadUint64(r, &e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to read event.Timestamp: %w", err))
	}
	b, err := util.ReadBytes32(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.BlobHash: %w", err))
	}
	e.BlobHash, err = hashing.HashValueFromBytes(b)
	if err != nil {
		panic(fmt.Errorf("failed to convert HashValue from bytes: %w", err))
	}

	for i := 0; r.Len() != 0; i++ {
		if err := util.ReadUint32(r, &e.FieldSizes[i]); err != nil {
			panic(fmt.Errorf("failed to read event.FieldSizes: %w", err))
		}
	}
}
