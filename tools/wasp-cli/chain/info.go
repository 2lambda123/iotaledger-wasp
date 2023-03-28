// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/cliclients"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/iotaledger/wasp/tools/wasp-cli/util"
	"github.com/iotaledger/wasp/tools/wasp-cli/waspcmd"
)

func initInfoCmd() *cobra.Command {
	var node string
	var chain string
	var debug bool
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show information about the chain",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			chainID := config.GetChain(chain)
			client := cliclients.WaspClient(node, debug)

			chainInfo, _, err := client.ChainsApi.
				GetChainInfo(context.Background(), chainID.String()).
				Execute() //nolint:bodyclose // false positive
			log.Check(err)

			committeeInfo, _, err := client.ChainsApi.
				GetCommitteeInfo(context.Background(), chainID.String()).
				Execute() //nolint:bodyclose // false positive
			log.Check(err)

			printNodesRowHdr := []string{"PubKey", "PeeringURL", "Alive", "Committee", "Access", "AccessAPI"}
			printNodesRowFmt := func(n apiclient.CommitteeNode, isCommitteeNode, isAccessNode bool) []string {
				return []string{
					n.Node.PublicKey,
					n.Node.PeeringURL,
					strconv.FormatBool(n.Node.IsAlive),
					strconv.FormatBool(isCommitteeNode),
					strconv.FormatBool(isAccessNode),
					n.AccessAPI,
				}
			}
			printNodes := func(label string, nodes []apiclient.CommitteeNode, isCommitteeNode, isAccessNode bool) {
				if nodes == nil {
					log.Printf("%s: N/A\n", label)
				}
				log.Printf("%s: %v\n", label, len(nodes))
				rows := make([][]string, 0)
				for _, n := range nodes {
					rows = append(rows, printNodesRowFmt(n, isCommitteeNode, isAccessNode))
				}
				log.PrintTable(printNodesRowHdr, rows)
			}

			log.Printf("Chain ID: %s\n", chainInfo.ChainID)
			log.Printf("EVM Chain ID: %d\n", chainInfo.EvmChainId)
			log.Printf("Active: %v\n", chainInfo.IsActive)

			if chainInfo.IsActive {
				log.Printf("State address: %v\n", committeeInfo.StateAddress)
				printNodes("Committee nodes", committeeInfo.CommitteeNodes, true, false)
				printNodes("Access nodes", committeeInfo.AccessNodes, false, true)
				printNodes("Candidate nodes", committeeInfo.CandidateNodes, false, false)

				contracts, _, err := client.ChainsApi.GetContracts(context.Background(), chainID.String()).Execute() //nolint:bodyclose // false positive
				log.Check(err)
				log.Printf("#Contracts: %d\n", len(contracts))

				log.Printf("Owner: %s\n", chainInfo.ChainOwnerId)
				log.Printf("Gas fee: gas units * (%d/%d)\n", chainInfo.GasFeePolicy.GasPerToken.A, chainInfo.GasFeePolicy.GasPerToken.B)
				log.Printf("Validator fee share: %d%%\n", chainInfo.GasFeePolicy.ValidatorFeeShare)

				if chainInfo.CustomMetadata != nil && len(*chainInfo.CustomMetadata) > 0 {
					customMetadata, err := base64.StdEncoding.DecodeString(*chainInfo.CustomMetadata)
					log.Check(err)
					log.Printf("Custom metadata (hex): %s\n", hex.EncodeToString(customMetadata))
				} else {
					log.Printf("Custom metadata: (empty)\n")
				}
			}
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	util.WithDebugFlag(cmd, &debug)
	return cmd
}
