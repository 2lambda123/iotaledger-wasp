// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import {ScFuncContext} from './context';
import {concat, hexEncode, stringToBytes, uint64Encode, WasmDecoder, WasmEncoder} from "./wasmtypes";

export interface IEventHandlers {
    callHandler(data: Uint8Array): void;

    id(): u32;
}

let nextID: u32 = 0;

export function eventHandlersGenerateID(): u32 {
    nextID++;
    return nextID;
}

export interface IEvent {
	topic(): Uint8Array;
	payload(): Uint8Array;
	decodePayload(payload: Uint8Array): void;
}

function encode(e: IEvent): Uint8Array {
	return new Uint8Array([...e.topic(), ...e.payload()]);
}

// TODO 
// func DecodePayloadTopic(payload []byte) []byte {
// 	dec := wasmtypes.NewWasmDecoder(payload)
// 	topic := wasmtypes.StringDecode(dec)
// 	return []byte(topic)
// }

// func DecodeEventTopic(e Event) []byte {
// 	dec := wasmtypes.NewWasmDecoder(e.Topic())
// 	return []byte(wasmtypes.StringDecode(dec))
// }