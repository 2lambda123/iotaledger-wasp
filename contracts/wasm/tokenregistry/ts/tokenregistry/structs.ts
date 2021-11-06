// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "wasmlib";

export class Token {
    created: i64 = 0;  // creation timestamp
    description: string = "";  // description what minted token represents
    mintedBy: wasmlib.ScAgentID = new wasmlib.ScAgentID();  // original minter
    owner: wasmlib.ScAgentID = new wasmlib.ScAgentID();  // current owner
    supply: i64 = 0;  // amount of tokens originally minted
    updated: i64 = 0;  // last update timestamp
    userDefined: string = "";  // any user defined text

    static fromBytes(bytes: u8[]): Token {
        let decode = new wasmlib.BytesDecoder(bytes);
        let data = new Token();
        data.created = decode.int64();
        data.description = decode.string();
        data.mintedBy = decode.agentID();
        data.owner = decode.agentID();
        data.supply = decode.int64();
        data.updated = decode.int64();
        data.userDefined = decode.string();
        decode.close();
        return data;
    }

    bytes(): u8[] {
        return new wasmlib.BytesEncoder().
		    int64(this.created).
		    string(this.description).
		    agentID(this.mintedBy).
		    agentID(this.owner).
		    int64(this.supply).
		    int64(this.updated).
		    string(this.userDefined).
            data();
    }
}

export class ImmutableToken {
    objID: i32;
    keyID: wasmlib.Key32;

    constructor(objID: i32, keyID: wasmlib.Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    exists(): boolean {
        return wasmlib.exists(this.objID, this.keyID, wasmlib.TYPE_BYTES);
    }

    value(): Token {
        return Token.fromBytes(wasmlib.getBytes(this.objID, this.keyID, wasmlib.TYPE_BYTES));
    }
}

export class MutableToken {
    objID: i32;
    keyID: wasmlib.Key32;

    constructor(objID: i32, keyID: wasmlib.Key32) {
        this.objID = objID;
        this.keyID = keyID;
    }

    exists(): boolean {
        return wasmlib.exists(this.objID, this.keyID, wasmlib.TYPE_BYTES);
    }

    setValue(value: Token): void {
        wasmlib.setBytes(this.objID, this.keyID, wasmlib.TYPE_BYTES, value.bytes());
    }

    value(): Token {
        return Token.fromBytes(wasmlib.getBytes(this.objID, this.keyID, wasmlib.TYPE_BYTES));
    }
}
