// // Copyright 2020 IOTA Stiftung
// // SPDX-License-Identifier: Apache-2.0

use crate::*;
use isc::{offledgerrequest::*, waspclient::*};
use std::sync::{mpsc, Arc, RwLock};
use std::time::Duration;

pub trait IClientService {
    fn call_view_by_hname(
        &self,
        chain_id: &ScChainID,
        contract_hname: &ScHname,
        function_hname: &ScHname,
        args: &[u8],
    ) -> errors::Result<Vec<u8>>;
    fn post_request(
        &self,
        chain_id: &ScChainID,
        contract_hname: &ScHname,
        function_hname: &ScHname,
        args: &[u8],
        allowance: &ScAssets,
        key_pair: &keypair::KeyPair,
        nonce: u64,
    ) -> errors::Result<ScRequestID>;
    fn subscribe_events(
        &self,
        tx: mpsc::Sender<Vec<String>>,
        done: Arc<RwLock<bool>>,
    ) -> errors::Result<()>;
    fn wait_until_request_processed(
        &self,
        chain_id: &ScChainID,
        req_id: &ScRequestID,
        timeout: Duration,
    ) -> errors::Result<()>;
}

#[derive(Clone, PartialEq)]
pub struct WasmClientService {
    client: waspclient::WaspClient,
    websocket: Option<websocket::Client>,
    event_port: String,
    last_err: errors::Result<()>,
}

impl IClientService for WasmClientService {
    fn call_view_by_hname(
        &self,
        chain_id: &ScChainID,
        contract_hname: &ScHname,
        function_hname: &ScHname,
        args: &[u8],
    ) -> errors::Result<Vec<u8>> {
        let params = ScDict::new(args);

        return self.client.call_view_by_hname(
            chain_id,
            contract_hname,
            function_hname,
            &params,
            None,
        );
    }

    fn post_request(
        &self,
        chain_id: &ScChainID,
        contract_hname: &ScHname,
        function_hname: &ScHname,
        args: &[u8],
        allowance: &ScAssets,
        key_pair: &keypair::KeyPair,
        nonce: u64,
    ) -> errors::Result<ScRequestID> {
        let params = ScDict::new(args);
        let mut req: offledgerrequest::OffLedgerRequestData =
            offledgerrequest::OffLedgerRequest::new(
                chain_id,
                contract_hname,
                function_hname,
                &params,
                nonce,
            );
        req.with_allowance(&allowance);
        req.sign(key_pair);
        self.client.post_offledger_request(&chain_id, &req)?;
        return Ok(req.id());
    }

    fn subscribe_events(
        &self,
        tx: mpsc::Sender<Vec<String>>,
        done: Arc<RwLock<bool>>,
    ) -> errors::Result<()> {
        self.websocket.clone().unwrap().subscribe(tx, done); // TODO remove clone
        return Ok(());
    }

    fn wait_until_request_processed(
        &self,
        chain_id: &ScChainID,
        req_id: &ScRequestID,
        timeout: Duration,
    ) -> errors::Result<()> {
        return self
            .client
            .wait_until_request_processed(&chain_id, req_id, timeout);
    }
}

impl WasmClientService {
    pub fn new(wasp_api: &str, event_port: &str) -> Self {
        return WasmClientService {
            client: waspclient::WaspClient::new(wasp_api),
            websocket: Some(websocket::Client::new(event_port).unwrap()),
            event_port: event_port.to_string(),
            last_err: Ok(()),
        };
    }
}

impl Default for WasmClientService {
    fn default() -> Self {
        return WasmClientService {
            client: waspclient::WaspClient::new("127.0.0.1:19090"),
            event_port: "127.0.0.1:15550".to_string(),
            websocket: None, // TODO set an empty object
            last_err: Ok(()),
        };
    }
}

#[cfg(test)]
mod tests {
    use crate::isc::{keypair, waspclient};
    use crate::{WasmClientContext, WasmClientService};

    #[test]
    fn service_default() {
        let service = WasmClientService::default();
        let default_service = WasmClientService {
            client: waspclient::WaspClient::new("127.0.0.1:19090"),
            websocket: None,
            event_port: "127.0.0.1:15550".to_string(),
            last_err: Ok(()),
        };
        assert!(default_service.event_port == service.event_port);
        assert!(default_service.last_err == Ok(()));
    }

    const MYCHAIN: &str = "tst1pqqf4qxh2w9x7rz2z4qqcvd0y8n22axsx82gqzmncvtsjqzwmhnjs438rhk";
    const MYSEED: &str = "0xa580555e5b84a4b72bbca829b4085a4725941f3b3702525f36862762d76c21f3";

    fn setup_client() -> WasmClientContext {
        let svc = WasmClientService::new("127.0.0.1:19090", "127.0.0.1:15550");
        let ctx =
            WasmClientContext::new(&svc, &wasmlib::chain_id_from_string(MYCHAIN), "testwasmlib");
        ctx.sign_requests(&keypair::KeyPair::from_sub_seed(
            &wasmlib::bytes_from_string(MYSEED),
            0,
        ));
        assert!(ctx.error.read().unwrap().is_ok());
        return ctx;
    }

    use testwasmlib::*;
    #[test]
    fn test_call_view() {
        let ctx = setup_client();

        let v = testwasmlib::ScFuncs::get_random(&ctx);
        // v.func.call();
        // expect(ctx.Err == null).toBeTruthy();
        // let rnd = v.results.random().value();
        // console.log("Rnd: " + rnd);
        // expect(rnd != 0n).toBeTruthy();
    }
}
