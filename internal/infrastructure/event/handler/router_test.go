package handler

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/podengo-project/idmsvc-backend/internal/test/mock/infrastructure/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.NotNil(t, router)
}

func TestRouterAdd(t *testing.T) {
	mockHandler := event.NewEventable(t)
	router := NewRouter()
	assert.Panics(t, func() {
		router.Add("", mockHandler)
	})
	assert.Panics(t, func() {
		router.Add("test-topic", nil)
	})
	assert.NotPanics(t, func() {
		router.Add("test-topic", mockHandler)
	})
}

func TestOnMessage(t *testing.T) {
	var (
		err error
		msg *kafka.Message
	)
	msg = &kafka.Message{}
	router := NewRouter()
	assert.NotPanics(t, func() {
		err := router.OnMessage(msg)
		assert.Error(t, err)
	})

	err = router.OnMessage(nil)
	require.Error(t, err)
	assert.Equal(t, "'msg' cannot be nil", err.Error())

	msg = &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: nil,
		},
	}
	err = router.OnMessage(msg)
	require.Error(t, err)
	assert.Equal(t, "'topic' is nil or empty", err.Error())

	msg = &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: pointy.String(""),
		},
	}
	err = router.OnMessage(msg)
	require.Error(t, err)
	assert.Equal(t, "'topic' is nil or empty", err.Error())

	msg = &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: pointy.String("non-existing-topic"),
		},
	}
	err = router.OnMessage(msg)
	require.Error(t, err)
	assert.Equal(t, "not found handler for topic 'non-existing-topic'", err.Error())

	msg = &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: pointy.String("demo-topic"),
		},
	}
	mockHandler := event.NewEventable(t)
	mockHandler.On("OnMessage", msg).
		Return(nil)
	router.Add("demo-topic", mockHandler)
	err = router.OnMessage(msg)
	assert.NoError(t, err)
}
