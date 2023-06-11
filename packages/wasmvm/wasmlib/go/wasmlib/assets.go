// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmlib

import (
	"sort"

	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type TokenAmounts map[wasmtypes.ScTokenID]wasmtypes.ScBigInt

type ScAssets struct {
	BaseTokens   uint64
	NativeTokens TokenAmounts
	NftIDs       map[wasmtypes.ScNftID]bool
}

func NewScAssets(buf []byte) *ScAssets {
	assets := &ScAssets{}
	if len(buf) == 0 {
		return assets
	}

	dec := wasmtypes.NewWasmDecoder(buf)
	empty := wasmtypes.BoolDecode(dec)
	if empty {
		return assets
	}

	assets.BaseTokens = wasmtypes.Uint64Decode(dec)

	size := wasmtypes.Uint16Decode(dec)
	if size > 0 {
		assets.NativeTokens = make(TokenAmounts, size)
		for ; size > 0; size-- {
			tokenID := wasmtypes.TokenIDDecode(dec)
			assets.NativeTokens[tokenID] = wasmtypes.BigIntDecode(dec)
		}
	}

	size = wasmtypes.Uint16Decode(dec)
	if size > 0 {
		assets.NftIDs = make(map[wasmtypes.ScNftID]bool)
	}
	for ; size > 0; size-- {
		nftID := wasmtypes.NftIDDecode(dec)
		assets.NftIDs[nftID] = true
	}
	return assets
}

func (a *ScAssets) Balances() ScBalances {
	return ScBalances{assets: a}
}

func (a *ScAssets) Bytes() []byte {
	if a == nil {
		return []byte{}
	}

	enc := wasmtypes.NewWasmEncoder()
	empty := a.IsEmpty()
	wasmtypes.BoolEncode(enc, empty)
	if empty {
		return enc.Buf()
	}

	wasmtypes.Uint64Encode(enc, a.BaseTokens)

	wasmtypes.Uint16Encode(enc, uint16(len(a.NativeTokens)))
	for _, tokenID := range a.TokenIDs() {
		wasmtypes.TokenIDEncode(enc, *tokenID)
		wasmtypes.BigIntEncode(enc, a.NativeTokens[*tokenID])
	}

	wasmtypes.Uint16Encode(enc, uint16(len(a.NftIDs)))
	for nftID := range a.NftIDs {
		wasmtypes.NftIDEncode(enc, nftID)
	}
	return enc.Buf()
}

func (a *ScAssets) IsEmpty() bool {
	if a.BaseTokens != 0 {
		return false
	}
	for _, val := range a.NativeTokens {
		if !val.IsZero() {
			return false
		}
	}
	return len(a.NftIDs) == 0
}

func (a *ScAssets) TokenIDs() []*wasmtypes.ScTokenID {
	tokenIDs := make([]*wasmtypes.ScTokenID, 0, len(a.NativeTokens))
	for key := range a.NativeTokens {
		// need a local copy to avoid referencing the single key var multiple times
		tokenID := key
		tokenIDs = append(tokenIDs, &tokenID)
	}
	sort.Slice(tokenIDs, func(i, j int) bool {
		return string(tokenIDs[i].Bytes()) < string(tokenIDs[j].Bytes())
	})
	return tokenIDs
}

type ScBalances struct {
	assets *ScAssets
}

func (b *ScBalances) Balance(tokenID *wasmtypes.ScTokenID) wasmtypes.ScBigInt {
	if len(b.assets.NativeTokens) == 0 {
		return wasmtypes.NewScBigInt()
	}
	return b.assets.NativeTokens[*tokenID]
}

func (b *ScBalances) Bytes() []byte {
	if b == nil {
		return []byte{}
	}
	return b.assets.Bytes()
}

func (b *ScBalances) BaseTokens() uint64 {
	return b.assets.BaseTokens
}

func (b *ScBalances) IsEmpty() bool {
	return b.assets.IsEmpty()
}

func (b *ScBalances) NftIDs() map[wasmtypes.ScNftID]bool {
	return b.assets.NftIDs
}

func (b *ScBalances) TokenIDs() []*wasmtypes.ScTokenID {
	return b.assets.TokenIDs()
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScTransfer struct {
	ScBalances
}

// create a new transfer object ready to add token transfers
func NewScTransfer() *ScTransfer {
	return &ScTransfer{ScBalances{assets: NewScAssets(nil)}}
}

// create a new transfer object from a balances object
func ScTransferFromBalances(balances *ScBalances) *ScTransfer {
	transfer := ScTransferFromBaseTokens(balances.BaseTokens())
	for _, tokenID := range balances.TokenIDs() {
		transfer.Set(tokenID, balances.Balance(tokenID))
	}
	for nftID := range balances.NftIDs() {
		transfer.AddNFT(&nftID)
	}
	return transfer
}

// create a new transfer object and initialize it with the specified amount of base tokens
func ScTransferFromBaseTokens(amount uint64) *ScTransfer {
	transfer := NewScTransfer()
	transfer.assets.BaseTokens = amount
	return transfer
}

// create a new transfer object and initialize it with the specified NFT
func ScTransferFromNFT(nftID *wasmtypes.ScNftID) *ScTransfer {
	transfer := NewScTransfer()
	transfer.AddNFT(nftID)
	return transfer
}

// create a new transfer object and initialize it with the specified token transfer
func ScTransferFromTokens(tokenID *wasmtypes.ScTokenID, amount wasmtypes.ScBigInt) *ScTransfer {
	transfer := NewScTransfer()
	transfer.Set(tokenID, amount)
	return transfer
}

func (t *ScTransfer) AddNFT(nftID *wasmtypes.ScNftID) {
	if t.assets.NftIDs == nil {
		t.assets.NftIDs = make(map[wasmtypes.ScNftID]bool)
	}
	t.assets.NftIDs[*nftID] = true
}

func (t *ScTransfer) Bytes() []byte {
	if t == nil {
		return []byte{}
	}
	return t.assets.Bytes()
}

// set the specified tokenID amount in the transfers object
// note that this will overwrite any previous amount for the specified tokenID
func (t *ScTransfer) Set(tokenID *wasmtypes.ScTokenID, amount wasmtypes.ScBigInt) {
	if t.assets.NativeTokens == nil {
		t.assets.NativeTokens = make(TokenAmounts)
	}
	t.assets.NativeTokens[*tokenID] = amount
}
