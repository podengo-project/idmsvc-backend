package event

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/exp/slog"
)

func getHeaderString(headers []kafka.Header) string {
	var output []string = make([]string, len(headers))
	for i, header := range headers {
		output[i] = fmt.Sprintf("%s: %s", header.Key, string(header.Value))
	}
	return fmt.Sprintf("{%s}", strings.Join(output, ", "))
}

func logEventMessageInfo(msg *kafka.Message, text string) {
	if msg == nil || text == "" {
		return
	}
	slog.Info(
		text,
		slog.Group("EventInfoMessage",
			slog.String("Topic", *msg.TopicPartition.Topic),
			slog.String("Key", string(msg.Key)),
			slog.String("Headers", getHeaderString(msg.Headers)),
		),
	)
}

func logEventMessageError(msg *kafka.Message, err error) {
	if msg == nil || err == nil {
		return
	}

	slog.Error(
		"Error processing event message",
		slog.Group("EventErrorMessage",
			slog.Any("Headers", msg.Headers),
			slog.String("Payload", string(msg.Value)),
			slog.String("Error", err.Error()),
		),
	)
}
