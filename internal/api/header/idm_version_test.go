package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewXRHIDMVersion(t *testing.T) {
	assert.Nil(t, NewXRHIDMVersion("", "", "", ""))
	assert.Nil(t, NewXRHIDMVersion("0.7", "", "", ""))
	assert.Nil(t, NewXRHIDMVersion("0.7", "4.10.0-8.el9_1", "", ""))
	assert.Nil(t, NewXRHIDMVersion("0.7", "4.10.0-8.el9_1", "rhel", ""))

	data := NewXRHIDMVersion("0.7", "4.10.0-8.el9_1", "rhel", "9.1")
	assert.NotNil(t, data)
	assert.Equal(t, "0.7", data.IPAHCCVersion)
	assert.Equal(t, "4.10.0-8.el9_1", data.IPAVersion)
	assert.Equal(t, "rhel", data.OSReleaseID)
	assert.Equal(t, "9.1", data.OSReleaseVersionID)
}

func TestNewXRHIDMVersionWithHeader(t *testing.T) {
	// Empty string
	assert.Nil(t, NewXRHIDMVersionWithHeader(""))
	// Invalid base64 content
	assert.Nil(t, NewXRHIDMVersionWithHeader("??"))
	// Invalid json
	assert.Nil(t, NewXRHIDMVersionWithHeader("ewo="))
	// Valid header representation
	data := NewXRHIDMVersionWithHeader("eyJpcGEtaGNjIjoiMC43IiwiaXBhIjoiNC4xMC4wLTguZWw5XzEiLCJvcy1yZWxlYXNlLWlkIjoicmhlbCIsIm9zLXJlbGVhc2UtdmVyc2lvbi1pZCI6IjkuMSJ9Cg==")
	assert.NotNil(t, data)
	assert.Equal(t, "0.7", data.IPAHCCVersion)
	assert.Equal(t, "4.10.0-8.el9_1", data.IPAVersion)
	assert.Equal(t, "rhel", data.OSReleaseID)
	assert.Equal(t, "9.1", data.OSReleaseVersionID)
}

func TestEncodeXRHIDMVersion(t *testing.T) {
	const b64Header = "eyJpcGEtaGNjIjoiMC43IiwiaXBhIjoiNC4xMC4wLTguZWw5XzEiLCJvcy1yZWxlYXNlLWlkIjoicmhlbCIsIm9zLXJlbGVhc2UtdmVyc2lvbi1pZCI6IjkuMSJ9"
	assert.Equal(t, "", EncodeXRHIDMVersion(nil))
	assert.Equal(t, b64Header, EncodeXRHIDMVersion(NewXRHIDMVersionWithHeader(b64Header)))
}
