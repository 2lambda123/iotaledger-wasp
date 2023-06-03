// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package gpa

import (
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/util"
)

// MsgWrapper can be used to compose an algorithm out of other abstractions.
// These messages are meant to wrap and route the messages of the sub-algorithms.
type MsgWrapper struct {
	msgType       byte
	subsystemFunc func(subsystem byte, index int) (GPA, error) // Resolve a subsystem GPA based on its code and index.
}

func NewMsgWrapper(msgType byte, subsystemFunc func(subsystem byte, index int) (GPA, error)) *MsgWrapper {
	return &MsgWrapper{msgType, subsystemFunc}
}

func (w *MsgWrapper) WrapMessage(subsystem byte, index int, msg Message) Message {
	return &WrappingMsg{w.msgType, subsystem, index, msg}
}

func (w *MsgWrapper) WrapMessages(subsystem byte, index int, msgs OutMessages) OutMessages {
	if msgs == nil {
		return nil
	}
	wrapped := NoMessages()
	msgs.MustIterate(func(msg Message) {
		wrapped.Add(w.WrapMessage(subsystem, index, msg))
	})
	return wrapped
}

func (w *MsgWrapper) DelegateInput(subsystem byte, index int, input Input) (GPA, OutMessages, error) {
	sub, err := w.subsystemFunc(subsystem, index)
	if err != nil {
		return nil, nil, err
	}
	return sub, w.WrapMessages(subsystem, index, sub.Input(input)), nil
}

func (w *MsgWrapper) DelegateMessage(msg *WrappingMsg) (GPA, OutMessages, error) {
	sub, err := w.subsystemFunc(msg.Subsystem(), msg.Index())
	if err != nil {
		return nil, nil, err
	}
	return sub, w.WrapMessages(msg.Subsystem(), msg.Index(), sub.Message(msg.Wrapped())), nil
}

func (w *MsgWrapper) UnmarshalMessage(data []byte) (Message, error) {
	mu := marshalutil.New(data)
	msgType, err := mu.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("cannot decode MsgWrapper::msgType: %v", msgType)
	}
	if msgType != w.msgType {
		return nil, fmt.Errorf("invalid MsgWrapper::msgType, got %v, expected %v", msgType, w.msgType)
	}
	subsystem, err := mu.ReadByte()
	if err != nil {
		return nil, err
	}
	indexU16, err := mu.ReadUint16()
	if err != nil {
		return nil, err
	}
	index := int(indexU16)
	wrappedBin, err := util.UnmarshallBytes(mu)
	if err != nil {
		return nil, err
	}
	subGPA, err := w.subsystemFunc(subsystem, index)
	if err != nil {
		return nil, err
	}
	wrapped, err := subGPA.UnmarshalMessage(wrappedBin)
	if err != nil {
		return nil, err
	}
	return &WrappingMsg{msgType, subsystem, index, wrapped}, nil
}

// The message that contains another, and its routing info.
type WrappingMsg struct {
	msgType   byte
	subsystem byte
	index     int
	wrapped   Message
}

var _ Message = &WrappingMsg{}

func NewWrappingMsg(msgType, subsystem byte, index int, wrapped Message) *WrappingMsg {
	return &WrappingMsg{msgType: msgType, subsystem: subsystem, index: index, wrapped: wrapped}
}

func (m *WrappingMsg) Subsystem() byte {
	return m.subsystem
}

func (m *WrappingMsg) Index() int {
	return m.index
}

func (m *WrappingMsg) Wrapped() Message {
	return m.wrapped
}

func (m *WrappingMsg) Recipient() NodeID {
	return m.wrapped.Recipient()
}

func (m *WrappingMsg) SetSender(sender NodeID) {
	m.wrapped.SetSender(sender)
}

func (m *WrappingMsg) MarshalBinary() ([]byte, error) {
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(m.msgType)
	mu.WriteByte(m.subsystem)
	mu.WriteUint16(uint16(m.index))
	bin, err := m.wrapped.MarshalBinary()
	if err != nil {
		return nil, err
	}
	util.MarshallBytes(mu, bin)
	return mu.Bytes(), nil
}

func (m *WrappingMsg) UnmarshalBinary(data []byte) error {
	panic("this message is un-marshaled by the gpa.MsgWrapper")
}
