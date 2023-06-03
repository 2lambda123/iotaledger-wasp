package codec

import (
	"errors"
	"time"

	"github.com/iotaledger/wasp/packages/util"
)

func DecodeTime(b []byte, def ...time.Time) (time.Time, error) {
	if b == nil {
		if len(def) == 0 {
			return time.Time{}, errors.New("cannot decode nil bytes")
		}
		return def[0], nil
	}
	nanos, err := util.Int64FromBytes(b)
	if err != nil {
		return time.Time{}, err
	}
	if nanos == 0 {
		return time.Time{}, nil
	}
	return time.Unix(0, nanos), nil
}

func MustDecodeTime(b []byte, def ...time.Time) time.Time {
	t, err := DecodeTime(b, def...)
	if err != nil {
		panic(err)
	}
	return t
}

func EncodeTime(value time.Time) []byte {
	if value.IsZero() {
		return make([]byte, 8)
	}
	return util.Int64ToBytes(value.UnixNano())
}
