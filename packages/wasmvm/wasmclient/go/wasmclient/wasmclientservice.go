// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmclient

// for some reason we cannot use the import name mangos, so we rename those packages
// for some other reason if the third mamgos import is missing things won't work
import (
	"context"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/clients/apiclient"
	"github.com/iotaledger/wasp/clients/apiextensions"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/publisher"
	"github.com/iotaledger/wasp/packages/publisher/publisherws"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type ContractEvent struct {
	ChainID    string
	ContractID string
	Data       string
}

type IClientService interface {
	CallViewByHname(chainID wasmtypes.ScChainID, hContract, hFunction wasmtypes.ScHname, args []byte) ([]byte, error)
	PostRequest(chainID wasmtypes.ScChainID, hContract, hFunction wasmtypes.ScHname, args []byte, allowance *wasmlib.ScAssets, keyPair *cryptolib.KeyPair, nonce uint64) (wasmtypes.ScRequestID, error)
	SubscribeEvents(msg chan ContractEvent, done chan bool) error
	WaitUntilRequestProcessed(chainID wasmtypes.ScChainID, reqID wasmtypes.ScRequestID, timeout time.Duration) error
}

type WasmClientService struct {
	waspClient *apiclient.APIClient
	eventPort  string
}

var _ IClientService = new(WasmClientService)

func NewWasmClientService(waspAPI, eventPort string) *WasmClientService {
	client, err := apiextensions.WaspAPIClientByHostName(waspAPI)
	if err != nil {
		panic(err.Error())
	}

	return &WasmClientService{waspClient: client, eventPort: eventPort}
}

func DefaultWasmClientService() *WasmClientService {
	return NewWasmClientService("http://localhost:19090", "ws://localhost:19090/ws")
}

func (sc *WasmClientService) CallViewByHname(chainID wasmtypes.ScChainID, hContract, hFunction wasmtypes.ScHname, args []byte) ([]byte, error) {
	iscChainID := cvt.IscChainID(&chainID)
	iscContract := cvt.IscHname(hContract)
	iscFunction := cvt.IscHname(hFunction)
	params, err := dict.FromBytes(args)
	if err != nil {
		return nil, err
	}
	res, _, err := sc.waspClient.RequestsApi.CallView(context.Background()).ContractCallViewRequest(apiclient.ContractCallViewRequest{
		ContractHName: iscContract.String(),
		FunctionHName: iscFunction.String(),
		ChainId:       iscChainID.String(),
		Arguments:     apiextensions.JSONDictToAPIJSONDict(params.JSONDict()),
	}).Execute()
	if err != nil {
		return nil, err
	}

	decodedParams, err := apiextensions.APIJsonDictToDict(*res)
	if err != nil {
		return nil, err
	}

	return decodedParams.Bytes(), nil
}

func (sc *WasmClientService) PostRequest(chainID wasmtypes.ScChainID, hContract, hFunction wasmtypes.ScHname, args []byte, allowance *wasmlib.ScAssets, keyPair *cryptolib.KeyPair, nonce uint64) (reqID wasmtypes.ScRequestID, err error) {
	iscChainID := cvt.IscChainID(&chainID)
	iscContract := cvt.IscHname(hContract)
	iscFunction := cvt.IscHname(hFunction)
	params, err := dict.FromBytes(args)
	if err != nil {
		return reqID, err
	}
	req := isc.NewOffLedgerRequest(iscChainID, iscContract, iscFunction, params, nonce)
	iscAllowance := cvt.IscAllowance(allowance)
	req.WithAllowance(iscAllowance)
	signed := req.Sign(keyPair)
	reqID = cvt.ScRequestID(signed.ID())

	_, err = sc.waspClient.RequestsApi.OffLedger(context.Background()).OffLedgerRequest(apiclient.OffLedgerRequest{
		ChainId: iscChainID.String(),
		Request: iotago.EncodeHex(signed.Bytes()),
	}).Execute()

	return reqID, err
}

func (sc *WasmClientService) subscribe(ctx context.Context, ws *websocket.Conn, topic string) error {
	msg := publisherws.SubscriptionCommand{
		Command: publisherws.CommandSubscribe,
		Topic:   topic,
	}
	err := wsjson.Write(ctx, ws, msg)
	if err != nil {
		return err
	}
	return wsjson.Read(ctx, ws, &msg)
}

func (sc *WasmClientService) SubscribeEvents(msgChannel chan ContractEvent, done chan bool) error {
	ctx := context.Background()
	ws, _, err := websocket.Dial(ctx, sc.eventPort, nil)
	if err != nil {
		return err
	}

	err = sc.subscribe(ctx, ws, "chains")
	if err != nil {
		return err
	}

	err = sc.subscribe(ctx, ws, publisher.ISCEventKindSmartContract)
	if err != nil {
		return err
	}

	//err = sc.subscribe(ctx, ws, publisher.ISCEventKindNewBlock)
	//if err != nil {
	//	return err
	//}
	//
	//err = sc.subscribe(ctx, ws, publisher.ISCEventKindReceipt)
	//if err != nil {
	//	return err
	//}

	go func() {
		for {
			evt := publisher.ISCEvent{}
			err := wsjson.Read(ctx, ws, &evt)
			if err != nil {
				close(msgChannel)
				return
			}
			if evt.Content != nil {
				items := evt.Content.([]interface{})
				for _, item := range items {
					parts := strings.Split(item.(string), ": ")
					// contract tst1pqqf4qxh2w9x7rz2z4qqcvd0y8n22axsx82gqzmncvtsjqzwmhnjs438rhk | vm (contract): 89703a45: testwasmlib.test|1671671237|tst1pqqf4qxh2w9x7rz2z4qqcvd0y8n22axsx82gqzmncvtsjqzwmhnjs438rhk|Lala
					event := ContractEvent{
						ChainID:    evt.ChainID,
						ContractID: parts[0],
						Data:       parts[1],
					}
					msgChannel <- event
				}
			}
		}
	}()

	go func() {
		<-done
		ws.Close(websocket.StatusNormalClosure, "intentional close")
	}()

	return nil
}

func (sc *WasmClientService) WaitUntilRequestProcessed(chainID wasmtypes.ScChainID, reqID wasmtypes.ScRequestID, timeout time.Duration) error {
	iscChainID := cvt.IscChainID(&chainID)
	iscReqID := cvt.IscRequestID(&reqID)

	_, _, err := sc.waspClient.RequestsApi.
		WaitForRequest(context.Background(), iscChainID.String(), iscReqID.String()).
		TimeoutSeconds(int32(timeout.Seconds())).
		Execute()

	return err
}
