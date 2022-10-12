package wallet

import (
	"github.com/spf13/cobra"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/tools/wasp-cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
)

var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "Show the wallet address",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		wallet := Load()
		log.Printf("Address index %d\n", addressIndex)
		log.Verbosef("  Private key: %s\n", wallet.KeyPair.GetPrivateKey().String())
		log.Verbosef("  Public key:  %s\n", wallet.KeyPair.GetPublicKey().String())
		log.Printf("  Address:     %s\n", wallet.Address().Bech32(parameters.L1().Protocol.Bech32HRP))
	},
}

type BalanceModel struct {
	AddressIndex int                 `json:"AddressIndex"`
	Address      string              `json:"Address"`
	BaseTokens   uint64              `json:"BaseTokens"`
	NativeTokens iotago.NativeTokens `json:"NativeTokens"`

	OutputMap      iotago.OutputSet `json:"-"`
	VerboseOutputs map[uint16]string
}

func (b *BalanceModel) AsJSON() ([]byte, error) {
	return log.DefaultJSONFormatter(b)
}

func (b *BalanceModel) AsText() (string, error) {
	balanceTemplate := `Address index: {{.AddressIndex}}
Address: {{.Address}}
NativeTokens: 
	Base tokens: {{.BaseTokens}}

{{range $i, $out := .NativeTokens}}
	{{$i.ID}} {{$out.Amount}}
{{end}}`

	return log.ParseCLIOutputTemplate(b, balanceTemplate)
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Show the wallet balance",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		wallet := Load()
		address := wallet.Address()

		outs, err := config.L1Client().OutputMap(address)
		log.Check(err)

		balance := isc.FungibleTokensFromOutputMap(outs)

		model := &BalanceModel{
			Address:      address.Bech32(parameters.L1().Protocol.Bech32HRP),
			AddressIndex: addressIndex,
			NativeTokens: balance.Tokens,
			BaseTokens:   balance.BaseTokens,
			OutputMap:    outs,
		}
		if log.VerboseFlag {
			model.VerboseOutputs = map[uint16]string{}

			for i, out := range outs {
				tokens := isc.FungibleTokensFromOutput(out)
				model.VerboseOutputs[i.Index()] = tokens.String()
			}
		}

		log.PrintCLIOutput(model)
	},
}
