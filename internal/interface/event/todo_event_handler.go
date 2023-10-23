package event

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type TodoInterface interface {
	TodoCreated(msg *kafka.Message) error
}
