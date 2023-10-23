package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicTranslationGetInternal(t *testing.T) {
	data := &TopicTranslation{
		internalToReal: map[string]string{
			"internal": "real",
		},
		realToInternal: map[string]string{
			"real": "internal",
		},
	}
	assert.Equal(t, "", data.GetInternal("nonexisting"))
	assert.Equal(t, "internal", data.GetInternal("real"))
	data = nil
	assert.Equal(t, "", data.GetInternal("real"))
}

func TestTopicTranslationGetReal(t *testing.T) {
	data := &TopicTranslation{
		internalToReal: map[string]string{
			"internal": "real",
		},
		realToInternal: map[string]string{
			"real": "internal",
		},
	}
	assert.Equal(t, "", data.GetReal("nonexisting"))
	assert.Equal(t, "real", data.GetReal("internal"))
	data = nil
	assert.Equal(t, "", data.GetReal("internal"))
}
