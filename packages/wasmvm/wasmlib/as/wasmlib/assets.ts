// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import {ScDict} from './dict';
import {ScTokenID, tokenIDDecode, tokenIDEncode, tokenIDFromBytes} from './wasmtypes/sctokenid';
import {uint64Decode, uint64Encode} from './wasmtypes/scuint64';
import {bigIntDecode, bigIntEncode, ScBigInt} from './wasmtypes/scbigint';
import {nftIDDecode, nftIDEncode, ScNftID} from './wasmtypes/scnftid';
import {WasmDecoder, WasmEncoder} from './wasmtypes/codec';
import {uint16Decode, uint16Encode} from './wasmtypes/scuint16';
import {boolDecode, boolEncode} from './wasmtypes/scbool';

export class ScAssets {
    baseTokens: u64 = 0;
    nativeTokens: Map<string, ScBigInt> = new Map();
    nftIDs: Set<ScNftID> = new Set();

    public constructor(buf: Uint8Array | null) {
        if (buf === null || buf.length == 0) {
            return this;
        }

        const dec = new WasmDecoder(buf);
        const empty = boolDecode(dec);
        if (empty) {
            return this;
        }

        this.baseTokens = uint64Decode(dec);

        let size = dec.vluDecode(32);
        for (let i: u64 = 0; i < size; i++) {
            const tokenID = tokenIDDecode(dec);
            const amount = bigIntDecode(dec);
            this.nativeTokens.set(ScDict.toKey(tokenID.id), amount);
        }

        size = dec.vluDecode(32);
        for (let i: u64 = 0; i < size; i++) {
            const nftID = nftIDDecode(dec);
            this.nftIDs.add(nftID);
        }
    }

    public balances(): ScBalances {
        return new ScBalances(this);
    }

    public isEmpty(): bool {
        if (this.baseTokens != 0) {
            return false;
        }
        const values = this.nativeTokens.values();
        for (let i = 0; i < values.length; i++) {
            if (!values[i].isZero()) {
                return false;
            }
        }
        return this.nftIDs.size == 0;
    }

    public toBytes(): Uint8Array {
        const enc = new WasmEncoder();
        const empty = this.isEmpty();
        boolEncode(enc, empty);
        if (empty) {
            return enc.buf();
        }

        uint64Encode(enc, this.baseTokens);

        const tokenIDs = this.tokenIDs();
        enc.vluEncode(tokenIDs.length as u64);
        for (let i = 0; i < tokenIDs.length; i++) {
            const tokenID = tokenIDs[i];
            tokenIDEncode(enc, tokenID);
            const mapKey = ScDict.toKey(tokenID.id);
            const amount = this.nativeTokens.get(mapKey);
            bigIntEncode(enc, amount);
        }

        enc.vluEncode(this.nftIDs.size as u64);
        const arr = this.nftIDs.values();
        for (let i = 0; i < arr.length; i++) {
            const nftID = arr[i];
            nftIDEncode(enc, nftID);
        }
        return enc.buf();
    }

    public tokenIDs(): ScTokenID[] {
        const tokenIDs: ScTokenID[] = [];
        const keys = this.nativeTokens.keys().sort();
        for (let i = 0; i < keys.length; i++) {
            const keyBytes = ScDict.fromKey(keys[i]);
            const tokenID = tokenIDFromBytes(keyBytes);
            tokenIDs.push(tokenID);
        }
        return tokenIDs;
    }
}

export class ScBalances {
    assets: ScAssets;

    constructor(assets: ScAssets) {
        this.assets = assets;
    }

    public balance(tokenID: ScTokenID): ScBigInt {
        const mapKey = ScDict.toKey(tokenID.id);
        if (!this.assets.nativeTokens.has(mapKey)) {
            return new ScBigInt();
        }
        return this.assets.nativeTokens.get(mapKey);
    }

    public baseTokens(): u64 {
        return this.assets.baseTokens;
    }

    public isEmpty(): bool {
        return this.assets.isEmpty();
    }

    public nftIDs(): Set<ScNftID> {
        return this.assets.nftIDs;
    }

    public toBytes(): Uint8Array {
        return this.assets.toBytes();
    }

    public tokenIDs(): ScTokenID[] {
        return this.assets.tokenIDs();
    }
}

export class ScTransfer extends ScBalances {
    public constructor() {
        super(new ScAssets(null));
    }

    public static fromBalances(balances: ScBalances): ScTransfer {
        const transfer = ScTransfer.baseTokens(balances.baseTokens());
        const tokenIDs = balances.tokenIDs();
        for (let i = 0; i < tokenIDs.length; i++) {
            const tokenID = tokenIDs[i];
            transfer.set(tokenID, balances.balance(tokenID));
        }
        const nftIDs = balances.nftIDs().values();
        for (let i = 0; i < nftIDs.length; i++) {
            transfer.addNFT(nftIDs[i]);
        }
        return transfer;
    }

    public static baseTokens(amount: u64): ScTransfer {
        const transfer = new ScTransfer();
        transfer.assets.baseTokens = amount;
        return transfer;
    }

    public static nft(nftID: ScNftID): ScTransfer {
        const transfer = new ScTransfer();
        transfer.addNFT(nftID);
        return transfer;
    }

    public static tokens(tokenID: ScTokenID, amount: ScBigInt): ScTransfer {
        const transfer = new ScTransfer();
        transfer.set(tokenID, amount);
        return transfer;
    }

    public addNFT(nftID: ScNftID): void {
        this.assets.nftIDs.add(nftID);
    }

    public set(tokenID: ScTokenID, amount: ScBigInt): void {
        const mapKey = ScDict.toKey(tokenID.id);
        this.assets.nativeTokens.set(mapKey, amount);
    }
}
