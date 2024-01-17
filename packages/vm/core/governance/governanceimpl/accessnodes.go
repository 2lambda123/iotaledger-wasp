// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// This file provides implementation for the governance SC, the ChainNode
// management functions.
//
// State of the SC (the ChainNodes part):
//
//	VarAccessNodeCandidates:  map[pubKey] => AccessNodeInfo    // A set of Access Node Info.
//	VarAccessNodes:           map[pubKey] => byte[0]           // A set of nodes.
//	VarValidatorNodes:        pubKey[]                         // An ordered list of nodes.
package governanceimpl

import (
	"encoding/base64"

	"github.com/samber/lo"

	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
)

// Can only be invoked by the access node owner (verified via the Certificate field).
func addCandidateNode(ctx isc.Sandbox, ani *governance.AccessNodeInfo) dict.Dict {
	ani = governance.AccessNodeInfoWithValidatorAddress(ctx, ani)
	pubKeyStr := base64.StdEncoding.EncodeToString(ani.NodePubKey)

	state := ctx.State()
	governance.AccessNodeCandidatesMap(state).SetAt(ani.NodePubKey, ani.Bytes())
	ctx.Log().LogInfof("Governance::AddCandidateNode: accessNodeCandidate added, pubKey=%s", pubKeyStr)

	if ctx.ChainOwnerID().Equals(ctx.Request().SenderAccount()) {
		governance.AccessNodesMap(state).SetAt(ani.NodePubKey, codec.Bool.Encode(true))
		ctx.Log().LogInfof("Governance::AddCandidateNode: accessNode added, pubKey=%s", pubKeyStr)
	}

	return nil
}

// Can only be invoked by the access node owner (verified via the Certificate field).
//
// It is possible that after executing `revokeAccessNode(...)` a node will stay
// in the list of validators, and will be absent in the candidate or an access node set.
// The node is removed from the list of access nodes immediately, but the validator rotation
// must be initiated by the chain owner explicitly.
func revokeAccessNode(ctx isc.Sandbox, ani *governance.AccessNodeInfo) dict.Dict {
	ani = governance.AccessNodeInfoWithValidatorAddress(ctx, ani)
	state := ctx.State()
	governance.AccessNodeCandidatesMap(state).DelAt(ani.NodePubKey)
	governance.AccessNodesMap(state).DelAt(ani.NodePubKey)
	return nil
}

// Can only be invoked by the chain owner.
func changeAccessNodes(ctx isc.Sandbox, req governance.ChangeAccessNodesRequest) dict.Dict {
	ctx.RequireCallerIsChainOwner()
	state := ctx.State()
	accessNodeCandidates := governance.AccessNodeCandidatesMap(state)
	accessNodes := governance.AccessNodesMap(state)
	ctx.Log().LogDebugf("changeAccessNodes: actions len: %d", len(req))
	for pubKey, action := range req {
		switch action {
		case governance.ChangeAccessNodeActionRemove:
			accessNodes.DelAt(pubKey[:])
		case governance.ChangeAccessNodeActionAccept:
			// TODO should the list of candidates be checked? we are just adding any pubkey
			accessNodes.SetAt(pubKey[:], codec.Bool.Encode(true))
			// TODO should the node be removed from the list of candidates? // accessNodeCandidates.DelAt(pubKey)
		case governance.ChangeAccessNodeActionDrop:
			accessNodes.DelAt(pubKey[:])
			accessNodeCandidates.DelAt(pubKey[:])
		}
	}
	return nil
}

func getChainNodes(ctx isc.SandboxView) *governance.GetChainNodesResponse {
	res := &governance.GetChainNodesResponse{
		AccessNodeCandidates: make(map[cryptolib.PublicKeyKey]*governance.AccessNodeInfo),
		AccessNodes:          make(map[cryptolib.PublicKeyKey]struct{}),
	}
	stateR := ctx.StateR()
	governance.AccessNodeCandidatesMapR(stateR).Iterate(func(key, value []byte) bool {
		ani := lo.Must(governance.AccessNodeInfoFromBytes(key, value))
		res.AccessNodeCandidates[lo.Must(cryptolib.PublicKeyFromBytes(key)).AsKey()] = ani
		return true
	})
	governance.AccessNodesMapR(stateR).IterateKeys(func(key []byte) bool {
		res.AccessNodes[lo.Must(cryptolib.PublicKeyFromBytes(key)).AsKey()] = struct{}{}
		return true
	})
	return res
}
