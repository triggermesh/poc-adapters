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

// Package fixedwidthtojson implements a CloudEvents adapter ...
package fixedwidthtojson

import (
	"context"
	"strings"

	"go.uber.org/zap"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	targetce "github.com/triggermesh/triggermesh/pkg/targets/adapter/cloudevents"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)
	replier, err := targetce.New(env.Component, logger.Named("replier"),
		targetce.ReplierWithStatefulHeaders(env.BridgeIdentifier),
		targetce.ReplierWithStaticResponseType("io.triggermesh.dataweavetransformation.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &fwadapter{
		sink:     env.Sink,
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*fwadapter)(nil)

type fwadapter struct {
	sink     string
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *fwadapter) Start(ctx context.Context) error {
	a.logger.Info("Starting fixedwidthToJSON Transformation Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *fwadapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	sd := string(event.Data())
	ss := strings.Split(sd, "\n")
	var spaceleft []int
	var fields []string
	var fieldLineNumber []int
	for i, s := range ss {
		split := strings.Split(s, "  ")
		var space int
		for _, sl := range split {
			if sl == "" {
				space++
			} else {
				fields = append(fields, sl)
				fieldLineNumber = append(fieldLineNumber, i)
				if space > 0 {
					spaceleft = append(spaceleft, (space*2)+2)
				} else {
					spaceleft = append(spaceleft, 0)
				}
				space = 0
			}
		}
	}

	fwj := &fixedwidthJSONRepresentation{}
	for i, f := range fields {
		fv := []Field{}
		fv = append(fv, Field{
			Value:      f,
			Spaceleft:  spaceleft[i],
			LineNumber: fieldLineNumber[i],
		})
		fwj.Fields = append(fwj.Fields, fv...)
	}

	err := event.SetData(cloudevents.ApplicationJSON, fwj)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "setting ce data")
	}

	event.SetType(event.Type() + ".response")
	a.logger.Infof("responding with transformed event: %v", event.Type())
	if a.sink != "" {
		if result := a.ceClient.Send(ctx, event); !cloudevents.IsACK(result) {
			a.logger.Errorf("Error sending event to sink: %v", result)
		}
		return nil, cloudevents.ResultACK
	}

	return &event, cloudevents.ResultACK
}
