package blocklog

import (
	"io"

	iotago "github.com/iotaledger/iota.go/v4"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util/rwutil"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
)

type unprocessableRequestRecord struct {
	outputID iotago.OutputID
	req      isc.Request
}

func unprocessableRequestRecordFromBytes(data []byte) (*unprocessableRequestRecord, error) {
	return rwutil.ReadFromBytes(data, new(unprocessableRequestRecord))
}

func mustUnprocessableRequestRecordFromBytes(data []byte) *unprocessableRequestRecord {
	rec, err := unprocessableRequestRecordFromBytes(data)
	if err != nil {
		panic(err)
	}
	return rec
}

func (rec *unprocessableRequestRecord) Bytes() []byte {
	return rwutil.WriteToBytes(rec)
}

func (rec *unprocessableRequestRecord) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadN(rec.outputID[:])
	rec.req = isc.RequestFromReader(rr)
	return rr.Err
}

func (rec *unprocessableRequestRecord) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteN(rec.outputID[:])
	ww.Write(rec.req)
	return ww.Err
}

func (s *StateWriter) newUnprocessableRequestsArray() *collections.Array {
	return collections.NewArray(s.state, prefixNewUnprocessableRequests)
}

func (s *StateWriter) unprocessableMap() *collections.Map {
	return collections.NewMap(s.state, prefixUnprocessableRequests)
}

func (s *StateReader) unprocessableMap() *collections.ImmutableMap {
	return collections.NewMapReadOnly(s.state, prefixUnprocessableRequests)
}

// save request reference / address of the sender
func (s *StateWriter) SaveUnprocessable(req isc.OnLedgerRequest, blockIndex uint32, outputIndex uint16) {
	rec := unprocessableRequestRecord{
		// TransactionID is unknown yet, will be filled next block
		outputID: iotago.OutputIDFromTransactionIDAndIndex(iotago.TransactionID{}, outputIndex),
		req:      req,
	}
	s.unprocessableMap().SetAt(req.ID().Bytes(), rec.Bytes())
	s.newUnprocessableRequestsArray().Push(req.ID().Bytes())
}

func (s *StateWriter) updateUnprocessableRequestsOutputID(anchorTxID iotago.TransactionID) {
	newReqs := s.newUnprocessableRequestsArray()
	allReqs := s.unprocessableMap()
	n := newReqs.Len()
	for i := uint32(0); i < n; i++ {
		k := newReqs.GetAt(i)
		rec := mustUnprocessableRequestRecordFromBytes(allReqs.GetAt(k))
		rec.outputID = iotago.OutputIDFromTransactionIDAndIndex(anchorTxID, rec.outputID.Index())
		allReqs.SetAt(k, rec.Bytes())
	}
	newReqs.Erase()
}

func (s *StateReader) GetUnprocessable(reqID isc.RequestID) (req isc.Request, outputID iotago.OutputID, err error) {
	recData := s.unprocessableMap().GetAt(reqID.Bytes())
	rec, err := unprocessableRequestRecordFromBytes(recData)
	if err != nil {
		return nil, iotago.OutputID{}, err
	}
	return rec.req, rec.outputID, nil
}

func (s *StateReader) HasUnprocessable(reqID isc.RequestID) bool {
	return s.unprocessableMap().HasAt(reqID.Bytes())
}

func (s *StateWriter) RemoveUnprocessable(reqID isc.RequestID) {
	s.unprocessableMap().DelAt(reqID.Bytes())
}

// ---- entrypoints

// view used to check if a given requestID exists on the unprocessable list
func viewHasUnprocessable(ctx isc.SandboxView, reqID isc.RequestID) bool {
	state := NewStateReaderFromSandbox(ctx)
	return state.HasUnprocessable(reqID)
}

var (
	ErrUnprocessableAlreadyExist = coreerrors.Register("request does not exist on the unprocessable list").Create()
	ErrUnprocessableUnexpected   = coreerrors.Register("unexpected error getting unprocessable request from the state").Create()
	ErrUnprocessableWrongSender  = coreerrors.Register("unprocessable request sender does not match the retry sender").Create()
)

func retryUnprocessable(ctx isc.Sandbox, reqID isc.RequestID) dict.Dict {
	state := NewStateReaderFromSandbox(ctx)
	exists := state.HasUnprocessable(reqID)
	if !exists {
		panic(ErrUnprocessableAlreadyExist)
	}
	rec, outputID, err := state.GetUnprocessable(reqID)
	if err != nil {
		panic(ErrUnprocessableUnexpected)
	}
	recSender := rec.SenderAccount()
	if rec.SenderAccount() == nil || !recSender.Equals(ctx.Request().SenderAccount()) {
		panic(ErrUnprocessableWrongSender)
	}
	ctx.Privileged().RetryUnprocessable(rec, outputID)
	return nil
}

func UnprocessableRequestsAddedInBlock(block state.Block) ([]isc.Request, error) {
	state := NewStateReaderFromBlockMutations(block)
	var respErr error
	requests := []isc.Request{}
	state.unprocessableMap().Iterate(func(_, recData []byte) bool {
		rec, err := unprocessableRequestRecordFromBytes(recData)
		if err != nil {
			respErr = err
			return false
		}
		requests = append(requests, rec.req)
		return true
	})
	return requests, respErr
}

func HasUnprocessableRequestBeenRemovedInBlock(block state.Block, requestID isc.RequestID) bool {
	keyBytes := Contract.Hname().Bytes()
	keyBytes = append(keyBytes, collections.MapElemKey(prefixUnprocessableRequests, requestID.Bytes())...)
	_, wasRemoved := block.Mutations().Dels[kv.Key(keyBytes)]
	return wasRemoved
}
