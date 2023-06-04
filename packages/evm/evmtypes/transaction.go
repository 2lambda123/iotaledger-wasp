// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package evmtypes

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/iotaledger/wasp/packages/util"
)

func EncodeTransaction(tx *types.Transaction) []byte {
	return util.BytesFromWriter(tx.EncodeRLP)
}

func DecodeTransaction(b []byte) (*types.Transaction, error) {
	tx := new(types.Transaction)
	err := tx.DecodeRLP(rlp.NewStream(bytes.NewReader(b), 0))
	return tx, err
}
