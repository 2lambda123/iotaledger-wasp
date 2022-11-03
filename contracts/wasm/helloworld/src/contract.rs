// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]

use wasmlib::*;
use crate::*;

pub struct HelloWorldCall {
	pub func: ScFunc,
}

pub struct GetHelloWorldCall {
	pub func: ScView,
	pub results: ImmutableGetHelloWorldResults,
}

pub struct ScFuncs {
}

impl ScFuncs {
    pub fn hello_world(_ctx: &dyn ScFuncCallContext) -> HelloWorldCall {
        HelloWorldCall {
            func: ScFunc::new(HSC_NAME, HFUNC_HELLO_WORLD),
        }
    }

    pub fn get_hello_world(_ctx: &dyn ScViewCallContext) -> GetHelloWorldCall {
        let mut f = GetHelloWorldCall {
            func: ScView::new(HSC_NAME, HVIEW_GET_HELLO_WORLD),
            results: ImmutableGetHelloWorldResults { proxy: Proxy::nil() },
        };
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }
}
