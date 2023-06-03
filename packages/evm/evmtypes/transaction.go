// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package evmtypes

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/util"
)

func EncodeTransaction(tx *types.Transaction) []byte {
	mu := new(marshalutil.MarshalUtil)
	util.MarshallWriter(mu, tx.EncodeRLP)
	return mu.Bytes()
}

func DecodeTransaction(b []byte) (*types.Transaction, error) {
	tx := new(types.Transaction)
	err := tx.DecodeRLP(rlp.NewStream(bytes.NewReader(b), 0))
	return tx, err
}
