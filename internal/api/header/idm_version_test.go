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
	headerValue := `{"ipa-hcc":"0.7","ipa":"4.10.0-8.el9_1","os-release-id":"rhel","os-release-version-id":"9.1"}`
	// Empty string
	assert.Nil(t, NewXRHIDMVersionWithHeader(""))
	// Invalid base64 content
	assert.Nil(t, NewXRHIDMVersionWithHeader("??"))
	// Invalid json
	assert.Nil(t, NewXRHIDMVersionWithHeader("ewo="))
	// Valid header representation
	data := NewXRHIDMVersionWithHeader(headerValue)
	assert.NotNil(t, data)
	assert.Equal(t, "0.7", data.IPAHCCVersion)
	assert.Equal(t, "4.10.0-8.el9_1", data.IPAVersion)
	assert.Equal(t, "rhel", data.OSReleaseID)
	assert.Equal(t, "9.1", data.OSReleaseVersionID)
}

func TestEncodeXRHIDMVersion(t *testing.T) {
	headerValue := `{"ipa-hcc":"0.7","ipa":"4.10.0-8.el9_1","os-release-id":"rhel","os-release-version-id":"9.1"}`
	assert.Equal(t, "", EncodeXRHIDMVersion(nil))
	assert.Equal(t, headerValue, EncodeXRHIDMVersion(NewXRHIDMVersionWithHeader(headerValue)))
}
