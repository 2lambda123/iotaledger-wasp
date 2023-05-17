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

func (e *InitializeEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
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

func (e *IncCounterEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
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

func (e *IncCounterAndRepeatOnceEvent) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
}
