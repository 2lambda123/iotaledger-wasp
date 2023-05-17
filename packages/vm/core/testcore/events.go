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

func (e *TestManyEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
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

func (e *TestSingleEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
}
