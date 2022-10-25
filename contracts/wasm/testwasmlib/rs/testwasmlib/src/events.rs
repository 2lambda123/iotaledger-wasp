// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]
#![allow(unused_mut)]

use wasmlib::*;

pub struct TestWasmLibEvents {
}

impl TestWasmLibEvents {

	pub fn test(&self,
        address: &ScAddress,
        name: &str,
    ) {
		let mut evt = EventEncoder::new("testwasmlib.test");
		evt.encode(&address_to_string(&address));
		evt.encode(&string_to_string(&name));
		evt.emit();
	}
}
