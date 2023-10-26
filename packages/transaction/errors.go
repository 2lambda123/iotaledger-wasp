package transaction

import "github.com/iotaledger/hive.go/ierrors"

var (
	ErrNotEnoughBaseTokens                  = ierrors.New("not enough base tokens")
	ErrNotEnoughBaseTokensForStorageDeposit = ierrors.New("not enough base tokens for storage deposit")
	ErrNotEnoughNativeTokens                = ierrors.New("not enough native tokens")
)
