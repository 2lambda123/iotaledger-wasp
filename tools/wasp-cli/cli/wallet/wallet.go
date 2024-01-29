package wallet

import (
	"github.com/spf13/viper"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

var AddressIndex int

type Wallet struct {
	KeyPair      *cryptolib.KeyPair
	AddressIndex int
}

func Load() *Wallet {
	seedHex := viper.GetString("wallet.seed")
	useLegacyDerivation := viper.GetBool("wallet.useLegacyDerivation")
	useCoinType := viper.IsSet("wallet.coinType")
	coinType := viper.GetUint32("wallet.coinType")
	if seedHex == "" {
		log.Fatal("call `init` first")
	}

	masterSeed, err := iotago.DecodeHex(seedHex)
	log.Check(err)

	subSeed := cryptolib.SubSeed(masterSeed, uint32(AddressIndex), useLegacyDerivation)
	if useCoinType {
		subSeed = cryptolib.SubSeed(masterSeed, uint32(AddressIndex), useLegacyDerivation, coinType)
	}

	kp := cryptolib.KeyPairFromSeed(subSeed)

	return &Wallet{KeyPair: kp, AddressIndex: AddressIndex}
}

func (w *Wallet) PrivateKey() *cryptolib.PrivateKey {
	return w.KeyPair.GetPrivateKey()
}

func (w *Wallet) Address() iotago.Address {
	return w.KeyPair.GetPublicKey().AsEd25519Address()
}
