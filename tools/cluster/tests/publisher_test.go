package tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"

	// without this import it won't work, no messages will be received by the client socket...
	_ "go.nanomsg.org/mangos/v3/transport/all"

	"github.com/iotaledger/wasp/client/chainclient"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
)

type nanoClientTest struct {
	id       int
	messages []string
}

func (c *nanoClientTest) start(t *testing.T, url string) {
	sock, err := sub.NewSocket()
	require.NoError(t, err)

	err = sock.Dial(url)
	require.NoError(t, err)
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	require.NoError(t, err)

	for {
		msg, err := sock.Recv()
		require.NoError(t, err)
		c.messages = append(c.messages, string(msg))
	}
}

func assertMessages(t *testing.T, messages []string, expectedFinalCounter int) {
	t.Logf("assertMessages: |messages|=%v, expectedFinalCounter=%v", len(messages), expectedFinalCounter)
	inccounterEventRegx := regexp.MustCompile(`.*incCounter: counter = (\d+)$`)
	counter := 1
	for _, msg := range messages {
		rs := inccounterEventRegx.FindStringSubmatch(msg)
		if len(rs) == 0 {
			continue
		}
		counterFromEvent, err := strconv.ParseInt(rs[1], 10, 64)
		require.NoError(t, err)
		require.EqualValues(t, counter, counterFromEvent)
		counter++
	}
	require.EqualValues(t, expectedFinalCounter, counter-1)
}

// TODO the TODOs on this test indicate that there is a race condition with the "await request endpoint", needs to be debugged
// executed in cluster_test.go
func testNanoPublisher(t *testing.T, e *chainEnv) {
	e.deployWasmInccounter(0)
	// spawn many NANOMSG nanoClients and subscribe to everything from the node
	nanoClients := make([]nanoClientTest, 10)
	nanoURL := fmt.Sprintf("tcp://127.0.0.1:%d", e.Clu.Config.NanomsgPort(0))
	for i := range nanoClients {
		nanoClients[i] = nanoClientTest{id: i, messages: []string{}}
		go nanoClients[i].start(t, nanoURL)
	}

	// deposit funds for offledger requests
	keyPair, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	accountsClient := e.Chain.SCClient(accounts.Contract.Hname(), keyPair)
	reqTx, err := accountsClient.PostRequest(accounts.FuncDeposit.Name, chainclient.PostRequestParams{
		Transfer: isc.NewFungibleBaseTokens(1_000_000),
	})
	require.NoError(t, err)

	_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, reqTx, 30*time.Second)
	require.NoError(t, err)

	// send 100 requests
	numRequests := 100
	myClient := e.Chain.SCClient(incHname, keyPair)

	reqIDs := make([]isc.RequestID, numRequests)
	for i := 0; i < numRequests; i++ {
		req, err := myClient.PostOffLedgerRequest(incrementFuncName, chainclient.PostRequestParams{Nonce: uint64(i + 1)})
		require.NoError(t, err)

		// ---
		// TODO shouldn't be needed
		_, err = e.Chain.CommitteeMultiClient().WaitUntilRequestProcessedSuccessfully(e.Chain.ChainID, req.ID(), 30*time.Second)
		require.NoError(t, err)
		// ---

		reqIDs[i] = req.ID()
	}

	waitUntil(t, e.counterEquals(int64(numRequests)), util.MakeRange(0, 1), 60*time.Second, "requests counted - A")

	// assert all clients received the correct number of messages, in the correct order.
	for _, client := range nanoClients {
		assertMessages(t, client.messages, numRequests)
	}

	for i := 0; i < numRequests; i++ {
		req, err := myClient.PostOffLedgerRequest(incrementFuncName, chainclient.PostRequestParams{Nonce: uint64(i + 101)})
		require.NoError(t, err)

		// ---
		// TODO shouldn't be needed
		_, err = e.Chain.CommitteeMultiClient().WaitUntilRequestProcessedSuccessfully(e.Chain.ChainID, req.ID(), 30*time.Second)
		require.NoError(t, err)
		// ---
	}

	waitUntil(t, e.counterEquals(int64(numRequests*2)), util.MakeRange(0, 1), 60*time.Second, "requests counted - B")

	// assert all clients received the correct number of messages, in the correct order.
	for _, client := range nanoClients {
		assertMessages(t, client.messages, numRequests*2)
	}
}
