package event

// Tools related with schema validation
// http://json-schema.org
// http://json-schema.org/latest/json-schema-core.html
// http://json-schema.org/latest/json-schema-validation.html
//
// Fancy online tools
// https://www.liquid-technologies.com/online-json-to-schema-converter
// https://app.quicktype.io/
//
// To generate the code from the schemas is used:
// https://github.com/atombender/go-jsonschema
//
// To validate the schemas against a data structure is
// used: https://github.com/qri-io/jsonschema
//
// Regular expression tools, useful when pattern attribute could be used:
// https://www.regexpal.com
// https://regex101.com/
//
// Just to mention that 'pattern' does not work into the validation.
// The regular expressions are different in ECMA Script and GoLang,
// but maybe this library could make the differences work:
// https://github.com/dlclark/regexp2#ecmascript-compatibility-mode
// anyway that would be something to be added as a PR to some of the
// above libraries; however, into the above library it is mentioned
// that it deal with ASCII, but not very well with unicode. That
// is a concern, more when using for message validation that could
// be a source of bugs and vulnerabilities.
//

// One topic has one schema assocaited (1-1)

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/qri-io/jsonschema"
)

const (
	// Topic constants
	TopicTodoCreated = "platform.idmsvc.todo-created"
)

// FIXME Refactor this to make it more dynamic and reduce work for the developer
var AllowedTopics = []string{
	TopicTodoCreated,
	// TODO Add here new topics
}

// https://pkg.go.dev/embed

// Embed message schemas here

//go:embed "todo_created.event.json"
var schemaEventTocoCreated string

// TODO Embed here new event schema string contents

var (
	schemaKey2JsonSpec = map[string]string{
		TopicTodoCreated: schemaEventTocoCreated,
		// TODO Add here new event schemas
	}
)

type Schema jsonschema.Schema

type TopicSchema map[string](*Schema)

// GetSchemaMap return a SchemaMap associated to one topic.
// topic the topic which want to retrieve the SchemaMap.
// Return a SchemaMap associated to the topic or nil if the
// topic is not found.
func (ts *TopicSchema) GetSchema(topic string) *Schema {
	if value, ok := (*ts)[topic]; ok {
		return value
	}
	return nil
}

// ValidateMessage check the msg is accomplish the schema defined for it.
// msg is a reference to a kafka.Message struct.
// Return nil if the check is success else an error reference is filled.
func (ts *TopicSchema) ValidateMessage(msg *kafka.Message) error {
	var (
		s       *Schema
		schemas map[string]*Schema
	)
	schemas = *ts
	if len(schemas) == 0 {
		return fmt.Errorf("schemas is empty")
	}
	if msg == nil {
		return fmt.Errorf("msg cannot be nil")
	}
	if msg.TopicPartition.Topic == nil || *msg.TopicPartition.Topic == "" {
		return fmt.Errorf("topic cannot be nil nor empty string")
	}
	topic := *msg.TopicPartition.Topic
	if s = ts.GetSchema(topic); s == nil {
		return fmt.Errorf("topic not found: '%s'", topic)
	}

	return s.ValidateBytes(msg.Value)
}

// LoadSchemas unmarshall all the embedded schemas and
// return all them in the output schemas variable.
// See also LoadSchemaFromString.
// schemas is a hashmap map[string]*gojsonschema.Schema that
// can be used to immediately validate schemas against
// unmarshalled schemas.
// Return the resulting list of schemas, or nil if an
// an error happens.
func LoadSchemas() (TopicSchema, error) {
	var (
		output TopicSchema = TopicSchema{}
		s      *Schema
		err    error
	)

	for topic, schema := range schemaKey2JsonSpec {
		if s, err = LoadSchemaFromString(schema); err != nil {
			// FIXME Refactor to cover this path: schemaKey2JsonSpec should be injected to the function
			return nil, fmt.Errorf("error unmarshalling for topic '%s': %w", topic, err)
		}
		output[topic] = s
	}

	return output, nil
}

// LoadSchemaFromString unmarshall a schema from
// its string representation in json format.
// schemas is a string representation in json format
// for gojsonschema.Schema.
// Return the resulting list of schemas, or nil if an
// an error happens.
func LoadSchemaFromString(schema string) (*Schema, error) {
	var err error
	var output *Schema
	rs := &jsonschema.Schema{}
	if err = json.Unmarshal([]byte(schema), rs); err != nil {
		return nil, fmt.Errorf("error unmarshalling schema '%s': %w", schema, err)
	}
	output = (*Schema)(rs)
	return output, nil
}

// ValidateBytes validate that a slice of bytes which
// represent an event message match the Schema.
// data is a byte slice with the event message representation.
// Return nil if check is success, else a filled error.
func (s *Schema) ValidateBytes(data []byte) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}
	jsSchema := (*jsonschema.Schema)(s)
	parseErrs, err := jsSchema.ValidateBytes(context.Background(), data)
	if err != nil {
		return err
	}
	if len(parseErrs) > 0 {
		return s.prepareParseErrorList(parseErrs)
	}
	return nil
}

// Validate check that data interface accomplish the Schema.
// data is any type, it cannot be nil.
// Return nil if the check is success, else a filled error.
func (s *Schema) Validate(data interface{}) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}
	jsSchema := (*jsonschema.Schema)(s)
	vs := jsSchema.Validate(context.Background(), data)
	if len(*vs.Errs) > 0 {
		return s.prepareParseErrorList(*vs.Errs)
	}
	return nil
}

func (s *Schema) prepareParseErrorList(parseErrs []jsonschema.KeyError) error {
	var errorList []string = []string{}
	for _, item := range parseErrs {
		errorList = append(errorList, fmt.Sprintf(
			"%s: %s = %s",
			item.Message,
			item.PropertyPath,
			item.InvalidValue,
		))
	}
	return fmt.Errorf("error validating schema: %s", strings.Join(errorList, ", "))
}
