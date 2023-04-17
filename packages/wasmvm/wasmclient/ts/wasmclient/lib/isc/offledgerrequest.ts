// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import {Blake2b} from '@iota/crypto.js';
import * as wasmlib from 'wasmlib';
import {concat} from 'wasmlib';
import {KeyPair} from './keypair';

export class OffLedgerSignatureScheme {
    keyPair: KeyPair;
    signature: Uint8Array;

    public constructor(keyPair: KeyPair) {
        this.keyPair = keyPair;
        this.signature = new Uint8Array(0);
    }
}

export class OffLedgerRequest {
    chainID: wasmlib.ScChainID;
    contract: wasmlib.ScHname;
    entryPoint: wasmlib.ScHname;
    params: Uint8Array;
    signatureScheme: OffLedgerSignatureScheme = new OffLedgerSignatureScheme(new KeyPair(new Uint8Array(0)));
    nonce: u64;
    allowance: wasmlib.ScAssets = new wasmlib.ScAssets(new Uint8Array(0));
    gasBudget: u64 = 2n ** 64n - 1n;

    public constructor(chainID: wasmlib.ScChainID, contract: wasmlib.ScHname, entryPoint: wasmlib.ScHname, params: Uint8Array, nonce: u64) {
        this.chainID = chainID;
        this.contract = contract;
        this.entryPoint = entryPoint;
        this.params = params;
        this.nonce = nonce;
    }

    public bytes(): Uint8Array {
        const sig = this.signatureScheme.signature;
        const data = concat(this.essence(), wasmlib.uint16ToBytes(sig.length as u16));
        return concat(data, sig);
    }

    public essence(): Uint8Array {
        const oneByte = new Uint8Array(1);
        oneByte[0] = 1; // requestKindTagOffLedgerISC
        let data = concat(oneByte, this.chainID.toBytes());
        data = concat(data, this.contract.toBytes());
        data = concat(data, this.entryPoint.toBytes());
        data = concat(data, this.params);
        data = concat(data, wasmlib.uint64ToBytes(this.nonce));
        data = concat(data, wasmlib.uint64ToBytes(this.gasBudget));
        const pubKey = this.signatureScheme.keyPair.publicKey;
        oneByte[0] = pubKey.length as u8;
        data = concat(data, oneByte);
        data = concat(data, pubKey);
        data = concat(data, this.allowance.toBytes());
        return data;
    }

    public ID(): wasmlib.ScRequestID {
        // req id is hash of req bytes with output index zero
        const hash = Blake2b.sum256(this.bytes());
        const reqId = new wasmlib.ScRequestID();
        reqId.id.set(hash, 0);
        return reqId;
    }

    public sign(keyPair: KeyPair): OffLedgerRequest {
        const req = new OffLedgerRequest(this.chainID, this.contract, this.entryPoint, this.params, this.nonce);
        req.signatureScheme = new OffLedgerSignatureScheme(keyPair);
        req.signatureScheme.signature = keyPair.sign(req.essence());
        return req;
    }

    public withAllowance(allowance: wasmlib.ScAssets): void {
        this.allowance = allowance;
    }
}
