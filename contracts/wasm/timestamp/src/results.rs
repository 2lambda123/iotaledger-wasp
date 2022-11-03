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
pub struct ImmutableGetTimestampResults {
	pub(crate) proxy: Proxy,
}

impl ImmutableGetTimestampResults {
    // last official timestamp generated
    pub fn timestamp(&self) -> ScImmutableUint64 {
		ScImmutableUint64::new(self.proxy.root(RESULT_TIMESTAMP))
	}
}

#[derive(Clone)]
pub struct MutableGetTimestampResults {
	pub(crate) proxy: Proxy,
}

impl MutableGetTimestampResults {
    // last official timestamp generated
    pub fn timestamp(&self) -> ScMutableUint64 {
		ScMutableUint64::new(self.proxy.root(RESULT_TIMESTAMP))
	}
}
