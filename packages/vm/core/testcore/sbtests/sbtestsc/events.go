package sbtestsc

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &GenericDataEvent{}

type GenericDataEvent struct {
	Counter uint32
}

func (e *GenericDataEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, "GenericDataEvent"); err != nil {
		panic(fmt.Errorf("failed to write GenericDataEvent: %w", err))
	}
	return w.Bytes()
}

func (e *GenericDataEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint32(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *GenericDataEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}

	if err := util.ReadUint32(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to read event.Counter: %w", err))
	}
}

var _ isc.Event = &GenericDataEvent{}

type TestEvent struct {
	Message string
}

func (e *TestEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, "TestEvent"); err != nil {
		panic(fmt.Errorf("failed to write TestEvent: %w", err))
	}
	return w.Bytes()
}

func (e *TestEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, e.Message); err != nil {
		panic(fmt.Errorf("failed to write event.Message: %w", err))
	}
	return w.Bytes()
}

func (e *TestEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}

	str, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Message: %w", err))
	}
	e.Message = str
}
