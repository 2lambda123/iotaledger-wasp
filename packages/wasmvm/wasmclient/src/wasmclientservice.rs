// // Copyright 2020 IOTA Stiftung
// // SPDX-License-Identifier: Apache-2.0

use std::time::Duration;

use wasmlib::*;

use isc::offledgerrequest::*;

use crate::*;
use crate::keypair::KeyPair;
use crate::waspclient::WaspClient;

#[derive(Clone, PartialEq)]
pub struct WasmClientService {
    client: WaspClient,
    last_err: errors::Result<()>,
    web_socket: String,
}

impl WasmClientService {
    pub fn new(wasp_api: &str) -> Self {
        return WasmClientService {
            client: WaspClient::new(wasp_api),
            last_err: Ok(()),
            web_socket: wasp_api.replace("http:", "ws:") + "/ws",
        };
    }

    pub fn call_view_by_hname(
        &self,
        chain_id: &ScChainID,
        contract_hname: &ScHname,
        function_hname: &ScHname,
        args: &[u8],
    ) -> errors::Result<Vec<u8>> {
        return self.client.call_view_by_hname(
            chain_id,
            contract_hname,
            function_hname,
            args,
        );
    }

    pub fn post_request(
        &self,
        chain_id: &ScChainID,
        h_contract: &ScHname,
        h_function: &ScHname,
        args: &[u8],
        allowance: &ScAssets,
        key_pair: &KeyPair,
        nonce: u64,
    ) -> errors::Result<ScRequestID> {
        let mut req: OffLedgerRequestData =
            OffLedgerRequest::new(
                chain_id,
                h_contract,
                h_function,
                args,
                nonce,
            );
        req.with_allowance(&allowance);
        let signed = req.sign(key_pair);
        let res = self.client.post_offledger_request(&chain_id, &signed);
        if let Err(e) = res {
            return Err(e);
        }
        Ok(signed.id())
    }

    pub fn wait_until_request_processed(
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

impl Default for WasmClientService {
    fn default() -> Self {
        return WasmClientService {
            client: WaspClient::new("http://localhost:19090"),
            last_err: Ok(()),
            web_socket: String::from("ws://localhost:19090/ws"),
        };
    }
}
