// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package tstemplates

var eventsTs = map[string]string{
	// *******************************
	"events.ts": `
$#emit importWasmLib
$#emit importWasmTypes
import * as sc from './index';

$#set TypeName $Package$+Events
export class $TypeName {
$#each events eventFunc
}
`,
	// *******************************
	"eventFunc": `
$#set endFunc ): void {
$#if event eventSetEndFunc

$#each eventComment _eventComment
    $evtName($endFunc
$#each event eventParam
$#if event eventEndFunc2
		let enc = new wasmtypes.WasmEncoder();
		wasmtypes.stringEncode(enc, "$hscName.$evtName");
		const timestamp = new wasmlib.ScFuncContext().timestamp();
		wasmtypes.uint64Encode(enc, timestamp);
$#each event eventEncode
		new wasmlib.ScFuncContext().event(enc.buf());
    }
`,
	// *******************************
	"eventParam": `
$#each fldComment _eventParamComment
        $fldName: $fldLangType,
`,
	// *******************************
	"eventEncode": `
		wasmtypes.$fldType$+Encode(enc, $fldName);
`,
	// *******************************
	"eventSetEndFunc": `
$#set endFunc $nil
`,
	// *******************************
	"eventEndFunc2": `
    ): void {
`,
	// *******************************
	"eventName": `
			$fldName,
`,
}
