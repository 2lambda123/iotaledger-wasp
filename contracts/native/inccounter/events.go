package inccounter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &InitializeEvent{}

type InitializeEvent struct {
	Timestamp uint64
	Counter   uint32
}

func (e *InitializeEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, []byte("InitializeEvent")); err != nil {
		panic(fmt.Errorf("failed to write InitializeEvent: %w", err))
	}
	return w.Bytes()
}

func (e *InitializeEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteUint32(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *InitializeEvent) DecodePayload(payload []byte) {
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
	if err := util.ReadUint32(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
}

var _ isc.Event = &InitializeEvent{}

type IncCounterEvent struct {
	Timestamp uint64
	Counter   uint32
}

func (e *IncCounterEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, []byte(("IncCounterEvent"))); err != nil {
		panic(fmt.Errorf("failed to write IncCounterEvent: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteUint32(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterEvent) DecodePayload(payload []byte) {
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
	if err := util.ReadUint32(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
}

type IncCounterAndRepeatOnceEvent struct {
	Timestamp uint64
	Counter   uint32
}

func (e *IncCounterAndRepeatOnceEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, []byte("IncCounterAndRepeatOnceEvent")); err != nil {
		panic(fmt.Errorf("failed to write IncCounterAndRepeatOnceEvent: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterAndRepeatOnceEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteUint32(&w, e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
	return w.Bytes()
}

func (e *IncCounterAndRepeatOnceEvent) DecodePayload(payload []byte) {
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
	if err := util.ReadUint32(r, &e.Counter); err != nil {
		panic(fmt.Errorf("failed to write event.Counter: %w", err))
	}
}
