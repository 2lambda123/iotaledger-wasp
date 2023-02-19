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
	webSocket  string
}

var _ IClientService = new(WasmClientService)

func NewWasmClientService(waspAPI string) *WasmClientService {
	client, err := apiextensions.WaspAPIClientByHostName(waspAPI)
	if err != nil {
		panic(err.Error())
	}
	return &WasmClientService{
		waspClient: client,
		webSocket:  strings.Replace(waspAPI, "http:", "ws:", 1) + "/ws",
	}
}

func (sc *WasmClientService) CallViewByHname(chainID wasmtypes.ScChainID, hContract, hFunction wasmtypes.ScHname, args []byte) ([]byte, error) {
	params, err := dict.FromBytes(args)
	if err != nil {
		return nil, err
	}

	res, _, err := sc.waspClient.RequestsApi.CallView(context.Background()).ContractCallViewRequest(apiclient.ContractCallViewRequest{
		ChainId:       cvt.IscChainID(&chainID).String(),
		ContractHName: cvt.IscHname(hContract).String(),
		FunctionHName: cvt.IscHname(hFunction).String(),
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

func (sc *WasmClientService) SubscribeEvents(msgChannel chan ContractEvent, done chan bool) error {
	ctx := context.Background()
	ws, _, err := websocket.Dial(ctx, sc.webSocket, nil)
	if err != nil {
		return err
	}
	err = eventSubscribe(ctx, ws, "chains")
	if err != nil {
		return err
	}
	err = eventSubscribe(ctx, ws, publisher.ISCEventKindSmartContract)
	if err != nil {
		return err
	}

	go eventLoop(ctx, ws, msgChannel)

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

func eventLoop(ctx context.Context, ws *websocket.Conn, msgChannel chan ContractEvent) {
	for {
		evt := publisher.ISCEvent{}
		err := wsjson.Read(ctx, ws, &evt)
		if err != nil {
			close(msgChannel)
			return
		}
		items := evt.Content.([]interface{})
		for _, item := range items {
			parts := strings.Split(item.(string), ": ")
			event := ContractEvent{
				ChainID:    evt.ChainID,
				ContractID: parts[0],
				Data:       parts[1],
			}
			msgChannel <- event
		}
	}
}

func eventSubscribe(ctx context.Context, ws *websocket.Conn, topic string) error {
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
