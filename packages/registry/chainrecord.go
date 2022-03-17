// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"fmt"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/hive.go/marshalutil"
	"github.com/iotaledger/wasp/packages/database/textdb"
	"github.com/iotaledger/wasp/packages/iscp"
)

// ChainRecord represents chain the node is participating in
// TODO optimize, no need for a persistent structure, simple activity tag is enough
type ChainRecord struct {
	ChainID *iscp.ChainID
	Active  bool
}

type ChainRecordText struct {
	ChainID string `json:"chainid" yaml:"chainid"`
	Active  bool   `json:"active" yaml:"active"`
}

func NewChainRecordText(rec *ChainRecord) *ChainRecordText {
	return &ChainRecordText{ChainID: rec.ChainID.Base58(), Active: rec.Active}
}

func FromMarshalUtil(mu *marshalutil.MarshalUtil) (*ChainRecord, error) {
	ret := &ChainRecord{}
	aliasAddr, err := ledgerstate.AliasAddressFromMarshalUtil(mu)
	if err != nil {
		return nil, err
	}
	ret.ChainID = iscp.NewChainID(aliasAddr)

	ret.Active, err = mu.ReadBool()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// CommitteeRecordFromBytes
func ChainRecordFromBytes(data []byte) (*ChainRecord, error) {
	return FromMarshalUtil(marshalutil.New(data))
}

func (rec *ChainRecord) Bytes() []byte {
	mu := marshalutil.New().WriteBytes(rec.ChainID.Bytes()).
		WriteBool(rec.Active)
	return mu.Bytes()
}

func (rec *ChainRecord) String() string {
	ret := "ChainID: " + rec.ChainID.String() + "\n"
	ret += fmt.Sprintf("      Active: %v\n", rec.Active)
	return ret
}

func (rec *ChainRecord) toText(m textdb.Marshaller) ([]byte, error) {
	obj := NewChainRecordText(rec)
	return m.Marshal(obj)
}

func ChainRecordFromText(in []byte, m textdb.Marshaller) (*ChainRecord, error) {
	obj := ChainRecordText{}
	err := m.Unmarshal(in, &obj)
	if err != nil {
		return nil, err
	}
	chID, err := iscp.ChainIDFromBase58(obj.ChainID)
	if err != nil {
		return nil, err
	}
	return &ChainRecord{ChainID: chID, Active: obj.Active}, nil
}
