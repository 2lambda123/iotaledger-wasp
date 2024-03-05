// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package bp

import (
	"fmt"
	"io"
	"time"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

type BatchProposal struct {
	l1API                   iotago.API                // Transient, for deserialization.
	nodeIndex               uint16                    // Just for a double-check.
	baseCO                  *isc.ChainOutputs         // Represents the consensus input received by the node[nodeIndex].
	baseBlockID             iotago.BlockID            // Represents the consensus input received by the node[nodeIndex].
	strongParents           iotago.BlockIDs           // Proposed TIPS to attach to.
	reattachTX              *iotago.SignedTransaction // The transaction to reattach, if any.
	dssTIndexProposal       util.BitVector            // DSS Index proposal for a TX.
	dssBIndexProposal       util.BitVector            // DSS Index proposal for a Block.
	timeData                time.Time                 // Our view of time.
	validatorFeeDestination isc.AgentID               // Proposed destination for fees.
	requestRefs             []*isc.RequestRef         // Requests we propose to include into the execution.
}

// This case is for deserialization.
func EmptyBatchProposal(l1API iotago.API) *BatchProposal {
	return &BatchProposal{l1API: l1API}
}

func NewBatchProposal(
	l1API iotago.API,
	nodeIndex uint16,
	baseBlock *iotago.Block,
	strongParents iotago.BlockIDs,
	baseCO *isc.ChainOutputs,
	reattachTX *iotago.SignedTransaction,
	dssTIndexProposal util.BitVector,
	dssBIndexProposal util.BitVector,
	timeData time.Time,
	validatorFeeDestination isc.AgentID,
	requestRefs []*isc.RequestRef,
) *BatchProposal {
	var baseBlockID iotago.BlockID
	var err error
	if baseBlock != nil {
		baseBlockID, err = baseBlock.ID()
		if err != nil {
			panic("cannot extract block id")
		}
	}
	return &BatchProposal{
		l1API:                   l1API,
		nodeIndex:               nodeIndex,
		baseCO:                  baseCO,
		baseBlockID:             baseBlockID,
		strongParents:           strongParents,
		reattachTX:              reattachTX,
		dssTIndexProposal:       dssTIndexProposal,
		dssBIndexProposal:       dssBIndexProposal,
		timeData:                timeData,
		validatorFeeDestination: validatorFeeDestination,
		requestRefs:             requestRefs,
	}
}

func (b *BatchProposal) Bytes() []byte {
	return rwutil.WriteToBytes(b)
}

func (b *BatchProposal) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	b.nodeIndex = rr.ReadUint16()

	if rr.ReadBool() {
		b.baseCO = new(isc.ChainOutputs)
		rr.Read(b.baseCO)
	} else {
		b.baseCO = nil
	}

	rr.ReadN(b.baseBlockID[:])

	spCount := rr.ReadInt8()
	b.strongParents = make(iotago.BlockIDs, spCount)
	for _, sp := range b.strongParents {
		rr.ReadN(sp[:])
	}

	if rr.ReadBool() {
		txBytes := rr.ReadBytes()
		b.reattachTX = new(iotago.SignedTransaction)
		b.reattachTX.API = b.l1API
		_, err := b.reattachTX.API.Decode(txBytes, b.reattachTX)
		if err != nil {
			return err
		}
	} else {
		b.reattachTX = nil
	}
	b.dssTIndexProposal = util.NewFixedSizeBitVector(0)
	rr.Read(b.dssTIndexProposal)
	b.dssBIndexProposal = util.NewFixedSizeBitVector(0)
	rr.Read(b.dssBIndexProposal)
	b.timeData = time.Unix(0, rr.ReadInt64())
	b.validatorFeeDestination = isc.AgentIDFromReader(rr)
	size := rr.ReadSize16()
	b.requestRefs = make([]*isc.RequestRef, size)
	for i := range b.requestRefs {
		b.requestRefs[i] = new(isc.RequestRef)
		rr.ReadN(b.requestRefs[i].ID[:])
		rr.ReadN(b.requestRefs[i].Hash[:])
	}
	return rr.Err
}

func (b *BatchProposal) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteUint16(b.nodeIndex)

	ww.WriteBool(b.baseCO != nil)
	if b.baseCO != nil {
		ww.Write(b.baseCO)
	}

	ww.WriteN(b.baseBlockID[:])

	ww.WriteInt8(int8(len(b.strongParents)))
	for _, sp := range b.strongParents {
		ww.WriteN(sp[:])
	}

	ww.WriteBool(b.reattachTX != nil)
	if b.reattachTX != nil {
		bs, err := b.reattachTX.API.Encode(b.reattachTX)
		if err != nil {
			panic(fmt.Errorf("cannot encode the TX: %v", err))
		}
		ww.WriteBytes(bs)
	}
	ww.Write(b.dssTIndexProposal)
	ww.Write(b.dssBIndexProposal)
	ww.WriteInt64(b.timeData.UnixNano())
	ww.Write(b.validatorFeeDestination)
	ww.WriteSize16(len(b.requestRefs))
	for i := range b.requestRefs {
		ww.WriteN(b.requestRefs[i].ID[:])
		ww.WriteN(b.requestRefs[i].Hash[:])
	}
	return ww.Err
}
