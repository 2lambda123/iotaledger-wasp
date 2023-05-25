// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package rstemplates

var eventhandlersRs = map[string]string{
	// *******************************
	"eventhandlers.rs": `
use std::collections::HashMap;
use wasmlib::*;

use crate::*;

pub struct $PkgName$+EventHandlers {
    my_id: u32,
    $pkg_name$+_handlers: HashMap<&'static str, fn(evt: &$PkgName$+EventHandlers, msg: &Vec<u8>)>,

$#each events eventHandlerMember
}

impl IEventHandlers for $PkgName$+EventHandlers {
    fn call_handler(&self, topic: &str, data: &Vec<u8>) {
        if let Some(handler) = self.$pkg_name$+_handlers.get(topic) {
            handler(self, data);
        }
    }

    fn id(&self) -> u32 {
        self.my_id
    }
}

unsafe impl Send for $PkgName$+EventHandlers {}
unsafe impl Sync for $PkgName$+EventHandlers {}

impl $PkgName$+EventHandlers {
    pub fn new() -> $PkgName$+EventHandlers {
        let mut handlers: HashMap<&str, fn(evt: &$PkgName$+EventHandlers, msg: &Vec<u8>)> = HashMap::new();
$#each events eventHandler
        return $PkgName$+EventHandlers {
            my_id: EventHandlers::generate_id(),
            $pkg_name$+_handlers: handlers,
$#each events eventHandlerMemberInit
        };
    }
$#each events eventFuncSignature
}
$#each events eventClass
`,
	// *******************************
	"eventHandlerMember": `
    $evt_name: Box<dyn Fn(&Event$EvtName)>,
`,
	// *******************************
	"eventHandlerMemberInit": `
            $evt_name: Box::new(|_e| {}),
`,
	// *******************************
	"eventFuncSignature": `

    pub fn on_$pkg_name$+_$evt_name<F>(&mut self, handler: F)
        where F: Fn(&Event$EvtName) + 'static {
        self.$evt_name = Box::new(handler);
    }
`,
	// *******************************
	"eventHandler": `
        handlers.insert("$hscName.$evtName", |e, m| { (e.$evt_name)(&Event$EvtName::new(m)); });
`,
	// *******************************
	"eventClass": `

pub struct Event$EvtName {
    pub timestamp: u64,
$#each event eventClassField
}

impl Event$EvtName {
    pub fn new(msg: &Vec<u8>) -> Event$EvtName {
        let mut dec = WasmDecoder::new(msg);
        let _topic = string_decode(&mut dec);
        Event$EvtName {
            timestamp: uint64_decode(&mut dec),
$#each event eventHandlerField
        }
    }

    pub fn encode(self) -> Vec<u8> {
        let mut enc = WasmEncoder::new();
        // topic
        string_encode(&mut enc, "$hscName.$evtName");

        // payload
        uint64_encode(&mut enc, self.timestamp);
$#each event eventEncode
        return enc.buf();
    }
}
`,
	// *******************************
	"eventClassField": `
    pub $fld_name: $fldLangType,
`,
	// *******************************
	"eventArgs": `
        $fld_name: $fldLangType,
`,
	// *******************************
	"eventHandlerField": `
            $fld_name: $fld_type$+_decode(&mut dec),
`,
	// *******************************
	"eventEncode": `
        $fld_type$+_encode(&mut enc, $fldRef$+self.$fld_name);
`,
}
