package trie

import (
	"fmt"

	"golang.org/x/crypto/blake2b"

	"github.com/iotaledger/hive.go/serializer/v2/marshalutil"
)

func concat(par ...[]byte) []byte {
	mu := new(marshalutil.MarshalUtil)
	for _, p := range par {
		mu.WriteBytes(p)
	}
	return mu.Bytes()
}

func blake2b160(data []byte) (ret [HashSizeBytes]byte) {
	hash, _ := blake2b.New(HashSizeBytes, nil)
	if _, err := hash.Write(data); err != nil {
		panic(err)
	}
	copy(ret[:], hash.Sum(nil))
	return
}

func assertf(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf("assertion failed:: "+format, args...))
	}
}

func assertNoError(err error) {
	assertf(err == nil, "error: %v", err)
}
