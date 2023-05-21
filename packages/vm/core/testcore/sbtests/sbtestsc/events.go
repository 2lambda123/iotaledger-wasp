package sbtestsc

import (
	"bytes"
	"fmt"
	"time"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &GenericDataEvent{}

type GenericDataEvent struct {
	Timestamp uint64
	Counter   uint64
}

func (e *GenericDataEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, []byte("GenericDataEvent")); err != nil {
		panic(fmt.Errorf("failed to write GenericDataEvent: %w", err))
	}
	return w.Bytes()
}

func (e *GenericDataEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteUint64(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *GenericDataEvent) DecodePayload(payload []byte) {
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
	if err := util.ReadUint64(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to read event.Counter: %w", err))
	}
}

var _ isc.Event = &GenericDataEvent{}

type TestEvent struct {
	Timestamp uint64
	Message   string
}

func (e *TestEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, []byte("TestEvent")); err != nil {
		panic(fmt.Errorf("failed to write TestEvent: %w", err))
	}
	return w.Bytes()
}

func (e *TestEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteString16(&w, e.Message); err != nil {
		panic(fmt.Errorf("failed to write event.Message: %w", err))
	}
	return w.Bytes()
}

func (e *TestEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadBytes8(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if !bytes.Equal(topic, isc.DecodeEventTopic(e)) {
		panic("decode by unmatched event type")
	}
	if err = util.ReadUint64(r, &e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to read event.Timestamp: %w", err))
	}
	str, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Message: %w", err))
	}
	e.Message = str
}
