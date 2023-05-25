// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package gotemplates

var eventhandlersGo = map[string]string{
	// *******************************
	"eventhandlers.go": `
package $package

$#emit importWasmLibAndWasmTypes

var $pkgName$+Handlers = map[string]func(*$PkgName$+EventHandlers, []byte){
$#each events eventHandler
}

type $PkgName$+EventHandlers struct {
	myID uint32
$#each events eventHandlerMember
}

var _ wasmlib.IEventHandlers = new($PkgName$+EventHandlers)

func New$PkgName$+EventHandlers() *$PkgName$+EventHandlers {
	return &$PkgName$+EventHandlers{ myID: wasmlib.EventHandlersGenerateID() }
}

func (h *$PkgName$+EventHandlers) CallHandler(data []byte) {
	dec := wasmtypes.NewWasmDecoder(data)
	topic := wasmtypes.StringDecode(dec)
	handler := $pkgName$+Handlers[topic]
	if handler != nil {
		handler(h, data)
	}
}

func (h *$PkgName$+EventHandlers) ID() uint32 {
	return h.myID
}
$#each events eventFuncSignature
$#each events eventClass
`,
	// *******************************
	"eventHandlerMember": `
	$evtName func(e *$EvtName$+Event)
`,
	// *******************************
	"eventFuncSignature": `

func (h *$PkgName$+EventHandlers) On$PkgName$EvtName(handler func(e *$EvtName$+Event)) {
	h.$evtName = handler
}
`,
	// *******************************
	"eventHandler": `
	"$hscName.$evtName": func(evt *$PkgName$+EventHandlers, msg []byte) { evt.on$PkgName$EvtName$+Thunk(msg) },
`,
	// *******************************
	"eventClass": `

type Event$EvtName struct {
	Timestamp uint64
$#each event eventClassField
}

func (h *$PkgName$+EventHandlers) on$PkgName$EvtName$+Thunk(msg []byte) {
	if h.$evtName == nil {
		return
	}
	e := &$EvtName$+Event{}
	e.DecodePayload(msg)
	h.$evtName(e)
}
`,
	// *******************************
	"eventClassField": `
	$FldName $fldLangType
`,
}
