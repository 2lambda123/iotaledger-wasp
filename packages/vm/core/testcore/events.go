package testcore

import (
	"bytes"
	"fmt"
	"time"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &TestManyEvent{}

type TestManyEvent struct {
	Timestamp uint64
	I         uint32
}

func (e *TestManyEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, []byte("TestManyEvent")); err != nil {
		panic(fmt.Errorf("failed to write TestManyEvent: %w", err))
	}
	return w.Bytes()
}

func (e *TestManyEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteUint32(&w, e.I); err != nil {
		panic(fmt.Errorf("failed to write event.I: %w", err))
	}
	return w.Bytes()
}

func (e *TestManyEvent) DecodePayload(payload []byte) {
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
	Timestamp uint64
	Message   string
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
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteString16(&w, e.Message); err != nil {
		panic(fmt.Errorf("failed to write event.Message: %w", err))
	}
	return w.Bytes()
}

func (e *TestSingleEvent) DecodePayload(payload []byte) {
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
