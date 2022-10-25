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
pub struct ImmutableDonateParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableDonateParams {
    // feedback for the person you donate to
    pub fn feedback(&self) -> ScImmutableString {
        ScImmutableString::new(self.proxy.root(PARAM_FEEDBACK))
    }
}

#[derive(Clone)]
pub struct MutableDonateParams {
    pub(crate) proxy: Proxy,
}

impl MutableDonateParams {
    // feedback for the person you donate to
    pub fn feedback(&self) -> ScMutableString {
        ScMutableString::new(self.proxy.root(PARAM_FEEDBACK))
    }
}

#[derive(Clone)]
pub struct ImmutableInitParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableInitParams {
    pub fn owner(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_OWNER))
    }
}

#[derive(Clone)]
pub struct MutableInitParams {
    pub(crate) proxy: Proxy,
}

impl MutableInitParams {
    pub fn owner(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_OWNER))
    }
}

#[derive(Clone)]
pub struct ImmutableWithdrawParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableWithdrawParams {
    // amount to withdraw
    pub fn amount(&self) -> ScImmutableUint64 {
        ScImmutableUint64::new(self.proxy.root(PARAM_AMOUNT))
    }
}

#[derive(Clone)]
pub struct MutableWithdrawParams {
    pub(crate) proxy: Proxy,
}

impl MutableWithdrawParams {
    // amount to withdraw
    pub fn amount(&self) -> ScMutableUint64 {
        ScMutableUint64::new(self.proxy.root(PARAM_AMOUNT))
    }
}

#[derive(Clone)]
pub struct ImmutableDonationParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableDonationParams {
    pub fn nr(&self) -> ScImmutableUint32 {
        ScImmutableUint32::new(self.proxy.root(PARAM_NR))
    }
}

#[derive(Clone)]
pub struct MutableDonationParams {
    pub(crate) proxy: Proxy,
}

impl MutableDonationParams {
    pub fn nr(&self) -> ScMutableUint32 {
        ScMutableUint32::new(self.proxy.root(PARAM_NR))
    }
}
