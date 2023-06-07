package chainmanager

import (
	"errors"
	"fmt"
	"io"

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
	return util.WriterToBytes(msg), nil
}

func (msg *msgCmtLog) UnmarshalBinary(data []byte) error {
	_, err := util.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgCmtLog) Read(r io.Reader) error {
	rr := util.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeCmtLog {
		return errors.New("unexpected message type")
	}
	rr.ReadN(msg.committeeAddr[:])
	wrappedMsgData := rr.ReadBytes()
	if rr.Err == nil {
		msg.wrapped, rr.Err = cmt_log.UnmarshalMessage(wrappedMsgData)
	}
	return rr.Err
}

func (msg *msgCmtLog) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeCmtLog)
	ww.WriteN(msg.committeeAddr[:])
	ww.WriteMarshaled(msg.wrapped)
	return ww.Err
}
