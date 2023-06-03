// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package am_dist

import (
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
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
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(msgTypeAccess)
	mu.WriteUint32(uint32(m.senderLClock))
	mu.WriteUint32(uint32(m.receiverLClock))
	mu.WriteUint32(uint32(len(m.accessForChains)))
	for i := range m.accessForChains {
		util.MarshalBytes(mu, m.accessForChains[i].Bytes())
	}
	mu.WriteUint32(uint32(len(m.serverForChains)))
	for i := range m.serverForChains {
		util.MarshalBytes(mu, m.serverForChains[i].Bytes())
	}
	return mu.Bytes(), nil
}

func (m *msgAccess) UnmarshalBinary(data []byte) error {
	mu := marshalutil.New(data)
	if msgType, err := mu.ReadByte(); err != nil || msgType != msgTypeAccess {
		if err != nil {
			return err
		}
		return fmt.Errorf("unexpected message type: %v", msgType)
	}
	//
	// senderLClock
	u32, err := mu.ReadUint32()
	if err != nil {
		return err
	}
	m.senderLClock = int(u32)
	//
	// receiverLClock
	u32, err = mu.ReadUint32()
	if err != nil {
		return err
	}
	m.receiverLClock = int(u32)
	//
	// accessForChains
	u32, err = mu.ReadUint32()
	if err != nil {
		return err
	}
	m.accessForChains = make([]isc.ChainID, u32)
	for i := range m.accessForChains {
		val, err2 := util.UnmarshalBytes(mu)
		if err2 != nil {
			return err2
		}
		chainID, err2 := isc.ChainIDFromBytes(val)
		if err2 != nil {
			return err2
		}
		m.accessForChains[i] = chainID
	}
	//
	// serverForChains
	u32, err = mu.ReadUint32()
	if err != nil {
		return err
	}
	m.serverForChains = make([]isc.ChainID, u32)
	for i := range m.serverForChains {
		val, err := util.UnmarshalBytes(mu)
		if err != nil {
			return err
		}
		chainID, err := isc.ChainIDFromBytes(val)
		if err != nil {
			return err
		}
		m.serverForChains[i] = chainID
	}
	return nil
}
