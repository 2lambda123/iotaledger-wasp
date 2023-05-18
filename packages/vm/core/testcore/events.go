package testcore

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &TestManyEvent{}

type TestManyEvent struct {
	I uint32
}

func (e *TestManyEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, "TestManyEvent"); err != nil {
		panic(fmt.Errorf("failed to write TestManyEvent: %w", err))
	}
	return w.Bytes()
}

func (e *TestManyEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint32(&w, e.I); err != nil {
		panic(fmt.Errorf("failed to write event.I: %w", err))
	}
	return w.Bytes()
}

func (e *TestManyEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}

	if err := util.ReadUint32(r, &e.I); err != nil {
		panic(fmt.Errorf("failed to read event.I: %w", err))
	}
}

type DEvent interface {
	Decode(b []byte)
}

func Decode(e DEvent, b []byte) {
	e.Decode(b)
}

var _ isc.Event = &TestManyEvent{}

type TestSingleEvent struct {
	Message string
}

func (e *TestSingleEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, "TestSingleEvent"); err != nil {
		panic(fmt.Errorf("failed to write TestSingleEvent: %w", err))
	}
	return w.Bytes()
}

func (e *TestSingleEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteString16(&w, e.Message); err != nil {
		panic(fmt.Errorf("failed to write event.Message: %w", err))
	}
	return w.Bytes()
}

func (e *TestSingleEvent) DecodePayload(payload []byte) {
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
