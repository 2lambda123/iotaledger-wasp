// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use crate::*;
use wasmlib::*;

pub trait IEventHandler {
    fn call_handler(&self, topic: &str, params: &[&str]);
    fn equal(&self, h: &dyn IEventHandler) -> bool;
}

pub struct WasmClientContext {
    pub chain_id: ScChainID,
    pub event_handlers: Vec<Box<dyn IEventHandler>>,
    pub key_pair: Option<keypair::KeyPair>,
    pub req_id: ScRequestID,
    pub sc_name: String,
    pub sc_hname: ScHname,
    pub svc_client: WasmClientService, //TODO Maybe  use 'dyn IClientService' for 'svc_client' instead of a struct
}

impl WasmClientContext {
    pub fn new(
        svc_client: &WasmClientService,
        chain_id: &wasmlib::ScChainID,
        sc_name: &str,
    ) -> WasmClientContext {
        WasmClientContext {
            svc_client: svc_client.clone(),
            sc_name: sc_name.to_string(),
            sc_hname: ScHname::new(sc_name),
            chain_id: chain_id.clone(),
            event_handlers: Vec::new(),
            key_pair: None,
            req_id: request_id_from_bytes(&[]),
        }
    }

    pub fn default() -> WasmClientContext {
        WasmClientContext {
            svc_client: WasmClientService::default(),
            sc_name: String::new(),
            sc_hname: ScHname(0),
            chain_id: chain_id_from_bytes(&[]),
            event_handlers: Vec::new(),
            key_pair: None,
            req_id: request_id_from_bytes(&[]),
        }
    }

    pub fn current_chain_id(&self) -> ScChainID {
        return self.chain_id;
    }

    pub fn init_func_call_context(&'static self) {
        wasmlib::host::connect_host(self);
    }

    pub fn init_view_call_context(&'static self, _contract_hname: &ScHname) -> ScHname {
        wasmlib::host::connect_host(self);
        return self.sc_hname;
    }

    pub fn register(&mut self, handler: Box<dyn IEventHandler>) -> errors::Result<()> {
        let handler_iterator = self.event_handlers.iter();
        for h in handler_iterator {
            if handler.equal(h.as_ref()) {
                return Ok(());
            }
        }
        self.event_handlers.push(handler);
        if self.event_handlers.len() > 1 {
            return Ok(());
        }
        return self.start_event_handlers();
    }

    // overrides default contract name
    pub fn service_contract_name(&mut self, contract_name: &str) {
        self.sc_hname = wasmlib::ScHname::new(contract_name);
    }

    pub fn sign_requests(&mut self, key_pair: &keypair::KeyPair) {
        self.key_pair = Some(key_pair.clone());
    }

    pub fn unregister(&mut self, handler: Box<dyn IEventHandler>) {
        self.event_handlers.retain(|h| {
            if handler.equal(h.as_ref()) {
                return false;
            } else {
                return true;
            }
        });

        if self.event_handlers.len() == 0 {
            self.stop_event_handlers();
        }
    }

    pub fn wait_request(&mut self, req_id: Option<&ScRequestID>) -> errors::Result<()> {
        let r_id;
        match req_id {
            Some(id) => r_id = id,
            None => r_id = &self.req_id,
        }
        return self.svc_client.wait_until_request_processed(
            &self.chain_id,
            &r_id,
            std::time::Duration::new(60, 0),
        );
    }

    pub fn start_event_handlers(&self) -> errors::Result<()> {
        todo!()
    }

    pub fn stop_event_handlers(&self) {
        todo!()
    }
}
