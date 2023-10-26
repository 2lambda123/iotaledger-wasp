package codec

import "github.com/iotaledger/hive.go/ierrors"

func DecodeBool(b []byte, def ...bool) (bool, error) {
	if b == nil {
		if len(def) == 0 {
			return false, ierrors.New("cannot decode nil bool")
		}
		return def[0], nil
	}
	if len(b) != 1 {
		return false, ierrors.New("invalid bool size")
	}
	if (b[0] & 0xfe) != 0x00 {
		return false, ierrors.New("invalid bool value")
	}
	return b[0] != 0, nil
}

func MustDecodeBool(b []byte, def ...bool) bool {
	ret, err := DecodeBool(b, def...)
	if err != nil {
		panic(err)
	}
	return ret
}

func EncodeBool(value bool) []byte {
	if value {
		return []byte{1}
	}
	return []byte{0}
}
