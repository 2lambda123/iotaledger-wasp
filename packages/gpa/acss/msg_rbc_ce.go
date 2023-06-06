// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"bytes"

	"go.dedis.ch/kyber/v3/suites"

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
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteBytes(m.data)
	return w.Bytes(), nil
}

func (m *msgRBCCEPayload) UnmarshalBinary(data []byte) (err error) {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)
	m.data = rr.ReadBytes()
	return nil
}
