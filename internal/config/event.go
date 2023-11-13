package config

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/event"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"golang.org/x/exp/slog"
)

// TopicMap is used to map between real and internal topics, this is
// it could be that the name we indicate for the topics into the
// clowderapp resource be different from the real created in kafka,
// so this type allow to preproce the mappings, and use them when
// needed to translate them into the producer and consumer functions
type TopicTranslation struct {
	internalToReal map[string]string
	realToInternal map[string]string
}

// GetInternal translate the name of a real topic to
// the internal topic name. This will be used by the
// consumers.
func (tm *TopicTranslation) GetInternal(realTopic string) string {
	if tm == nil {
		return ""
	}
	if val, ok := tm.realToInternal[realTopic]; ok {
		return val
	}
	return ""
}

// GetReal translate the name of an internal topic
// to the real topic name. This will be used by the
// producers.
// Returns empty string when the topic is not found
// into the translation map.
func (tm *TopicTranslation) GetReal(internalTopic string) string {
	if tm == nil {
		return ""
	}
	if val, ok := tm.internalToReal[internalTopic]; ok {
		return val
	}
	return ""
}

// It store the mapping between the internal topic managed by
// the service and the real topic managed by kafka
var TopicTranslationConfig *TopicTranslation = nil

// NewDefaultTopicMap Build a default topic map that map
// all the allowed topics to itselfs
// Return A TopicMap initialized as default values
func NewTopicTranslationWithDefaults() *TopicTranslation {
	var tm *TopicTranslation = &TopicTranslation{
		internalToReal: make(map[string]string),
		realToInternal: make(map[string]string),
	}
	for _, topic := range event.AllowedTopics {
		tm.internalToReal[topic] = topic
		tm.realToInternal[topic] = topic
	}
	return tm
}

// NewTopicTranslationWithClowder Build a topic map based into the
// clowder configuration.
func NewTopicTranslationWithClowder(cfg *clowder.AppConfig) *TopicTranslation {
	if cfg == nil {
		return NewTopicTranslationWithDefaults()
	}

	var tm *TopicTranslation = &TopicTranslation{
		internalToReal: make(map[string]string),
		realToInternal: make(map[string]string),
	}
	for _, topic := range cfg.Kafka.Topics {
		tm.internalToReal[topic.RequestedName] = topic.Name
		tm.realToInternal[topic.Name] = topic.RequestedName
		slog.Debug("internalToReal",
			slog.String(topic.RequestedName, topic.Name))
		slog.Debug("realToInternal",
			slog.String(topic.Name, topic.RequestedName))
	}
	return tm
}
