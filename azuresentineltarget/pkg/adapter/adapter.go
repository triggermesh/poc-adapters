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

// Package AZURESENTINELTARGET implements a CloudEvents adapter that...
package azuresentineltarget

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/google/uuid"
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
		targetce.ReplierWithStaticResponseType("io.triggermesh.azuresentineltarget.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &azuresentineltargetadapter{
		client:         http.DefaultClient,
		clientID:       env.ClientID,
		tenantID:       env.TenantID,
		azureCreds:     env.ClientSecret,
		subscriptionID: env.SubscriptionID,
		resourceGroup:  env.ResourceGroup,
		workspace:      env.Workspace,
		clientSecret:   env.ClientSecret,

		sink:     env.Sink,
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*azuresentineltargetadapter)(nil)

func (a *azuresentineltargetadapter) Start(ctx context.Context) error {
	a.logger.Info("Starting AZURESENTINELTARGET Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *azuresentineltargetadapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	ee := &expectedEvent{}
	if err := event.DataAs(ee); err != nil {
		a.logger.Errorf("Error decoding event: %v", err)
		return nil, nil
	}

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		a.logger.Errorf("Error creating Azure authorizer: %v", err)
		return nil, nil
	}

	incident := createIncident(*ee)
	reqBody, err := json.Marshal(*incident)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "marshaling request for retrieving an access token")
	}

	rURL := `https://management.azure.com/subscriptions/` + a.subscriptionID + `/resourceGroups/` + a.resourceGroup + `/providers/Microsoft.OperationalInsights/workspaces/` + a.workspace + `/providers/Microsoft.SecurityInsights/incidents/` + uuid.New().String() + `?api-version=2020-01-01`
	request, err := http.NewRequest(http.MethodPut, rURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "creating request token")
	}

	request.Header.Set("Content-Type", "application/json")
	req, err := autorest.Prepare(request,
		authorizer.WithAuthorization())
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "preparing request")
	}

	res, err := autorest.Send(req)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "sending request")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "reading response body ")
	}

	if res.StatusCode != 201 {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "invalid response from Azure: "+string(body))
	}

	return a.replier.Ok(&event, body)
}

func createIncident(ee expectedEvent) *Incident {
	i := &Incident{}
	alertProductNames := []string{}
	alertProductNames = append(alertProductNames, ee.Event.GUID)
	alertProductNames = append(alertProductNames, ee.Event.Name)
	alertProductNames = append(alertProductNames, ee.Event.Severity)
	alertProductNames = append(alertProductNames, ee.Event.ShortDescription)

	// i.Properties.ProviderIncidentId = ee.Event.GUID

	fmt.Printf("+%v", ee.Event)

	i.Properties.Title = ee.Event.Name
	i.Properties.Description = ee.Event.ShortDescription
	i.Properties.AdditionalData.AlertProductNames = alertProductNames

	i.Properties.Owner.AssignedTo = ee.Resource.Name

	// incidentLabelType := IncidentLabel{
	// 	LabelName: ee.Provider.AccountID,
	// 	LabelType: IncidentLabelType[
	// 		{
	// 		Name: "Account",
	// 		Type: "Azure",
	// 		},
	// 	]
	// 	},
	// }

	// i.Properties.Labels = incidentLabelType

	// i.Properties.Labels = []struct {
	// 	LabelName string `json:"labelName"`
	// 	LabelType string `json:"labelType"`
	// }{
	// 	{
	// 		LabelName: ee.Provider.AccountID,
	// 		LabelType: "accountID",
	// 	},
	// }
	// alertProductNames := []string{}
	// alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Platform)
	// alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].AccountID)
	// alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Region)
	// alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Service)
	// alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Type+":"+ee.Event.Event.Resources[0].Name+":"+ee.Event.Event.Resources[0].GUID)
	// i.Properties.Title = ee.Event.Event.Metadata.Name
	// i.Properties.Description = ee.Event.Event.Metadata.ShortDescription
	// i.Properties.AdditionalData.AlertProductNames = alertProductNames
	// i.Properties.Labels = []struct {
	// 	LabelName string `json:"labelName"`
	// 	LabelType string `json:"labelType"`
	// }{
	// 	{
	// 		LabelName: "accountID",
	// 		LabelType: ee.Provider.AccountID,
	// 	},
	// }
	i.Properties.Severity = "High"
	i.Properties.Status = "Active"
	return i
}
