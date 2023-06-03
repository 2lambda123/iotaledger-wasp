// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package evmtypes

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/util"
)

// EncodeReceipt serializes the receipt in RLP format
func EncodeReceipt(receipt *types.Receipt) []byte {
	mu := new(marshalutil.MarshalUtil)
	util.MarshalWriter(mu, receipt.EncodeRLP)
	return mu.Bytes()
}

func DecodeReceipt(data []byte) (*types.Receipt, error) {
	receipt := new(types.Receipt)
	err := receipt.DecodeRLP(rlp.NewStream(bytes.NewReader(data), 0))
	return receipt, err
}

func BloomFilter(bloom types.Bloom, addresses []common.Address, topics [][]common.Hash) bool {
	return bloomMatchesAddresses(bloom, addresses) && bloomMatchesAllEvents(bloom, topics)
}

func bloomMatchesAddresses(bloom types.Bloom, addresses []common.Address) bool {
	if len(addresses) == 0 {
		return true
	}
	for _, addr := range addresses {
		if types.BloomLookup(bloom, addr) {
			return true
		}
	}
	return false
}

func bloomMatchesAllEvents(bloom types.Bloom, events [][]common.Hash) bool {
	for _, topics := range events {
		if !bloomMatchesAnyTopic(bloom, topics) {
			return false
		}
	}
	return true
}

func bloomMatchesAnyTopic(bloom types.Bloom, topics []common.Hash) bool {
	if len(topics) == 0 {
		return true
	}
	for _, topic := range topics {
		if types.BloomLookup(bloom, topic) {
			return true
		}
	}
	return false
}
