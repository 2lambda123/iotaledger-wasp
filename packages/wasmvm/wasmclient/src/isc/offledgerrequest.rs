// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use crate::keypair;
use wasmlib::*;

//TODO generalize this trait
pub trait OffLedgerRequest<'a> {
    fn new(
        chain_id: &ScChainID,
        contract: &ScHname,
        entry_point: &ScHname,
        params: &ScDict,
        signature_scheme: Option<&'a OffLedgerSignatureScheme>,
        nonce: u64,
    ) -> Self;
    fn with_nonce(&mut self, nonce: u64) -> &Self;
    fn with_gas_budget(&mut self, gas_budget: u64) -> &Self;
    fn with_allowance(&mut self, allowance: &ScAssets) -> &Self;
    fn sign(&mut self, key: &keypair::KeyPair) -> &Self;
}

#[derive(Clone)]
pub struct OffLedgerRequestData<'a> {
    chain_id: ScChainID,
    contract: ScHname,
    entry_point: ScHname,
    params: ScDict,
    signature_scheme: Option<&'a OffLedgerSignatureScheme>, // None if unsigned
    nonce: u64,
    allowance: ScAssets,
    gas_budget: u64,
}

pub struct OffLedgerSignatureScheme {
    key_pair: keypair::KeyPair,
    signature: Vec<u8>,
}

impl OffLedgerSignatureScheme {
    pub fn clone(&self) -> Self {
        todo!()
    }
}

impl<'a> OffLedgerRequest<'a> for OffLedgerRequestData<'a> {
    fn new(
        chain_id: &ScChainID,
        contract: &ScHname,
        entry_point: &ScHname,
        params: &ScDict,
        signature_scheme: Option<&'a OffLedgerSignatureScheme>,
        nonce: u64,
    ) -> Self {
        return OffLedgerRequestData {
            chain_id: chain_id.clone(),
            contract: contract.clone(),
            entry_point: entry_point.clone(),
            params: params.clone(),
            signature_scheme: match signature_scheme {
                Some(signature_scheme_val) => Some(signature_scheme_val.to_owned()),
                None => None,
            },
            nonce: nonce,
            allowance: ScAssets::new(&Vec::new()),
            gas_budget: super::gas::MAX_GAS_PER_REQUEST,
        };
    }
    fn with_nonce(&mut self, nonce: u64) -> &Self {
        self.nonce = nonce;
        return self;
    }
    fn with_gas_budget(&mut self, gas_budget: u64) -> &Self {
        self.gas_budget = gas_budget;
        return self;
    }
    fn with_allowance(&mut self, allowance: &ScAssets) -> &Self {
        self.allowance = allowance.clone();
        return self;
    }
    fn sign(&mut self, _key: &keypair::KeyPair) -> &Self {
        todo!()
    }
}

impl<'a> OffLedgerRequestData<'a> {
    pub fn id(&self) -> ScRequestID {
        todo!()
    }
    pub fn essence(&self) -> Vec<u8> {
        let mut data: Vec<u8> = vec![1];
        data.append(self.chain_id.to_bytes().as_mut());
        data.append(self.contract.to_bytes().as_mut());
        data.append(self.entry_point.to_bytes().as_mut());
        data.append(self.params.to_bytes().as_mut());
        data.append(wasmlib::uint64_to_bytes(self.nonce).as_mut());
        data.append(wasmlib::uint64_to_bytes(self.gas_budget).as_mut());
        let scheme = match self.signature_scheme {
            Some(val) => val.clone(),
            None => {
                panic!("none")
            }
        };
        let mut public_key = scheme.key_pair.public_key.to_bytes().to_vec();
        data.push(public_key.len() as u8);
        data.append(&mut public_key);
        data.append(self.allowance.to_bytes().as_mut());
        return data;
    }
    pub fn to_bytes(&self) -> Vec<u8> {
        todo!()
    }
}
