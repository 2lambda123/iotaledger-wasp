package tests

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/client/chainclient"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/utxodb"
	"github.com/iotaledger/wasp/packages/vm/core/blocklog"
	"github.com/iotaledger/wasp/packages/vm/core/testcore"
)

// executed in cluster_test.go
func testSpamOnledger(t *testing.T, e *chainEnv) {
	testutil.RunHeavy(t)
	// in the privtangle setup, with 1s milestones, this test takes ~50m to process 10k requests
	const numRequests = 10_000

	// send requests from many different wallets to speed things up
	numAccounts := 1000
	numRequestsPerAccount := numRequests / numAccounts
	errCh := make(chan error, numRequests)
	txCh := make(chan iotago.Transaction, numRequests)
	for i := 0; i < numAccounts; i++ {
		keyPair, _, err := e.Clu.NewKeyPairWithFunds()
		createWalletRetries := 0
		if err != nil {
			if createWalletRetries >= 5 {
				t.Fatal("failed to create wallet, got an error 5 times, %w", err)
			}
			// wait and re-try
			createWalletRetries++
			i--
			time.Sleep(1 * time.Second)
			continue
		}
		go func() {
			chainClient := e.Chain.SCClient(incHname, keyPair)
			retries := 0
			for i := 0; i < numRequestsPerAccount; i++ {
				tx, err := chainClient.PostRequest(incrementFuncName)
				if err != nil {
					if retries >= 5 {
						errCh <- fmt.Errorf("failed to issue tx, an error 5 times, %w", err)
						break
					}
					// wait and re-try the tx
					retries++
					i--
					time.Sleep(1 * time.Second)
					continue
				}
				retries = 0
				errCh <- err
				txCh <- *tx
				time.Sleep(1 * time.Second) // give time for the indexer to get the new UTXOs (so we don't issue conflicting txs)
			}
		}()
	}

	// wait for all requests to be sent
	for i := 0; i < numRequests; i++ {
		err := <-errCh
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < numRequests; i++ {
		tx := <-txCh
		_, err := e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, &tx, 30*time.Second)
		require.NoError(t, err)
	}

	waitUntil(t, e.counterEqualsCondition(int64(numRequests)), []int{0}, 5*time.Minute)

	res, err := e.Chain.Cluster.WaspClient(0).CallView(e.Chain.ChainID, blocklog.Contract.Hname(), blocklog.ViewGetEventsForBlock.Name, dict.Dict{})
	require.NoError(t, err)
	events, err := testcore.EventsViewResultToStringArray(res)
	require.NoError(t, err)
	println(events)
}

// executed in cluster_test.go
func testSpamOffLedger(t *testing.T, e *chainEnv) {
	testutil.RunHeavy(t)

	// we need to cap the limit of parallel requests, otherwise some reqs will fail due to local tcp limits: `dial tcp 127.0.0.1:9090: socket: too many open files`
	const maxParallelRequests = 700
	const numRequests = 100_000

	// deposit funds for offledger requests
	keyPair, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)

	e.DepositFunds(utxodb.FundsFromFaucetAmount, keyPair)

	myClient := e.Chain.SCClient(incHname, keyPair)

	durationsMutex := sync.Mutex{}
	processingDurationsSum := uint64(0)
	maxProcessingDuration := uint64(0)

	maxChan := make(chan int, maxParallelRequests)
	reqSuccessChan := make(chan uint64, numRequests)
	reqErrorChan := make(chan error, 1)

	go func() {
		for i := 0; i < numRequests; i++ {
			maxChan <- i
			nonce := uint64(i + 1)
			go func() {
				// send the request
				req, er := myClient.PostOffLedgerRequest(incrementFuncName, chainclient.PostRequestParams{Nonce: nonce})
				if er != nil {
					reqErrorChan <- er
					return
				}
				reqSentTime := time.Now()
				// wait for the request to be processed
				_, err = e.Chain.CommitteeMultiClient().WaitUntilRequestProcessedSuccessfully(e.Chain.ChainID, req.ID(), 5*time.Minute)
				if err != nil {
					reqErrorChan <- err
					return
				}
				processingDuration := uint64(time.Since(reqSentTime).Seconds())
				reqSuccessChan <- nonce
				<-maxChan

				durationsMutex.Lock()
				defer durationsMutex.Unlock()
				processingDurationsSum += processingDuration
				if processingDuration > maxProcessingDuration {
					maxProcessingDuration = processingDuration
				}
			}()
		}
	}()

	n := 0
	for {
		select {
		case <-reqSuccessChan:
			n++
		case e := <-reqErrorChan:
			// no request should fail
			fmt.Printf("ERROR sending offledger request, err: %v\n", e)
			t.Fatal(e)
		}
		if n == numRequests {
			break
		}
	}

	waitUntil(t, e.counterEqualsCondition(int64(numRequests)), []int{0}, 5*time.Minute)

	res, err := e.Chain.Cluster.WaspClient(0).CallView(e.Chain.ChainID, blocklog.Contract.Hname(), blocklog.ViewGetEventsForBlock.Name, dict.Dict{})
	require.NoError(t, err)
	events, err := testcore.EventsViewResultToStringArray(res)
	require.NoError(t, err)
	require.Regexp(t, fmt.Sprintf("counter = %d", numRequests), events[len(events)-1])
	avgProcessingDuration := processingDurationsSum / numRequests
	fmt.Printf("avg processing duration: %ds\n max: %ds\n", avgProcessingDuration, maxProcessingDuration)
}

// executed in cluster_test.go
func testSpamCallViewWasm(t *testing.T, e *chainEnv) {
	testutil.RunHeavy(t)

	wallet, _, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(t, err)
	client := e.Chain.SCClient(incHname, wallet)
	{
		// increment counter once
		tx, err := client.PostRequest(incrementFuncName)
		require.NoError(t, err)
		_, err = e.Chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(e.Chain.ChainID, tx, 30*time.Second)
		require.NoError(t, err)
	}

	const n = 200
	ch := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			r, err := client.CallView("getCounter", nil)
			if err != nil {
				ch <- err
				return
			}

			v, err := codec.DecodeInt64(r.MustGet(varCounter))
			if err == nil && v != 1 {
				err = errors.New("v != 1")
			}
			ch <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-ch
		if err != nil {
			t.Error(err)
		}
	}
}
