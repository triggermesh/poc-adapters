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

package mongodbtarget

import (
	"context"
	"os"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	adaptertest "knative.dev/eventing/pkg/adapter/v2/test"
	logtesting "knative.dev/pkg/logging/testing"
)

const (
	tCloudEventID     = "ce-abcd-0123"
	tCloudEventType   = "ce.test.type"
	tCloudEventSource = "ce.test.source"

	tCollection = "test"
	tDatabase   = "test"

	tQuery        = ".foo | .."
	tInsert       = `{"database":"test","collection": "test","mapStrVal":{"test":"testvalue","test2":"test3"}}`
	tExpectedJSON = "{\"foo\":\"richard@triggermesh.com\"}"

	tFalseJSON       = "not json"
	tExpectedFailure = "{\"Code\":\"request-parsing\",\"Description\":\"[xml] found bytes, but failed to unmarshal: EOF not json\",\"Details\":null}"
)

// requires the enviroment variable `MONGODB_SERVER_URL` to contain a valid mongodb connection string

func TestInsert(t *testing.T) {
	ctx := context.Background()
	serverUrl := os.Getenv("MONGODB_SERVER_URL")
	require.NotEmpty(t, serverUrl, "MONGODB_SERVER_URL must be set")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverUrl))
	require.NotNil(t, client, "client should not be nil")
	require.Nil(t, err, "error should be nil")
	testCases := map[string]struct {
		inEvent cloudevents.Event
		mClient *mongo.Client
	}{
		"Consume event of type io.triggermesh.mongodb.insert": {
			inEvent: newCloudEvent(tInsert, "io.triggermesh.mongodb.insert"),
			mClient: client,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {

			ceClient := adaptertest.NewTestClient()
			mA := &mongodbAdapter{
				logger:   logtesting.TestLogger(t),
				ceClient: ceClient,
				mclient:  tc.mClient,
			}

			mA.dispatch(tc.inEvent)

			// find the inserted values
			episodesFiltered := findInsertedValues("testvalue", client, t)
			assert.Equal(t, 1, len(episodesFiltered), "should contain 1")
			assert.Equal(t, "testvalue", string(episodesFiltered[0]["test"].(string)), "should contain `testvalue`")

		})
	}
	// cleanup
	client.Database(tDatabase).Collection(tCollection).Drop(ctx)
}

func findInsertedValues(value string, c *mongo.Client, t *testing.T) []bson.M {
	ctx := context.Background()
	var err error

	collection := c.Database("test").Collection("test")
	filterCursor, err := collection.Find(ctx, bson.M{"test": value})
	require.Nil(t, err, "error should be nil")

	var episodesFiltered []bson.M
	err = filterCursor.All(ctx, &episodesFiltered)
	require.Nil(t, err, "error should be nil")

	return episodesFiltered
}

func newCloudEvent(data, cetype string) cloudevents.Event {
	event := cloudevents.NewEvent()

	if err := event.SetData(cloudevents.ApplicationJSON, []byte(data)); err != nil {
		// not expected
		panic(err)
	}

	event.SetID(tCloudEventID)
	event.SetSource(tCloudEventSource)
	event.SetType(cetype)

	return event
}

// func TestSink(t *testing.T) {
// 	testCases := map[string]struct {
// 		inEvent     cloudevents.Event
// 		expectEvent cloudevents.Event
// 		query       string
// 	}{
// 		"sink ok": {
// 			inEvent:     newCloudEvent(t, tJSON, cloudevents.ApplicationJSON),
// 			expectEvent: newCloudEvent(t, tExpectedJSON, cloudevents.ApplicationJSON),
// 			query:       tQuery,
// 		},
// 	}
// 	for name, tc := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				body, err := ioutil.ReadAll(r.Body)
// 				assert.NoError(t, err)
// 				assert.Equal(t, tExpectedJSON, string(body))
// 				fmt.Fprintf(w, "OK")
// 			}))
// 			defer svr.Close()

// 			env := &envAccessor{
// 				EnvConfig: adapter.EnvConfig{
// 					Component: tCloudEventSource,
// 					Sink:      svr.URL,
// 				},
// 				Query: tc.query,
// 			}
// 			ctx := context.Background()
// 			c, err := cloudevents.NewClientHTTP()
// 			assert.NoError(t, err)
// 			a := NewAdapter(ctx, env, c)

// 			go func() {
// 				if err := a.Start(ctx); err != nil {
// 					assert.FailNow(t, "could not start test adapter")
// 				}
// 			}()

// 			response := sendCE(t, &tc.inEvent, c, svr.URL)
// 			assert.NotEqual(t, cloudevents.IsUndelivered(response), response)
// 		})
// 	}
// }
