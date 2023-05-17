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
	wasmlib.ScFuncContext{}.Event(evt.Encode())
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

func (e *$EvtName$+Event) Encode() []byte {
	return append(e.Topic(), e.Payload()...)
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
	"eventSetEndFunc": `
$#set endFunc $nil
`,
	// *******************************
	"eventEndFunc2": `
) {
`,
}
