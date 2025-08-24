package handler

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/event"
)

type eventRouter struct {
	routes map[string]event.Eventable
}

type EventRouter interface {
	Add(topic string, handler event.Eventable)
}

type EventRouterHandler interface {
	event.Eventable
	EventRouter
}

func NewRouter() EventRouterHandler {
	return &eventRouter{
		routes: map[string]event.Eventable{},
	}
}

func (h *eventRouter) Add(topic string, handler event.Eventable) {
	if topic == "" {
		panic("'topic' cannot be empty")
	}
	if handler == nil {
		panic("'handler' cannot be nil")
	}
	h.routes[topic] = handler
}

func (h *eventRouter) OnMessage(msg *kafka.Message) error {
	if msg == nil {
		return fmt.Errorf("'msg' cannot be nil")
	}
	if topic := msg.TopicPartition.Topic; topic != nil && *topic != "" {
		if handler, ok := h.routes[*topic]; ok {
			return handler.OnMessage(msg)
		} else {
			return fmt.Errorf("not found handler for topic '%s'", *topic)
		}
	} else {
		return fmt.Errorf("'topic' is nil or empty")
	}
}
