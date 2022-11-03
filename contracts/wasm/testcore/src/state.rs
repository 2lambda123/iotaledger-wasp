// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]
#![allow(unused_imports)]

use wasmlib::*;

use crate::*;

#[derive(Clone)]
pub struct ImmutableTestCoreState {
	pub(crate) proxy: Proxy,
}

impl ImmutableTestCoreState {
    pub fn counter(&self) -> ScImmutableUint64 {
		ScImmutableUint64::new(self.proxy.root(STATE_COUNTER))
	}

    pub fn ints(&self) -> MapStringToImmutableInt64 {
		MapStringToImmutableInt64 { proxy: self.proxy.root(STATE_INTS) }
	}

    pub fn strings(&self) -> MapStringToImmutableString {
		MapStringToImmutableString { proxy: self.proxy.root(STATE_STRINGS) }
	}
}

#[derive(Clone)]
pub struct MutableTestCoreState {
	pub(crate) proxy: Proxy,
}

impl MutableTestCoreState {
    pub fn as_immutable(&self) -> ImmutableTestCoreState {
		ImmutableTestCoreState { proxy: self.proxy.root("") }
	}

    pub fn counter(&self) -> ScMutableUint64 {
		ScMutableUint64::new(self.proxy.root(STATE_COUNTER))
	}

    pub fn ints(&self) -> MapStringToMutableInt64 {
		MapStringToMutableInt64 { proxy: self.proxy.root(STATE_INTS) }
	}

    pub fn strings(&self) -> MapStringToMutableString {
		MapStringToMutableString { proxy: self.proxy.root(STATE_STRINGS) }
	}
}
