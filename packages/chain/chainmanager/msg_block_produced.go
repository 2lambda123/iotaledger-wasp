package chainmanager

import (
	"errors"
	"fmt"
	"io"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util"
)

// This message is used to inform access nodes on new blocks
// produced so that they can update their active state faster.
type msgBlockProduced struct {
	gpa.BasicMessage
	tx    *iotago.Transaction
	block state.Block
}

var _ gpa.Message = &msgCmtLog{}

func NewMsgBlockProduced(recipient gpa.NodeID, tx *iotago.Transaction, block state.Block) gpa.Message {
	return &msgBlockProduced{
		BasicMessage: gpa.NewBasicMessage(recipient),
		tx:           tx,
		block:        block,
	}
}

func (msg *msgBlockProduced) String() string {
	txID, err := msg.tx.ID()
	if err != nil {
		panic(fmt.Errorf("cannot extract TX ID: %w", err))
	}
	return fmt.Sprintf(
		"{chainMgr.msgBlockProduced, stateIndex=%v, l1Commitment=%v, tx.ID=%v}",
		msg.block.StateIndex(), msg.block.L1Commitment(), txID.ToHex(),
	)
}

func (msg *msgBlockProduced) MarshalBinary() (ret []byte, err error) {
	return util.WriterToBytes(msg), nil
}

func (msg *msgBlockProduced) UnmarshalBinary(data []byte) error {
	_, err := util.ReaderFromBytes(data, msg)
	return err
}

func (msg *msgBlockProduced) Read(r io.Reader) error {
	rr := util.NewReader(r)
	msgType := rr.ReadByte()
	if rr.Err == nil && msgType != msgTypeBlockProduced {
		return errors.New("unexpected message type")
	}
	msg.tx = new(iotago.Transaction)
	rr.ReadSerialized(msg.tx)
	msg.block = util.ReadFromBytes(rr, state.BlockFromBytes)
	return rr.Err
}

func (msg *msgBlockProduced) Write(w io.Writer) error {
	ww := util.NewWriter(w)
	ww.WriteByte(msgTypeBlockProduced)
	ww.WriteSerialized(msg.tx)
	ww.WriteFromBytes(msg.block)
	return ww.Err
}
