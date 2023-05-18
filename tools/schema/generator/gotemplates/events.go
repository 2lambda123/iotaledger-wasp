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
$#each event eventName
	}
	wasmlib.ScFuncContext{}.Event(isc.Encode(evt))
}

var _ isc.Event = &$EvtName$+Event{}

type $EvtName$+Event struct {
$#each event eventDefParam
}

func (e *$EvtName$+Event) Topic() []byte {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.StringEncode(enc, "$package.$evtName")
	return enc.Buf()
}

func (e *$EvtName$+Event) Payload() []byte {
	enc := wasmtypes.NewWasmEncoder()
$#each event eventEncode
	return enc.Buf()
}

func (e *$EvtName$+Event) DecodePayload(payload []byte) {
	dec := wasmtypes.NewWasmDecoder(payload)
	topic := wasmtypes.StringDecode(dec)
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}
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
	$fldName $fldLangType
`,
	// *******************************
	"eventName": `
		$fldName,
`,
	// *******************************
	"eventEncode": `
	wasmtypes.$FldType$+Encode(enc, e.$fldName)
`,
	// *******************************
	"eventDecode": `
	e.$fldName = wasmtypes.$FldType$+Decode(dec)
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
