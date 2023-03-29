package chain

import (
	"context"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/clients/apiextensions"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/cliclients"
	"github.com/iotaledger/wasp/tools/wasp-cli/cli/config"
	"github.com/iotaledger/wasp/tools/wasp-cli/log"
	"github.com/iotaledger/wasp/tools/wasp-cli/waspcmd"
)

func initBlockCmd() *cobra.Command {
	var node string
	var chain string
	cmd := &cobra.Command{
		Use:   "block [index]",
		Short: "Get information about a block given its index, or latest block if missing",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			bi := fetchBlockInfo(args, node, chain)
			log.Printf("Block index: %d\n", bi.BlockIndex)
			log.Printf("Timestamp: %s\n", bi.Timestamp.UTC().Format(time.RFC3339))
			log.Printf("Total requests: %d\n", bi.TotalRequests)
			log.Printf("Successful requests: %d\n", bi.NumSuccessfulRequests)
			log.Printf("Off-ledger requests: %d\n", bi.NumOffLedgerRequests)
			log.Printf("\n")
			logRequestsInBlock(bi.BlockIndex, node, chain)
			log.Printf("\n")
			logEventsInBlock(bi.BlockIndex, node, chain)
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func fetchBlockInfo(args []string, node, chain string) *apiclient.BlockInfoResponse {
	client := cliclients.WaspClient(node)

	if len(args) == 0 {
		blockInfo, _, err := client.
			CorecontractsApi.
			BlocklogGetLatestBlockInfo(context.Background(), config.GetChain(chain).String()).
			Execute() //nolint:bodyclose // false positive

		log.Check(err)
		return blockInfo
	}

	index, err := strconv.Atoi(args[0])
	log.Check(err)

	blockInfo, _, err := client.
		CorecontractsApi.
		BlocklogGetBlockInfo(context.Background(), config.GetChain(chain).String(), uint32(index)).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)

	return blockInfo
}

func logRequestsInBlock(index uint32, node, chain string) {
	client := cliclients.WaspClient(node)
	receipts, _, err := client.CorecontractsApi.
		BlocklogGetRequestReceiptsOfBlock(context.Background(), config.GetChain(chain).String(), index).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)

	for i, receipt := range receipts.Receipts {
		r := receipt
		logReceipt(&r, i)
	}
}

func logReceipt(receipt *apiclient.RequestReceiptResponse, index ...int) {
	req := receipt.Request

	kind := "on-ledger"
	if req.IsOffLedger {
		kind = "off-ledger"
	}

	args, err := apiextensions.APIJsonDictToDict(req.Params)
	log.Check(err)

	var argsTree interface{} = "(empty)"
	if len(args) > 0 {
		argsTree = args
	}

	errMsg := "(empty)"
	if receipt.Error != nil {
		errMsg = receipt.Error.ErrorMessage
	}

	tree := []log.TreeItem{
		{K: "Kind", V: kind},
		{K: "Sender", V: req.SenderAccount},
		{K: "Contract Hname", V: req.CallTarget.ContractHName},
		{K: "Function Hname", V: req.CallTarget.FunctionHName},
		{K: "Arguments", V: argsTree},
		{K: "Error", V: errMsg},
	}
	if len(index) > 0 {
		log.Printf("Request #%d (%s):\n", index[0], req.RequestId)
	} else {
		log.Printf("Request %s:\n", req.RequestId)
	}
	log.PrintTree(tree, 2, 2)
}

func logEventsInBlock(index uint32, node, chain string) {
	client := cliclients.WaspClient(node)
	events, _, err := client.CorecontractsApi.
		BlocklogGetEventsOfBlock(context.Background(), config.GetChain(chain).String(), index).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)
	logEvents(events)
}

func initRequestCmd() *cobra.Command {
	var node string
	var chain string
	cmd := &cobra.Command{
		Use:   "request <request-id>",
		Short: "Get information about a request given its ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chain = defaultChainFallback(chain)

			reqID, err := isc.RequestIDFromString(args[0])
			log.Check(err)

			client := cliclients.WaspClient(node)
			receipt, _, err := client.RequestsApi.
				GetReceipt(context.Background(), config.GetChain(chain).String(), reqID.String()).
				Execute() //nolint:bodyclose // false positive

			log.Check(err)

			log.Printf("Request found in block %d\n\n", receipt.BlockIndex)
			logResolvedReceipt(receipt)

			log.Printf("\n")
			logEventsInRequest(reqID, node, chain)
			log.Printf("\n")
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	withChainFlag(cmd, &chain)
	return cmd
}

func logResolvedReceipt(receipt *apiclient.ReceiptResponse, index ...int) {
	reqBytes, err := iotago.DecodeHex(receipt.Request)
	log.Check(err)
	req, err := isc.NewRequestFromBytes(reqBytes)
	log.Check(err)

	kind := "on-ledger"
	if req.IsOffLedger() {
		kind = "off-ledger"
	}

	var argsTree interface{} = "(empty)"
	if len(req.Params()) > 0 {
		argsTree = req.Params()
	}

	errMsg := "(empty)"
	if receipt.Error != nil {
		errMsg = receipt.Error.Message
	}

	tree := []log.TreeItem{
		{K: "Kind", V: kind},
		{K: "Sender", V: req.SenderAccount().String()},
		{K: "Contract Hname", V: req.CallTarget().Contract.String()},
		{K: "Function Hname", V: req.CallTarget().EntryPoint.String()},
		{K: "Arguments", V: argsTree},
		{K: "Error", V: errMsg},
	}
	if len(index) > 0 {
		log.Printf("Request #%d (%s):\n", index[0], req.ID())
	} else {
		log.Printf("Request %s:\n", req.ID())
	}
	log.PrintTree(tree, 2, 2)
}

func logEventsInRequest(reqID isc.RequestID, node, chain string) {
	client := cliclients.WaspClient(node)
	events, _, err := client.CorecontractsApi.
		BlocklogGetEventsOfRequest(context.Background(), config.GetChain(chain).String(), reqID.String()).
		Execute() //nolint:bodyclose // false positive

	log.Check(err)
	logEvents(events)
}

func logEvents(ret *apiclient.EventsResponse) {
	header := []string{"event"}
	rows := make([][]string, len(ret.Events))

	for i, event := range ret.Events {
		rows[i] = []string{event}
	}

	log.Printf("Total %d events\n", len(ret.Events))
	log.PrintTable(header, rows)
}
