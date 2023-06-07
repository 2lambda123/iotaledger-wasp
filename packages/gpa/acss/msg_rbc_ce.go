// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package acss

import (
	"io"

	"github.com/iotaledger/wasp/packages/util/rwutil"
	"go.dedis.ch/kyber/v3/suites"
)

// This message is used as a payload of the RBC:
//
// > RBC(C||E)
type msgRBCCEPayload struct {
	suite suites.Suite
	data  []byte
}

func (msg *msgRBCCEPayload) MarshalBinary() ([]byte, error) {
	return rwutil.WriterToBytes(msg), nil
}

func (msg *msgRBCCEPayload) UnmarshalBinary(data []byte) error {
	_, err := rwutil.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgRBCCEPayload) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	msg.data = rr.ReadBytes()
	return rr.Err
}

func (msg *msgRBCCEPayload) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteBytes(msg.data)
	return ww.Err
}
