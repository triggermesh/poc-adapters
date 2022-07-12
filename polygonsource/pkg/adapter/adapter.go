/*
Copyright 2022 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package POLYGONSOURCE implements a CloudEvents adapter that...
package polygonsource

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/google/uuid"
	"go.uber.org/zap"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	targetce "github.com/triggermesh/triggermesh/pkg/targets/adapter/cloudevents"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

type Result struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

type PolyTransactionResults struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []Result `json:"result"`
}

type envAccessor struct {
	pkgadapter.EnvConfig

	PollingTimeout time.Duration `envconfig:"POLYGONSOURCE_POLLING_TIMEOUT" default:"5s"`

	WalletAddress string `envconfig:"POLYGONSOURCE_WALLET_ADDRESS"`

	PolyScanAPIKey string `envconfig:"POLYGONSOURCE_API_KEY"`

	// IngoreFirstBatch is a bool flag that tells the adapter if it should ignore the first batch of transactions.
	IngoreFirstBatch bool `envconfig:"POLYGONSOURCE_IGNORE_FIRST_BATCH" default:"false"`

	// BridgeIdentifier is the name of the bridge workflow this target is part of
	BridgeIdentifier string `envconfig:"EVENTS_BRIDGE_IDENTIFIER"`
	// CloudEvents responses parametrization
	CloudEventPayloadPolicy string `envconfig:"EVENTS_PAYLOAD_POLICY" default:"error"`
	// Sink defines the target sink for the events. If no Sink is defined the
	// events are replied back to the sender.
	Sink string `envconfig:"K_SINK"`
}

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	replier, err := targetce.New(env.Component, logger.Named("replier"),
		targetce.ReplierWithStatefulHeaders(env.BridgeIdentifier),
		targetce.ReplierWithStaticResponseType("io.triggermesh.polygonsource.response"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &polygonsourceadapter{
		WalletAddress:    env.WalletAddress,
		pollingTimeout:   env.PollingTimeout,
		polyScanAPIKey:   env.PolyScanAPIKey,
		ignoreFirstBatch: env.IngoreFirstBatch,

		sink:     env.Sink,
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*polygonsourceadapter)(nil)

type polygonsourceadapter struct {
	WalletAddress    string
	pollingTimeout   time.Duration
	polyScanAPIKey   string
	walletState      []Result
	ignoreFirstBatch bool

	sink     string
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *polygonsourceadapter) Start(ctx context.Context) error {
	a.logger.Info("Starting POLYGONSOURCE Adapter")
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			a.poll(ctx)
			a.sleep()
		}
	}
	return nil
}

func (a *polygonsourceadapter) poll(ctx context.Context) (err error) {
	url := "https://api.polygonscan.com/api?module=account&action=txlist&address=" + a.WalletAddress + "&apikey=" + a.polyScanAPIKey
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	ptr := &PolyTransactionResults{}

	err = json.Unmarshal(body, ptr)
	if err != nil {
		return err
	}

	if a.ignoreFirstBatch {
		for _, tx := range ptr.Result {
			a.logger.Debug("Ignoring first batch of transactions")
			a.walletState = append(a.walletState, tx)
		}
		return nil
	}

	// append to walletState if not already in there
	for _, tx := range ptr.Result {
		if !a.isTxInWalletState(tx) {
			a.logger.Debug("adding tx to walletState")
			a.walletState = append(a.walletState, tx)
			// send new event
			a.emitCE(ctx, tx)
		}
	}

	return nil
}

func (a *polygonsourceadapter) emitCE(ctx context.Context, r Result) protocol.Result {
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType("io.triggermesh.polygonsource.response")
	event.SetSource(a.WalletAddress)
	event.SetData(cloudevents.ApplicationJSON, r)

	if result := a.ceClient.Send(ctx, event); result.Error != nil {
		a.logger.Errorf("Error sending event: %v", result.Error)
		return result
	}

	return nil
}

func (a *polygonsourceadapter) isTxInWalletState(tx Result) bool {
	for _, t := range a.walletState {
		if t.Hash == tx.Hash {
			return true
		}
	}
	return false
}

func (a *polygonsourceadapter) sleep() {
	time.Sleep(time.Second * a.pollingTimeout)
}
