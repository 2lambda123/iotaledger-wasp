// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

import * as wasmlib from "wasmlib";
import * as sc from "./index";

export class ImmutableAllowanceResults extends wasmlib.ScMapID {

    amount(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, sc.idxMap[sc.IdxResultAmount]);
	}
}

export class MutableAllowanceResults extends wasmlib.ScMapID {

    amount(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, sc.idxMap[sc.IdxResultAmount]);
	}
}

export class ImmutableBalanceOfResults extends wasmlib.ScMapID {

    amount(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, sc.idxMap[sc.IdxResultAmount]);
	}
}

export class MutableBalanceOfResults extends wasmlib.ScMapID {

    amount(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, sc.idxMap[sc.IdxResultAmount]);
	}
}

export class ImmutableTotalSupplyResults extends wasmlib.ScMapID {

    supply(): wasmlib.ScImmutableInt64 {
		return new wasmlib.ScImmutableInt64(this.mapID, sc.idxMap[sc.IdxResultSupply]);
	}
}

export class MutableTotalSupplyResults extends wasmlib.ScMapID {

    supply(): wasmlib.ScMutableInt64 {
		return new wasmlib.ScMutableInt64(this.mapID, sc.idxMap[sc.IdxResultSupply]);
	}
}
