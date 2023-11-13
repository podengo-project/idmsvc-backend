package event

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"golang.org/x/exp/slog"
)

// Adapted from: https://github.com/RedHatInsights/playbook-dispatcher/blob/master/internal/response-consumer/main.go#L21

// Start initiate a kafka run loop consumer given the
// configuration and the event handler for the received
// messages.
// config a reference to an initialized KafkaConfig. It cannot be nil.
// handler is the event handler which receive the read messages.
func Start(ctx context.Context, config *config.Kafka, handler Eventable) {
	var (
		err      error
		consumer *kafka.Consumer
	)

	if consumer, err = NewConsumer(config); err != nil {
		slog.Error("error creating consumer", slog.Any("error", err.Error()))
		return
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			slog.Error("error closing consumer", slog.Any("error", err.Error()))
		}
	}()

	start := NewConsumerEventLoop(ctx, consumer, handler)
	start()
}
