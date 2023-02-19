// // Copyright 2020 IOTA Stiftung
// // SPDX-License-Identifier: Apache-2.0

use std::time::Duration;

use reqwest::{blocking, StatusCode};
use wasmlib::*;

use codec::*;

use crate::*;

pub const ISC_EVENT_KIND_NEW_BLOCK: &str = "new_block";
pub const ISC_EVENT_KIND_RECEIPT: &str = "receipt";
pub const ISC_EVENT_KIND_SMART_CONTRACT: &str = "contract";
pub const ISC_EVENT_KIND_ERROR: &str = "error";

const READ_TIMEOUT: Duration = Duration::from_millis(10000);

#[derive(Clone, PartialEq, Default)]
pub struct WasmClientService {
    wasp_api: String,
}

impl WasmClientService {
    pub fn new(wasp_api: &str) -> Self {
        return WasmClientService {
            wasp_api: String::from(wasp_api),
        };
    }

    pub fn call_view_by_hname(
        &self,
        chain_id: &ScChainID,
        contract_hname: &ScHname,
        function_hname: &ScHname,
        args: &[u8],
    ) -> Result<Vec<u8>> {
        let url = format!("{}/requests/callview", self.wasp_api);
        let client = blocking::Client::builder()
            .timeout(READ_TIMEOUT)
            .build()
            .unwrap();
        let body = APICallViewRequest {
            arguments: json_encode(args),
            chain_id: chain_id.to_string(),
            contract_hname: contract_hname.to_string(),
            function_hname: function_hname.to_string(),
        };
        let res = client.post(url).json(&body).send();
        match res {
            Ok(v) => match v.status() {
                StatusCode::OK => {
                    match v.json::<JsonResponse>() {
                        Ok(json_obj) => {
                            return Ok(json_decode(json_obj));
                        }
                        Err(e) => {
                            return Err(format!("parse post response failed: {}", e.to_string()));
                        }
                    };
                }
                failed_status_code => {
                    let status_code = failed_status_code.as_u16();
                    match v.json::<JsonError>() {
                        Ok(err_msg) => {
                            return Err(format!("{status_code}: {}", err_msg.message));
                        }
                        Err(e) => return Err(e.to_string()),
                    }
                }
            },
            Err(e) => {
                return Err(format!("call() request failed: {}", e.to_string()));
            }
        }
    }

    pub fn post_request(
        &self,
        chain_id: &ScChainID,
        h_contract: &ScHname,
        h_function: &ScHname,
        args: &[u8],
        allowance: &ScAssets,
        key_pair: &keypair::KeyPair,
        nonce: u64,
    ) -> Result<ScRequestID> {
        let mut req =
            offledgerrequest::OffLedgerRequest::new(
                chain_id,
                h_contract,
                h_function,
                args,
                nonce,
            );
        req.with_allowance(&allowance);
        let signed = req.sign(key_pair);

        let url = format!("{}/requests/offledger", self.wasp_api);
        let client = blocking::Client::new();
        let body = APIOffLedgerRequest {
            chain_id: chain_id.to_string(),
            request: hex_encode(&signed.to_bytes()),
        };
        let res = client.post(url).json(&body).send();
        match res {
            Ok(v) => match v.status() {
                StatusCode::OK => {}
                StatusCode::ACCEPTED => {}
                failed_status_code => {
                    let status_code = failed_status_code.as_u16();
                    match v.json::<JsonError>() {
                        Ok(err_msg) => {
                            return Err(format!("{status_code}: {}", err_msg.message));
                        }
                        Err(e) => return Err(e.to_string()),
                    }
                }
            },
            Err(e) => {
                return Err(format!("post() request failed: {}", e.to_string()));
            }
        }
        Ok(signed.id())
    }

    pub fn wait_until_request_processed(
        &self,
        chain_id: &ScChainID,
        req_id: &ScRequestID,
        timeout: Duration,
    ) -> Result<()> {
        let url = format!(
            "{}/chains/{}/requests/{}/wait",
            self.wasp_api,
            chain_id.to_string(),
            req_id.to_string()
        );
        let client = blocking::Client::builder()
            .timeout(timeout)
            .build()
            .unwrap();
        let res = client.get(url).header("Content-Type", "application/json").send();
        match res {
            Ok(v) => match v.status() {
                StatusCode::OK => {
                    return Ok(());
                }
                failed_status_code => {
                    let status_code = failed_status_code.as_u16();
                    match v.text() {
                        Ok(err_msg) => {
                            return Err(format!("{status_code}: {err_msg}"));
                        }
                        Err(e) => return Err(e.to_string()),
                    }
                }
            },
            Err(e) => {
                return Err(format!("request failed: {}", e.to_string()));
            }
        }
    }
}
