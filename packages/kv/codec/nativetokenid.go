package codec

import (
	"errors"

	iotago "github.com/iotaledger/iota.go/v4"
)

func DecodeNativeTokenID(b []byte, def ...iotago.NativeTokenID) (iotago.NativeTokenID, error) {
	if len(b) != iotago.NativeTokenIDLength {
		if len(def) == 0 {
			return iotago.NativeTokenID{}, errors.New("wrong data length")
		}
		return def[0], nil
	}
	var ret iotago.NativeTokenID
	copy(ret[:], b)
	return ret, nil
}

func EncodeNativeTokenID(value iotago.NativeTokenID) []byte {
	return value[:]
}
