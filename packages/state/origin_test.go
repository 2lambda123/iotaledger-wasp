package state

import (
	"testing"

	"github.com/iotaledger/hive.go/core/kvstore/mapdb"
	"github.com/stretchr/testify/require"
)

func TestOrigin(t *testing.T) {
	store := InitChainStore(mapdb.NewMapDB())
	l1commitment := OriginL1Commitment()
	require.True(t, l1commitment.Equals(store.LatestBlock().L1Commitment()))
}
