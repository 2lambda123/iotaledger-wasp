// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

import * as wasmtypes from "../wasmtypes";

export class CallRequest {
    // caller assets that the call is allowed to access
    allowance : Uint8Array = new Uint8Array(0);
    contract  : wasmtypes.ScHname = new wasmtypes.ScHname(0);
    function  : wasmtypes.ScHname = new wasmtypes.ScHname(0);
    params    : Uint8Array = new Uint8Array(0);

    static fromBytes(buf: Uint8Array): CallRequest {
        const dec = new wasmtypes.WasmDecoder(buf);
        const data = new CallRequest();
        data.allowance = wasmtypes.bytesDecode(dec);
        data.contract  = wasmtypes.hnameDecode(dec);
        data.function  = wasmtypes.hnameDecode(dec);
        data.params    = wasmtypes.bytesDecode(dec);
        dec.close();
        return data;
    }

    bytes(): Uint8Array {
        const enc = new wasmtypes.WasmEncoder();
        wasmtypes.bytesEncode(enc, this.allowance);
        wasmtypes.hnameEncode(enc, this.contract);
        wasmtypes.hnameEncode(enc, this.function);
        wasmtypes.bytesEncode(enc, this.params);
        return enc.buf();
    }
}

export class ImmutableCallRequest extends wasmtypes.ScProxy {

    exists(): bool {
        return this.proxy.exists();
    }

    value(): CallRequest {
        return CallRequest.fromBytes(this.proxy.get());
    }
}

export class MutableCallRequest extends wasmtypes.ScProxy {

    delete(): void {
        this.proxy.delete();
    }

    exists(): bool {
        return this.proxy.exists();
    }

    setValue(value: CallRequest): void {
        this.proxy.set(value.bytes());
    }

    value(): CallRequest {
        return CallRequest.fromBytes(this.proxy.get());
    }
}

export class DeployRequest {
    description : string = "";
    name        : string = "";
    params      : Uint8Array = new Uint8Array(0);
    progHash    : wasmtypes.ScHash = new wasmtypes.ScHash();

    static fromBytes(buf: Uint8Array): DeployRequest {
        const dec = new wasmtypes.WasmDecoder(buf);
        const data = new DeployRequest();
        data.description = wasmtypes.stringDecode(dec);
        data.name        = wasmtypes.stringDecode(dec);
        data.params      = wasmtypes.bytesDecode(dec);
        data.progHash    = wasmtypes.hashDecode(dec);
        dec.close();
        return data;
    }

    bytes(): Uint8Array {
        const enc = new wasmtypes.WasmEncoder();
        wasmtypes.stringEncode(enc, this.description);
        wasmtypes.stringEncode(enc, this.name);
        wasmtypes.bytesEncode(enc, this.params);
        wasmtypes.hashEncode(enc, this.progHash);
        return enc.buf();
    }
}

export class ImmutableDeployRequest extends wasmtypes.ScProxy {

    exists(): bool {
        return this.proxy.exists();
    }

    value(): DeployRequest {
        return DeployRequest.fromBytes(this.proxy.get());
    }
}

export class MutableDeployRequest extends wasmtypes.ScProxy {

    delete(): void {
        this.proxy.delete();
    }

    exists(): bool {
        return this.proxy.exists();
    }

    setValue(value: DeployRequest): void {
        this.proxy.set(value.bytes());
    }

    value(): DeployRequest {
        return DeployRequest.fromBytes(this.proxy.get());
    }
}

export class PostRequest {
    // caller assets that the call is allowed to access
    allowance : Uint8Array = new Uint8Array(0);
    chainID   : wasmtypes.ScChainID = new wasmtypes.ScChainID();
    contract  : wasmtypes.ScHname = new wasmtypes.ScHname(0);
    delay     : u32 = 0;
    function  : wasmtypes.ScHname = new wasmtypes.ScHname(0);
    params    : Uint8Array = new Uint8Array(0);
    // assets that are transferred into caller account
    transfer  : Uint8Array = new Uint8Array(0);

    static fromBytes(buf: Uint8Array): PostRequest {
        const dec = new wasmtypes.WasmDecoder(buf);
        const data = new PostRequest();
        data.allowance = wasmtypes.bytesDecode(dec);
        data.chainID   = wasmtypes.chainIDDecode(dec);
        data.contract  = wasmtypes.hnameDecode(dec);
        data.delay     = wasmtypes.uint32Decode(dec);
        data.function  = wasmtypes.hnameDecode(dec);
        data.params    = wasmtypes.bytesDecode(dec);
        data.transfer  = wasmtypes.bytesDecode(dec);
        dec.close();
        return data;
    }

    bytes(): Uint8Array {
        const enc = new wasmtypes.WasmEncoder();
        wasmtypes.bytesEncode(enc, this.allowance);
        wasmtypes.chainIDEncode(enc, this.chainID);
        wasmtypes.hnameEncode(enc, this.contract);
        wasmtypes.uint32Encode(enc, this.delay);
        wasmtypes.hnameEncode(enc, this.function);
        wasmtypes.bytesEncode(enc, this.params);
        wasmtypes.bytesEncode(enc, this.transfer);
        return enc.buf();
    }
}

export class ImmutablePostRequest extends wasmtypes.ScProxy {

    exists(): bool {
        return this.proxy.exists();
    }

    value(): PostRequest {
        return PostRequest.fromBytes(this.proxy.get());
    }
}

export class MutablePostRequest extends wasmtypes.ScProxy {

    delete(): void {
        this.proxy.delete();
    }

    exists(): bool {
        return this.proxy.exists();
    }

    setValue(value: PostRequest): void {
        this.proxy.set(value.bytes());
    }

    value(): PostRequest {
        return PostRequest.fromBytes(this.proxy.get());
    }
}

export class SendRequest {
    address  : wasmtypes.ScAddress = new wasmtypes.ScAddress();
    transfer : Uint8Array = new Uint8Array(0);

    static fromBytes(buf: Uint8Array): SendRequest {
        const dec = new wasmtypes.WasmDecoder(buf);
        const data = new SendRequest();
        data.address  = wasmtypes.addressDecode(dec);
        data.transfer = wasmtypes.bytesDecode(dec);
        dec.close();
        return data;
    }

    bytes(): Uint8Array {
        const enc = new wasmtypes.WasmEncoder();
        wasmtypes.addressEncode(enc, this.address);
        wasmtypes.bytesEncode(enc, this.transfer);
        return enc.buf();
    }
}

export class ImmutableSendRequest extends wasmtypes.ScProxy {

    exists(): bool {
        return this.proxy.exists();
    }

    value(): SendRequest {
        return SendRequest.fromBytes(this.proxy.get());
    }
}

export class MutableSendRequest extends wasmtypes.ScProxy {

    delete(): void {
        this.proxy.delete();
    }

    exists(): bool {
        return this.proxy.exists();
    }

    setValue(value: SendRequest): void {
        this.proxy.set(value.bytes());
    }

    value(): SendRequest {
        return SendRequest.fromBytes(this.proxy.get());
    }
}

export class TransferRequest {
    agentID  : wasmtypes.ScAgentID = wasmtypes.agentIDFromBytes(new Uint8Array(0));
    create   : bool = false;
    transfer : Uint8Array = new Uint8Array(0);

    static fromBytes(buf: Uint8Array): TransferRequest {
        const dec = new wasmtypes.WasmDecoder(buf);
        const data = new TransferRequest();
        data.agentID  = wasmtypes.agentIDDecode(dec);
        data.create   = wasmtypes.boolDecode(dec);
        data.transfer = wasmtypes.bytesDecode(dec);
        dec.close();
        return data;
    }

    bytes(): Uint8Array {
        const enc = new wasmtypes.WasmEncoder();
        wasmtypes.agentIDEncode(enc, this.agentID);
        wasmtypes.boolEncode(enc, this.create);
        wasmtypes.bytesEncode(enc, this.transfer);
        return enc.buf();
    }
}

export class ImmutableTransferRequest extends wasmtypes.ScProxy {

    exists(): bool {
        return this.proxy.exists();
    }

    value(): TransferRequest {
        return TransferRequest.fromBytes(this.proxy.get());
    }
}

export class MutableTransferRequest extends wasmtypes.ScProxy {

    delete(): void {
        this.proxy.delete();
    }

    exists(): bool {
        return this.proxy.exists();
    }

    setValue(value: TransferRequest): void {
        this.proxy.set(value.bytes());
    }

    value(): TransferRequest {
        return TransferRequest.fromBytes(this.proxy.get());
    }
}
