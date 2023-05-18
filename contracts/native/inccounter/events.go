package inccounter

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &InitializeEvent{}

type InitializeEvent struct {
	Counter uint32
}

func (e *InitializeEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, "InitializeEvent"); err != nil {
		panic(fmt.Errorf("failed to write InitializeEvent: %w", err))
	}
	return w.Bytes()
}

func (e *InitializeEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint32(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *InitializeEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}
	if err := util.ReadUint32(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
}

var _ isc.Event = &InitializeEvent{}

type IncCounterEvent struct {
	Counter uint32
}

func (e *IncCounterEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, "IncCounterEvent"); err != nil {
		panic(fmt.Errorf("failed to write IncCounterEvent: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint32(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}

	if err := util.ReadUint32(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
}

type IncCounterAndRepeatOnceEvent struct {
	Counter uint32
}

func (e *IncCounterAndRepeatOnceEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, "IncCounterAndRepeatOnceEvent"); err != nil {
		panic(fmt.Errorf("failed to write IncCounterAndRepeatOnceEvent: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterAndRepeatOnceEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint32(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterAndRepeatOnceEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}

	if err := util.ReadUint32(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
}
