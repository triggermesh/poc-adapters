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

// Package jsontoxmltransformation implements a CloudEvents adapter that accepts a JSON payload
// and returns a correspoding event containg the same payload represented in XML.
package jsontoxmltransformation

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/itchyny/gojq"
	"vimagination.zapto.org/json2xml"

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
		targetce.ReplierWithStaticResponseType("io.triggermesh.jsontoxmltransformation.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &jqadapter{
		sink:     env.Sink,
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*jqadapter)(nil)

type jqadapter struct {
	query *gojq.Query

	sink     string
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *jqadapter) Start(ctx context.Context) error {
	a.logger.Info("Starting JSONToXMLTransformation Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *jqadapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	r := io.Reader(bytes.NewReader(event.Data()))
	json := json.NewDecoder(r)
	op := bytes.NewBuffer(nil)
	iow := io.Writer(op)
	xml := xml.NewEncoder(iow)
	if err := json2xml.Convert(json, xml); err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "converting the data")
	}

	var x interface{}
	xml.Encode(x)
	if err := event.SetData(cloudevents.ApplicationXML, op.Bytes()); err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, nil)
	}

	if a.sink != "" {
		if result := a.ceClient.Send(ctx, event); !cloudevents.IsACK(result) {
			return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, result, "sending the cloudevent to the sink")
		}
		return nil, cloudevents.ResultACK
	}

	return &event, cloudevents.ResultACK
}
