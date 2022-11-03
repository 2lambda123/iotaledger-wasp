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
pub struct ImmutableAllowanceResults {
	pub(crate) proxy: Proxy,
}

impl ImmutableAllowanceResults {
    pub fn amount(&self) -> ScImmutableUint64 {
		ScImmutableUint64::new(self.proxy.root(RESULT_AMOUNT))
	}
}

#[derive(Clone)]
pub struct MutableAllowanceResults {
	pub(crate) proxy: Proxy,
}

impl MutableAllowanceResults {
    pub fn amount(&self) -> ScMutableUint64 {
		ScMutableUint64::new(self.proxy.root(RESULT_AMOUNT))
	}
}

#[derive(Clone)]
pub struct ImmutableBalanceOfResults {
	pub(crate) proxy: Proxy,
}

impl ImmutableBalanceOfResults {
    pub fn amount(&self) -> ScImmutableUint64 {
		ScImmutableUint64::new(self.proxy.root(RESULT_AMOUNT))
	}
}

#[derive(Clone)]
pub struct MutableBalanceOfResults {
	pub(crate) proxy: Proxy,
}

impl MutableBalanceOfResults {
    pub fn amount(&self) -> ScMutableUint64 {
		ScMutableUint64::new(self.proxy.root(RESULT_AMOUNT))
	}
}

#[derive(Clone)]
pub struct ImmutableTotalSupplyResults {
	pub(crate) proxy: Proxy,
}

impl ImmutableTotalSupplyResults {
    pub fn supply(&self) -> ScImmutableUint64 {
		ScImmutableUint64::new(self.proxy.root(RESULT_SUPPLY))
	}
}

#[derive(Clone)]
pub struct MutableTotalSupplyResults {
	pub(crate) proxy: Proxy,
}

impl MutableTotalSupplyResults {
    pub fn supply(&self) -> ScMutableUint64 {
		ScMutableUint64::new(self.proxy.root(RESULT_SUPPLY))
	}
}
