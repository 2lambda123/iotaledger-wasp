// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package gotemplates

var eventsGo = map[string]string{
	// *******************************
	"events.go": `
package $package

$#emit importWasmLibAndWasmTypes

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
$#each event eventParam
$#if event eventEndFunc2
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.StringEncode(enc, "$package.$evtName")
$#each event eventEncode
	buf := enc.Buf()
	wasmlib.ScFuncContext{}.Event(buf)
}
`,
	// *******************************
	"eventParam": `
$#each fldComment _eventParamComment
	$fldName $fldLangType,
`,
	// *******************************
	"eventEncode": `
	wasmtypes.$FldType$+Encode(enc, $fldName)
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
