package header

import (
	"encoding/json"
)

type XRHIDMVersion struct {
	IPAHCCVersion      string `json:"ipa-hcc"`
	IPAVersion         string `json:"ipa"`
	OSReleaseID        string `json:"os-release-id"`
	OSReleaseVersionID string `json:"os-release-version-id"`
}

// NewXRHIDMVersion create a new XRHIDMVersion from the arguments.
// IpaHccVersion is the ipa-hcc client version.
// IpaVersion is the rhel-idm version.
// OSReleaseID identifies the operating
// system. See: https://www.freedesktop.org/software/systemd/man/os-release.html#ID=
// OSReleaseVersionID identifies the operating system
// version. See: https://www.freedesktop.org/software/systemd/man/os-release.html#VERSION_ID=
// Return a filled structure or nil if something is wrong.
func NewXRHIDMVersion(IpaHccVersion string, IpaVersion string, OSReleaseID string, OSReleaseVersionID string) *XRHIDMVersion {
	if IpaHccVersion == "" || IpaVersion == "" || OSReleaseID == "" || OSReleaseVersionID == "" {
		return nil
	}
	return &XRHIDMVersion{
		IPAHCCVersion:      IpaHccVersion,
		IPAVersion:         IpaVersion,
		OSReleaseID:        OSReleaseID,
		OSReleaseVersionID: OSReleaseVersionID,
	}
}

// NewXRHIDMVersionWithHeader create a new XRHIDMVersion from the arguments.
// header is the string with the X-RH-IDM-Version content.
// Return a filled structure or nil if something is wrong.
func NewXRHIDMVersionWithHeader(header string) *XRHIDMVersion {
	if header == "" {
		return nil
	}
	output := &XRHIDMVersion{}
	if err := json.Unmarshal([]byte(header), output); err != nil {
		return nil
	}
	return output
}

// EncodeXRHIDMVersion encode a base64 x-rh-idm-version header value
// from a XRHIDMVersion.
// data is the reference to the XRHIDMVersion information.
// Return the base64 encoded header value.
func EncodeXRHIDMVersion(data *XRHIDMVersion) string {
	if data == nil {
		return ""
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
