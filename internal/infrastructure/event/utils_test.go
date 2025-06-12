package event

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/podengo-project/idmsvc-backend/internal/api/event"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/event/message"
	"github.com/stretchr/testify/assert"
	"go.openly.dev/pointy"
)

func TestLogEventMessageInfo(t *testing.T) {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	opts := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Do not print date/time or we can't verfiy
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue("NODATE")
			}

			return a
		},
	}

	h := slog.NewJSONHandler(bufWriter, &opts)
	slog.SetDefault(slog.New(h))

	buf.Reset()
	logEventMessageInfo(nil, "")
	bufWriter.Flush()
	assert.Equal(t, "", buf.String())

	buf.Reset()
	logEventMessageInfo(
		&kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   string(message.HdrType),
					Value: []byte(message.HdrTypeIntrospect),
				},
			},
			TopicPartition: kafka.TopicPartition{
				Topic: pointy.String(event.TopicTodoCreated),
			},
			Value: []byte(`{"uuid":"5e759032-5124-11ed-a029-482ae3863d30","url":"https://example.test"}`),
		}, "")
	bufWriter.Flush()
	assert.Equal(t, "", buf.String())

	buf.Reset()
	logEventMessageInfo(nil, "Any additional message")
	bufWriter.Flush()
	assert.Equal(t, "", buf.String())

	buf.Reset()
	logEventMessageInfo(
		&kafka.Message{
			Headers: []kafka.Header{
				// {
				// 	Key:   string(message.HdrType),
				// 	Value: []byte(message.HdrTypeIntrospect),
				// },
			},
			TopicPartition: kafka.TopicPartition{
				Topic: pointy.String(event.TopicTodoCreated),
			},
			Value: []byte(`{"uuid":"c810053e-512b-11ed-9d9c-482ae3863d30","url":"https://example.test"}`),
		},
		"Some message",
	)
	bufWriter.Flush()
	assert.Equal(t, "{\"time\":\"NODATE\",\"level\":\"INFO\",\"msg\":\"Some message\",\"EventInfoMessage\":{\"Topic\":\"platform.idmsvc.todo-created\",\"Key\":\"\",\"Headers\":\"{}\"}}\n", buf.String())
}

func TestLogEventMessageError(t *testing.T) {
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	opts := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Do not print date/time or we can't verfiy
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue("NODATE")
			}

			return a
		},
	}

	h := slog.NewJSONHandler(bufWriter, &opts)
	slog.SetDefault(slog.New(h))
	logEventMessageError(nil, nil)
	assert.Equal(t, "", buf.String())

	buf.Reset()
	logEventMessageError(
		&kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   string(message.HdrType),
					Value: []byte(message.HdrTypeIntrospect),
				},
			},
			TopicPartition: kafka.TopicPartition{
				Topic: pointy.String(event.TopicTodoCreated),
			},
			Value: []byte(`{"uuid":"28304ad2-512c-11ed-bd7e-482ae3863d30","url":"https://example.test"}`),
		}, nil)
	bufWriter.Flush()
	assert.Equal(t, "", buf.String())

	buf.Reset()
	logEventMessageError(nil, fmt.Errorf("Any error message"))
	bufWriter.Flush()
	assert.Equal(t, "", buf.String())

	buf.Reset()
	logEventMessageError(
		&kafka.Message{
			Headers: []kafka.Header{
				{
					Key:   string(message.HdrType),
					Value: []byte(message.HdrTypeIntrospect),
				},
			},
			TopicPartition: kafka.TopicPartition{
				Topic: pointy.String(event.TopicTodoCreated),
			},
			Value: []byte(`{"uuid":"28304ad2-512c-11ed-bd7e-482ae3863d30","url":"https://example.test"}`),
		},
		fmt.Errorf("Any error message"),
	)
	bufWriter.Flush()
	assert.Equal(t, "{\"time\":\"NODATE\",\"level\":\"ERROR\",\"msg\":\"Error processing event message\",\"EventErrorMessage\":{\"Headers\":[{\"Key\":\"Type\",\"Value\":\"SW50cm9zcGVjdA==\"}],\"Payload\":\"{\\\"uuid\\\":\\\"28304ad2-512c-11ed-bd7e-482ae3863d30\\\",\\\"url\\\":\\\"https://example.test\\\"}\",\"Error\":\"Any error message\"}}\n", buf.String())
}

func TestGetHeaderString(t *testing.T) {
	// func getHeaderString(headers []kafka.Header) string {
	// 	var output []string = make([]string, len(headers))
	// 	for i, header := range headers {
	// 		output[i] = fmt.Sprintf("%s: %s", header.Key, string(header.Value))
	// 	}
	// 	return fmt.Sprintf("{%s}", strings.Join(output, "\n"))
	// }
	var result string

	headers := [][]kafka.Header{
		{},
		{
			kafka.Header{
				Key:   "Header1",
				Value: []byte("Value1"),
			},
		},
		{
			kafka.Header{
				Key:   "Header1",
				Value: []byte("Value1"),
			},
			kafka.Header{
				Key:   "Header2",
				Value: []byte("Value2"),
			},
		},
	}

	result = getHeaderString(headers[0])
	assert.Equal(t, "{}", result)

	result = getHeaderString(headers[1])
	assert.Equal(t, "{Header1: Value1}", result)

	result = getHeaderString(headers[2])
	assert.Equal(t, "{Header1: Value1, Header2: Value2}", result)
}
