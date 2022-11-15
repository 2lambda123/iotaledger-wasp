// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package evmimpl

import (
	"math/big"

	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/evm/iscmagic"
)

// handler for ISCAccounts::getL2BalanceBaseTokens
func (h *magicContractViewHandler) GetL2BalanceBaseTokens(agentID iscmagic.ISCAgentID) uint64 {
	r := h.ctx.CallView(accounts.Contract.Hname(), accounts.ViewBalanceBaseToken.Hname(), dict.Dict{
		accounts.ParamAgentID: codec.EncodeAgentID(agentID.MustUnwrap()),
	})
	return codec.MustDecodeUint64(r.MustGet(accounts.ParamBalance))
}

// handler for ISCAccounts::getL2BalanceNativeTokens
func (h *magicContractViewHandler) GetL2BalanceNativeTokens(nativeTokenID iscmagic.NativeTokenID, agentID iscmagic.ISCAgentID) *big.Int {
	r := h.ctx.CallView(accounts.Contract.Hname(), accounts.ViewBalanceNativeToken.Hname(), dict.Dict{
		accounts.ParamNativeTokenID: codec.EncodeNativeTokenID(nativeTokenID.Unwrap()),
		accounts.ParamAgentID:       codec.EncodeAgentID(agentID.MustUnwrap()),
	})
	return codec.MustDecodeBigIntAbs(r.MustGet(accounts.ParamBalance))
}

// handler for ISCAccounts::getL2NFTs
func (h *magicContractViewHandler) GetL2NFTs(agentID iscmagic.ISCAgentID) []iscmagic.NFTID {
	r := h.ctx.CallView(accounts.Contract.Hname(), accounts.ViewAccountNFTs.Hname(), dict.Dict{
		accounts.ParamAgentID: codec.EncodeAgentID(agentID.MustUnwrap()),
	})
	arr := collections.NewArray16(r, accounts.ParamNFTIDs)
	ret := make([]iscmagic.NFTID, arr.MustLen())
	for i := uint16(0); i < arr.MustLen(); i++ {
		copy(ret[i][:], arr.MustGetAt(i))
	}
	return ret
}
