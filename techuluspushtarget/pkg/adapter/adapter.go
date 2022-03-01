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

// Package techuluspushtarget implements a CloudEvents adapter that accepts a JSON payload
// and returns a correspoding event containg the same payload represented in XML.
package techuluspushtarget

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
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
	// Refrence to a valid techulus push API KEY
	APIKey string `envconfig:"API_KEY" required:"true"`
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
		targetce.ReplierWithStaticResponseType("io.triggermesh.techuluspush.target.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &pushAdapter{
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*pushAdapter)(nil)

type pushAdapter struct {
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

type eventPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *pushAdapter) Start(ctx context.Context) error {
	a.logger.Info("Starting techulus push target Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *pushAdapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	ep := &eventPayload{}
	if err := event.DataAs(ep); err != nil {
		a.logger.Errorw("Failed to unmarshal event data", zap.Error(err))
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "Error creating request")
	}

	url := "https://push.techulus.com/api/v1/notify/83316c2c-5069-439e-b30a-78a349af080f?title=" + urlify(ep.Title) + "&body=" + urlify(ep.Body)
	fmt.Println(url)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		if err != nil {
			a.logger.Errorw("Failed to create request", zap.Error(err))
			return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "Error creating request")
		}
	}
	res, err := client.Do(req)
	if err != nil {
		a.logger.Errorw("Failed to send request", zap.Error(err))
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "Error sending request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		a.logger.Errorw("Failed to read response body", zap.Error(err))
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "Error reading response body")
	}
	fmt.Println(string(body))

	event.SetType(event.Type() + ".response")
	err = event.SetData("application/json", body)
	if err != nil {
		a.logger.Errorw("Failed to set event data", zap.Error(err))
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "Error setting event data")
	}
	a.logger.Infof("responding with transformed event: %v", event.Type())

	return &event, cloudevents.ResultACK
}

func urlify(s string) string {
	sr := strings.ReplaceAll(s, " ", "%20")
	return sr
}
