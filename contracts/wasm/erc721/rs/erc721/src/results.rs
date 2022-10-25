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
pub struct ImmutableBalanceOfResults {
    pub(crate) proxy: Proxy,
}

impl ImmutableBalanceOfResults {
    // amount of tokens owned by owner
    pub fn amount(&self) -> ScImmutableUint64 {
        ScImmutableUint64::new(self.proxy.root(RESULT_AMOUNT))
    }
}

#[derive(Clone)]
pub struct MutableBalanceOfResults {
    pub(crate) proxy: Proxy,
}

impl MutableBalanceOfResults {
    // amount of tokens owned by owner
    pub fn amount(&self) -> ScMutableUint64 {
        ScMutableUint64::new(self.proxy.root(RESULT_AMOUNT))
    }
}

#[derive(Clone)]
pub struct ImmutableGetApprovedResults {
    pub(crate) proxy: Proxy,
}

impl ImmutableGetApprovedResults {
    pub fn approved(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(RESULT_APPROVED))
    }
}

#[derive(Clone)]
pub struct MutableGetApprovedResults {
    pub(crate) proxy: Proxy,
}

impl MutableGetApprovedResults {
    pub fn approved(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(RESULT_APPROVED))
    }
}

#[derive(Clone)]
pub struct ImmutableIsApprovedForAllResults {
    pub(crate) proxy: Proxy,
}

impl ImmutableIsApprovedForAllResults {
    pub fn approval(&self) -> ScImmutableBool {
        ScImmutableBool::new(self.proxy.root(RESULT_APPROVAL))
    }
}

#[derive(Clone)]
pub struct MutableIsApprovedForAllResults {
    pub(crate) proxy: Proxy,
}

impl MutableIsApprovedForAllResults {
    pub fn approval(&self) -> ScMutableBool {
        ScMutableBool::new(self.proxy.root(RESULT_APPROVAL))
    }
}

#[derive(Clone)]
pub struct ImmutableNameResults {
    pub(crate) proxy: Proxy,
}

impl ImmutableNameResults {
    pub fn name(&self) -> ScImmutableString {
        ScImmutableString::new(self.proxy.root(RESULT_NAME))
    }
}

#[derive(Clone)]
pub struct MutableNameResults {
    pub(crate) proxy: Proxy,
}

impl MutableNameResults {
    pub fn name(&self) -> ScMutableString {
        ScMutableString::new(self.proxy.root(RESULT_NAME))
    }
}

#[derive(Clone)]
pub struct ImmutableOwnerOfResults {
    pub(crate) proxy: Proxy,
}

impl ImmutableOwnerOfResults {
    pub fn owner(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(RESULT_OWNER))
    }
}

#[derive(Clone)]
pub struct MutableOwnerOfResults {
    pub(crate) proxy: Proxy,
}

impl MutableOwnerOfResults {
    pub fn owner(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(RESULT_OWNER))
    }
}

#[derive(Clone)]
pub struct ImmutableSymbolResults {
    pub(crate) proxy: Proxy,
}

impl ImmutableSymbolResults {
    pub fn symbol(&self) -> ScImmutableString {
        ScImmutableString::new(self.proxy.root(RESULT_SYMBOL))
    }
}

#[derive(Clone)]
pub struct MutableSymbolResults {
    pub(crate) proxy: Proxy,
}

impl MutableSymbolResults {
    pub fn symbol(&self) -> ScMutableString {
        ScMutableString::new(self.proxy.root(RESULT_SYMBOL))
    }
}

#[derive(Clone)]
pub struct ImmutableTokenURIResults {
    pub(crate) proxy: Proxy,
}

impl ImmutableTokenURIResults {
    pub fn token_uri(&self) -> ScImmutableString {
        ScImmutableString::new(self.proxy.root(RESULT_TOKEN_URI))
    }
}

#[derive(Clone)]
pub struct MutableTokenURIResults {
    pub(crate) proxy: Proxy,
}

impl MutableTokenURIResults {
    pub fn token_uri(&self) -> ScMutableString {
        ScMutableString::new(self.proxy.root(RESULT_TOKEN_URI))
    }
}
