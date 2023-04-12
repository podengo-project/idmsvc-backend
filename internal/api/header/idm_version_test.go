package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewXRHIDMVersion(t *testing.T) {
	assert.Nil(t, NewXRHIDMVersion("", ""))
	assert.Nil(t, NewXRHIDMVersion("0.7", ""))
	assert.Nil(t, NewXRHIDMVersion("", "4.10.0-8.el9_1"))

	data := NewXRHIDMVersion("0.7", "4.10.0-8.el9_1")
	assert.NotNil(t, data)
	assert.Equal(t, "0.7", data.IPAHCCVersion)
	assert.Equal(t, "4.10.0-8.el9_1", data.IPAVersion)
}

func TestNewXRHIDMVersionWithHeader(t *testing.T) {
	// Empty string
	assert.Nil(t, NewXRHIDMVersionWithHeader(""))
	// Invalid base64 content
	assert.Nil(t, NewXRHIDMVersionWithHeader("??"))
	// Invalid json
	assert.Nil(t, NewXRHIDMVersionWithHeader("ewo="))
	// Valid header representation
	data := NewXRHIDMVersionWithHeader("eyJpcGEtaGNjIjogIjAuNyIsICJpcGEiOiAiNC4xMC4wLTguZWw5XzEifQo=")
	assert.NotNil(t, data)
	assert.Equal(t, "0.7", data.IPAHCCVersion)
	assert.Equal(t, "4.10.0-8.el9_1", data.IPAVersion)
}

func TestEncodeXRHIDMVersion(t *testing.T) {
	const b64Header = "eyJpcGEtaGNjIjoiMC43IiwiaXBhIjoiNC4xMC4wLTguZWw5XzEifQ=="
	assert.Equal(t, "", EncodeXRHIDMVersion(nil))
	assert.Equal(t, b64Header, EncodeXRHIDMVersion(NewXRHIDMVersionWithHeader(b64Header)))
}
