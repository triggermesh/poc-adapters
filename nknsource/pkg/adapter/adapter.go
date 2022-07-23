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

// Package NKNSOURCE implements a CloudEvents adapter that...
package nknsource

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	nkn "github.com/nknorg/nkn-sdk-go"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/uuid"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	targetce "github.com/triggermesh/triggermesh/pkg/targets/adapter/cloudevents"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

type envAccessor struct {
	pkgadapter.EnvConfig
	// BridgeIdentifier is the name of the bridge workflow this target is part of
	BridgeIdentifier string `envconfig:"EVENTS_BRIDGE_IDENTIFIER"`
	// CloudEvents responses parametrization
	CloudEventPayloadPolicy string `envconfig:"EVENTS_PAYLOAD_POLICY" default:"error"`
	// Sink defines the target sink for the events. If no Sink is defined the
	// events are replied back to the sender.
	Sink string `envconfig:"K_SINK" default:"http://localhost:8081`

	// Seed pasprase of a Wallet to use for recieving
	Seed string `envconfig:"K_SEED"`

	// Identifier
	Identifier []byte `envconfig:"K_IDENTIFIER" default:"[0 0 0 0 0 0 0 0]"`
}

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	replier, err := targetce.New(env.Component, logger.Named("replier"),
		targetce.ReplierWithStatefulHeaders(env.BridgeIdentifier),
		targetce.ReplierWithStaticResponseType("io.triggermesh.nkntarget.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	seed, err := hex.DecodeString(env.Seed)
	if err != nil {
		logger.Panicf("Error decoding seed from hex: %v", err)
	}

	account, err := nkn.NewAccount(seed)
	if err != nil {
		logger.Panicf("Error creating NKN account from seed: %v", err)
	}

	client, err := nkn.NewClient(account, "any string", nil)
	if err != nil {
		logger.Panicf("Error creating NKN client: %v", err)
	}

	return &nkntargetadapter{
		nknClient:  client,
		nknAccount: account,
		Identifier: env.Identifier,

		sink:     env.Sink,
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*nkntargetadapter)(nil)

type nkntargetadapter struct {
	nknClient  *nkn.Client
	nknAccount *nkn.Account
	Identifier []byte

	sink     string
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

type nknResponse struct {
	StatusCode int
	Error      string
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *nkntargetadapter) Start(ctx context.Context) error {
	a.logger.Info("NKNSource Adapter Started at address ", a.nknClient.Address())

	for {
		msg := <-a.nknClient.OnMessage.C
		err := a.emitCe(ctx, msg)
		if err != nil {
			a.logger.Errorf("Error emitting cloudevent: %v", err)
			nknr := nknResponse{
				StatusCode: 200,
				Error:      err.Error(),
			}
			msg.Reply(nknr)
		}

		nknr := nknResponse{
			StatusCode: 200,
			Error:      "",
		}
		b, err := json.Marshal(nknr)
		if err != nil {
			a.logger.Errorf("Error marshalling nkn response: %v", err)
			msg.Reply(err)
		}
		msg.Reply([]byte(b))
	}
	return nil
}

func (a *nkntargetadapter) emitCe(ctx context.Context, data *nkn.Message) error {
	event := cloudevents.NewEvent()
	event.SetID(string(uuid.NewUUID()))
	event.SetType("io.triggermesh.nkntarget.event")
	event.SetSource("a.nknClient.Address()")

	ctxwt := cloudevents.ContextWithTarget(context.Background(), a.sink)
	if result := a.ceClient.Send(ctxwt, event); !cloudevents.IsACK(result) {
		return fmt.Errorf("failed to send event: %s", result.Error())
	}
	return nil
}
