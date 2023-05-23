package blob

import (
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

var Processor = Contract.Processor(nil,
	FuncStoreBlob.WithHandler(storeBlob),
	ViewGetBlobField.WithHandler(getBlobField),
	ViewGetBlobInfo.WithHandler(getBlobInfo),
	ViewListBlobs.WithHandler(listBlobs),
)

func SetInitialState(state kv.KVStore) {
	// does not do anything
}

// storeBlob treats parameters as names of fields and field values
// it stores it in the state in deterministic binary representation
// Returns hash of the blob
func storeBlob(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf("blob.storeBlob.begin")
	state := ctx.State()
	params := ctx.Params()
	// calculate a deterministic hash of all blob fields
	blobHash, kSorted, values := mustGetBlobHash(params.Dict)

	directory := GetDirectory(state)
	ctx.Requiref(!directory.HasAt(blobHash[:]),
		"blob.storeBlob.fail: blob with hash %s already exists", blobHash.String())

	// get a record by blob hash
	blbValues := GetBlobValues(state, blobHash)
	blbSizes := GetBlobSizes(state, blobHash)

	totalSize := uint32(0)
	totalSizeWithKeys := uint32(0)

	// save record of the blob. In parallel save record of sizes of blob fields
	sizes := make([]uint32, len(kSorted))
	for i, k := range kSorted {
		size := uint32(len(values[i]))
		blbValues.SetAt([]byte(k), values[i])
		blbSizes.SetAt([]byte(k), EncodeSize(size))
		sizes[i] = size
		totalSize += size
		totalSizeWithKeys += size + uint32(len(k))
	}

	ret := dict.New()
	ret.Set(ParamHash, codec.EncodeHashValue(blobHash))

	directory.SetAt(blobHash[:], EncodeSize(totalSize))

	evt := &StoreBlobEvent{
		BlobHash:   blobHash,
		FieldSizes: sizes,
	}
	ctx.Event(evt.Encode())
	return ret
}

// getBlobInfo return lengths of all fields in the blob
func getBlobInfo(ctx isc.SandboxView) dict.Dict {
	ctx.Log().Debugf("blob.getBlobInfo.begin")

	blobHash := ctx.Params().MustGetHashValue(ParamHash)

	blbSizes := GetBlobSizesR(ctx.StateR(), blobHash)
	ret := dict.New()
	blbSizes.Iterate(func(field []byte, value []byte) bool {
		ret.Set(kv.Key(field), value)
		return true
	})
	return ret
}

func getBlobField(ctx isc.SandboxView) dict.Dict {
	ctx.Log().Debugf("blob.getBlobField.begin")
	state := ctx.StateR()

	params := ctx.Params()
	blobHash := params.MustGetHashValue(ParamHash)
	field := params.MustGetBytes(ParamField)

	blobValues := GetBlobValuesR(state, blobHash)
	ctx.Requiref(blobValues.Len() != 0, "blob with hash %s has not been found", blobHash.String())
	value := blobValues.GetAt(field)
	ctx.Requiref(value != nil, "'blob field %s value not found", string(field))
	ret := dict.New()
	ret.Set(ParamBytes, value)
	return ret
}

func listBlobs(ctx isc.SandboxView) dict.Dict {
	ctx.Log().Debugf("blob.listBlobs.begin")
	ret := dict.New()
	GetDirectoryR(ctx.StateR()).Iterate(func(hash []byte, totalSize []byte) bool {
		ret.Set(kv.Key(hash), totalSize)
		return true
	})
	return ret
}
