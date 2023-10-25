// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscutils

import (
	"testing"

	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
)

func TestNFTLibrary(t *testing.T) {
	env := evmtest.InitEVM(t)
	ethKey, _ := env.Chain.NewEthereumAccountWithL2Funds()

	nftTest := env.DeployContract(ethKey, NFTTestContractABI, NFTTestContractBytecode)
	_ = nftTest
}
