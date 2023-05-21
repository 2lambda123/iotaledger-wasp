// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmlib

import (
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type IEventHandlers interface {
	CallHandler(data []byte)
	ID() uint32
}

var nextID = uint32(0)

func EventHandlersGenerateID() uint32 {
	nextID++
	return nextID
}

// a copy of isc.Event
type Event interface {
	Topic() []byte
	Payload() []byte
	DecodePayload(payload []byte)
}

func Encode(e Event) []byte {
	return append(e.Topic(), e.Payload()...)
}

func DecodePayloadTopic(payload []byte) []byte {
	dec := wasmtypes.NewWasmDecoder(payload)
	// FIXME topic should be bytes
	topic := wasmtypes.StringDecode(dec)
	return []byte(topic)
}

func DecodeEventTopic(e Event) []byte {
	dec := wasmtypes.NewWasmDecoder(e.Topic())
	return []byte(wasmtypes.StringDecode(dec))
}
