// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscutils

import (
	_ "embed"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// If you change any of the .sol files, you must recompile them.  You will need
// the `solc` binary installed in your system. Then, simply run `go generate`
// in this directory.

//go:generate solc --abi --bin --bin-runtime --overwrite prng_test.sol -o build/.
var (
	//go:embed build/PRNGTest.abi
	PRNGTestContractABI string
	//go:embed build/PRNGTest.bin
	PRNGTestContractBytecodeHex string
	PRNGTestContractBytecode    = common.FromHex(strings.TrimSpace(PRNGTestContractBytecodeHex))
	//deployed bytecode and runtime bytecode are different, see: https://ethereum.stackexchange.com/questions/13086/whats-the-difference-between-solcs-bin-bytecode-versus-bin-runtime
	//go:embed build/PRNGTest.bin-runtime
	PRNGTestContractRuntimeBytecodeHex string
	PRNGTestContractRuntimeBytecode    = common.FromHex(strings.TrimSpace(PRNGTestContractRuntimeBytecodeHex))
)

//go:generate solc --abi --bin --bin-runtime --overwrite nft_test.sol @iota=../../../packages/vm/core/evm -o build/.
var (
	//go:embed build/NFTTest.abi
	NFTTestContractABI string
	//go:embed build/NFTTest.bin
	NFTTestContractBytecodeHex string
	NFTTestContractBytecode    = common.FromHex(strings.TrimSpace(NFTTestContractBytecodeHex))
	//deployed bytecode and runtime bytecode are different, see: https://ethereum.stackexchange.com/questions/13086/whats-the-difference-between-solcs-bin-bytecode-versus-bin-runtime
	//go:embed build/NFTTest.bin-runtime
	NFTTestContractRuntimeBytecodeHex string
	NFTTestContractRuntimeBytecode    = common.FromHex(strings.TrimSpace(NFTTestContractRuntimeBytecodeHex))
)
