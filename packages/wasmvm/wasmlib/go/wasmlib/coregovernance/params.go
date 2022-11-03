// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the schema definition file instead

package coregovernance

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"

type ImmutableAddAllowedStateControllerAddressParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableAddAllowedStateControllerAddressParams) StateControllerAddress() wasmtypes.ScImmutableAddress {
	return wasmtypes.NewScImmutableAddress(s.proxy.Root(ParamStateControllerAddress))
}

type MutableAddAllowedStateControllerAddressParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableAddAllowedStateControllerAddressParams) StateControllerAddress() wasmtypes.ScMutableAddress {
	return wasmtypes.NewScMutableAddress(s.proxy.Root(ParamStateControllerAddress))
}

type ImmutableAddCandidateNodeParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableAddCandidateNodeParams) AccessNodeInfoAccessAPI() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.proxy.Root(ParamAccessNodeInfoAccessAPI))
}

func (s ImmutableAddCandidateNodeParams) AccessNodeInfoCertificate() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.proxy.Root(ParamAccessNodeInfoCertificate))
}

func (s ImmutableAddCandidateNodeParams) AccessNodeInfoForCommittee() wasmtypes.ScImmutableBool {
	return wasmtypes.NewScImmutableBool(s.proxy.Root(ParamAccessNodeInfoForCommittee))
}

func (s ImmutableAddCandidateNodeParams) AccessNodeInfoPubKey() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.proxy.Root(ParamAccessNodeInfoPubKey))
}

type MutableAddCandidateNodeParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableAddCandidateNodeParams) AccessNodeInfoAccessAPI() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.proxy.Root(ParamAccessNodeInfoAccessAPI))
}

func (s MutableAddCandidateNodeParams) AccessNodeInfoCertificate() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.proxy.Root(ParamAccessNodeInfoCertificate))
}

func (s MutableAddCandidateNodeParams) AccessNodeInfoForCommittee() wasmtypes.ScMutableBool {
	return wasmtypes.NewScMutableBool(s.proxy.Root(ParamAccessNodeInfoForCommittee))
}

func (s MutableAddCandidateNodeParams) AccessNodeInfoPubKey() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.proxy.Root(ParamAccessNodeInfoPubKey))
}

type MapBytesToImmutableUint8 struct {
	proxy wasmtypes.Proxy
}

func (m MapBytesToImmutableUint8) GetUint8(key []byte) wasmtypes.ScImmutableUint8 {
	return wasmtypes.NewScImmutableUint8(m.proxy.Key(wasmtypes.BytesToBytes(key)))
}

type ImmutableChangeAccessNodesParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableChangeAccessNodesParams) ChangeAccessNodesActions() MapBytesToImmutableUint8 {
	return MapBytesToImmutableUint8{proxy: s.proxy.Root(ParamChangeAccessNodesActions)}
}

type MapBytesToMutableUint8 struct {
	proxy wasmtypes.Proxy
}

func (m MapBytesToMutableUint8) Clear() {
	m.proxy.ClearMap()
}

func (m MapBytesToMutableUint8) GetUint8(key []byte) wasmtypes.ScMutableUint8 {
	return wasmtypes.NewScMutableUint8(m.proxy.Key(wasmtypes.BytesToBytes(key)))
}

type MutableChangeAccessNodesParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableChangeAccessNodesParams) ChangeAccessNodesActions() MapBytesToMutableUint8 {
	return MapBytesToMutableUint8{proxy: s.proxy.Root(ParamChangeAccessNodesActions)}
}

type ImmutableDelegateChainOwnershipParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableDelegateChainOwnershipParams) ChainOwner() wasmtypes.ScImmutableAgentID {
	return wasmtypes.NewScImmutableAgentID(s.proxy.Root(ParamChainOwner))
}

type MutableDelegateChainOwnershipParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableDelegateChainOwnershipParams) ChainOwner() wasmtypes.ScMutableAgentID {
	return wasmtypes.NewScMutableAgentID(s.proxy.Root(ParamChainOwner))
}

type ImmutableRemoveAllowedStateControllerAddressParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableRemoveAllowedStateControllerAddressParams) StateControllerAddress() wasmtypes.ScImmutableAddress {
	return wasmtypes.NewScImmutableAddress(s.proxy.Root(ParamStateControllerAddress))
}

type MutableRemoveAllowedStateControllerAddressParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableRemoveAllowedStateControllerAddressParams) StateControllerAddress() wasmtypes.ScMutableAddress {
	return wasmtypes.NewScMutableAddress(s.proxy.Root(ParamStateControllerAddress))
}

type ImmutableRevokeAccessNodeParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableRevokeAccessNodeParams) AccessNodeInfoCertificate() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.proxy.Root(ParamAccessNodeInfoCertificate))
}

func (s ImmutableRevokeAccessNodeParams) AccessNodeInfoPubKey() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.proxy.Root(ParamAccessNodeInfoPubKey))
}

type MutableRevokeAccessNodeParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableRevokeAccessNodeParams) AccessNodeInfoCertificate() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.proxy.Root(ParamAccessNodeInfoCertificate))
}

func (s MutableRevokeAccessNodeParams) AccessNodeInfoPubKey() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.proxy.Root(ParamAccessNodeInfoPubKey))
}

type ImmutableRotateStateControllerParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableRotateStateControllerParams) StateControllerAddress() wasmtypes.ScImmutableAddress {
	return wasmtypes.NewScImmutableAddress(s.proxy.Root(ParamStateControllerAddress))
}

type MutableRotateStateControllerParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableRotateStateControllerParams) StateControllerAddress() wasmtypes.ScMutableAddress {
	return wasmtypes.NewScMutableAddress(s.proxy.Root(ParamStateControllerAddress))
}

type ImmutableSetChainInfoParams struct {
	proxy wasmtypes.Proxy
}

// default maximum size of a blob
func (s ImmutableSetChainInfoParams) MaxBlobSize() wasmtypes.ScImmutableUint32 {
	return wasmtypes.NewScImmutableUint32(s.proxy.Root(ParamMaxBlobSize))
}

// default maximum size of a single event
func (s ImmutableSetChainInfoParams) MaxEventSize() wasmtypes.ScImmutableUint16 {
	return wasmtypes.NewScImmutableUint16(s.proxy.Root(ParamMaxEventSize))
}

// default maximum number of events per request
func (s ImmutableSetChainInfoParams) MaxEventsPerReq() wasmtypes.ScImmutableUint16 {
	return wasmtypes.NewScImmutableUint16(s.proxy.Root(ParamMaxEventsPerReq))
}

type MutableSetChainInfoParams struct {
	proxy wasmtypes.Proxy
}

// default maximum size of a blob
func (s MutableSetChainInfoParams) MaxBlobSize() wasmtypes.ScMutableUint32 {
	return wasmtypes.NewScMutableUint32(s.proxy.Root(ParamMaxBlobSize))
}

// default maximum size of a single event
func (s MutableSetChainInfoParams) MaxEventSize() wasmtypes.ScMutableUint16 {
	return wasmtypes.NewScMutableUint16(s.proxy.Root(ParamMaxEventSize))
}

// default maximum number of events per request
func (s MutableSetChainInfoParams) MaxEventsPerReq() wasmtypes.ScMutableUint16 {
	return wasmtypes.NewScMutableUint16(s.proxy.Root(ParamMaxEventsPerReq))
}

type ImmutableSetFeePolicyParams struct {
	proxy wasmtypes.Proxy
}

func (s ImmutableSetFeePolicyParams) FeePolicyBytes() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.proxy.Root(ParamFeePolicyBytes))
}

type MutableSetFeePolicyParams struct {
	proxy wasmtypes.Proxy
}

func (s MutableSetFeePolicyParams) FeePolicyBytes() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.proxy.Root(ParamFeePolicyBytes))
}
