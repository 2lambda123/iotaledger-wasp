// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as isc from './isc';
import * as wasmlib from 'wasmlib';
import {WebSocket} from 'ws';

type ClientCallBack = (msg: string) => void;

export class WasmClientService {
    private callbacks: ClientCallBack[] = [];
    private eventPort: string;
    private ws: WebSocket;
    private subscribers: any[] = [];
    private waspClient: isc.WaspClient;

    public constructor(waspAPI: string, eventPort: string) {
        this.waspClient = new isc.WaspClient(waspAPI);
        this.eventPort = eventPort;
        this.ws = new WebSocket(eventPort, {
            perMessageDeflate: false
        });
    }

    public static DefaultWasmClientService(): WasmClientService {
        return new WasmClientService('http://localhost:19090', 'ws://localhost:19090/ws');
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

    public subscribeEvents(who: any, callback: (msg: string) => void): isc.Error {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const self = this;
        this.callbacks.push(callback);
        this.subscribers.push(who);
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
                    console.log(msg);
                } catch (ex) {
                    console.log(`Failed to parse expected JSON message: ${data} ${ex}`);
                    return;
                }

                if (!msg.Kind) {
                    return;
                }

                const items: string[] = msg.Content;
                for (const item of items) {
                    const parts = item.split(': ');
                    for (let i = 0; i < self.callbacks.length; i++) {
                        self.callbacks[i](parts[1]);
                    }
                }
            });
        }
        return null;
    }

    public unsubscribeEvents(who: any): void {
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
