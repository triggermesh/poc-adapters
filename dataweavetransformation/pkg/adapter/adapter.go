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

package dataweavetransformation

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
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

type envAccessor struct {
	pkgadapter.EnvConfig
	// Spell defines the Dataweave spell to use on the incoming data at the event payload.
	Spell string `envconfig:"DW_SPELL" required:"true"`
	// IncomingContentType defines the expected content type of the incoming data.
	IncomingContentType string `envconfig:"INCOMING_CONTENT_TYPE" default:"application/json"`
	// OutputContentType defines the content the cloudevent to be sent with the transformed data.
	OutputContentType string `envconfig:"OUTPUT_CONTENT_TYPE" default:"application/json"`
	// BridgeIdentifier is the name of the bridge workflow this target is part of.
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
		targetce.ReplierWithStaticResponseType("io.triggermesh.dataweavetransformation.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &adapter{
		spell:               env.Spell,
		incomingContentType: env.IncomingContentType,
		outputContentType:   env.OutputContentType,

		sink:     env.Sink,
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*adapter)(nil)

type adapter struct {
	spell               string
	incomingContentType string
	outputContentType   string

	sink     string
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *adapter) Start(ctx context.Context) error {
	a.logger.Info("Starting Dataweave Transformation Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *adapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	var err error
	tmpfile, err := ioutil.TempFile("/app", "*.json")
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "creating the file")
	}

	if _, err := tmpfile.Write(event.Data()); err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "writing to the file")
	}

	cn := strings.Replace(tmpfile.Name(), "/app/", "", 1)
	out, err := exec.Command("/app/.dw/bin/dw", "-i", "payload", cn, a.spell).Output()
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "executing the spell")
	}

	err = os.Remove(tmpfile.Name())
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "removing the file")
	}

	if err := event.SetData(a.outputContentType, out); err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, nil)
	}

	return &event, cloudevents.ResultACK
}
