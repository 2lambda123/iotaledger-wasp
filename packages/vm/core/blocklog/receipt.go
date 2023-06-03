package blocklog

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/subrealm"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/util"
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
	return RequestReceiptFromMarshalUtil(marshalutil.New(data))
}

func RequestReceiptFromMarshalUtil(mu *marshalutil.MarshalUtil) (*RequestReceipt, error) {
	ret := &RequestReceipt{}

	var err error

	if ret.GasBudget, err = mu.ReadUint64(); err != nil {
		return nil, fmt.Errorf("cannot read GasBudget: %w", err)
	}
	if ret.GasBurned, err = mu.ReadUint64(); err != nil {
		return nil, fmt.Errorf("cannot read GasBurned: %w", err)
	}
	if ret.GasFeeCharged, err = mu.ReadUint64(); err != nil {
		return nil, fmt.Errorf("cannot read GasFeeCharged: %w", err)
	}
	if ret.SDCharged, err = mu.ReadUint64(); err != nil {
		return nil, fmt.Errorf("cannot read SDCharged: %w", err)
	}
	if ret.Request, err = isc.NewRequestFromMarshalUtil(mu); err != nil {
		return nil, fmt.Errorf("cannot read Request: %w", err)
	}

	if isError, err2 := mu.ReadBool(); err2 != nil {
		return nil, fmt.Errorf("cannot read isError: %w", err2)
	} else if !isError {
		return ret, nil
	}

	if ret.Error, err = isc.UnresolvedVMErrorFromMarshalUtil(mu); err != nil {
		return nil, fmt.Errorf("cannot read Error: %w", err)
	}

	return ret, nil
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
	mu := marshalutil.New()

	mu.WriteUint64(r.GasBudget).
		WriteUint64(r.GasBurned).
		WriteUint64(r.GasFeeCharged).
		WriteUint64(r.SDCharged)

	r.Request.WriteToMarshalUtil(mu)

	if r.Error == nil {
		mu.WriteBool(false)
	} else {
		mu.WriteBool(true)
		mu.WriteBytes(r.Error.Bytes())
	}

	return mu.Bytes()
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
	copy(ret[:4], util.Uint32To4Bytes(blockIndex))
	copy(ret[4:6], util.Uint16To2Bytes(requestIndex))
	return ret
}

func (k RequestLookupKey) BlockIndex() uint32 {
	return util.MustUint32From4Bytes(k[:4])
}

func (k RequestLookupKey) RequestIndex() uint16 {
	return util.MustUint16From2Bytes(k[4:6])
}

func (k RequestLookupKey) Bytes() []byte {
	return k[:]
}

func (k *RequestLookupKey) Write(w io.Writer) error {
	_, err := w.Write(k[:])
	return err
}

func (k *RequestLookupKey) Read(r io.Reader) error {
	n, err := r.Read(k[:])
	if err != nil || n != 6 {
		return io.EOF
	}
	return nil
}

// endregion ///////////////////////////////////////////////////////////

// region RequestLookupKeyList //////////////////////////////////////////////

// RequestLookupKeyList a list of RequestLookupReference of requests with colliding isc.RequestLookupDigest
type RequestLookupKeyList []RequestLookupKey

func RequestLookupKeyListFromBytes(data []byte) (RequestLookupKeyList, error) {
	rdr := bytes.NewReader(data)
	var size uint16
	if err := util.ReadUint16(rdr, &size); err != nil {
		return nil, err
	}
	ret := make(RequestLookupKeyList, size)
	for i := uint16(0); i < size; i++ {
		if err := ret[i].Read(rdr); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (ll RequestLookupKeyList) Bytes() []byte {
	if len(ll) > math.MaxUint16 {
		panic("RequestLookupKeyList::Write: too long")
	}
	mu := new(marshalutil.MarshalUtil)
	mu.WriteUint16(uint16(len(ll)))
	for i := range ll {
		mu.WriteBytes(ll[i][:])
	}
	return mu.Bytes()
}

// endregion /////////////////////////////////////////////////////////////
