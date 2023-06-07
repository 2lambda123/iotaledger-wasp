package blocklog

import (
	"fmt"
	"io"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/subrealm"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util/rwutil"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

// region RequestReceipt /////////////////////////////////////////////////////

// RequestReceipt represents log record of processed request on the chain
type RequestReceipt struct {
	Request       isc.Request            `json:"request"`
	Error         *isc.UnresolvedVMError `json:"error"`
	GasBudget     uint64                 `json:"gasBudget"`
	GasBurned     uint64                 `json:"gasBurned"`
	GasFeeCharged uint64                 `json:"gasFeeCharged"`
	SDCharged     uint64                 `json:"storageDepositCharged"`
	// not persistent
	BlockIndex   uint32       `json:"blockIndex"`
	RequestIndex uint16       `json:"requestIndex"`
	GasBurnLog   *gas.BurnLog `json:"-"`
}

func RequestReceiptFromBytes(data []byte) (*RequestReceipt, error) {
	rr := rwutil.NewBytesReader(data)
	ret := &RequestReceipt{
		GasBudget:     rr.ReadUint64(),
		GasBurned:     rr.ReadUint64(),
		GasFeeCharged: rr.ReadUint64(),
		SDCharged:     rr.ReadUint64(),
		Request:       rwutil.FromMarshalUtil(rr, isc.NewRequestFromMarshalUtil),
	}
	hasError := rr.ReadBool()
	if hasError {
		ret.Error = rwutil.FromMarshalUtil(rr, isc.UnresolvedVMErrorFromMarshalUtil)
	}
	return ret, rr.Err
}

func RequestReceiptsFromBlock(block state.Block) ([]*RequestReceipt, error) {
	var respErr error
	receipts := []*RequestReceipt{}
	kvStore := subrealm.NewReadOnly(block.MutationsReader(), kv.Key(Contract.Hname().Bytes()))
	kvStore.Iterate(kv.Key(prefixRequestReceipts+"."), func(key kv.Key, value []byte) bool { // TODO: Nicer way to construct the key?
		receipt, err := RequestReceiptFromBytes(value)
		if err != nil {
			respErr = fmt.Errorf("cannot deserialize requestReceipt: %w", err)
			return true
		}
		receipts = append(receipts, receipt)
		return true
	})
	if respErr != nil {
		return nil, respErr
	}
	return receipts, nil
}

func (r *RequestReceipt) Bytes() []byte {
	ww := rwutil.NewBytesWriter()

	ww.WriteUint64(r.GasBudget).
		WriteUint64(r.GasBurned).
		WriteUint64(r.GasFeeCharged).
		WriteUint64(r.SDCharged).
		WriteToMarshalUtil(r.Request)

	hasError := r.Error != nil
	ww.WriteBool(hasError)
	if hasError {
		ww.WriteN(r.Error.Bytes())
	}

	return ww.Bytes()
}

func (r *RequestReceipt) WithBlockData(blockIndex uint32, requestIndex uint16) *RequestReceipt {
	r.BlockIndex = blockIndex
	r.RequestIndex = requestIndex
	return r
}

func (r *RequestReceipt) String() string {
	ret := fmt.Sprintf("ID: %s\n", r.Request.ID().String())
	ret += fmt.Sprintf("Err: %v\n", r.Error)
	ret += fmt.Sprintf("Block/Request index: %d / %d\n", r.BlockIndex, r.RequestIndex)
	ret += fmt.Sprintf("Gas budget / burned / fee charged: %d / %d /%d\n", r.GasBudget, r.GasBurned, r.GasFeeCharged)
	ret += fmt.Sprintf("Storage deposit charged: %d\n", r.SDCharged)
	ret += fmt.Sprintf("Call data: %s\n", r.Request)
	return ret
}

func (r *RequestReceipt) Short() string {
	prefix := "tx"
	if r.Request.IsOffLedger() {
		prefix = "api"
	}

	ret := fmt.Sprintf("%s/%s", prefix, r.Request.ID())

	if r.Error != nil {
		ret += fmt.Sprintf(": Err: %v", r.Error)
	}

	return ret
}

func (r *RequestReceipt) LookupKey() RequestLookupKey {
	return NewRequestLookupKey(r.BlockIndex, r.RequestIndex)
}

func (r *RequestReceipt) ToISCReceipt(resolvedError *isc.VMError) *isc.Receipt {
	return &isc.Receipt{
		Request:       r.Request.Bytes(),
		Error:         r.Error,
		GasBudget:     r.GasBudget,
		GasBurned:     r.GasBurned,
		GasFeeCharged: r.GasFeeCharged,
		BlockIndex:    r.BlockIndex,
		RequestIndex:  r.RequestIndex,
		ResolvedError: resolvedError.Error(),
	}
}

// endregion  /////////////////////////////////////////////////////////////

// region RequestLookupKey /////////////////////////////////////////////

// RequestLookupReference globally unique reference to the request: block index and index of the request within block
type RequestLookupKey [6]byte

func NewRequestLookupKey(blockIndex uint32, requestIndex uint16) RequestLookupKey {
	ret := RequestLookupKey{}
	copy(ret[:4], codec.EncodeUint32(blockIndex))
	copy(ret[4:6], codec.EncodeUint16(requestIndex))
	return ret
}

func (k RequestLookupKey) BlockIndex() uint32 {
	return codec.MustDecodeUint32(k[:4])
}

func (k RequestLookupKey) RequestIndex() uint16 {
	return codec.MustDecodeUint16(k[4:6])
}

func (k RequestLookupKey) Bytes() []byte {
	return k[:]
}

func (k *RequestLookupKey) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadN(k[:])
	return rr.Err
}

func (k *RequestLookupKey) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteN(k[:])
	return ww.Err
}

// endregion ///////////////////////////////////////////////////////////

// region RequestLookupKeyList //////////////////////////////////////////////

// RequestLookupKeyList a list of RequestLookupReference of requests with colliding isc.RequestLookupDigest
type RequestLookupKeyList []RequestLookupKey

func RequestLookupKeyListFromBytes(data []byte) (RequestLookupKeyList, error) {
	rr := rwutil.NewBytesReader(data)
	size := rr.ReadSize()
	ll := make(RequestLookupKeyList, size)
	for i := 0; i < size; i++ {
		rr.ReadN(ll[i][:])
	}
	return ll, rr.Err
}

func (ll RequestLookupKeyList) Bytes() []byte {
	ww := rwutil.NewBytesWriter()
	ww.WriteSize(len(ll))
	for i := range ll {
		ww.WriteN(ll[i][:])
	}
	return ww.Bytes()
}

// endregion /////////////////////////////////////////////////////////////
