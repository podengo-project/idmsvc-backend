package event

import (
	"context"
	"fmt"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/podengo-project/idmsvc-backend/internal/api/event"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/event/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
)

func TestNewConsumer(t *testing.T) {
	var (
		consumer *kafka.Consumer
		err      error
	)

	type TestCase struct {
		Given    *config.Kafka
		Expected error
	}

	testCases := []TestCase{
		// When config is nil
		{
			Given:    nil,
			Expected: fmt.Errorf("config cannot be nil"),
		},
		// Failing kafka.NewConsumer
		{
			Given: &config.Kafka{
				Bootstrap: struct{ Servers string }{},
			},
			Expected: fmt.Errorf("Configuration property \"auto.offset.reset\" cannot be set to empty value"),
		},
		// Fail SubscribeTopics because unknown group
		{
			Given: &config.Kafka{
				Bootstrap: struct{ Servers string }{
					Servers: "localhost:9092",
				},
				Auto: struct {
					Offset struct{ Reset string }
					Commit struct{ Interval struct{ Ms int } }
				}{
					Offset: struct{ Reset string }{
						Reset: "latest",
					},
				},
			},
			Expected: fmt.Errorf("Local: Unknown group"),
		},
		// Fail SubscribeTopics because unknown group
		{
			Given: &config.Kafka{
				Bootstrap: struct{ Servers string }{
					Servers: "localhost:9092",
				},
				Auto: struct {
					Offset struct{ Reset string }
					Commit struct{ Interval struct{ Ms int } }
				}{
					Offset: struct{ Reset string }{
						Reset: "latest",
					},
				},
			},
			Expected: fmt.Errorf("Local: Unknown group"),
		},
		// Success return
		{
			Given: &config.Kafka{
				Bootstrap: struct{ Servers string }{
					Servers: "localhost:9092",
				},
				Auto: struct {
					Offset struct{ Reset string }
					Commit struct{ Interval struct{ Ms int } }
				}{
					Offset: struct{ Reset string }{
						Reset: "latest",
					},
				},
				Group: struct{ Id string }{
					Id: "main",
				},
				Topics: []string{
					event.TopicTodoCreated,
				},
			},
			Expected: nil,
		},
		// Success return with sasl
		{
			Given: &config.Kafka{
				Bootstrap: struct{ Servers string }{
					Servers: "localhost:9092",
				},
				Auto: struct {
					Offset struct{ Reset string }
					Commit struct{ Interval struct{ Ms int } }
				}{
					Offset: struct{ Reset string }{
						Reset: "latest",
					},
				},
				Group: struct{ Id string }{
					Id: "main",
				},
				Topics: []string{
					event.TopicTodoCreated,
				},
				Sasl: struct {
					Username  string
					Password  string `json:"-"`
					Mechanism string
					Protocol  string
				}{
					Username:  "myusername",
					Password:  "mypassword",
					Mechanism: "SCRAM-SHA-512",
					Protocol:  "sasl_plaintext",
				},
				Capath: "",
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		consumer, err = NewConsumer(testCase.Given)
		if testCase.Expected != nil {
			assert.Nil(t, consumer)
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Error(), err.Error())
		} else {
			assert.NotNil(t, consumer)
			assert.NoError(t, err)
			if err != nil {
				assert.Equal(t, "", err.Error())
			}
		}
	}
}

type MockEventable struct {
	mock.Mock
}

func (m *MockEventable) OnMessage(msg *kafka.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func TestProcessConsumedMessage(t *testing.T) {
	type TestCaseGiven struct {
		Schemas event.TopicSchema
		Message *kafka.Message
		Handler Eventable
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected error
	}

	msgValid := &kafka.Message{
		Key: []byte("this-is-my-key"),
		TopicPartition: kafka.TopicPartition{
			Topic: pointy.String(event.TopicTodoCreated),
		},
		Headers: []kafka.Header{
			{
				Key:   "Type",
				Value: []byte(message.HdrTypeIntrospect),
			},
		},
		Value: []byte(`{"id":12345,"title":"todo title","description":"todo description"}`),
	}
	msgNoValid := &kafka.Message{
		Key: []byte("this-is-my-key"),
		TopicPartition: kafka.TopicPartition{
			Topic: pointy.String(event.TopicTodoCreated),
		},
		Headers: []kafka.Header{
			{
				Key:   string(message.HdrType),
				Value: []byte(message.HdrTypeIntrospect),
			},
		},
		Value: []byte(`{}`),
	}
	mockOnMessageFailure := &MockEventable{}
	mockOnMessageFailure.On("OnMessage", msgValid).Return(fmt.Errorf("Error in handler"))
	mockOnMessageSuccess := &MockEventable{}
	mockOnMessageSuccess.On("OnMessage", msgValid).Return(nil)

	schemas, err := event.LoadSchemas()
	require.NoError(t, err)
	require.NotNil(t, schemas)

	testCases := []TestCase{
		// nil arguments return error
		{
			Name: "force error for nil arguments",
			Given: TestCaseGiven{
				Schemas: nil,
				Message: nil,
				Handler: nil,
			},
			Expected: fmt.Errorf("schemas, msg or handler is nil"),
		},
		// nil topic
		{
			Name: "force error when topic is nil",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic: nil,
					},
				},
				Handler: mockOnMessageFailure,
			},
			Expected: fmt.Errorf("Topic cannot be nil"),
		},
		// Wrong topic
		{
			Name: "force error when topic does not exist",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic: pointy.String("AnyNonExistingTopic"),
					},
				},
				Handler: mockOnMessageFailure,
			},
			Expected: fmt.Errorf("Topic mapping not found for: AnyNonExistingTopic"),
		},
		// Invalid message
		{
			Name: "force error when message is not validated",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic: pointy.String(event.TopicTodoCreated),
					},
				},
				Handler: mockOnMessageFailure,
			},
			Expected: fmt.Errorf("data cannot be nil"),
		},
		// Error validating message schema
		{
			Name: "force error when validating message schema",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: msgNoValid,
				Handler: mockOnMessageFailure,
			},
			Expected: fmt.Errorf("error validating schema: \"id\" value is required: / = map[], \"title\" value is required: / = map[], \"description\" value is required: / = map[]"),
		},
		// Valid message but failure on handler
		{
			Name: "force error when the handler return error",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: msgValid,
				Handler: mockOnMessageFailure,
			},
			Expected: fmt.Errorf("Error in handler"),
		},
		// Valid message handled
		{
			Name: "success case where the message is handled",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: msgValid,
				Handler: mockOnMessageSuccess,
			},
			Expected: nil,
		},
	}

	config.TopicTranslationConfig = config.NewTopicTranslationWithDefaults()

	for _, testCase := range testCases {
		t.Logf("Testing case: '%s'", testCase.Name)
		result := processConsumedMessage(
			testCase.Given.Schemas,
			testCase.Given.Message,
			testCase.Given.Handler)
		if testCase.Expected != nil {
			require.Error(t, result)
			assert.Equal(t, testCase.Expected.Error(), result.Error())
		} else {
			assert.NoError(t, result)
		}
	}
}

func TestNewConsumerEventLoop(t *testing.T) {
	var (
		result   func()
		consumer *kafka.Consumer
		cfg      config.Kafka
		h        Eventable
		err      error
	)
	assert.PanicsWithErrorf(t, "consumer cannot be nil", func() {
		NewConsumerEventLoop(context.Background(), nil, nil)
	}, "consumer cannot be nil")

	cfg = config.Kafka{}
	cfg.Auto.Offset.Reset = "latest"
	cfg.Topics = []string{event.TopicTodoCreated}
	cfg.Group.Id = "unit-tests"
	consumer, err = NewConsumer(&cfg)
	require.NotNil(t, consumer)
	require.NoError(t, err)
	assert.PanicsWithErrorf(t, "handler cannot be nil", func() {
		NewConsumerEventLoop(context.Background(), consumer, nil)
	}, "consumer cannot be nil")

	h = &MockEventable{}
	assert.NotPanics(t, func() {
		result = NewConsumerEventLoop(context.Background(), consumer, h)
	})
	assert.NotNil(t, result)
}
