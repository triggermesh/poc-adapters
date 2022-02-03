/*
Copyright (c) 2022 TriggerMesh Inc.

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

// Package mongodbtarget implements an adapter that connects to a MongoDB database
// and allows a user to insert, query, and update documents via cloudevents.
package mongodbtarget

import (
	"context"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"

	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
)

// NewTarget returns the adapter implementation.
func NewTarget(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.ServerURL))
	if err != nil {
		return nil
	}

	return &mongodbAdapter{
		mclient: client,

		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*mongodbAdapter)(nil)

type mongodbAdapter struct {
	mclient *mongo.Client

	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Returns if stopCh is closed or Send() returns an error.
func (a *mongodbAdapter) Start(ctx context.Context) error {
	a.logger.Info("Starting MongoDB adapter")
	return a.ceClient.StartReceiver(ctx, a.dispatch)
}

func (a *mongodbAdapter) dispatch(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	ctx := context.Background()
	a.logger.Debug("Processing event")
	switch typ := event.Type(); typ {
	case "io.triggermesh.mongodb.insert":
		if err := a.insert(event, ctx); err != nil {
			return a.reportError("error invoking function: ", err)
		}

	case "io.triggermesh.mongodb.query.kv":
		resp, err := a.kvQuery(event, ctx)
		if err != nil {
			return a.reportError("error invoking kvQuery function: ", err)
		}

		return resp, nil
	case "io.triggermesh.mongodb.update":
		if err := a.update(event, ctx); err != nil {
			return a.reportError("error invoking update function: ", err)
		}

	default:
		return a.reportError("event type not supported ", nil)
	}

	return nil, cloudevents.ResultACK
}
func (a *mongodbAdapter) kvQuery(e cloudevents.Event, ctx context.Context) (*cloudevents.Event, cloudevents.Result) {
	qpd := &QueryPayload{}
	if err := e.DataAs(qpd); err != nil {
		return nil, err
	}

	collection := a.mclient.Database(qpd.Database).Collection(qpd.Collection)
	filterCursor, err := collection.Find(ctx, bson.M{qpd.Key: qpd.Value})
	if err != nil {
		return a.reportError("error finding in collection: ", err)
	}

	var episodesFiltered []bson.M
	if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
		return a.reportError("error filtering by cursor: ", err)
	}

	a.logger.Debug(episodesFiltered)
	responseEvent := cloudevents.NewEvent(cloudevents.VersionV1)
	err = responseEvent.SetData("application/json", episodesFiltered)
	if err != nil {
		return a.reportError("error generating response event: ", err)
	}

	responseEvent.SetType("io.triggermesh.mongodb.query.kv.result")
	responseEvent.SetSource("io.triggermesh.mongodb")
	responseEvent.SetSubject("query-result")
	responseEvent.SetDataContentType(cloudevents.ApplicationJSON)

	return &responseEvent, nil
}

func (a *mongodbAdapter) insert(e cloudevents.Event, ctx context.Context) error {
	a.logger.Infof("Inserting data..")
	ipd := &InsertPayload{}
	if err := e.DataAs(ipd); err != nil {
		return err
	}

	collection := a.mclient.Database(ipd.Database).Collection(ipd.Collection)
	if ipd.MapStrVal != nil {
		res, err := collection.InsertOne(ctx, ipd.MapStrVal)
		if err != nil {
			return err
		}

		a.logger.Infof("Posted data id:")
		a.logger.Info(res.InsertedID)
	}

	return nil
}

func (a *mongodbAdapter) update(e cloudevents.Event, ctx context.Context) error {
	a.logger.Infof("updating data..")
	up := &UpdatePayload{}
	if err := e.DataAs(up); err != nil {
		return err
	}

	collection := a.mclient.Database(up.Database).Collection(up.Collection)
	result, err := collection.UpdateOne(
		ctx,
		bson.M{up.SearchKey: up.SearchValue},
		bson.D{{Key: "$set", Value: bson.D{{Key: up.UpdateKey, Value: up.UpdateValue}}}},
	)
	if err != nil {
		return err
	}

	a.logger.Infof("Updated %v Documents!\n", result.ModifiedCount)
	return nil
}

func (a *mongodbAdapter) reportError(msg string, err error) (*cloudevents.Event, cloudevents.Result) {
	a.logger.Errorw(msg, zap.Error(err))

	responseEvent := cloudevents.NewEvent(cloudevents.VersionV1)
	responseEvent.SetType("io.triggermesh.mongodb.error")
	responseEvent.SetSource("io.triggermesh.mongodb")
	responseEvent.SetSubject("error")
	responseEvent.SetDataContentType(cloudevents.ApplicationJSON)
	if err := responseEvent.SetData(cloudevents.ApplicationJSON, msg); err != nil {
		a.logger.Errorw("could not set error response data")
	}

	return &responseEvent, cloudevents.NewHTTPResult(http.StatusInternalServerError, msg)
}
