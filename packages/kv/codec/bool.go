package codec

import (
	"bytes"
	"errors"

	"github.com/iotaledger/wasp/packages/util"
)

func DecodeBool(b []byte, def ...bool) (bool, error) {
	if b == nil {
		if len(def) == 0 {
			return false, errors.New("cannot decode nil bytes")
		}
		return def[0], nil
	}
	return util.ReadBool(bytes.NewReader(b))
}

func MustDecodeBool(b []byte, def ...bool) bool {
	ret, err := DecodeBool(b, def...)
	if err != nil {
		panic(err)
	}
	return ret
}

func EncodeBool(value bool) []byte {
	w := new(bytes.Buffer)
	_ = util.WriteBool(w, value)
	return w.Bytes()
}
