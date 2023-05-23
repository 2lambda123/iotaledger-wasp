package accounts

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &FoundryCreateNewEvent{}

type FoundryCreateNewEvent struct {
	Timestamp    uint64
	SerialNumber uint32
}

func (e *FoundryCreateNewEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncFoundryCreateNew.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncFoundryCreateNew.Hname(): %w", err))
	}
	return w.Bytes()
}

func (e *FoundryCreateNewEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteUint32(&w, e.SerialNumber); err != nil {
		panic(fmt.Errorf("failed to write event.SerialNumber: %w", err))
	}
	return w.Bytes()
}

func (e *FoundryCreateNewEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadBytes8(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if !bytes.Equal(topic, isc.DecodeEventTopic(e)) {
		panic("decode by unmatched event type")
	}
	if err := util.ReadUint64(r, &e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to read event.Timestamp: %w", err))
	}
	if err := util.ReadUint32(r, &e.SerialNumber); err != nil {
		panic(fmt.Errorf("failed to read event.SerialNumber: %w", err))
	}
}
