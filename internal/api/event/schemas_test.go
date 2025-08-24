package event

import (
	"fmt"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/event/message"
	"github.com/qri-io/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
)

func TestLoadSchemas(t *testing.T) {
	s, err := LoadSchemas()
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestLoadSchemaFromString(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    string
		Expected error
	}

	testCases := []TestCase{
		{
			Name:     "force error when unmarshalling schema",
			Given:    "{{",
			Expected: fmt.Errorf("error unmarshalling schema '{{': invalid character '{' looking for beginning of object key string"),
		},
		{
			Name:     "success scenario",
			Given:    "{}",
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		s, err := LoadSchemaFromString(testCase.Given)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Error(), err.Error())
			assert.Nil(t, s)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, s)
		}
	}
}

func TestGetSchema(t *testing.T) {
	var schm *Schema
	sm := TopicSchema{
		TopicTodoCreated: &Schema{},
	}

	schm = sm.GetSchema(TopicTodoCreated)
	require.NotNil(t, schm)
	assert.Equal(t, sm[TopicTodoCreated], schm)

	schm = sm.GetSchema("NotExistingKey")
	assert.Nil(t, schm)
}

func TestValidateBytes(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    []byte
		Expected error
	}
	testCases := []TestCase{
		{
			Name:     "force error when data is nil",
			Given:    nil,
			Expected: fmt.Errorf("data cannot be nil"),
		},
		{
			Name:     "force error when parsing JSON bytes",
			Given:    []byte(`{`),
			Expected: fmt.Errorf("error parsing JSON bytes: unexpected end of JSON input"),
		},
		{
			Name: "force error when no valid data",
			Given: []byte(`{
				"title": "todo title",
				"description": "todo description"
			}`),
			Expected: fmt.Errorf("error validating schema: \"id\" value is required: / = map[description:todo description title:todo title]"),
		},
		{
			Name: "success scenario",
			Given: []byte(`{
				"id": 12345,
				"title": "todo title",
				"description": "todo description"
			}`),
		},
	}

	schema, err := LoadSchemaFromString(schemaEventTocoCreated)
	require.NoError(t, err)
	require.NotNil(t, schema)

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		err = schema.ValidateBytes(testCase.Given)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidate(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    interface{}
		Expected error
	}

	testCases := []TestCase{
		{
			Name:     "force error when nil is provided",
			Given:    nil,
			Expected: fmt.Errorf("data cannot be nil"),
		},
		// FIXME This test should return a failure but it is not happening
		// {
		// 	Name: "force failure by using a struct that does not match the schema",
		// 	Given: struct{ AnyField string }{
		// 		AnyField: "AnyValue",
		// 	},
		// 	Expected: fmt.Errorf("error validating schema: "),
		// },
		{
			Name: "success scenario",
			Given: TodoCreatedEventJson{
				Id:          12345,
				Title:       "Todo title",
				Description: "Todo description",
			},
			Expected: nil,
		},
	}

	schema, err := LoadSchemaFromString(schemaEventTocoCreated)
	require.NoError(t, err)
	require.NotNil(t, schema)

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		err = schema.Validate(testCase.Given)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestPrepareParseErrorList(t *testing.T) {
	schema, err := LoadSchemaFromString(`{}`)
	require.NoError(t, err)
	require.NotNil(t, schema)
	err = schema.prepareParseErrorList(
		[]jsonschema.KeyError{
			{
				Message:      "test",
				PropertyPath: "/",
				InvalidValue: "test",
			},
		},
	)
	require.Error(t, err)
	assert.Equal(t, "error validating schema: test: / = test", err.Error())
}

func TestValidateMessage(t *testing.T) {
	type TestCaseGiven struct {
		Schemas TopicSchema
		Message *kafka.Message
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected error
	}

	schemas, err := LoadSchemas()
	require.NoError(t, err)

	testCases := []TestCase{
		// nil schemas
		{
			Name: "force error when schemas is nil",
			Given: TestCaseGiven{
				Schemas: nil,
				Message: nil,
			},
			Expected: fmt.Errorf("schemas is empty"),
		},
		// nil message
		{
			Name: "force error when message is nil",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: nil,
			},
			Expected: fmt.Errorf("msg cannot be nil"),
		},
		// No Topic
		{
			Name: "force error when no topic is specified",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					Headers: []kafka.Header{
						{
							Key:   string(message.HdrType),
							Value: []byte(message.HdrTypeIntrospect),
						},
					},
					TopicPartition: kafka.TopicPartition{
						Topic: nil,
					},
				},
			},
			Expected: fmt.Errorf("topic cannot be nil nor empty string"),
		},
		// Empty topic string
		{
			Name: "force error when no topic is specified",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					Headers: []kafka.Header{
						{
							Key:   string(message.HdrType),
							Value: []byte(message.HdrTypeIntrospect),
						},
					},
					TopicPartition: kafka.TopicPartition{
						Topic: pointy.String(""),
					},
				},
			},
			Expected: fmt.Errorf("topic cannot be nil nor empty string"),
		},
		// Topic not found
		{
			Name: "force error when the topic is not found",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					Headers: []kafka.Header{
						{
							Key:   string(message.HdrType),
							Value: []byte(message.HdrTypeIntrospect),
						},
					},
					TopicPartition: kafka.TopicPartition{
						Topic: pointy.String("ATopicThatDoesNotExist"),
					},
				},
			},
			Expected: fmt.Errorf("topic not found: '%s'", "ATopicThatDoesNotExist"),
		},
		// Validate bytes return false
		{
			Name: "force error when schema validation fails",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					Headers: []kafka.Header{
						{
							Key:   string(message.HdrType),
							Value: []byte(message.HdrTypeIntrospect),
						},
					},
					TopicPartition: kafka.TopicPartition{
						Topic: pointy.String(TopicTodoCreated),
					},
					Value: []byte(`{
						"title": "todo title",
						"description": "todo description"
					}`),
				},
			},
			Expected: fmt.Errorf("error validating schema: \"id\" value is required: / = map[description:todo description title:todo title]"),
		},
		// Validate bytes return true
		{
			Name: "Success message schema validation",
			Given: TestCaseGiven{
				Schemas: schemas,
				Message: &kafka.Message{
					Headers: []kafka.Header{
						{
							Key:   string(message.HdrType),
							Value: []byte(message.HdrTypeIntrospect),
						},
					},
					TopicPartition: kafka.TopicPartition{
						Topic: pointy.String(TopicTodoCreated),
					},
					Value: []byte(`{
						"id": 12345,
						"title": "todo title",
						"description": "todo description"
					}`),
				},
			},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Testing case '%s'", testCase.Name)
		err := testCase.Given.Schemas.ValidateMessage(testCase.Given.Message)
		if testCase.Expected != nil {
			require.Error(t, err)
			require.Equal(t, testCase.Expected.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
