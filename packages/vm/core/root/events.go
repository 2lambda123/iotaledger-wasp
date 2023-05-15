package root

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
)

var _ isc.Event = &DeployContractEvent{}

type DeployContractEvent struct {
	Name        string
	Hname       isc.Hname
	ProgramHash hashing.HashValue
	Description string
}

func (e *DeployContractEvent) Encode() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncDeployContract.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncDeployContract.Hname(): %w", err))
	}
	if err := util.WriteString16(&w, e.Name); err != nil {
		panic(fmt.Errorf("failed to write event.Name: %w", err))
	}
	if err := util.WriteBytes32(&w, e.ProgramHash.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.ProgramHash: %w", err))
	}
	if err := util.WriteString16(&w, e.Description); err != nil {
		panic(fmt.Errorf("failed to write event.Description: %w", err))
	}
	return w.Bytes()
}

var _ isc.Event = &GrantDeployPermissionEvent{}

type GrantDeployPermissionEvent struct {
	AgentID isc.AgentID
}

func (e *GrantDeployPermissionEvent) Encode() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncGrantDeployPermission.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncGrantDeployPermission.Hname(): %w", err))
	}
	if err := util.WriteBytes8(&w, e.AgentID.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.AgentID: %w", err))
	}
	return w.Bytes()
}

var _ isc.Event = &RevokeDeployPermissionEvent{}

type RevokeDeployPermissionEvent struct {
	AgentID isc.AgentID
}

func (e *RevokeDeployPermissionEvent) Encode() []byte {
	w := bytes.Buffer{}
	if err := util.WriteBytes8(&w, FuncRevokeDeployPermission.Hname().Bytes()); err != nil {
		panic(fmt.Errorf("failed to write FuncRevokeDeployPermission.Hname(): %w", err))
	}
	if err := util.WriteBytes8(&w, e.AgentID.Bytes()); err != nil {
		panic(fmt.Errorf("failed to write event.AgentID: %w", err))
	}
	return w.Bytes()
}
