// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package bp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

func TestBatchProposal1Serialization(t *testing.T) {
	var reqRefs []*isc.RequestRef
	for i := uint64(0); i < 5; i++ {
		req := isc.NewOffLedgerRequest(isc.RandomChainID(), isc.NewMessage(3, 14, dict.New()), i, 200).Sign(cryptolib.NewKeyPair())
		reqRefs = append(reqRefs, &isc.RequestRef{
			ID:   req.ID(),
			Hash: hashing.PseudoRandomHash(nil),
		})
	}

	batchProposal1 := NewBatchProposal(
		iotago.LatestAPI(iotago.NewV3SnapshotProtocolParameters()),
		10,
		nil,
		iotago.BlockIDs{},
		isc.RandomChainOutputs(),
		nil,
		util.NewFixedSizeBitVector(11),
		util.NewFixedSizeBitVector(11),
		time.Now(),
		isc.NewRandomAgentID(),
		reqRefs,
	)

	b := rwutil.WriteToBytes(batchProposal1)
	batchProposal2, err := rwutil.ReadFromBytes(b, new(BatchProposal))
	require.NoError(t, err)
	require.Equal(t, batchProposal1.nodeIndex, batchProposal2.nodeIndex)
	require.Equal(t, batchProposal1.baseCO, batchProposal2.baseCO)
	require.Equal(t, batchProposal1.dssTIndexProposal, batchProposal2.dssTIndexProposal)
	require.Equal(t, batchProposal1.dssBIndexProposal, batchProposal2.dssBIndexProposal)
	require.Equal(t, batchProposal1.timeData.UnixNano(), batchProposal2.timeData.UnixNano())
	require.Equal(t, batchProposal1.validatorFeeDestination, batchProposal2.validatorFeeDestination)
	require.Equal(t, batchProposal1.requestRefs, batchProposal2.requestRefs)
}
