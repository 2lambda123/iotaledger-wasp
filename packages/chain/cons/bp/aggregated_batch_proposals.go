// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package bp

import (
	"bytes"
	"sort"
	"time"

	"github.com/iotaledger/hive.go/logger"
	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// Here we store just an aggregated info.
type AggregatedBatchProposals struct {
	shouldBeSkipped           bool
	batchProposalSet          batchProposalSet
	decidedDSStIndexProposals map[gpa.NodeID][]int
	decidedDSSbIndexProposals map[gpa.NodeID][]int
	decidedBaseBlockID        *iotago.BlockID
	decidedBaseCO             *isc.ChainOutputs
	decidedReattachTX         *iotago.SignedTransaction
	decidedRequestRefs        []*isc.RequestRef
	aggregatedTime            time.Time
}

func AggregateBatchProposals(inputs map[gpa.NodeID][]byte, nodeIDs []gpa.NodeID, f int, l1API iotago.API, log *logger.Logger) *AggregatedBatchProposals {
	bps := batchProposalSet{}
	//
	// Parse and validate the batch proposals. Skip the invalid ones.
	for nid := range inputs {
		batchProposal := EmptyBatchProposal(l1API)
		batchProposal, err := rwutil.ReadFromBytes(inputs[nid], batchProposal)
		if err != nil {
			log.Warnf("cannot decode BatchProposal from %v: %v", nid, err)
			continue
		}
		if int(batchProposal.nodeIndex) >= len(nodeIDs) || nodeIDs[batchProposal.nodeIndex] != nid {
			log.Warnf("invalid nodeIndex=%v in batchProposal from %v", batchProposal.nodeIndex, nid)
			continue
		}
		bps[nid] = batchProposal
	}
	//
	// Store the aggregated values.
	if len(bps) == 0 {
		log.Debugf("Cant' aggregate batch proposal: have 0 batch proposals.")
		return &AggregatedBatchProposals{shouldBeSkipped: true}
	}
	aggregatedTime := bps.aggregatedTime(f)
	decidedBaseCO := bps.decidedBaseAnchorOutput(f)
	abp := &AggregatedBatchProposals{
		batchProposalSet:          bps,
		decidedDSStIndexProposals: bps.decidedDSStIndexProposals(),
		decidedDSSbIndexProposals: bps.decidedDSSbIndexProposals(),
		decidedBaseCO:             decidedBaseCO,
		decidedBaseBlockID:        bps.decidedBaseBlockID(f),
		decidedReattachTX:         bps.decidedReattachTX(f),
		decidedRequestRefs:        bps.decidedRequestRefs(f, decidedBaseCO),
		aggregatedTime:            aggregatedTime,
	}
	if abp.decidedBaseCO == nil && abp.decidedReattachTX == nil {
		log.Debugf("Cant' aggregate batch proposal: decidedBaseCO and decidedReattachTX are both nil.")
		abp.shouldBeSkipped = true
	}
	if abp.decidedBaseCO != nil && abp.decidedReattachTX != nil {
		log.Debugf("Cant' aggregate batch proposal: decidedBaseCO and decidedReattachTX are both non-nil.")
		abp.shouldBeSkipped = true
	}
	if abp.decidedBaseCO != nil && len(abp.decidedRequestRefs) == 0 {
		log.Debugf("Cant' aggregate batch proposal: decidedBaseCO is non-nil, but there is no decided requests.")
		abp.shouldBeSkipped = true
	}
	if abp.aggregatedTime.IsZero() {
		log.Debugf("Cant' aggregate batch proposal: aggregatedTime is zero")
		abp.shouldBeSkipped = true
	}
	return abp
}

func (abp *AggregatedBatchProposals) ShouldBeSkipped() bool {
	return abp.shouldBeSkipped
}

func (abp *AggregatedBatchProposals) ShouldBuildNewTX() bool {
	return !abp.shouldBeSkipped && abp.decidedBaseCO != nil
}

func (abp *AggregatedBatchProposals) DecidedReattachTX() *iotago.SignedTransaction {
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	if abp.decidedReattachTX == nil {
		panic("trying to use reattach TX id when no TX was decided to be reused")
	}
	return abp.decidedReattachTX
}

func (abp *AggregatedBatchProposals) DecidedDSStIndexProposals() map[gpa.NodeID][]int {
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	return abp.decidedDSStIndexProposals
}

func (abp *AggregatedBatchProposals) DecidedDSSbIndexProposals() map[gpa.NodeID][]int {
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	return abp.decidedDSSbIndexProposals
}

func (abp *AggregatedBatchProposals) DecidedBaseCO() *isc.ChainOutputs { // TODO: Use it as one of the parents, if non-nil.
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	return abp.decidedBaseCO
}

func (abp *AggregatedBatchProposals) DecidedStrongParents(randomness hashing.HashValue) iotago.BlockIDs {
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	return abp.batchProposalSet.decidedStrongParents(abp.aggregatedTime, randomness)
}

func (abp *AggregatedBatchProposals) AggregatedTime() time.Time {
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	return abp.aggregatedTime
}

func (abp *AggregatedBatchProposals) ValidatorFeeTarget(randomness hashing.HashValue) isc.AgentID {
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	return abp.batchProposalSet.selectedFeeDestination(abp.aggregatedTime, randomness)
}

func (abp *AggregatedBatchProposals) DecidedRequestRefs() []*isc.RequestRef {
	if abp.shouldBeSkipped {
		panic("trying to use aggregated proposal marked to be skipped")
	}
	return abp.decidedRequestRefs
}

// TODO should this be moved to the VM?
func (abp *AggregatedBatchProposals) OrderedRequests(requests []isc.Request, randomness hashing.HashValue) []isc.Request {
	type sortStruct struct {
		key hashing.HashValue
		ref *isc.RequestRef
		req isc.Request
	}

	sortBuf := make([]*sortStruct, len(abp.decidedRequestRefs))
	for i := range abp.decidedRequestRefs {
		ref := abp.decidedRequestRefs[i]
		var found isc.Request
		for j := range requests {
			if ref.IsFor(requests[j]) {
				found = requests[j]
				break
			}
		}
		if found == nil {
			panic("request was not provided by mempool")
		}
		sortBuf[i] = &sortStruct{
			key: hashing.HashDataBlake2b(ref.ID.Bytes(), ref.Hash[:], randomness[:]),
			ref: ref,
			req: found,
		}
	}
	sort.Slice(sortBuf, func(i, j int) bool {
		return bytes.Compare(sortBuf[i].key[:], sortBuf[j].key[:]) < 0
	})

	// Make sure the requests are sorted such way, that the nonces per account are increasing.
	// This is needed to handle several requests per batch for the VMs that expect the in-order nonces.
	// We make a second pass here to tain the overall ordering of requests (module account) without
	// making requests from a single account grouped together while sorting.
	for i := range sortBuf {
		oi, ok := sortBuf[i].req.(isc.OffLedgerRequest)
		if !ok {
			continue
		}
		for j := i + 1; j < len(sortBuf); j++ {
			oj, ok := sortBuf[j].req.(isc.OffLedgerRequest)
			if !ok {
				continue
			}
			if oi.SenderAccount().Equals(oj.SenderAccount()) && oi.Nonce() > oj.Nonce() {
				sortBuf[i], sortBuf[j] = sortBuf[j], sortBuf[i]
				oi = oj
			}
		}
	}

	sorted := make([]isc.Request, len(abp.decidedRequestRefs))
	for i := range sortBuf {
		sorted[i] = sortBuf[i].req
	}
	return sorted
}
