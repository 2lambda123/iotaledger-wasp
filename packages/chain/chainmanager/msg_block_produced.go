package chainmanager

import (
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2"
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
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
	mu := new(marshalutil.MarshalUtil)
	mu.WriteByte(msgTypeBlockProduced)
	txBytes, err := msg.tx.Serialize(serializer.DeSeriModeNoValidation, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize tx: %w", err)
	}
	util.MarshallBytes(mu, txBytes)
	util.MarshallBytes(mu, msg.block.Bytes())
	return mu.Bytes(), nil
}

func (msg *msgBlockProduced) UnmarshalBinary(data []byte) error {
	mu := marshalutil.New(data)

	msgType, err := mu.ReadByte()
	if err != nil {
		return fmt.Errorf("cannot read msgType byte: %w", err)
	}
	if msgType != msgTypeBlockProduced {
		return fmt.Errorf("unexpected msgType: %v", msgType)
	}

	txBytes, err := util.UnmarshallBytes(mu)
	if err != nil {
		return fmt.Errorf("cannot read tx bytes: %w", err)
	}
	msg.tx = &iotago.Transaction{}
	_, err = msg.tx.Deserialize(txBytes, serializer.DeSeriModeNoValidation, nil)
	if err != nil {
		return fmt.Errorf("cannot deserialize tx: %w", err)
	}

	blockBytes, err := util.UnmarshallBytes(mu)
	if err != nil {
		return fmt.Errorf("cannot read block bytes: %w", err)
	}
	msg.block, err = state.BlockFromBytes(blockBytes)
	if err != nil {
		return fmt.Errorf("cannot deserialize block: %w", err)
	}
	return nil
}
