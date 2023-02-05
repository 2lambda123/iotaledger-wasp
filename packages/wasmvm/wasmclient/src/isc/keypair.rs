// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use crypto::{
    hashes::{blake2b::Blake2b256, Digest},
    keys::slip10::{Curve, Seed},
    signatures::ed25519,
};
// use crypto::keys::slip10::Key;
// use crypto::macs::hmac::HMAC_SHA512;
use wasmlib::*;

pub struct KeyPair {
    private_key: ed25519::SecretKey,
    pub public_key: ed25519::PublicKey,
}

// const SEED_KEY: [u8;64] = [
// 0x6a, 0x09, 0xe6, 0x67,
// 0xf3, 0xbc, 0xc9, 0x08,
// 0xbb, 0x67, 0xae, 0x85,
// 0x84, 0xca, 0xa7, 0x3b,
// 0x3c, 0x6e, 0xf3, 0x72,
// 0xfe, 0x94, 0xf8, 0x2b,
// 0xa5, 0x4f, 0xf5, 0x3a,
// 0x5f, 0x1d, 0x36, 0xf1,
// 0x51, 0x0e, 0x52, 0x7f,
// 0xad, 0xe6, 0x82, 0xd1,
// 0x9b, 0x05, 0x68, 0x8c,
// 0x2b, 0x3e, 0x6c, 0x1f,
// 0x1f, 0x83, 0xd9, 0xab,
// 0xfb, 0x41, 0xbd, 0x6b,
// 0x5b, 0xe0, 0xcd, 0x19,
// 0x13, 0x7e, 0x21, 0x79,
// ];

impl KeyPair {
    pub fn new(seed_bytes: &[u8]) -> KeyPair {
        if seed_bytes.len() == 0 {
            return KeyPair {
                private_key: ed25519::SecretKey::from_bytes([0; ed25519::SECRET_KEY_LENGTH]),
                public_key: ed25519::PublicKey::try_from_bytes([0; ed25519::PUBLIC_KEY_LENGTH])
                    .unwrap(),
            };
        }
        // let mut i = [0; 64];
        // HMAC_SHA512(seed_bytes, &SEED_KEY, &mut i);
        let seed = Seed::from_bytes(seed_bytes);
        let key = seed.to_master_key(Curve::Ed25519);
        let vec: Vec<u8> = key.into();
        println!("{}", bytes_to_string(&vec));
        return KeyPair {
            private_key: key.secret_key(),
            public_key: key.secret_key().public_key(),
        };
    }

    pub fn address(&self) -> ScAddress {
        let mut addr: Vec<u8> = Vec::with_capacity(wasmlib::SC_LENGTH_ED25519);
        addr[0] = wasmlib::SC_ADDRESS_ED25519;
        let hash = Blake2b256::digest(self.public_key.to_bytes());
        addr.copy_from_slice(&hash[..]);
        return wasmlib::address_from_bytes(&addr);
    }

    pub fn sign(&self, data: &[u8]) -> Vec<u8> {
        return self.private_key.sign(data).to_bytes().to_vec();
    }

    pub fn sub_seed(seed: &[u8], n: u64) -> Vec<u8> {
        let index_bytes = uint64_to_bytes(n);
        let mut hash_of_index_bytes = Blake2b256::digest(index_bytes.to_owned());
        for i in 0..seed.len() {
            hash_of_index_bytes[i] ^= seed[i];
        }
        hash_of_index_bytes.to_vec()
    }

    pub fn from_sub_seed(seed: &[u8], n: u64) -> KeyPair {
        let sub_seed = KeyPair::sub_seed(seed, n);
        return KeyPair::new(&sub_seed);
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

// impl Debug for KeyPair {
//     fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
//         f.debug_tuple("KeyPair").field(&self.public_key).finish()
//     }
// }

impl PartialEq for KeyPair {
    fn eq(&self, other: &Self) -> bool {
        return (self.public_key == other.public_key)
            && (self.private_key.as_slice() == other.private_key.as_slice());
    }
}

#[cfg(test)]
mod tests {
    use crate::keypair::KeyPair;
    use wasmlib::{bytes_from_string, bytes_to_string};

    const MYSEED: &str = "0xa580555e5b84a4b72bbca829b4085a4725941f3b3702525f36862762d76c21f3";

    #[test]
    fn keypair_sub_seed_0() {
        let my_seed = bytes_from_string(&MYSEED);
        let sub_seed = KeyPair::sub_seed(&my_seed, 0);
        println!("Seed: {}", bytes_to_string(&sub_seed));
        assert!(
            bytes_to_string(&sub_seed)
                == "0x24642f47bd363fbd4e05f13ed6c60b04c8a4cf1d295f76fc16917532bc4cd0af"
        );
    }

    #[test]
    fn keypair_sub_seed_1() {
        let my_seed = bytes_from_string(&MYSEED);
        let sub_seed = KeyPair::sub_seed(&my_seed, 1);
        println!("Seed: {}", bytes_to_string(&sub_seed));
        assert!(
            bytes_to_string(&sub_seed)
                == "0xb83d28550d9ee5651796eeb36027e737f0d79495b56d3d8931c716f2141017c8"
        );
    }

    #[test]
    fn keypair_construct() {
        let my_seed = bytes_from_string(&MYSEED);
        let pair = KeyPair::new(&my_seed);
        println!("Publ: {}", bytes_to_string(&pair.public_key.to_bytes()));
        println!("Priv: {}", bytes_to_string(&pair.private_key.to_bytes()));
        assert!(
            bytes_to_string(&pair.public_key.to_bytes())
                == "0x30adc0bd555d56ed51895528e47dcb403e36e0026fe49b6ae59e9adcea5f9a87"
        );
        assert!(bytes_to_string(&pair.private_key.to_bytes()) == "0xa580555e5b84a4b72bbca829b4085a4725941f3b3702525f36862762d76c21f330adc0bd555d56ed51895528e47dcb403e36e0026fe49b6ae59e9adcea5f9a87");
    }

    #[test]
    fn keypair_from_sub_seed_0() {
        let my_seed = bytes_from_string(&MYSEED);
        let pair = KeyPair::from_sub_seed(&my_seed, 0);
        println!("Publ: {}", bytes_to_string(&pair.public_key.to_bytes()));
        println!("Priv: {}", bytes_to_string(&pair.private_key.to_bytes()));
        assert!(
            bytes_to_string(&pair.public_key.to_bytes())
                == "0x40a757d26f6ef94dccee5b4f947faa78532286fe18117f2150a80acf2a95a8e2"
        );
        assert!(bytes_to_string(&pair.private_key.to_bytes()) == "0x24642f47bd363fbd4e05f13ed6c60b04c8a4cf1d295f76fc16917532bc4cd0af40a757d26f6ef94dccee5b4f947faa78532286fe18117f2150a80acf2a95a8e2");
    }

    #[test]
    fn keypair_from_sub_seed_1() {
        let my_seed = bytes_from_string(&MYSEED);
        let pair = KeyPair::from_sub_seed(&my_seed, 1);
        println!("Publ: {}", bytes_to_string(&pair.public_key.to_bytes()));
        println!("Priv: {}", bytes_to_string(&pair.private_key.to_bytes()));
        assert!(
            bytes_to_string(&pair.public_key.to_bytes())
                == "0x120d2b26fc1b1d53bb916b8a277bcc2efa09e92c95be1a8fd5c6b3adbc795679"
        );
        assert!(bytes_to_string(&pair.private_key.to_bytes()) == "0xb83d28550d9ee5651796eeb36027e737f0d79495b56d3d8931c716f2141017c8120d2b26fc1b1d53bb916b8a277bcc2efa09e92c95be1a8fd5c6b3adbc795679");
    }
}
