// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "wasmlib"
import * as sc from "./index";

export class HelloWorldCall {
    func: wasmlib.ScFunc = new wasmlib.ScFunc(sc.HScName, sc.HFuncHelloWorld);
}

export class HelloWorldContext {
    state: sc.MutableHelloWorldState = new sc.MutableHelloWorldState();
}

export class GetHelloWorldCall {
    func: wasmlib.ScView = new wasmlib.ScView(sc.HScName, sc.HViewGetHelloWorld);
    results: sc.ImmutableGetHelloWorldResults = new sc.ImmutableGetHelloWorldResults();
}

export class GetHelloWorldContext {
    results: sc.MutableGetHelloWorldResults = new sc.MutableGetHelloWorldResults();
    state: sc.ImmutableHelloWorldState = new sc.ImmutableHelloWorldState();
}

export class ScFuncs {

    static helloWorld(ctx: wasmlib.ScFuncCallContext): HelloWorldCall {
        let f = new HelloWorldCall();
        return f;
    }

    static getHelloWorld(ctx: wasmlib.ScViewCallContext): GetHelloWorldCall {
        let f = new GetHelloWorldCall();
        f.func.setPtrs(null, f.results);
        return f;
    }
}
