// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package governance

import (
	"bytes"
	"fmt"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
)

// NodeOwnershipCertificate is a proof that a specified address is an owner of the specified node.
// It is implemented as a signature over the node pub key concatenated with the owner address.
type NodeOwnershipCertificate []byte

func NewNodeOwnershipCertificate(nodeKeyPair *cryptolib.KeyPair, ownerAddress iotago.Address) NodeOwnershipCertificate {
	mu := new(marshalutil.MarshalUtil)
	mu.WriteBytes(nodeKeyPair.GetPublicKey().AsBytes())
	mu.WriteBytes(isc.BytesFromAddress(ownerAddress))
	return nodeKeyPair.GetPrivateKey().Sign(mu.Bytes())
}

func NewNodeOwnershipCertificateFromBytes(data []byte) NodeOwnershipCertificate {
	return data
}

func (c NodeOwnershipCertificate) Verify(nodePubKey *cryptolib.PublicKey, ownerAddress iotago.Address) bool {
	mu := new(marshalutil.MarshalUtil)
	mu.WriteBytes(nodePubKey.AsBytes())
	mu.WriteBytes(isc.BytesFromAddress(ownerAddress))
	return nodePubKey.Verify(mu.Bytes(), c.Bytes())
}

func (c NodeOwnershipCertificate) Bytes() []byte {
	return c
}

// AccessNodeInfo conveys all the information that is maintained
// on the governance SC about a specific node.
type AccessNodeInfo struct {
	NodePubKey    []byte // Public Key of the node. Stored as a key in the SC State and Params.
	ValidatorAddr []byte // Address of the validator owning the node. Not sent via parameters.
	Certificate   []byte // Proof that Validator owns the Node.
	ForCommittee  bool   // true, if Node should be a candidate to a committee.
	AccessAPI     string // API URL, if any.
}

func NewAccessNodeInfoFromBytes(pubKey, value []byte) (*AccessNodeInfo, error) {
	var a AccessNodeInfo
	var err error
	r := bytes.NewReader(value)
	a.NodePubKey = pubKey // NodePubKey stored as a map key.
	if a.ValidatorAddr, err = util.ReadBytes(r); err != nil {
		return nil, fmt.Errorf("failed to read AccessNodeInfo.ValidatorAddr: %w", err)
	}
	if a.Certificate, err = util.ReadBytes(r); err != nil {
		return nil, fmt.Errorf("failed to read AccessNodeInfo.Certificate: %w", err)
	}
	if err2 := util.ReadBool(r, &a.ForCommittee); err2 != nil {
		return nil, fmt.Errorf("failed to read AccessNodeInfo.ForCommittee: %w", err2)
	}
	if a.AccessAPI, err = util.ReadString(r); err != nil {
		return nil, fmt.Errorf("failed to read AccessNodeInfo.AccessAPI: %w", err)
	}
	return &a, nil
}

func NewAccessNodeInfoListFromMap(infoMap *collections.ImmutableMap) ([]*AccessNodeInfo, error) {
	res := make([]*AccessNodeInfo, 0)
	var accErr error
	infoMap.Iterate(func(elemKey, value []byte) bool {
		var a *AccessNodeInfo
		if a, accErr = NewAccessNodeInfoFromBytes(elemKey, value); accErr != nil {
			return false
		}
		res = append(res, a)
		return true
	})
	if accErr != nil {
		return nil, fmt.Errorf("failed to iterate over AccessNodeInfo list: %w", accErr)
	}
	return res, nil
}

func (a *AccessNodeInfo) Bytes() []byte {
	mu := new(marshalutil.MarshalUtil)
	util.WriteBytesMu(mu, a.ValidatorAddr)
	util.WriteBytesMu(mu, a.Certificate)
	mu.WriteBool(a.ForCommittee)
	util.WriteStringMu(mu, a.AccessAPI)
	return mu.Bytes()
}

var errSenderMustHaveL1Address = coreerrors.Register("sender must have L1 address").Create()

func NewAccessNodeInfoFromAddCandidateNodeParams(ctx isc.Sandbox) *AccessNodeInfo {
	validatorAddr, _ := isc.AddressFromAgentID(ctx.Request().SenderAccount()) // Not from params, to have it validated.
	if validatorAddr == nil {
		panic(errSenderMustHaveL1Address)
	}
	params := ctx.Params()
	ani := AccessNodeInfo{
		NodePubKey:    params.MustGetBytes(ParamAccessNodeInfoPubKey),
		ValidatorAddr: isc.BytesFromAddress(validatorAddr),
		Certificate:   params.MustGetBytes(ParamAccessNodeInfoCertificate),
		ForCommittee:  params.MustGetBool(ParamAccessNodeInfoForCommittee, false),
		AccessAPI:     params.MustGetString(ParamAccessNodeInfoAccessAPI, ""),
	}
	return &ani
}

func (a *AccessNodeInfo) ToAddCandidateNodeParams() dict.Dict {
	d := dict.New()
	d.Set(ParamAccessNodeInfoForCommittee, codec.EncodeBool(a.ForCommittee))
	d.Set(ParamAccessNodeInfoPubKey, a.NodePubKey)
	d.Set(ParamAccessNodeInfoCertificate, a.Certificate)
	d.Set(ParamAccessNodeInfoAccessAPI, codec.EncodeString(a.AccessAPI))
	return d
}

func NewAccessNodeInfoFromRevokeAccessNodeParams(ctx isc.Sandbox) *AccessNodeInfo {
	validatorAddr, _ := isc.AddressFromAgentID(ctx.Request().SenderAccount()) // Not from params, to have it validated.
	if validatorAddr == nil {
		panic(errSenderMustHaveL1Address)
	}
	params := ctx.Params()
	ani := AccessNodeInfo{
		NodePubKey:    params.MustGetBytes(ParamAccessNodeInfoPubKey),
		ValidatorAddr: isc.BytesFromAddress(validatorAddr), // Not from params, to have it validated.
		Certificate:   params.MustGetBytes(ParamAccessNodeInfoCertificate),
	}
	return &ani
}

func (a *AccessNodeInfo) ToRevokeAccessNodeParams() dict.Dict {
	d := dict.New()
	d.Set(ParamAccessNodeInfoPubKey, a.NodePubKey)
	d.Set(ParamAccessNodeInfoCertificate, a.Certificate)
	return d
}

func (a *AccessNodeInfo) AddCertificate(nodeKeyPair *cryptolib.KeyPair, ownerAddress iotago.Address) *AccessNodeInfo {
	a.Certificate = NewNodeOwnershipCertificate(nodeKeyPair, ownerAddress).Bytes()
	return a
}

func (a *AccessNodeInfo) ValidateCertificate(ctx isc.Sandbox) bool {
	nodePubKey, err := cryptolib.NewPublicKeyFromBytes(a.NodePubKey)
	if err != nil {
		return false
	}
	validatorAddr, _, err := isc.AddressFromBytes(a.ValidatorAddr)
	if err != nil {
		return false
	}
	cert := NewNodeOwnershipCertificateFromBytes(a.Certificate)
	return cert.Verify(nodePubKey, validatorAddr)
}

// GetChainNodesRequest
type GetChainNodesRequest struct{}

func (req GetChainNodesRequest) AsDict() dict.Dict {
	return dict.New()
}

// GetChainNodesResponse
type GetChainNodesResponse struct {
	AccessNodeCandidates []*AccessNodeInfo      // Application info for the AccessNodes.
	AccessNodes          []*cryptolib.PublicKey // Public Keys of Access Nodes.
}

func NewGetChainNodesResponseFromDict(d dict.Dict) *GetChainNodesResponse {
	res := GetChainNodesResponse{
		AccessNodeCandidates: make([]*AccessNodeInfo, 0),
		AccessNodes:          make([]*cryptolib.PublicKey, 0),
	}

	ac := collections.NewMapReadOnly(d, ParamGetChainNodesAccessNodeCandidates)
	ac.Iterate(func(pubKey, value []byte) bool {
		ani, err := NewAccessNodeInfoFromBytes(pubKey, value)
		if err != nil {
			panic(fmt.Errorf("unable to decode access node info: %w", err))
		}
		res.AccessNodeCandidates = append(res.AccessNodeCandidates, ani)
		return true
	})

	an := collections.NewMapReadOnly(d, ParamGetChainNodesAccessNodes)
	an.Iterate(func(pubKeyBin, value []byte) bool {
		publicKey, err := cryptolib.NewPublicKeyFromBytes(pubKeyBin)
		if err != nil {
			panic(fmt.Errorf("unable to decode public key: %w", err))
		}
		res.AccessNodes = append(res.AccessNodes, publicKey)
		return true
	})
	return &res
}

//
//	ChangeAccessNodesRequest
//

type ChangeAccessNodeAction byte

const (
	ChangeAccessNodeActionRemove = ChangeAccessNodeAction(iota)
	ChangeAccessNodeActionAccept
	ChangeAccessNodeActionDrop
)

type ChangeAccessNodesRequest struct {
	actions map[cryptolib.PublicKeyKey]ChangeAccessNodeAction
}

func NewChangeAccessNodesRequest() *ChangeAccessNodesRequest {
	return &ChangeAccessNodesRequest{
		actions: make(map[cryptolib.PublicKeyKey]ChangeAccessNodeAction),
	}
}

func (req *ChangeAccessNodesRequest) Remove(pubKey *cryptolib.PublicKey) *ChangeAccessNodesRequest {
	req.actions[pubKey.AsKey()] = ChangeAccessNodeActionRemove
	return req
}

func (req *ChangeAccessNodesRequest) Accept(pubKey *cryptolib.PublicKey) *ChangeAccessNodesRequest {
	req.actions[pubKey.AsKey()] = ChangeAccessNodeActionAccept
	return req
}

func (req *ChangeAccessNodesRequest) Drop(pubKey *cryptolib.PublicKey) *ChangeAccessNodesRequest {
	req.actions[pubKey.AsKey()] = ChangeAccessNodeActionDrop
	return req
}

func (req *ChangeAccessNodesRequest) AsDict() dict.Dict {
	d := dict.New()
	actionsMap := collections.NewMap(d, ParamChangeAccessNodesActions)
	for pubKey, action := range req.actions {
		actionsMap.SetAt(pubKey[:], []byte{byte(action)})
	}
	return d
}
