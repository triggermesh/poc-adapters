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
	"strings"
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
	tCloudEventSource = "ce.test.source"

	tCollection = "test"
	tDatabase   = "test"

	tArbitraryInstert = `{"test":"testvalue","test2":"test3"}`
	tInsert           = `{"database":"test","collection": "test","mapStrVal":{"test":"testvalue","test2":"test3"}}`
	tUpdate           = `{"database":"test","collection": "test","searchKey":"test","searchValue":"testvalue","updateKey":"partstore","updateValue":"UP FOR GRABS"}`
	tQuery            = `{"database":"test","Collection": "test","key":"partstore","value":"UP FOR GRABS"}`
)

// requires the enviroment variable `MONGODB_SERVER_URL` to contain a valid mongodb connection string

// Insert the tInsert test data into a mongodb table via sending an arbitrary event type
func TestArbitraryEvent(t *testing.T) {
	ctx := context.Background()
	serverURL := os.Getenv("MONGODB_SERVER_URL")
	require.NotEmpty(t, serverURL, "MONGODB_SERVER_URL must be set")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverURL))
	require.NotNil(t, client, "client should not be nil")
	require.Nil(t, err, "error should be nil")
	testCases := map[string]struct {
		inEvent cloudevents.Event
		mClient *mongo.Client
	}{
		"Consume event of type io.triggermesh.mongodb.arbitrary": {
			inEvent: newCloudEvent(tArbitraryInstert, "io.triggermesh.arbitrary"),
			mClient: client,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			ceClient := adaptertest.NewTestClient()
			mA := &mongodbAdapter{
				defaultDatabase:   tDatabase,
				defaultCollection: tCollection,

				logger:   logtesting.TestLogger(t),
				ceClient: ceClient,
				mclient:  tc.mClient,
			}

			_, _ = mA.dispatch(tc.inEvent)

			// find the inserted values
			episodesFiltered := findInsertedValues("test", "testvalue", client, t)
			assert.Equal(t, 1, len(episodesFiltered), "should contain 1")
			assert.Equal(t, "testvalue", string(episodesFiltered[0]["test"].(string)), "should contain `testvalue`")

		})
	}
	_ = client.Database(tDatabase).Collection(tCollection).Drop(ctx)
}

// Insert the tInsert test data into a mongodb table via sending an event of type io.triggermesh.mongodb.insert
func TestInsert(t *testing.T) {
	ctx := context.Background()
	serverURL := os.Getenv("MONGODB_SERVER_URL")
	require.NotEmpty(t, serverURL, "MONGODB_SERVER_URL must be set")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverURL))
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

			_, _ = mA.dispatch(tc.inEvent)

			// find the inserted values
			episodesFiltered := findInsertedValues("test", "testvalue", client, t)
			assert.Equal(t, 1, len(episodesFiltered), "should contain 1")
			assert.Equal(t, "testvalue", string(episodesFiltered[0]["test"].(string)), "should contain `testvalue`")

		})
	}
}

// Update the previously inserted tInsert test data into a mongodb table via sending an event of type io.triggermesh.mongodb.update
func TestUpdate(t *testing.T) {
	ctx := context.Background()
	serverURL := os.Getenv("MONGODB_SERVER_URL")
	require.NotEmpty(t, serverURL, "MONGODB_SERVER_URL must be set")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverURL))
	require.NotNil(t, client, "client should not be nil")
	require.Nil(t, err, "error should be nil")
	testCases := map[string]struct {
		inEvent cloudevents.Event
		mClient *mongo.Client
	}{
		"Consume event of type io.triggermesh.mongodb.update": {
			inEvent: newCloudEvent(tUpdate, "io.triggermesh.mongodb.update"),
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

			_, _ = mA.dispatch(tc.inEvent)

			// find the updated values
			episodesFiltered := findInsertedValues("partstore", "UP FOR GRABS", client, t)
			assert.Equal(t, 1, len(episodesFiltered), "should contain 1")
			assert.Equal(t, "UP FOR GRABS", string(episodesFiltered[0]["partstore"].(string)), "should contain `testvalue`")

		})
	}
}

// Find the changed values in the mongodb table created by the tested io.triggermesh.mongodb.update event
// by testing the consumtion of an event of type io.triggermesh.mongodb.query.kv
func TestQueryKV(t *testing.T) {
	ctx := context.Background()
	serverURL := os.Getenv("MONGODB_SERVER_URL")
	require.NotEmpty(t, serverURL, "MONGODB_SERVER_URL must be set")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverURL))
	require.NotNil(t, client, "client should not be nil")
	require.Nil(t, err, "error should be nil")
	testCases := map[string]struct {
		inEvent cloudevents.Event
		mClient *mongo.Client
	}{
		"Consume event of type io.triggermesh.mongodb.query.kv": {
			inEvent: newCloudEvent(tQuery, "io.triggermesh.mongodb.query.kv"),
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

			e, _ := mA.dispatch(tc.inEvent)
			sv := e.String()
			found := strings.Contains(sv, "UP FOR GRABS")
			assert.True(t, found, "should contain `UP FOR GRABS`")
		})
	}
	_ = client.Database(tDatabase).Collection(tCollection).Drop(ctx)
}

func findInsertedValues(key, value string, c *mongo.Client, t *testing.T) []bson.M {
	ctx := context.Background()
	var err error

	collection := c.Database("test").Collection("test")
	filterCursor, err := collection.Find(ctx, bson.M{key: value})
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
