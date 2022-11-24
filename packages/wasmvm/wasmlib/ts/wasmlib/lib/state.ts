// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import {stateDelete, stateExists, stateGet, stateSet} from "./host";
import {IKvStore, Proxy} from "./wasmtypes/proxy";

export class ScImmutableState {
    exists(key: u8[]): bool {
        return stateExists(key);
    }

    get(key: u8[]): u8[] {
        const val = stateGet(key);
        return val === null ? [] : val;
    }
}

export class ScState implements IKvStore {
    public static proxy(): Proxy {
        return new Proxy(new ScState());
    }

    delete(key: u8[]): void {
        stateDelete(key);
    }

    exists(key: u8[]): bool {
        return stateExists(key);
    }

    get(key: u8[]): u8[] {
        const val = stateGet(key);
        return val === null ? [] : val;
    }

    public immutable(): ScImmutableState {
        return this;
    }

    set(key: u8[], value: u8[]): void {
        stateSet(key, value);
    }
}
