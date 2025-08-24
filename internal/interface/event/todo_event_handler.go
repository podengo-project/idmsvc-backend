package event

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type TodoInterface interface {
	TodoCreated(msg *kafka.Message) error
}
