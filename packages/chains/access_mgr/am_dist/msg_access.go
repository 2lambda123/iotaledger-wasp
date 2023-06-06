// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package am_dist

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

// Send by a node which has a chain enabled to a node it considers an access node.
type msgAccess struct {
	gpa.BasicMessage
	senderLClock    int
	receiverLClock  int
	accessForChains []isc.ChainID
	serverForChains []isc.ChainID
}

var _ gpa.Message = &msgAccess{}

func newMsgAccess(
	recipient gpa.NodeID,
	senderLClock, receiverLClock int,
	accessForChains []isc.ChainID,
	serverForChains []isc.ChainID,
) gpa.Message {
	return &msgAccess{
		BasicMessage:    gpa.NewBasicMessage(recipient),
		senderLClock:    senderLClock,
		receiverLClock:  receiverLClock,
		accessForChains: accessForChains,
		serverForChains: serverForChains,
	}
}

func (m *msgAccess) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeAccess)
	ww.WriteUint32(uint32(m.senderLClock))
	ww.WriteUint32(uint32(m.receiverLClock))
	ww.WriteUint32(uint32(len(m.accessForChains)))
	for i := range m.accessForChains {
		ww.WriteBytes(m.accessForChains[i].Bytes())
	}
	ww.WriteUint32(uint32(len(m.serverForChains)))
	for i := range m.serverForChains {
		ww.WriteBytes(m.serverForChains[i].Bytes())
	}
	return w.Bytes(), nil
}

func (m *msgAccess) UnmarshalBinary(data []byte) (err error) {
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	if msgType := rr.ReadByte(); msgType != msgTypeAccess {
		return fmt.Errorf("unexpected message type: %v", msgType)
	}

	// senderLClock
	m.senderLClock = int(rr.ReadUint32())

	// receiverLClock
	m.receiverLClock = int(rr.ReadUint32())

	// accessForChains
	m.accessForChains = make([]isc.ChainID, rr.ReadUint32())
	for i := range m.accessForChains {
		val := rr.ReadBytes()
		chainID, err2 := isc.ChainIDFromBytes(val)
		if err2 != nil {
			return err2
		}
		m.accessForChains[i] = chainID
	}

	// serverForChains
	m.serverForChains = make([]isc.ChainID, rr.ReadUint32())
	for i := range m.serverForChains {
		val := rr.ReadBytes()
		chainID, err2 := isc.ChainIDFromBytes(val)
		if err2 != nil {
			return err2
		}
		m.serverForChains[i] = chainID
	}
	return nil
}
