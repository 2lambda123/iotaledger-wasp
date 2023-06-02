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
	var ret bool
	err := util.ReadBool(bytes.NewReader(b), &ret)
	return ret, err
}

func MustDecodeBool(b []byte, def ...bool) bool {
	ret, err := DecodeBool(b, def...)
	if err != nil {
		panic(err)
	}
	return ret
}

func EncodeBool(value bool) []byte {
	buf := new(bytes.Buffer)
	_ = util.WriteBool(buf, value)
	return buf.Bytes()
}
