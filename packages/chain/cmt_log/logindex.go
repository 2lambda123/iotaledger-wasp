// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cmt_log

import (
	"encoding/binary"
)

// LogIndex starts from 1. 0 is used as a nil value.
type LogIndex uint32

func (li LogIndex) AsUint32() uint32 {
	return uint32(li)
}

func (li LogIndex) Bytes() []byte {
	liBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(liBytes, li.AsUint32())
	return liBytes
}

func (li LogIndex) IsNil() bool {
	return li == 0
}

func (li LogIndex) Next() LogIndex {
	return LogIndex(li.AsUint32() + 1)
}

func NilLogIndex() LogIndex {
	return LogIndex(0)
}

func MaxLogIndex(lis ...LogIndex) LogIndex {
	max := NilLogIndex()
	for _, li := range lis {
		if li > max {
			max = li
		}
	}
	return max
}
