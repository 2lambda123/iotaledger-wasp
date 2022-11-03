// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]

use wasmlib::*;
use crate::*;

pub struct NowCall {
	pub func: ScFunc,
}

pub struct GetTimestampCall {
	pub func: ScView,
	pub results: ImmutableGetTimestampResults,
}

pub struct ScFuncs {
}

impl ScFuncs {
    pub fn now(_ctx: &dyn ScFuncCallContext) -> NowCall {
        NowCall {
            func: ScFunc::new(HSC_NAME, HFUNC_NOW),
        }
    }

    pub fn get_timestamp(_ctx: &dyn ScViewCallContext) -> GetTimestampCall {
        let mut f = GetTimestampCall {
            func: ScView::new(HSC_NAME, HVIEW_GET_TIMESTAMP),
            results: ImmutableGetTimestampResults { proxy: Proxy::nil() },
        };
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }
}
