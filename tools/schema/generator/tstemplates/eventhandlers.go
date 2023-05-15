// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package tstemplates

var eventhandlersTs = map[string]string{
	// *******************************
	"eventhandlers.ts": `
$#emit importWasmLib
$#emit importWasmTypes

export class $PkgName$+EventHandlers implements wasmlib.IEventHandlers {
    private myID: u32;
    private $pkgName$+Handlers: Map<string, (evt: $PkgName$+EventHandlers, msg: Uint8Array) => void> = new Map();

    /* eslint-disable @typescript-eslint/no-empty-function */
$#each events eventHandlerMember
    /* eslint-enable @typescript-eslint/no-empty-function */

    public constructor() {
        this.myID = wasmlib.eventHandlersGenerateID();
$#each events eventHandler
    }

    public callHandler(data: Uint8Array): void {
        const dec = new wasmtypes.WasmDecoder(data);
	    const topic = wasmtypes.stringDecode(dec);
        const handler = this.$pkgName$+Handlers.get(topic);
        if (handler) {
            handler(this, data);
        }
    }

    public id(): u32 {
        return this.myID;
    }
$#each events eventFuncSignature
}
$#each events eventClass
`,
	// *******************************
	"eventHandler": `
        this.$pkgName$+Handlers.set('$hscName.$evtName', (evt: $PkgName$+EventHandlers, msg: Uint8Array) => evt.$evtName(new $EvtName$+Event(msg)));
`,
	// *******************************
	"eventHandlerMember": `
    $evtName: (evt: $EvtName$+Event) => void = () => {};
`,
	// *******************************
	"eventFuncSignature": `

    public on$PkgName$EvtName(handler: (evt: $EvtName$+Event) => void): void {
        this.$evtName = handler;
    }
`,
	// *******************************
	"eventClass": `

export class $EvtName$+Event {
    public readonly timestamp: u64;
$#each event eventClassField

    public constructor(msg: Uint8Array) {
        const dec = new wasmlib.WasmDecoder(msg);
        const topic =  wasmtypes.stringDecode(dec);
        this.timestamp = wasmtypes.uint64Decode(dec);
$#each event eventHandlerField
        dec.close();
    }
}
`,
	// *******************************
	"eventClassField": `
    public readonly $fldName: $fldLangType;
`,
	// *******************************
	"eventHandlerField": `
        this.$fldName = wasmtypes.$fldType$+Decode(dec);
`,
}
