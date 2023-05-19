// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package gotemplates

var eventsGo = map[string]string{
	// *******************************
	"events.go": `
package $package

$#emit importWasmLibAndWasmTypesISC

$#set TypeName $Package$+Events
type $TypeName struct{}
$#each events eventFunc
`,
	// *******************************
	"eventFunc": `
$#set endFunc ) {
$#if event eventSetEndFunc

$#each eventComment _eventComment
func (e $TypeName) $EvtName($endFunc
	// TODO event param as array type will crash 
$#each event eventParam
$#if event eventEndFunc2
	evt := &$EvtName$+Event{
		wasmlib.ScFuncContext{}.Timestamp(),
$#each event eventName
	}
	wasmlib.ScFuncContext{}.Event(isc.Encode(evt))
}

var _ isc.Event = &$EvtName$+Event{}

type $EvtName$+Event struct {
	Timestamp uint64
$#each event eventDefParam
}

func (e *$EvtName$+Event) Topic() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.StringEncode(enc, "$package.$evtName")
	return enc.Buf()
}

func (e *$EvtName$+Event) Payload() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.Uint64Encode(enc, wasmlib.ScFuncContext{}.Timestamp())
$#each event eventEncode
	return enc.Buf()
}

func (e *$EvtName$+Event) DecodePayload(payload []byte) {
	dec := wasmtypes.NewWasmDecoder(payload)
	topic := wasmtypes.StringDecode(dec)
	// FIXME an unknown space is in front of e.Topic()
	etopic := string(e.Topic())[1:]
	if topic != etopic {
		panic("decode by unmatched event type")
	}
	wasmtypes.Uint64Decode(dec)
$#each event eventDecode
}

`,
	// *******************************
	"eventParam": `
$#each fldComment _eventParamComment
	$fldName $fldLangType,
`,
	// *******************************
	"eventDefParam": `
	$FldName $fldLangType
`,
	// *******************************
	"eventName": `
		$fldName,
`,
	// *******************************
	"eventEncode": `
	wasmtypes.$FldType$+Encode(enc, e.$FldName)
`,
	// *******************************
	"eventDecode": `
	e.$FldName = wasmtypes.$FldType$+Decode(dec)
`,
	// *******************************
	"eventSetEndFunc": `
$#set endFunc $nil
`,
	// *******************************
	"eventEndFunc2": `
) {
`,
}
