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

// Package PATHSADAPTER implements a CloudEvents adapter that accepts javascript if expressions
// along with a new type and path to route the event to if the expression evaluates to true.
package pathsadapter

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	targetce "github.com/triggermesh/triggermesh/pkg/targets/adapter/cloudevents"

	"github.com/robertkrimen/otto"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

type envAccessor struct {
	// PathAContinueIf is a javascript if expression to evaluate to determine if event should continue to path A
	PathAContinueIf string `envconfig:"PATH_A_CONTINUE_IF"`
	// PathAContinuePath is a path to route event to if path A should continue
	PathAContinuePath string `envconfig:"PATH_A_CONTINUE_PATH"`
	// PathAContinueType is the event type to change the outgoing event type to if path A should continue
	PathAContinueType string `envconfig:"PATH_A_CONTINUE_TYPE"`

	// PathBContinueIf is a javascript if expression to evaluate to determine if event should continue to path B
	PathBContinueIf string `envconfig:"PATH_B_CONTINUE_IF"`
	// PathBContinuePath is a path to route event to if path B should continue
	PathBContinuePath string `envconfig:"PATH_B_CONTINUE_PATH"`
	// PathBContinueType is the event type to change the outgoing event type to if path B should continue
	PathBContinueType string `envconfig:"PATH_B_CONTINUE_TYPE"`

	// DefaultContinuePath is a path to route event to if no other paths should continue
	DefaultContinuePath string `envconfig:"DEFAULT_CONTINUE_PATH"`
	// DefaultContinueType is the event type to change the outgoing event type to if no other paths should continue
	DefaultContinueType string `envconfig:"DEFAULT_CONTINUE_TYPE"`

	pkgadapter.EnvConfig
	// BridgeIdentifier is the name of the bridge workflow this target is part of
	BridgeIdentifier string `envconfig:"EVENTS_BRIDGE_IDENTIFIER"`
	// CloudEvents responses parametrization
	CloudEventPayloadPolicy string `envconfig:"EVENTS_PAYLOAD_POLICY" default:"error"`
}

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	replier, err := targetce.New(env.Component, logger.Named("replier"),
		targetce.ReplierWithStatefulHeaders(env.BridgeIdentifier),
		targetce.ReplierWithStaticResponseType("io.triggermesh.paths.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &pathsadapter{
		pathAContinueIf:   env.PathAContinueIf,
		pathAContinueType: env.PathAContinueType,
		pathAContinuePath: env.PathAContinuePath,

		pathBContinueIf:   env.PathBContinueIf,
		pathBContinueType: env.PathBContinueType,
		pathBContinuePath: env.PathBContinuePath,

		defaultContinuePath: env.DefaultContinuePath,
		defaultContinueType: env.DefaultContinueType,

		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}

}

var _ pkgadapter.Adapter = (*pathsadapter)(nil)

type pathsadapter struct {
	pathAContinueIf   string
	pathAContinuePath string
	pathAContinueType string

	pathBContinueIf   string
	pathBContinuePath string
	pathBContinueType string

	defaultContinuePath string
	defaultContinueType string

	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *pathsadapter) Start(ctx context.Context) error {
	a.logger.Info("Starting Path Adapter")

	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *pathsadapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	vm := otto.New()
	vm.Set("event", string(event.Data()))
	vm.Set("response", 0)
	_, err := vm.Run(`
		event = JSON.parse(event);
		if ` + a.pathAContinueIf + ` {
			response = 1
		}
		if ` + a.pathBContinueIf + ` {
			response = 2
		}
	`)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "running script")
	}

	value, err := vm.Get("response")
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "getting response")
	}

	value_int, err := value.ToInteger()
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "getting response to integer")
	}

	if value_int == 0 {
		a.logger.Infof("Sending event to default path")
		event.SetType(a.defaultContinueType)
		ctx := cloudevents.ContextWithTarget(context.Background(), a.defaultContinuePath)
		if result := a.ceClient.Send(ctx, event); cloudevents.IsUndelivered(result) {
			log.Fatalf("failed to send, %v", result)
			return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, result, "failed to send event to default path")
		}
	}

	if value_int == 1 {
		a.logger.Infof("Sending event to path A")
		event.SetType(a.pathAContinueType)
		ctx := cloudevents.ContextWithTarget(context.Background(), a.pathAContinuePath)
		if result := a.ceClient.Send(ctx, event); cloudevents.IsUndelivered(result) {
			log.Fatalf("failed to send, %v", result)
			return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "getting response to path a")

		}
	}

	if value_int == 2 {
		a.logger.Infof("Sending event to path B")
		event.SetType(a.pathBContinueType)
		ctx := cloudevents.ContextWithTarget(context.Background(), a.pathBContinuePath)
		if result := a.ceClient.Send(ctx, event); cloudevents.IsUndelivered(result) {
			log.Fatalf("failed to send, %v", result)
			return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "getting response to path B")
		}
	}

	return &event, cloudevents.ResultACK
}
