// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]

use wasmlib::*;

use crate::*;

pub struct FinalizeAuctionCall {
    pub func: ScFunc,
    pub params: MutableFinalizeAuctionParams,
}

pub struct InitCall {
    pub func: ScInitFunc,
    pub params: MutableInitParams,
}

pub struct PlaceBidCall {
    pub func: ScFunc,
    pub params: MutablePlaceBidParams,
}

pub struct SetOwnerMarginCall {
    pub func: ScFunc,
    pub params: MutableSetOwnerMarginParams,
}

pub struct StartAuctionCall {
    pub func: ScFunc,
    pub params: MutableStartAuctionParams,
}

pub struct GetAuctionInfoCall {
    pub func: ScView,
    pub params: MutableGetAuctionInfoParams,
    pub results: ImmutableGetAuctionInfoResults,
}

pub struct ScFuncs {
}

impl ScFuncs {
    pub fn finalize_auction(_ctx: &dyn ScFuncCallContext) -> FinalizeAuctionCall {
        let mut f = FinalizeAuctionCall {
            func: ScFunc::new(HSC_NAME, HFUNC_FINALIZE_AUCTION),
            params: MutableFinalizeAuctionParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    pub fn init(_ctx: &dyn ScFuncCallContext) -> InitCall {
        let mut f = InitCall {
            func: ScInitFunc::new(HSC_NAME, HFUNC_INIT),
            params: MutableInitParams { proxy: Proxy::nil() },
        };
        ScInitFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    pub fn place_bid(_ctx: &dyn ScFuncCallContext) -> PlaceBidCall {
        let mut f = PlaceBidCall {
            func: ScFunc::new(HSC_NAME, HFUNC_PLACE_BID),
            params: MutablePlaceBidParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    pub fn set_owner_margin(_ctx: &dyn ScFuncCallContext) -> SetOwnerMarginCall {
        let mut f = SetOwnerMarginCall {
            func: ScFunc::new(HSC_NAME, HFUNC_SET_OWNER_MARGIN),
            params: MutableSetOwnerMarginParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    pub fn start_auction(_ctx: &dyn ScFuncCallContext) -> StartAuctionCall {
        let mut f = StartAuctionCall {
            func: ScFunc::new(HSC_NAME, HFUNC_START_AUCTION),
            params: MutableStartAuctionParams { proxy: Proxy::nil() },
        };
        ScFunc::link_params(&mut f.params.proxy, &f.func);
        f
    }

    pub fn get_auction_info(_ctx: &dyn ScViewCallContext) -> GetAuctionInfoCall {
        let mut f = GetAuctionInfoCall {
            func: ScView::new(HSC_NAME, HVIEW_GET_AUCTION_INFO),
            params: MutableGetAuctionInfoParams { proxy: Proxy::nil() },
            results: ImmutableGetAuctionInfoResults { proxy: Proxy::nil() },
        };
        ScView::link_params(&mut f.params.proxy, &f.func);
        ScView::link_results(&mut f.results.proxy, &f.func);
        f
    }
}
