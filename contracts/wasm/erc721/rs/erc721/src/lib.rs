// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

#![allow(dead_code)]
#![allow(unused_imports)]

use erc721::*;
use wasmlib::*;

use crate::consts::*;
use crate::events::*;
use crate::params::*;
use crate::results::*;
use crate::state::*;
use crate::typedefs::*;

mod consts;
mod contract;
mod events;
mod params;
mod results;
mod state;
mod typedefs;

mod erc721;

const EXPORT_MAP: ScExportMap = ScExportMap {
    names: &[
        FUNC_APPROVE,
        FUNC_BURN,
        FUNC_INIT,
        FUNC_MINT,
        FUNC_SAFE_TRANSFER_FROM,
        FUNC_SET_APPROVAL_FOR_ALL,
        FUNC_TRANSFER_FROM,
        VIEW_BALANCE_OF,
        VIEW_GET_APPROVED,
        VIEW_IS_APPROVED_FOR_ALL,
        VIEW_NAME,
        VIEW_OWNER_OF,
        VIEW_SYMBOL,
        VIEW_TOKEN_URI,
    ],
    funcs: &[
        func_approve_thunk,
        func_burn_thunk,
        func_init_thunk,
        func_mint_thunk,
        func_safe_transfer_from_thunk,
        func_set_approval_for_all_thunk,
        func_transfer_from_thunk,
    ],
    views: &[
        view_balance_of_thunk,
        view_get_approved_thunk,
        view_is_approved_for_all_thunk,
        view_name_thunk,
        view_owner_of_thunk,
        view_symbol_thunk,
        view_token_uri_thunk,
    ],
};

pub fn on_dispatch(index: i32) {
    EXPORT_MAP.dispatch(index);
}

pub struct ApproveContext {
    events:  Erc721Events,
    params: ImmutableApproveParams,
    state: MutableErc721State,
}

fn func_approve_thunk(ctx: &ScFuncContext) {
    ctx.log("erc721.funcApprove");
    let f = ApproveContext {
        events:  Erc721Events {},
        params: ImmutableApproveParams { proxy: params_proxy() },
        state: MutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    func_approve(ctx, &f);
    ctx.log("erc721.funcApprove ok");
}

pub struct BurnContext {
    events:  Erc721Events,
    params: ImmutableBurnParams,
    state: MutableErc721State,
}

fn func_burn_thunk(ctx: &ScFuncContext) {
    ctx.log("erc721.funcBurn");
    let f = BurnContext {
        events:  Erc721Events {},
        params: ImmutableBurnParams { proxy: params_proxy() },
        state: MutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    func_burn(ctx, &f);
    ctx.log("erc721.funcBurn ok");
}

pub struct InitContext {
    events:  Erc721Events,
    params: ImmutableInitParams,
    state: MutableErc721State,
}

fn func_init_thunk(ctx: &ScFuncContext) {
    ctx.log("erc721.funcInit");
    let f = InitContext {
        events:  Erc721Events {},
        params: ImmutableInitParams { proxy: params_proxy() },
        state: MutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.name().exists(), "missing mandatory name");
    ctx.require(f.params.symbol().exists(), "missing mandatory symbol");
    func_init(ctx, &f);
    ctx.log("erc721.funcInit ok");
}

pub struct MintContext {
    events:  Erc721Events,
    params: ImmutableMintParams,
    state: MutableErc721State,
}

fn func_mint_thunk(ctx: &ScFuncContext) {
    ctx.log("erc721.funcMint");
    let f = MintContext {
        events:  Erc721Events {},
        params: ImmutableMintParams { proxy: params_proxy() },
        state: MutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    func_mint(ctx, &f);
    ctx.log("erc721.funcMint ok");
}

pub struct SafeTransferFromContext {
    events:  Erc721Events,
    params: ImmutableSafeTransferFromParams,
    state: MutableErc721State,
}

fn func_safe_transfer_from_thunk(ctx: &ScFuncContext) {
    ctx.log("erc721.funcSafeTransferFrom");
    let f = SafeTransferFromContext {
        events:  Erc721Events {},
        params: ImmutableSafeTransferFromParams { proxy: params_proxy() },
        state: MutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.from().exists(), "missing mandatory from");
    ctx.require(f.params.to().exists(), "missing mandatory to");
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    func_safe_transfer_from(ctx, &f);
    ctx.log("erc721.funcSafeTransferFrom ok");
}

pub struct SetApprovalForAllContext {
    events:  Erc721Events,
    params: ImmutableSetApprovalForAllParams,
    state: MutableErc721State,
}

fn func_set_approval_for_all_thunk(ctx: &ScFuncContext) {
    ctx.log("erc721.funcSetApprovalForAll");
    let f = SetApprovalForAllContext {
        events:  Erc721Events {},
        params: ImmutableSetApprovalForAllParams { proxy: params_proxy() },
        state: MutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.approval().exists(), "missing mandatory approval");
    ctx.require(f.params.operator().exists(), "missing mandatory operator");
    func_set_approval_for_all(ctx, &f);
    ctx.log("erc721.funcSetApprovalForAll ok");
}

pub struct TransferFromContext {
    events:  Erc721Events,
    params: ImmutableTransferFromParams,
    state: MutableErc721State,
}

fn func_transfer_from_thunk(ctx: &ScFuncContext) {
    ctx.log("erc721.funcTransferFrom");
    let f = TransferFromContext {
        events:  Erc721Events {},
        params: ImmutableTransferFromParams { proxy: params_proxy() },
        state: MutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.from().exists(), "missing mandatory from");
    ctx.require(f.params.to().exists(), "missing mandatory to");
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    func_transfer_from(ctx, &f);
    ctx.log("erc721.funcTransferFrom ok");
}

pub struct BalanceOfContext {
    params: ImmutableBalanceOfParams,
    results: MutableBalanceOfResults,
    state: ImmutableErc721State,
}

fn view_balance_of_thunk(ctx: &ScViewContext) {
    ctx.log("erc721.viewBalanceOf");
    let f = BalanceOfContext {
        params: ImmutableBalanceOfParams { proxy: params_proxy() },
        results: MutableBalanceOfResults { proxy: results_proxy() },
        state: ImmutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.owner().exists(), "missing mandatory owner");
    view_balance_of(ctx, &f);
    ctx.results(&f.results.proxy.kv_store);
    ctx.log("erc721.viewBalanceOf ok");
}

pub struct GetApprovedContext {
    params: ImmutableGetApprovedParams,
    results: MutableGetApprovedResults,
    state: ImmutableErc721State,
}

fn view_get_approved_thunk(ctx: &ScViewContext) {
    ctx.log("erc721.viewGetApproved");
    let f = GetApprovedContext {
        params: ImmutableGetApprovedParams { proxy: params_proxy() },
        results: MutableGetApprovedResults { proxy: results_proxy() },
        state: ImmutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    view_get_approved(ctx, &f);
    ctx.results(&f.results.proxy.kv_store);
    ctx.log("erc721.viewGetApproved ok");
}

pub struct IsApprovedForAllContext {
    params: ImmutableIsApprovedForAllParams,
    results: MutableIsApprovedForAllResults,
    state: ImmutableErc721State,
}

fn view_is_approved_for_all_thunk(ctx: &ScViewContext) {
    ctx.log("erc721.viewIsApprovedForAll");
    let f = IsApprovedForAllContext {
        params: ImmutableIsApprovedForAllParams { proxy: params_proxy() },
        results: MutableIsApprovedForAllResults { proxy: results_proxy() },
        state: ImmutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.operator().exists(), "missing mandatory operator");
    ctx.require(f.params.owner().exists(), "missing mandatory owner");
    view_is_approved_for_all(ctx, &f);
    ctx.results(&f.results.proxy.kv_store);
    ctx.log("erc721.viewIsApprovedForAll ok");
}

pub struct NameContext {
    results: MutableNameResults,
    state: ImmutableErc721State,
}

fn view_name_thunk(ctx: &ScViewContext) {
    ctx.log("erc721.viewName");
    let f = NameContext {
        results: MutableNameResults { proxy: results_proxy() },
        state: ImmutableErc721State { proxy: state_proxy() },
    };
    view_name(ctx, &f);
    ctx.results(&f.results.proxy.kv_store);
    ctx.log("erc721.viewName ok");
}

pub struct OwnerOfContext {
    params: ImmutableOwnerOfParams,
    results: MutableOwnerOfResults,
    state: ImmutableErc721State,
}

fn view_owner_of_thunk(ctx: &ScViewContext) {
    ctx.log("erc721.viewOwnerOf");
    let f = OwnerOfContext {
        params: ImmutableOwnerOfParams { proxy: params_proxy() },
        results: MutableOwnerOfResults { proxy: results_proxy() },
        state: ImmutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    view_owner_of(ctx, &f);
    ctx.results(&f.results.proxy.kv_store);
    ctx.log("erc721.viewOwnerOf ok");
}

pub struct SymbolContext {
    results: MutableSymbolResults,
    state: ImmutableErc721State,
}

fn view_symbol_thunk(ctx: &ScViewContext) {
    ctx.log("erc721.viewSymbol");
    let f = SymbolContext {
        results: MutableSymbolResults { proxy: results_proxy() },
        state: ImmutableErc721State { proxy: state_proxy() },
    };
    view_symbol(ctx, &f);
    ctx.results(&f.results.proxy.kv_store);
    ctx.log("erc721.viewSymbol ok");
}

pub struct TokenURIContext {
    params: ImmutableTokenURIParams,
    results: MutableTokenURIResults,
    state: ImmutableErc721State,
}

fn view_token_uri_thunk(ctx: &ScViewContext) {
    ctx.log("erc721.viewTokenURI");
    let f = TokenURIContext {
        params: ImmutableTokenURIParams { proxy: params_proxy() },
        results: MutableTokenURIResults { proxy: results_proxy() },
        state: ImmutableErc721State { proxy: state_proxy() },
    };
    ctx.require(f.params.token_id().exists(), "missing mandatory tokenID");
    view_token_uri(ctx, &f);
    ctx.results(&f.results.proxy.kv_store);
    ctx.log("erc721.viewTokenURI ok");
}
