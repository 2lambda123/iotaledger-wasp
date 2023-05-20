package root

import (
	"bytes"
	"fmt"
	"time"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &DeployContractEvent{}

type DeployContractEvent struct {
	Timestamp   uint64
	Name        string
	Hname       isc.Hname
	ProgramHash hashing.HashValue
	Description string
}

func (e *DeployContractEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncDeployContract.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncDeployContract.Hname(): %w", err))
	}
	return w.Bytes()
}

func (e *DeployContractEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteString16(&w, e.Name); err != nil {
		panic(fmt.Errorf("failed to write event.Name: %w", err))
	}
	if err := util.WriteBytes8(&w, e.Hname.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.Hname: %w", err))
	}
	if err := util.WriteBytes32(&w, e.ProgramHash.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.ProgramHash: %w", err))
	}
	if err := util.WriteString16(&w, e.Description); err != nil {
		panic(fmt.Errorf("failed to write event.Description: %w", err))
	}
	return w.Bytes()
}

func (e *DeployContractEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}
	if err := util.ReadUint64(r, &e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to read event.Timestamp: %w", err))
	}
	str, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Name: %w", err))
	}
	e.Name = str

	b, err := util.ReadBytes8(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Hname: %w", err))
	}
	e.Hname, err = isc.HnameFromBytes(b)

	b, err = util.ReadBytes32(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.ProgramHash: %w", err))
	}
	e.ProgramHash, err = hashing.HashValueFromBytes(b)

	str, err = util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Description: %w", err))
	}
	e.Description = str
}

var _ isc.Event = &GrantDeployPermissionEvent{}

type GrantDeployPermissionEvent struct {
	Timestamp uint64
	AgentID   isc.AgentID
}

func (e *GrantDeployPermissionEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncGrantDeployPermission.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncGrantDeployPermission.Hname(): %w", err))
	}
	return w.Bytes()
}

func (e *GrantDeployPermissionEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteBytes8(&w, e.AgentID.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.AgentID: %w", err))
	}
	return w.Bytes()
}

func (e *GrantDeployPermissionEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}
	if err := util.ReadUint64(r, &e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to read event.Timestamp: %w", err))
	}
	agentIDBytes, err := util.ReadBytes8(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.AgentID: %w", err))
	}
	e.AgentID, err = isc.AgentIDFromBytes(agentIDBytes)
	if err != nil {
		panic(fmt.Errorf("failed to decode AgentID: %w", err))
	}
}

var _ isc.Event = &RevokeDeployPermissionEvent{}

type RevokeDeployPermissionEvent struct {
	Timestamp uint64
	AgentID   isc.AgentID
}

func (e *RevokeDeployPermissionEvent) Topic() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncRevokeDeployPermission.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncRevokeDeployPermission.Hname(): %w", err))
	}
	return w.Bytes()
}

func (e *RevokeDeployPermissionEvent) Payload() []byte {
	w := bytes.Buffer{}
	if err := util.WriteUint64(&w, uint64(time.Now().Unix())); err != nil {
		panic(fmt.Errorf("failed to write event.Timestamp: %w", err))
	}
	if err := util.WriteBytes8(&w, e.AgentID.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.AgentID: %w", err))
	}
	return w.Bytes()
}

func (e *RevokeDeployPermissionEvent) DecodePayload(payload []byte) {
	r := bytes.NewReader(payload)
	topic, err := util.ReadString16(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.Topic: %w", err))
	}
	if topic != string(e.Topic()) {
		panic("decode by unmatched event type")
	}
	if err := util.ReadUint64(r, &e.Timestamp); err != nil {
		panic(fmt.Errorf("failed to read event.Timestamp: %w", err))
	}
	agentIDBytes, err := util.ReadBytes8(r)
	if err != nil {
		panic(fmt.Errorf("failed to read event.AgentID: %w", err))
	}
	e.AgentID, err = isc.AgentIDFromBytes(agentIDBytes)
	if err != nil {
		panic(fmt.Errorf("failed to decode AgentID: %w", err))
	}
}
