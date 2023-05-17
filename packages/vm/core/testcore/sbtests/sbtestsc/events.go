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

func (e *GenericDataEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
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

func (e *TestEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
}
