package header

import (
	"encoding/base64"
	"encoding/json"
)

type XRHIDMVersion struct {
	IPAHCCVersion string `json:"ipa-hcc"`
	IPAVersion    string `json:"ipa"`
}

// NewXRHIDMVersion create a new XRHIDMVersion from the arguments.
// IpaHccVersion is the ipa-hcc client version.
// IpaVersion is the rhel-idm version.
// Return a filled structure or nil if something is wrong.
func NewXRHIDMVersion(IpaHccVersion string, IpaVersion string) *XRHIDMVersion {
	if IpaHccVersion == "" || IpaVersion == "" {
		return nil
	}
	return &XRHIDMVersion{
		IPAHCCVersion: IpaHccVersion,
		IPAVersion:    IpaVersion,
	}
}

// NewXRHIDMVersionWithHeader create a new XRHIDMVersion from the arguments.
// header is the string with the X-RH-IDM-Version content.
// Return a filled structure or nil if something is wrong.
func NewXRHIDMVersionWithHeader(header string) *XRHIDMVersion {
	if header == "" {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return nil
	}
	output := &XRHIDMVersion{}
	if err = json.Unmarshal(data, output); err != nil {
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
	return base64.StdEncoding.EncodeToString(jsonData)
}
