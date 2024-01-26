package event

import (
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/event"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewTopicTranslationWithDefaults(t *testing.T) {
	result := config.NewTopicTranslationWithDefaults()
	for _, topic := range event.AllowedTopics {
		assert.Equal(t, topic, result.GetInternal(topic))
		assert.Equal(t, topic, result.GetReal(topic))
	}
}

func TestNewTopicTranslationWithClowder(t *testing.T) {
	var (
		cfg    *clowder.AppConfig
		result *config.TopicTranslation
	)

	// When it is nil, it returns the defaults
	result = config.NewTopicTranslationWithClowder(nil)
	for _, topic := range event.AllowedTopics {
		assert.Equal(t, topic, result.GetInternal(topic))
		assert.Equal(t, topic, result.GetReal(topic))
	}

	// When Kafka is nil, it returns the defaults
	cfg = &clowder.AppConfig{
		Kafka: nil,
	}
	result = config.NewTopicTranslationWithClowder(cfg)
	for _, topic := range event.AllowedTopics {
		assert.Equal(t, topic, result.GetInternal(topic))
		assert.Equal(t, topic, result.GetReal(topic))
	}

	// Try the custom mapping is right
	cfg = &clowder.AppConfig{
		Kafka: &clowder.KafkaConfig{
			Topics: []clowder.TopicConfig{
				{
					RequestedName: "requested.topic.name",
					Name:          "real.topic.name",
				},
			},
		},
	}
	result = config.NewTopicTranslationWithClowder(cfg)
	assert.Equal(t, "real.topic.name", result.GetReal("requested.topic.name"))
	assert.Equal(t, "requested.topic.name", result.GetInternal("real.topic.name"))
}

func TestGetInternal(t *testing.T) {
	tt := config.NewTopicTranslationWithClowder(&clowder.AppConfig{
		Kafka: &clowder.KafkaConfig{
			Topics: []clowder.TopicConfig{
				{
					RequestedName: "requested.topic.name",
					Name:          "real.topic.name",
				},
			},
		},
	})
	assert.Equal(t, "", tt.GetInternal("ItDoesNotExist"))
	assert.Equal(t, "requested.topic.name", tt.GetInternal("real.topic.name"))
}

func TestGetReal(t *testing.T) {
	tt := config.NewTopicTranslationWithClowder(&clowder.AppConfig{
		Kafka: &clowder.KafkaConfig{
			Topics: []clowder.TopicConfig{
				{
					RequestedName: "requested.topic.name",
					Name:          "real.topic.name",
				},
			},
		},
	})
	assert.Equal(t, "", tt.GetReal("ItDoesNotExist"))
	assert.Equal(t, "real.topic.name", tt.GetReal("requested.topic.name"))
}
