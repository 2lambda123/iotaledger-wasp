// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscutils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/iotaledger/wasp/packages/vm/core/evm/evmtest"
)

func TestNFTLibrary(t *testing.T) {
	env := evmtest.InitEVM(t)
	ethKey, _ := env.Chain.NewEthereumAccountWithL2Funds()

	nftTest := env.DeployContract(ethKey, NFTTestContractABI, NFTTestContractBytecode)

	var value []byte
	//nftTest.CallFnExpectEvent(nil, "Minted", &value, "MintNFT", []byte{},
	//	iscmagic.ISCAgentID{Data: addr.Bytes()}, iscmagic.NFTID{}, false)

	res := nftTest.CallFnExpectEvent(nil, "Minted", &value, "MintTestNFT")
	//res, err := nftTest.CallFn(nil, "MintNFTWithReturnedDict")
	//fmt.Println("err:", err)
	data, _ := res.EVMReceipt.MarshalJSON()
	fmt.Println("DATA:", string(data))
	return
	require.NotEmpty(t, fmt.Sprintf("%x", value))
}
