// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]

use wasmlib::*;

use crate::*;

#[derive(Clone)]
pub struct MapAgentIDToImmutableUint64 {
    pub(crate) proxy: Proxy,
}

impl MapAgentIDToImmutableUint64 {
    pub fn get_uint64(&self, key: &ScAgentID) -> ScImmutableUint64 {
        ScImmutableUint64::new(self.proxy.key(&agent_id_to_bytes(key)))
    }
}

pub type ImmutableAllowancesForAgent = MapAgentIDToImmutableUint64;

#[derive(Clone)]
pub struct MapAgentIDToMutableUint64 {
    pub(crate) proxy: Proxy,
}

impl MapAgentIDToMutableUint64 {
    pub fn clear(&self) {
        self.proxy.clear_map();
    }

    pub fn get_uint64(&self, key: &ScAgentID) -> ScMutableUint64 {
        ScMutableUint64::new(self.proxy.key(&agent_id_to_bytes(key)))
    }
}

pub type MutableAllowancesForAgent = MapAgentIDToMutableUint64;
