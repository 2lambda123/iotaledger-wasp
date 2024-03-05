// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package apilib

import (
	"context"
	"math"
	"time"

	"github.com/iotaledger/iota.go/v4/hexutil"
	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/packages/cryptolib"
)

// RunDKG runs DKG procedure on specific Wasp hosts: generates new keys and puts corresponding committee records
// into nodes. In case of success, generated address is returned
func RunDKG(client *apiclient.APIClient, peerPubKeys []string, threshold uint16, timeout ...time.Duration) (*cryptolib.PublicKey, error) {
	to := uint32(60 * 1000)
	if len(timeout) > 0 {
		n := timeout[0].Milliseconds()
		if n < int64(math.MaxUint16) {
			to = uint32(n)
		}
	}

	dkShares, _, err := client.NodeApi.GenerateDKS(context.Background()).DKSharesPostRequest(apiclient.DKSharesPostRequest{
		Threshold:      uint32(threshold),
		TimeoutMS:      to,
		PeerIdentities: peerPubKeys,
	}).Execute()
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := hexutil.DecodeHex(dkShares.PublicKey)
	if err != nil {
		return nil, err
	}
	pubKey, err := cryptolib.PublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}
