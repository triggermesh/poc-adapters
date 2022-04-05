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

// Package TERAFORMTARGET implements a CloudEvents adapter that...
package teraformtarget

import (
	"context"
	"io/ioutil"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	"fmt"

	targetce "github.com/triggermesh/triggermesh/pkg/targets/adapter/cloudevents"

	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

type envAccessor struct {
	TerrafromPlan   string `envconfig:"TERAFORMTARGET_TERRAFORM_PLAN"`
	TerraformConfig string `envconfig:"TERAFORMTARGET_TERRAFORM_CONFIG"`
	WorkingDir      string `envconfig:"TERAFORMTARGET_WORKING_DIR"`
	Config          string `envconfig:"TERAFORMTARGET_CONFIG"`
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

	// os.Setenv("KUBECONFIG", env.WorkingDir+"/config")

	tmpfile, err := ioutil.TempFile(env.WorkingDir, "*.tf")
	if err != nil {
		logger.Fatal("error creating tempfile", zap.Error(err))
	}

	tmpfile.Write([]byte(env.TerrafromPlan))

	replier, err := targetce.New(env.Component, logger.Named("replier"),
		targetce.ReplierWithStatefulHeaders(env.BridgeIdentifier),
		targetce.ReplierWithStaticResponseType("io.triggermesh.teraformtarget.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &teraformtargetadapter{
		workingdir: env.WorkingDir,
		sink:       env.Sink,
		replier:    replier,
		ceClient:   ceClient,
		logger:     logger,
	}
}

var _ pkgadapter.Adapter = (*teraformtargetadapter)(nil)

type teraformtargetadapter struct {
	workingdir string
	sink       string
	replier    *targetce.Replier
	ceClient   cloudevents.Client
	logger     *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *teraformtargetadapter) Start(ctx context.Context) error {
	a.logger.Info("Starting TERAFORMTARGET Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *teraformtargetadapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	a.logger.Infof("Received event: %v", event)

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "error installing Terraform")

	}

	tf, err := tfexec.NewTerraform(a.workingdir, execPath)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "error running NewTerraform")
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "error running Init")
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "error running Show")
	}

	if err := tf.Apply(context.Background()); err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "error running Apply")
	}

	fmt.Println(state.FormatVersion)

	return nil, cloudevents.ResultACK
}
