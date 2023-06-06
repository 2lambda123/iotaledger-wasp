package chainmanager

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2"
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

func (msg *msgBlockProduced) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)
	ww := util.NewWriter(w)

	ww.WriteByte(msgTypeBlockProduced)

	// TX
	txBytes, err := msg.tx.Serialize(serializer.DeSeriModeNoValidation, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize tx: %w", err)
	}
	ww.WriteBytes(txBytes)

	// Block
	ww.WriteBytes(msg.block.Bytes())
	return w.Bytes(), nil
}

func (msg *msgBlockProduced) UnmarshalBinary(data []byte) error {
	var err error
	r := bytes.NewReader(data)
	rr := util.NewReader(r)

	// MsgType
	if msgType := rr.ReadByte(); msgType != msgTypeBlockProduced {
		return fmt.Errorf("unexpected msgType: %v", msgType)
	}

	// TX
	txBytes := rr.ReadBytes()
	tx := &iotago.Transaction{}
	_, err = tx.Deserialize(txBytes, serializer.DeSeriModeNoValidation, nil)
	if err != nil {
		return fmt.Errorf("cannot deserialize tx: %w", err)
	}
	msg.tx = tx

	// Block
	blockBytes := rr.ReadBytes()
	block, err := state.BlockFromBytes(blockBytes)
	if err != nil {
		return fmt.Errorf("cannot deserialize block: %w", err)
	}
	msg.block = block
	return nil
}
