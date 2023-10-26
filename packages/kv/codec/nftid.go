package codec

import (
	"github.com/iotaledger/hive.go/ierrors"
	iotago "github.com/iotaledger/iota.go/v3"
)

func DecodeNFTID(b []byte, def ...iotago.NFTID) (ret iotago.NFTID, err error) {
	if b == nil {
		if len(def) == 0 {
			return ret, ierrors.New("cannot decode nil NFTID")
		}
		return def[0], nil
	}
	if len(b) != iotago.NFTIDLength {
		return ret, ierrors.New("invalid NFTID size")
	}
	copy(ret[:], b)
	return ret, nil
}

func MustDecodeNFTID(b []byte) iotago.NFTID {
	r, err := DecodeNFTID(b)
	if err != nil {
		panic(err)
	}
	return r
}

func EncodeNFTID(nftID iotago.NFTID) []byte {
	return nftID[:]
}
