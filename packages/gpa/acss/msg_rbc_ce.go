// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"go.dedis.ch/kyber/v3/suites"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/util"
)

// This message is used as a payload of the RBC:
//
// > RBC(C||E)
type msgRBCCEPayload struct {
	suite suites.Suite
	data  []byte
}

func (m *msgRBCCEPayload) MarshalBinary() ([]byte, error) {
	mu := new(marshalutil.MarshalUtil)
	util.WriteBytesMu(mu, m.data)
	return mu.Bytes(), nil
}

func (m *msgRBCCEPayload) UnmarshalBinary(data []byte) (err error) {
	mu := marshalutil.New(data)
	m.data, err = util.ReadBytesMu(mu)
	return err
}
