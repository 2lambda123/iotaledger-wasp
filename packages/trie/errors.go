package trie

import "github.com/iotaledger/hive.go/ierrors"

var (
	ErrWrongNibble = ierrors.New("key16 byte must be less than 0x0F")
	ErrEmpty       = ierrors.New("encoded key16 can't be empty")
	ErrWrongFormat = ierrors.New("encoded key16 wrong format")
)
