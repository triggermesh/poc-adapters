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

// Package kafkatarget implements a CloudEvents adapter that..
package kafkatarget

import (
	"context"
	"encoding/json"
	"fmt"
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

// KafkaClient is a wrapper of the Confluent Kafka producer
// functions needed for the Confluent adapter.
type KafkaClient interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	CreateKafkaAdminClient() (KafkaAdminClient, error)
	Flush(timeoutMs int) int
	Close()
}

type envAccessor struct {
	pkgadapter.EnvConfig

	BootstrapServers    string `envconfig:"CONFLUENT_BOOTSTRAP_SERVERS" required:"true"`
	Topic               string `envconfig:"CONFLUENT_TOPIC" required:"true"`
	GroupID             string `envconfig:"CONFLUENT_GROUP_ID" required:"true" `
	SecurityMechanisms  string `envconfig:"CONFLUENT_SECURITY_MECANISMS" required:"false" default:"GSSAPI"`
	SecurityProtocol    string `envconfig:"CONFLUENT_SECURITY_PROTOCOL" required:"false" default:"SASL_SSL"`
	KerberosPrincipal   string `envconfig:"KERBEROS_PRINCIPAL" required:"true" `
	KerberosServiceName string `envconfig:"KERBEROS_SERVICE_NAME" required:"true" `
	SSLCALocation       string `envconfig:"SSL_CA_LOCATION" required:"true"`
	KerberosKeytab      string `envconfig:"KERBEROS_KEYTAB" required:"true"`

	// This set of variables are experimental and not graduated to the CRD.
	BrokerVersionFallback       string `envconfig:"CONFLUENT_BROKER_VERSION_FALLBACK" required:"false" default:"0.10.0.0"`
	APIVersionFallbackMs        string `envconfig:"CONFLUENT_API_VERSION_FALLBACK_MS" required:"false" default:"0"`
	CreateTopicIfMissing        bool   `envconfig:"CONFLUENT_CREATE_MISSING_TOPIC" default:"true"`
	FlushOnExitTimeoutMillisecs int    `envconfig:"CONFLUENT_FLUSH_ON_EXIT_TIMEOUT_MS" default:"10000"`
	CreateTopicTimeoutMillisecs int    `envconfig:"CONFLUENT_CREATE_TOPIC_TIMEOUT_MS" default:"10000"`
	NewTopicPartitions          int    `envconfig:"CONFLUENT_TOPIC_PARTITIONS" default:"1"`
	NewTopicReplicationFactor   int    `envconfig:"CONFLUENT_TOPIC_REPLICATION_FACTOR" default:"1"`

	DiscardCEContext bool `envconfig:"CONFLUENT_DISCARD_CE_CONTEXT"`
}

// NewAdapter adapter implementation
func NewAdapter(ctx context.Context, envAcc pkgadapter.EnvConfigAccessor, ceClient cloudevents.Client) pkgadapter.Adapter {
	env := envAcc.(*envAccessor)
	logger := logging.FromContext(ctx)

	kafkaClient, err := NewKafkaClient(&kafka.ConfigMap{
		"bootstrap.servers":          env.BootstrapServers,
		"group.id":                   env.GroupID,
		"fetch.min.bytes":            1000000,
		"fetch.max.bytes":            1000000000,
		"fetch.wait.max.ms":          1 * time.Second,
		"security.protocol":          env.SecurityProtocol,
		"sasl.mechanisms":            env.SecurityMechanisms,
		"sasl.kerberos.service.name": env.KerberosServiceName,
		"sasl.kerberos.principal":    env.KerberosPrincipal,
		"sasl.kerberos.keytab":       env.KerberosKeytab,
		"ssl.ca.location":            env.SSLCALocation,
	})

	if err != nil {
		logger.Panic(err)
	}

	return &kafkaAdapter{
		kafkaClient:               kafkaClient,
		topic:                     env.Topic,
		createTopicIfMissing:      env.CreateTopicIfMissing,
		flushTimeout:              env.FlushOnExitTimeoutMillisecs,
		topicTimeout:              env.CreateTopicTimeoutMillisecs,
		newTopicPartitions:        env.NewTopicPartitions,
		newTopicReplicationFactor: env.NewTopicReplicationFactor,

		discardCEContext: env.DiscardCEContext,

		ceClient: ceClient,
		logger:   logger,
	}
}

var _ pkgadapter.Adapter = (*kafkaAdapter)(nil)

type kafkaAdapter struct {
	kafkaClient KafkaClient
	topic       string

	createTopicIfMissing bool

	flushTimeout              int
	topicTimeout              int
	newTopicPartitions        int
	newTopicReplicationFactor int

	discardCEContext bool

	ceClient cloudevents.Client
	logger   *zap.SugaredLogger
}

// Start is a blocking function and will return if an error occurs
// or the context is cancelled.
func (a *kafkaAdapter) Start(ctx context.Context) error {
	a.logger.Info("Starting Kafka Target Adapter")

	defer func() {
		a.kafkaClient.Flush(a.flushTimeout)
		a.kafkaClient.Close()
	}()

	if a.createTopicIfMissing {
		if err := a.ensureTopic(ctx, a.topic); err != nil {
			return fmt.Errorf("failed ensuring Topic %s: %w", a.topic, err)
		}
	}

	if err := a.ceClient.StartReceiver(ctx, a.dispatch); err != nil {
		return fmt.Errorf("error starting the cloud events server: %w", err)
	}
	return nil
}
func (a *kafkaAdapter) dispatch(event cloudevents.Event) cloudevents.Result {
	var msgVal []byte

	if a.discardCEContext {
		msgVal = event.Data()
	} else {
		jsonEvent, err := json.Marshal(event)
		if err != nil {
			a.logger.Errorw("Error marshalling CloudEvent", zap.Error(err))
			return cloudevents.ResultNACK
		}
		msgVal = jsonEvent
	}

	km := &kafka.Message{
		Key:            []byte(event.ID()),
		TopicPartition: kafka.TopicPartition{Topic: &a.topic, Partition: kafka.PartitionAny},
		Value:          msgVal,
	}

	// librdkafka provides buffering, we set channel size to 1
	// to avoid blocking tests as they execute in the same thread
	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	if err := a.kafkaClient.Produce(km, deliveryChan); err != nil {
		a.logger.Errorw("Error producing Kafka message", zap.String("msg", km.String()), zap.Error(err))
		return cloudevents.ResultNACK
	}

	r := <-deliveryChan
	m := r.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		a.logger.Infof("Message delivery failed: %v", m.TopicPartition.Error)
		return cloudevents.ResultNACK
	}

	a.logger.Info("Delivered message to topic %s [%d] at offset %v",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

	return cloudevents.ResultACK
}

//ensureTopic creates a topic if missing
func (a *kafkaAdapter) ensureTopic(ctx context.Context, topic string) error {
	a.logger.Infof("Ensuring topic %q", topic)

	adminClient, err := a.kafkaClient.CreateKafkaAdminClient()
	if err != nil {
		return fmt.Errorf("error creating admin client from producer: %w", err)
	}
	defer adminClient.Close()

	ts := []kafka.TopicSpecification{{
		Topic:             topic,
		NumPartitions:     a.newTopicPartitions,
		ReplicationFactor: a.newTopicReplicationFactor}}

	m, err := adminClient.GetMetadata(&topic, false, a.topicTimeout)
	if err != nil {
		return fmt.Errorf("error retrieving topic %q metadata: %w", topic, err)
	}
	if m == nil {
		return fmt.Errorf("empty response requesting topic metadata for %q", a.topic)
	}

	t, ok := m.Topics[a.topic]
	if !ok {
		return fmt.Errorf("topic %q metadata response does not contain required information", a.topic)
	}

	switch t.Error.Code() {
	case kafka.ErrNoError:
		a.logger.Infof("Topic found: %q with %d partitions", t.Topic, len(t.Partitions))
		return nil
	case kafka.ErrUnknownTopic, kafka.ErrUnknownTopicOrPart:
		// topic does not exists, we need to create it.
	default:
		return fmt.Errorf("topic %q metadata returned inexpected status: %w", a.topic, t.Error)
	}

	a.logger.Infof("Creating topic %q", topic)
	results, err := adminClient.CreateTopics(ctx, ts, kafka.SetAdminOperationTimeout(time.Duration(a.topicTimeout)*time.Millisecond))
	if err != nil {
		return fmt.Errorf("error creating topic %q: %w", a.topic, err)
	}

	if len(results) != 1 {
		return fmt.Errorf("creating topic %s returned inexpected results: %+v", a.topic, results)
	}

	if results[0].Error.Code() != kafka.ErrNoError {
		return fmt.Errorf("failed to create topic %s: %w", a.topic, results[0].Error)
	}

	return nil
}
