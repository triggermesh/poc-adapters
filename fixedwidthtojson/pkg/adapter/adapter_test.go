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

package fixedwidthtojson

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	adaptertest "knative.dev/eventing/pkg/adapter/v2/test"
	logtesting "knative.dev/pkg/logging/testing"
)

const (
	tCloudEventID           = "ce-abcd-0123"
	tCloudEventType         = "ce.test.type"
	tCloudEventTypeResponse = "ce.test.type.response"
	tCloudEventSource       = "ce.test.source"

	tData = `NAME                STATE     TELEPHONE

	John Smith          WA        418-Y11-4111

	Mary Hartford       CA        319-Z19-4341

	Evan Nolan          IL        219-532-c301
	`

	tExpectedResponse = `{"fields":[{"value":"NAME","spaceLeft":0,"lineNumber":0},{"value":"STATE","spaceLeft":16,"lineNumber":0},{"value":" TELEPHONE","spaceLeft":4,"lineNumber":0},{"value":"\tJohn Smith","spaceLeft":0,"lineNumber":2},{"value":"WA","spaceLeft":10,"lineNumber":2},{"value":"418-Y11-4111","spaceLeft":8,"lineNumber":2},{"value":"\tMary Hartford","spaceLeft":0,"lineNumber":4},{"value":" CA","spaceLeft":6,"lineNumber":4},{"value":"319-Z19-4341","spaceLeft":8,"lineNumber":4},{"value":"\tEvan Nolan","spaceLeft":0,"lineNumber":6},{"value":"IL","spaceLeft":10,"lineNumber":6},{"value":"219-532-c301","spaceLeft":8,"lineNumber":6},{"value":"\t","spaceLeft":0,"lineNumber":7}]}`
)

func TestSink(t *testing.T) {
	testCases := map[string]struct {
		inEvent     cloudevents.Event
		expectEvent cloudevents.Event
	}{
		"sink ok": {
			inEvent:     newCloudEvent(t, tData, "application/text", tCloudEventType),
			expectEvent: newCloudEvent(t, tExpectedResponse, cloudevents.ApplicationJSON, tCloudEventTypeResponse),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ceClient := adaptertest.NewTestClient()
			ctx := context.Background()

			a := &fwadapter{
				logger:   logtesting.TestLogger(t),
				ceClient: ceClient,
				sink:     "http://localhost:8080",
			}

			e, r := a.dispatch(ctx, tc.inEvent)
			assert.Nil(t, e)
			assert.Equal(t, cloudevents.ResultACK, r)
			events := ceClient.Sent()
			require.Equal(t, 1, len(events))
			assert.Equal(t, tc.expectEvent.DataEncoded, events[0].DataEncoded)
		})
	}
}

func newCloudEvent(t *testing.T, data, contentType, eventType string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID(tCloudEventID)
	event.SetType(eventType)
	event.SetSource(tCloudEventSource)
	err := event.SetData(contentType, []byte(data))
	require.NoError(t, err)
	return event
}
