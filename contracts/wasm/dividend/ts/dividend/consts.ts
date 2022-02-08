// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmtypes from "wasmlib/wasmtypes";

export const ScName        = "dividend";
export const ScDescription = "Simple dividend smart contract";
export const HScName       = new wasmtypes.ScHname(0xcce2e239);

export const ParamAddress = "address";
export const ParamFactor  = "factor";
export const ParamOwner   = "owner";

export const ResultFactor = "factor";
export const ResultOwner  = "owner";

export const StateMemberList  = "memberList";
export const StateMembers     = "members";
export const StateOwner       = "owner";
export const StateTotalFactor = "totalFactor";

export const FuncDivide    = "divide";
export const FuncInit      = "init";
export const FuncMember    = "member";
export const FuncSetOwner  = "setOwner";
export const ViewGetFactor = "getFactor";
export const ViewGetOwner  = "getOwner";

export const HFuncDivide    = new wasmtypes.ScHname(0xc7878107);
export const HFuncInit      = new wasmtypes.ScHname(0x1f44d644);
export const HFuncMember    = new wasmtypes.ScHname(0xc07da2cb);
export const HFuncSetOwner  = new wasmtypes.ScHname(0x2a15fe7b);
export const HViewGetFactor = new wasmtypes.ScHname(0x0ee668fe);
export const HViewGetOwner  = new wasmtypes.ScHname(0x137107a6);
