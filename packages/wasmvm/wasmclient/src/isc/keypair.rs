// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use std::convert::TryInto;

use crypto::hashes::{blake2b::Blake2b256, Digest};
use crypto::signatures::ed25519;
use wasmlib::*;

pub struct KeyPair {
    pub public_key: ed25519::PublicKey,
    private_key: ed25519::SecretKey,
}

impl KeyPair {
    pub fn sign(&self, data: &[u8]) -> Vec<u8> {
        return self.private_key.sign(data).to_bytes().to_vec();
    }

    pub fn from_sub_seed(seed: &[u8], n: u64) -> KeyPair {
        let index_bytes = uint64_to_bytes(n);
        let mut hash_of_index_bytes = Blake2b256::digest(index_bytes.to_owned());
        for i in 0..seed.len() {
            hash_of_index_bytes[i] ^= seed[i];
        }
        let public_key =
            ed25519::PublicKey::try_from_bytes(hash_of_index_bytes.try_into().unwrap()).unwrap();
        let private_key = ed25519::SecretKey::from_bytes(hash_of_index_bytes.try_into().unwrap());
        return KeyPair {
            public_key: public_key,
            private_key: private_key,
        };
    }
}

impl Clone for KeyPair {
    fn clone(&self) -> Self {
        return KeyPair {
            private_key: ed25519::SecretKey::from_bytes(self.private_key.to_bytes()),
            public_key: self.public_key.clone(),
        };
    }
}
