// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package isc_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/stretchr/testify/require"
)

func TestHnameSerialize(t *testing.T) {
	data1 := make([]byte, isc.HnameLength)
	rand.Read(data1)
	hname, err := isc.HnameFromBytes(data1)
	require.NoError(t, err)
	data2 := hname.Bytes()
	require.True(t, bytes.Equal(data1, data2))
}
