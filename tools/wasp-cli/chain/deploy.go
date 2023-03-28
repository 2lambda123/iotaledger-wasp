// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/apilib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/cliclients"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/wallet"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	cliutil "github.com/iotaledger/wasp/tools/wasp-cli/util"
	"github.com/iotaledger/wasp/tools/wasp-cli/waspcmd"
)

func GetAllWaspNodes() []int {
	ret := []int{}
	for index := range viper.GetStringMap("wasp") {
		i, err := strconv.Atoi(index)
		log.Check(err)
		ret = append(ret, i)
	}
	return ret
}

func controllerAddrDefaultFallback(addr string) iotago.Address {
	if addr == "" {
		return wallet.Load().Address()
	}
	prefix, govControllerAddr, err := iotago.ParseBech32(addr)
	log.Check(err)
	if parameters.L1().Protocol.Bech32HRP != prefix {
		log.Fatalf("unexpected prefix. expected: %s, actual: %s", parameters.L1().Protocol.Bech32HRP, prefix)
	}
	return govControllerAddr
}

func initDeployCmd() *cobra.Command {
	var (
		node             string
		peers            []string
		quorum           int
		evmParams        evmDeployParams
		govControllerStr string
		chainName        string
		debug bool
	)

	cmd := &cobra.Command{
		Use:   "deploy --chain=<name>",
		Short: "Deploy a new chain",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chainName = defaultChainFallback(chainName)

			if !util.IsSlug(chainName) {
				log.Fatalf("invalid chain name: %s, must be in slug format, only lowercase and hypens, example: foo-bar", chainName)
			}

			l1Client := cliclients.L1Client()

			govController := controllerAddrDefaultFallback(govControllerStr)

			stateController := doDKG(node, peers, quorum, debug)

			par := apilib.CreateChainParams{
				Layer1Client:         l1Client,
				CommitteeAPIHosts:    config.NodeAPIURLs([]string{node}),
				N:                    uint16(len(node)),
				T:                    uint16(quorum),
				OriginatorKeyPair:    wallet.Load().KeyPair,
				Textout:              os.Stdout,
				GovernanceController: govController,
				InitParams: dict.Dict{
					origin.ParamChainOwner:   isc.NewAgentID(govController).Bytes(),
					origin.ParamEVMChainID:   codec.EncodeUint16(evmParams.ChainID),
					origin.ParamEVMBlockKeep: codec.EncodeInt32(evmParams.BlockKeepAmount),
				},
			}

			chainID, err := apilib.DeployChain(par, stateController, govController)
			log.Check(err)

			config.AddChain(chainName, chainID.String())

			activateChain(node, chainName, chainID, debug)
		},
	}

	waspcmd.WithWaspNodeFlag(cmd, &node)
	waspcmd.WithPeersFlag(cmd, &peers)
	evmParams.initFlags(cmd)
	cmd.Flags().StringVar(&chainName, "chain", "", "name of the chain)")
	log.Check(cmd.MarkFlagRequired("chain"))
	cmd.Flags().IntVar(&quorum, "quorum", 0, "quorum (default: 3/4s of the number of committee nodes)")
	cmd.Flags().StringVar(&govControllerStr, "gov-controller", "", "governance controller address")
	cliutil.WithDebugFlag(cmd, &debug)
	return cmd
}
