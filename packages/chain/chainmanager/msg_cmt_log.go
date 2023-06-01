package chainmanager

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/chain/cmt_log"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/util"
)

// gpa.Wrapper is not applicable here, because here the addressing
// is by CommitteeID, not by integer index.
type msgCmtLog struct {
	committeeAddr iotago.Ed25519Address
	wrapped       gpa.Message
}

var _ gpa.Message = &msgCmtLog{}

func NewMsgCmtLog(committeeAddr iotago.Ed25519Address, wrapped gpa.Message) gpa.Message {
	return &msgCmtLog{
		committeeAddr: committeeAddr,
		wrapped:       wrapped,
	}
}

func (msg *msgCmtLog) String() string {
	return fmt.Sprintf("{chainMgr.msgCmtLog, committeeAddr=%v, wrapped=%+v}", msg.committeeAddr.String(), msg.wrapped)
}

func (msg *msgCmtLog) Recipient() gpa.NodeID {
	return msg.wrapped.Recipient()
}

func (msg *msgCmtLog) SetSender(sender gpa.NodeID) {
	msg.wrapped.SetSender(sender)
}

func (msg *msgCmtLog) MarshalBinary() ([]byte, error) {
	w := bytes.NewBuffer([]byte{})
	if err := util.WriteByte(w, msgTypeCmtLog); err != nil {
		return nil, fmt.Errorf("cannot serialize msgType: %w", err)
	}
	committeeAddrBytes, err := msg.committeeAddr.Serialize(serializer.DeSeriModeNoValidation, nil)
	if err != nil {
		return nil, err
	}
	if err2 := util.WriteBytes(w, committeeAddrBytes); err2 != nil {
		return nil, err2
	}
	bin, err := msg.wrapped.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := util.WriteBytes(w, bin); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (msg *msgCmtLog) UnmarshalBinary(data []byte) error {
	var err error
	r := bytes.NewReader(data)
	//
	// MsgType
	msgType, err := util.ReadByte(r)
	if err != nil {
		return fmt.Errorf("cannot read msgType byte: %w", err)
	}
	if msgType != msgTypeCmtLog {
		return fmt.Errorf("unexpected msgType: %v", msgType)
	}
	//
	committeeAddrBytes, err := util.ReadBytes(r)
	if err != nil {
		return err
	}
	_, err = msg.committeeAddr.Deserialize(committeeAddrBytes, serializer.DeSeriModeNoValidation, nil)
	if err != nil {
		return err
	}
	wrappedMsgData, err := util.ReadBytes(r)
	if err != nil {
		return err
	}
	msg.wrapped, err = cmt_log.UnmarshalMessage(wrappedMsgData)
	if err != nil {
		return err
	}
	return nil
}
