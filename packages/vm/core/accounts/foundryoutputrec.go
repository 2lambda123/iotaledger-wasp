package accounts

import (
	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// foundryOutputRec contains information to reconstruct output
type foundryOutputRec struct {
	Amount      uint64 // always storage deposit
	TokenScheme iotago.TokenScheme
	Metadata    []byte
	BlockIndex  uint32
	OutputIndex uint16
}

func (f *foundryOutputRec) Bytes() []byte {
	mu := marshalutil.New()
	mu.WriteUint32(f.BlockIndex).
		WriteUint16(f.OutputIndex).
		WriteUint64(f.Amount)
	rwutil.WriteBytesToMarshalUtil(codec.EncodeTokenScheme(f.TokenScheme), mu)
	rwutil.WriteBytesToMarshalUtil(f.Metadata, mu)

	return mu.Bytes()
}

func foundryOutputRecFromMarshalUtil(mu *marshalutil.MarshalUtil) (*foundryOutputRec, error) {
	ret := &foundryOutputRec{}
	var err error
	if ret.BlockIndex, err = mu.ReadUint32(); err != nil {
		return nil, err
	}
	if ret.OutputIndex, err = mu.ReadUint16(); err != nil {
		return nil, err
	}
	if ret.Amount, err = mu.ReadUint64(); err != nil {
		return nil, err
	}
	schemeBin, err := rwutil.ReadBytesFromMarshalUtil(mu)
	if err != nil {
		return nil, err
	}
	if ret.TokenScheme, err = codec.DecodeTokenScheme(schemeBin); err != nil {
		return nil, err
	}
	if ret.Metadata, err = rwutil.ReadBytesFromMarshalUtil(mu); err != nil {
		return nil, err
	}
	return ret, nil
}

func mustFoundryOutputRecFromBytes(data []byte) *foundryOutputRec {
	ret, err := foundryOutputRecFromMarshalUtil(marshalutil.New(data))
	if err != nil {
		panic(err)
	}
	return ret
}
