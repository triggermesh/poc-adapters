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

// Package kafkasource implements a CloudEvents adapter that..
package kafkasource

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	pkgadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// EnvAccessorCtor for configuration parameters
func EnvAccessorCtor() pkgadapter.EnvConfigAccessor {
	return &envAccessor{}
}

type envAccessor struct {
	pkgadapter.EnvConfig

	BootstrapServers    []string `envconfig:"CONFLUENT_BOOTSTRAP_SERVERS" required:"true"`
	Topics              []string `envconfig:"CONFLUENT_TOPICS" required:"true"`
	GroupID             string   `envconfig:"CONFLUENT_GROUP_ID" required:"true"`
	SecurityMechanisms  string   `envconfig:"CONFLUENT_SECURITY_MECANISMS" required:"false" default:"GSSAPI"`
	SecurityProtocol    string   `envconfig:"CONFLUENT_SECURITY_PROTOCOL" required:"false" default:"SASL_SSL"`
	KerberosPrincipal   string   `envconfig:"KERBEROS_PRINCIPAL" required:"true"`
	KerberosServiceName string   `envconfig:"KERBEROS_SERVICE_NAME" required:"true"`
	SSLCALocation       string   `envconfig:"SSL_CA_LOCATION" required:"true"`
	KerberosKeytab      string   `envconfig:"KERBEROS_KEYTAB" required:"true"`
}

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	kafkaClient, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":          env.BootstrapServers,
		"group.id":                   env.GroupID,
		"fetch.min.bytes":            1000000,
		"fetch.max.bytes":            1000000000,
		"fetch.wait.max.ms":          1 * time.Second,
		"security.protocol":          env.SecurityProtocol,
		"security.mecanisms":         env.SecurityMechanisms,
		"sasl.kerberos.service.name": env.KerberosServiceName,
		"sasl.kerberos.principal":    env.KerberosPrincipal,
		"sasl.kerberos.keytab":       env.KerberosKeytab,
		"ssl.ca.location":            env.SSLCALocation,
	})

	if err != nil {
		logger.Panic(err)
	}

	return &kafkaAdapter{
		kafkaConsumer: kafkaClient,
		topics:        env.Topics,

		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*kafkaAdapter)(nil)

type kafkaAdapter struct {
	kafkaConsumer *kafka.Consumer
	topics        []string

	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *kafkaAdapter) Start(ctx context.Context) error {
	a.logger.Info("Starting Kafka Target Adapter")
	err := a.kafkaConsumer.SubscribeTopics(a.topics, nil)
	if err != nil {
		return err
	}

	run := true
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := a.kafkaConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				a.logger.Infof("Received message on topic %s", e.TopicPartition.Topic)
				err := a.emitEvent(ctx, string(e.Value), e.TopicPartition)
				if err != nil {
					a.logger.Errorf("Failed to emit event: %v", err)
				}
			case kafka.Error:
				a.logger.Error("Error: %v: %v", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				a.logger.Infof("Ignored %v\n", e)
			}
		}
	}

	a.kafkaConsumer.Close()
	return nil
}

func (a *kafkaAdapter) emitEvent(ctx context.Context, message string, topicPartition kafka.TopicPartition) error {
	event := cloudevents.NewEvent(cloudevents.VersionV1)
	event.SetType("io.triggermesh.kafka.event")
	event.SetSubject("/kafka/target/event")
	event.SetSource(*topicPartition.Topic)
	event.SetID(topicPartition.Offset.String())
	in := `{"message": ` + message + `}`
	if err := event.SetData(cloudevents.ApplicationJSON, in); err != nil {
		return fmt.Errorf("failed to set event data: %w", err)
	}

	if result := a.ceClient.Send(context.Background(), event); !cloudevents.IsACK(result) {
		return result
	}
	return nil
}
