// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as isc from './isc';
import * as wasmlib from 'wasmlib';
import {WebSocket} from 'ws';
import {WasmClientContext} from './wasmclientcontext';

export class ContractEvent {
    chainID = '';
    contractID = '';
    data = '';
}

type ClientCallBack = (event: ContractEvent) => void;

export class WasmClientService {
    private callbacks: ClientCallBack[] = [];
    private ws: WebSocket;
    private subscribers: WasmClientContext[] = [];
    private waspClient: isc.WaspClient;

    public constructor(waspAPI: string) {
        this.waspClient = new isc.WaspClient(waspAPI);
        const eventPort = waspAPI.replace('http:', 'ws:') + '/ws';
        this.ws = new WebSocket(eventPort, {
            perMessageDeflate: false
        });
    }

    public callViewByHname(chainID: wasmlib.ScChainID, hContract: wasmlib.ScHname, hFunction: wasmlib.ScHname, args: Uint8Array): [Uint8Array, isc.Error] {
        return this.waspClient.callViewByHname(chainID, hContract, hFunction, args);
    }

    public postRequest(chainID: wasmlib.ScChainID, hContract: wasmlib.ScHname, hFunction: wasmlib.ScHname, args: Uint8Array, allowance: wasmlib.ScAssets, keyPair: isc.KeyPair, nonce: u64): [wasmlib.ScRequestID, isc.Error] {
        const req = new isc.OffLedgerRequest(chainID, hContract, hFunction, args, nonce);
        req.withAllowance(allowance);
        const signed = req.sign(keyPair);
        const reqID = signed.ID();
        const err = this.waspClient.postOffLedgerRequest(chainID, signed);
        return [reqID, err];
    }

    subscribe(topic: string) {
        const msg = {
            command: 'subscribe',
            topic: topic,
        };
        const rawMsg = JSON.stringify(msg);
        this.ws.send(rawMsg);
    }

    public subscribeEvents(who: WasmClientContext, callback: ClientCallBack): isc.Error {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const self = this;
        this.subscribers.push(who);
        this.callbacks.push(callback);
        if (this.subscribers.length == 1) {
            this.ws.on('open', () => {
                self.subscribe('chains');
                self.subscribe('contract');
            });
            this.ws.on('error', (err) => {
                // callback(['error', err.toString()]);
            });
            this.ws.on('message', (data) => {
                let msg: any;
                try {
                    msg = JSON.parse(data.toString());
                    if (!msg.Kind) {
                        // filter out subscribe responses
                        return;
                    }
                    console.log(msg);
                } catch (ex) {
                    console.log(`Failed to parse expected JSON message: ${data} ${ex}`);
                    return;
                }

                const event = new ContractEvent();
                event.chainID = msg.ChainID;
                const items: string[] = msg.Content;
                for (const item of items) {
                    const parts = item.split(': ');
                    event.contractID = parts[0];
                    event.data = parts[1];
                    for (const callback of self.callbacks) {
                        callback(event);
                    }
                }
            });
        }
        return null;
    }

    public unsubscribeEvents(who: WasmClientContext): void {
        for (let i = 0; i < this.subscribers.length; i++) {
            if (this.subscribers[i] === who) {
                this.subscribers.splice(i, 1);
                this.callbacks.splice(i, 1);
                if (this.subscribers.length == 0) {
                    this.ws.close();
                }
                return;
            }
        }
    }

    public waitUntilRequestProcessed(chainID: wasmlib.ScChainID, reqID: wasmlib.ScRequestID, timeout: u32): isc.Error {
        return this.waspClient.waitUntilRequestProcessed(chainID, reqID, timeout);
    }
}
