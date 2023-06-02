// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package evmtypes

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
)

func readBytes(mu *marshalutil.MarshalUtil) (b []byte, err error) {
	var n uint32
	if n, err = mu.ReadUint32(); err != nil {
		return nil, err
	}
	return mu.ReadBytes(int(n))
}

func writeBytes(mu *marshalutil.MarshalUtil, b []byte) {
	mu.WriteUint32(uint32(len(b)))
	mu.WriteBytes(b)
}
