pub use crate::offledgerrequest::*;
pub use crate::receipt::*;
pub use std::time::*;
pub use wasmlib::*;

const DEFAULT_OPTIMISTIC_READ_TIMEOUT: Duration = Duration::from_millis(1100);

pub struct WaspClient {
    base_url: String,
    token: String,
}

impl WaspClient {
    pub fn new(base_url: &str) -> WaspClient {
        return WaspClient {
            base_url: base_url.to_string(),
            token: String::from(""),
        };
    }
    pub fn call_view_by_hname(
        &self,
        chain_id: &ScChainID,
        contract_hname: &ScHname,
        function_hname: &ScHname,
        args: ScDict,
        optimistic_read_timeout: Option<Duration>,
    ) -> Result<(), String> {
        let deadline = match optimistic_read_timeout {
            Some(duration) => duration,
            None => DEFAULT_OPTIMISTIC_READ_TIMEOUT,
        };

        let url = format!(
            "/chain/{}/contract/{}/callviewbyhname/{}",
            chain_id.to_string(),
            contract_hname.to_string(),
            function_hname.to_string()
        );
        let client = reqwest::blocking::Client::builder()
            .timeout(deadline)
            .build()
            .unwrap();
        let _ = client.post(url).body(args.to_bytes()).send();
        Ok(())
    }
    pub fn post_offledger_request(
        &self,
        chain_id: &ScChainID,
        req: &OffLedgerRequestData,
    ) -> Result<(), String> {
        todo!()
    }
    pub fn wait_until_request_processed(
        &self,
        chain_id: &ScChainID,
        req_id: &ScRequestID,
        timeout: Duration,
    ) -> Result<Receipt, String> {
        todo!()
    }
}

// fn send_request(method: &str, route: &str) -> Result<Vec<u8>, String> {
//     todo!()
// }
