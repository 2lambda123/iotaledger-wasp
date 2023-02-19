// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmlib from 'wasmlib';
import {SyncRequestClient} from './ts-sync-request';
import {OffLedgerRequest} from './offledgerrequest';
import {APICallViewRequest, APIOffLedgerRequest, Codec, JsonReq, JsonResp} from './codec';

export type Error = string | null;

export class WaspClient {
    baseURL: string;

    public constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    public callViewByHname(chainID: wasmlib.ScChainID, hContract: wasmlib.ScHname, hFunction: wasmlib.ScHname, args: Uint8Array): [Uint8Array, Error] {
        const url = this.baseURL + '/requests/callview';
        const req = new SyncRequestClient();
        req.addHeader('Content-Type', 'application/json');

        const callViewRequest: APICallViewRequest = {
            contractHName: hContract.toString(),
            functionHName: hFunction.toString(),
            chainId: chainID.toString(),
            arguments: Codec.jsonEncode(args),
        };

        try {
            const resp = req.post<APICallViewRequest, JsonResp>(url, callViewRequest);
            const result = Codec.jsonDecode(resp);
            return [result, null];
        } catch (error) {
            let message;
            if (error instanceof Error) message = error.message;
            else message = String(error);
            return [new Uint8Array(0), message];
        }
    }

    public postOffLedgerRequest(chainID: wasmlib.ScChainID, signed: OffLedgerRequest): Error {
        const url = this.baseURL + '/requests/offledger';
        const req = new SyncRequestClient();
        req.addHeader('Content-Type', 'application/json');

        const offLedgerRequest: APIOffLedgerRequest = {
            chainId: chainID.toString(),
            request: wasmlib.hexEncode(signed.bytes()),
        };

        try {
            req.post(url, offLedgerRequest);
            return null;
        } catch (error) {
            let message;
            if (error instanceof Error) message = error.message;
            else message = String(error);
            return message;
        }
    }

    public waitUntilRequestProcessed(chainID: wasmlib.ScChainID, reqID: wasmlib.ScRequestID, timeout: u32): Error {
        //TODO Timeout of the wait can be set with `/wait?timeoutSeconds=`. Max seconds are 60secs.
        const url = this.baseURL + '/chains/' + chainID.toString() + '/requests/' + reqID.toString() + '/wait';
        const response = new SyncRequestClient().get(url);
        return null;
    }
}
