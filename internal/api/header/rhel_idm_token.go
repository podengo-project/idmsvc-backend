package header

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

const XRHIDMRHELIDMRegisterToken = "X-Rh-Idm-RhelIdm-Register-Token"

type RhelIdmToken struct {
	Secret     *string    `json:"secret,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
}

// DecodeRhelIdmToken decode a base64 x-rh-idm-rhelidm-token
// data is the base64 coded string
// Returns a RhelIdmToken with the information and nil for
// a success call, else nil and an error filled with all the
// additional information.
func DecodeRhelIdmToken(data string) (*RhelIdmToken, error) {
	if data == "" {
		return nil, fmt.Errorf("'data' is empty")
	}
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	output := &RhelIdmToken{}
	err = json.Unmarshal(bytes, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// EncodeRhelIdmToken encode a RhelIdmToken into base64 representation
// data is the RhelIdmToken reference to be encoded.
// Returns a base64 string representation for the json structure provided
// and nil, else an empty string and an error with additional information.
func EncodeRhelIdmToken(data *RhelIdmToken) (string, error) {
	if data == nil {
		return "", fmt.Errorf("'data' is nil")
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	output := base64.StdEncoding.EncodeToString(bytes)
	return output, nil
}
