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
	"time"

	"github.com/Azure/go-autorest/autorest"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
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

type envAccessor struct {
	pkgadapter.EnvConfig
	SubscriptionID string `envconfig:"AZURE_SUBSCRIPTION_ID" required:"true"`
	ResourceGroup  string `envconfig:"AZURE_RESOURCE_GROUP" required:"true"`
	Workspace      string `envconfig:"AZURE_WORKSPACE" required:"true"`
	ClientSecret   string `envconfig:"AZURE_CLIENT_SECRET" required:"true"`
	ClientID       string `envconfig:"AZURE_CLIENT_ID" required:"true"`
	TenantID       string `envconfig:"AZURE_TENANT_ID" required:"true"`
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

	rURL := `https://management.azure.com/subscriptions/` + env.SubscriptionID + `/resourceGroups/` + env.ResourceGroup + `/providers/Microsoft.OperationalInsights/workspaces/` + env.Workspace + `/providers/Microsoft.OperationalInsights/sources/` + env.BridgeIdentifier + `/events?api-version=2018-11-01`

	replier, err := targetce.New(env.Component, logger.Named("replier"),
		targetce.ReplierWithStatefulHeaders(env.BridgeIdentifier),
		targetce.ReplierWithStaticResponseType("io.triggermesh.azuresentineltarget.error"),
		targetce.ReplierWithPayloadPolicy(targetce.PayloadPolicy(env.CloudEventPayloadPolicy)))
	if err != nil {
		logger.Panicf("Error creating CloudEvents replier: %v", err)
	}

	return &azuresentineltargetadapter{
		client:     http.DefaultClient,
		requestURL: rURL,
		incidendID: uuid.New(),
		clientID:   env.ClientID,
		tenantID:   env.TenantID,
		azureCreds: env.ClientSecret,

		sink:     env.Sink,
		replier:  replier,
		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*azuresentineltargetadapter)(nil)

type Incident struct {
	Properties struct {
		Severity       string `json:"severity"`
		Status         string `json:"status"`
		Title          string `json:"title"`
		Description    string `json:"description"`
		AdditionalData struct {
			AlertProductNames []string `json:"alertProductNames"`
		} `json:"additionalData"`
		Labels []struct {
			LabelName string `json:"labelName"`
			LabelType string `json:"labelType"`
		} `json:"labels"`
	} `json:"properties"`
}

type expectedEvent struct {
	Event struct {
		Event struct {
			Metadata struct {
				GUID             int         `json:"guid"`
				Name             string      `json:"name"`
				URL              interface{} `json:"url"`
				Severity         string      `json:"severity"`
				ShortDescription string      `json:"shortDescription"`
				LongDescription  string      `json:"longDescription"`
				Time             int         `json:"time"`
			} `json:"metadata"`
			Producer struct {
				Name string `json:"name"`
			} `json:"producer"`
			Reporter struct {
				Name string `json:"name"`
			} `json:"reporter"`
			Resources []struct {
				GUID      string `json:"guid"`
				Name      string `json:"name"`
				Region    string `json:"region"`
				Platform  string `json:"platform"`
				Service   string `json:"service"`
				Type      string `json:"type"`
				AccountID string `json:"accountId"`
				Package   string `json:"package"`
			} `json:"resources"`
		} `json:"event"`
		Decoration []struct {
			Decorator string    `json:"decorator"`
			Timestamp time.Time `json:"timestamp"`
			Payload   struct {
				Registry         string    `json:"registry"`
				Namespace        string    `json:"namespace"`
				Image            string    `json:"image"`
				Tag              string    `json:"tag"`
				Digests          []string  `json:"digests"`
				ImageLastUpdated time.Time `json:"imageLastUpdated"`
				TagLastUpdated   time.Time `json:"tagLastUpdated"`
				Description      string    `json:"description"`
				StarCount        int       `json:"starCount"`
				PullCount        int64     `json:"pullCount"`
			} `json:"payload"`
		} `json:"decoration"`
	} `json:"event"`
	Sourcetype string `json:"sourcetype"`
}

type azuresentineltargetadapter struct {
	client     *http.Client
	requestURL string
	incidendID uuid.UUID
	clientID   string
	tenantID   string
	azureCreds string

	sink     string
	replier  *targetce.Replier
	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *azuresentineltargetadapter) Start(ctx context.Context) error {
	a.logger.Info("Starting AZURESENTINELTARGET Adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *azuresentineltargetadapter) dispatch(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	// a.logger.Infof("Received event: %v", event)
	ee := &expectedEvent{}
	if err := event.DataAs(ee); err != nil {
		a.logger.Errorf("Error decoding event: %v", err)
		return nil, nil
	}

	i := &Incident{}
	alertProductNames := []string{}
	alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Platform)
	alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].AccountID)
	alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Region)
	alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Service)
	alertProductNames = append(alertProductNames, ee.Event.Event.Resources[0].Type+":"+ee.Event.Event.Resources[0].Name+":"+ee.Event.Event.Resources[0].GUID)
	i.Properties.Title = ee.Event.Event.Metadata.Name
	i.Properties.Description = ee.Event.Event.Metadata.ShortDescription
	i.Properties.AdditionalData.AlertProductNames = alertProductNames
	i.Properties.Labels = []struct {
		LabelName string `json:"labelName"`
		LabelType string `json:"labelType"`
	}{
		{
			LabelName: ee.Event.Event.Reporter.Name,
			LabelType: "User",
		},
	}

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		a.logger.Errorf("Error creating Azure authorizer: %v", err)
		return nil, nil
	}

	// config := auth.NewDeviceFlowConfig(a.clientID, a.tenantID)
	// autho, err := config.Authorizer()

	// cred, err := confidential.NewCredFromSecret(a.azureCreds)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not create a cred from a secret: %w", err)
	// }

	// confidentialClientApp, err := confidential.New(a.clientID, cred, confidential.WithAuthority("https://login.microsoftonline.com/Enter_The_Tenant_Name_Here"))

	reqBody, err := json.Marshal(*i)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "marshaling request for retrieving an access token")
	}

	request, err := http.NewRequest(http.MethodPost, a.requestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "creating request token")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer ")

	req, err := autorest.Prepare(request,
		authorizer.WithAuthorization())
	fmt.Println(authorizer.WithAuthorization())

	res, err := autorest.Send(req)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "sending request")
	}

	// res, err := a.client.Do(request)
	// if err != nil {
	// 	return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "processing auth request")
	// }

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "reading response body ")
	}

	fmt.Println(string(body))
	fmt.Println(res.StatusCode)

	// if res.StatusCode != http.StatusOK {
	// 	return a.replier.Error(&event, targetce.ErrorCodeAdapterProcess, err, "invalid response from Azure")
	// }

	fmt.Println(string(body))

	return nil, cloudevents.ResultACK
}
