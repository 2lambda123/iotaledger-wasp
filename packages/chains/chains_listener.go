// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chains

import (
	"github.com/iotaledger/wasp/packages/chain/chaintypes"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/state"
)

type chainsListener struct {
	parent        chaintypes.ChainListener
	accessNodesCB func(chainID isc.ChainID, accessNodes []*cryptolib.PublicKey)
}

func NewChainsListener(parent chaintypes.ChainListener, accessNodesCB func(chainID isc.ChainID, accessNodes []*cryptolib.PublicKey)) chaintypes.ChainListener {
	return &chainsListener{parent: parent, accessNodesCB: accessNodesCB}
}

func (cl *chainsListener) BlockApplied(chainID isc.ChainID, block state.Block, latestState kv.KVStoreReader) {
	cl.parent.BlockApplied(chainID, block, latestState)
}

func (cl *chainsListener) AccessNodesUpdated(chainID isc.ChainID, accessNodes []*cryptolib.PublicKey) {
	cl.accessNodesCB(chainID, accessNodes)
	cl.parent.AccessNodesUpdated(chainID, accessNodes)
}

func (cl *chainsListener) ServerNodesUpdated(chainID isc.ChainID, serverNodes []*cryptolib.PublicKey) {
	cl.parent.ServerNodesUpdated(chainID, serverNodes)
}
