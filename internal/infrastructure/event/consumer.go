package event

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	schema "github.com/podengo-project/idmsvc-backend/internal/api/event"
	"github.com/podengo-project/idmsvc-backend/internal/config"
)

// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
// https://docs.confluent.io/platform/current/clients/consumer.html#ak-consumer-configuration

// NewConsumer create a new consumer based on the configuration
// supplied.
// config Provide the necessary configuration to create the consumer.
func NewConsumer(config *config.Kafka) (*kafka.Consumer, error) {
	var (
		consumer *kafka.Consumer
		err      error
	)

	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	kafkaConfigMap := &kafka.ConfigMap{
		"bootstrap.servers":        config.Bootstrap.Servers,
		"group.id":                 config.Group.Id,
		"auto.offset.reset":        config.Auto.Offset.Reset,
		"auto.commit.interval.ms":  config.Auto.Commit.Interval.Ms,
		"go.logs.channel.enable":   false,
		"allow.auto.create.topics": true,
		// NOTE This could be useful when launching locally
		// "socket.timeout.ms":                  60000,
		// "socket.connection.setup.timeout.ms": 3000,
		// "session.timeout.ms":                 6000,
	}

	if config.Sasl.Username != "" {
		_ = kafkaConfigMap.SetKey("sasl.username", config.Sasl.Username)
		_ = kafkaConfigMap.SetKey("sasl.password", config.Sasl.Password)
		_ = kafkaConfigMap.SetKey("sasl.mechanism", config.Sasl.Mechanism)
		_ = kafkaConfigMap.SetKey("security.protocol", config.Sasl.Protocol)
		_ = kafkaConfigMap.SetKey("ssl.ca.location", config.Capath)
	}

	if consumer, err = kafka.NewConsumer(kafkaConfigMap); err != nil {
		return nil, err
	}

	if err = consumer.SubscribeTopics(config.Topics, nil); err != nil {
		return nil, err
	}
	slog.Info(
		"Consumer subscribed to topics",
		slog.String("topics", strings.Join(config.Topics, ",")),
	)

	return consumer, nil
}

func checkConsumerEventLoop(ctx context.Context, consumer *kafka.Consumer, handler Eventable) schema.TopicSchema {
	if consumer == nil {
		panic(fmt.Errorf("consumer cannot be nil"))
	}
	if handler == nil {
		panic(fmt.Errorf("handler cannot be nil"))
	}
	if schemas, err := schema.LoadSchemas(); err != nil {
		panic(err)
	} else {
		return schemas
	}
}

// NewConsumerEventLoop creates a consumer event loop, which is awaiting for
// new kafka messages and process them by the specified handler.
//
// consumer is an initialized kafka.Consumer. It cannot be nil.
// handler is the event handler which will dispatch the received messages.
// It cannot be nil.
//
// Return a function that represent the event loop or a panic if a failure
// happens.
// TODO Refactor this function to reduce the complexity
// nolint:gocognit
func NewConsumerEventLoop(ctx context.Context, consumer *kafka.Consumer, handler Eventable) func() {
	var (
		err     error
		msg     *kafka.Message
		schemas schema.TopicSchema
	)
	if schemas, err = schema.LoadSchemas(); err != nil {
		panic(err)
	}

	schemas = checkConsumerEventLoop(ctx, consumer, handler)

	return func() {
		slog.Info("Consumer loop awaiting to consume messages")
		for {
			// Message wait loop
			for {
				if msg, err = consumer.ReadMessage(1 * time.Second); err == nil {
					break
				}

				val, ok := err.(kafka.Error)
				if !ok || val.Code() != kafka.ErrTimedOut {
					slog.Error(
						"error awaiting to read a message",
						slog.Any("error", err),
					)
				}

				select {
				case <-ctx.Done():
					slog.Info("Context done for NewConsumerEventLoop")
					return
				default:
				}

				continue
			}

			// FIXME Serialize message into the database before call proccess
			if err = processConsumedMessage(schemas, msg, handler); err != nil {
				logEventMessageError(msg, err)
				// FIXME Update serialized message to flag as error processing
				continue
			}
		}
	}
}

func processConsumedMessage(schemas schema.TopicSchema, msg *kafka.Message, handler Eventable) error {
	var err error

	if schemas == nil || msg == nil || handler == nil {
		return fmt.Errorf("schemas, msg or handler is nil")
	}
	if msg.TopicPartition.Topic == nil {
		return fmt.Errorf("Topic cannot be nil")
	}

	internalTopic := config.TopicTranslationConfig.GetInternal(*msg.TopicPartition.Topic)
	if internalTopic == "" {
		return fmt.Errorf("Topic mapping not found for: %s", *msg.TopicPartition.Topic)
	}
	slog.Info(
		"Topic mapping",
		slog.String("topic_name", *msg.TopicPartition.Topic),
		slog.String("requested_topic_name", internalTopic),
	)
	*msg.TopicPartition.Topic = internalTopic
	logEventMessageInfo(msg, "Consuming message")

	if err = schemas.ValidateMessage(msg); err != nil {
		return err
	}

	// Dispatch message
	if err = handler.OnMessage(msg); err != nil {
		return err
	}
	return nil
}
