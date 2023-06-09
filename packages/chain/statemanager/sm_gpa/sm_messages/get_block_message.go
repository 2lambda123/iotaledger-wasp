package sm_messages

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
)

type GetBlockMessage struct {
	gpa.BasicMessage
	commitment *state.L1Commitment
}

var _ gpa.Message = new(GetBlockMessage)

func NewGetBlockMessage(commitment *state.L1Commitment, to gpa.NodeID) *GetBlockMessage {
	return &GetBlockMessage{
		BasicMessage: gpa.NewBasicMessage(to),
		commitment:   commitment,
	}
}

func NewEmptyGetBlockMessage() *GetBlockMessage { // `UnmarshalBinary` must be called afterwards
	return NewGetBlockMessage(&state.L1Commitment{}, gpa.NodeID{})
}

func (gbmT *GetBlockMessage) MarshalBinary() (data []byte, err error) {
	return append([]byte{MsgTypeGetBlockMessage}, gbmT.commitment.Bytes()...), nil
}

func (gbmT *GetBlockMessage) UnmarshalBinary(data []byte) error {
	if data[0] != MsgTypeGetBlockMessage {
		return fmt.Errorf("error creating get block message from bytes: wrong message type %v", data[0])
	}
	var err error
	gbmT.commitment, err = state.L1CommitmentFromBytes(data[1:])
	if err != nil {
		return fmt.Errorf("error creating get block message from bytes: %w", err)
	}
	return nil
}

func (gbmT *GetBlockMessage) GetL1Commitment() *state.L1Commitment {
	return gbmT.commitment
}
