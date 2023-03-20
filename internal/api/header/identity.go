package header

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// DecodeIdentity from a base64 representation and return
// the unmarshaled value.
// data is the base64 x-rh-identity header representation.
// Return the identity unmarshalled on success and nil, else
// nil and an error.
func DecodeIdentity(data string) (*identity.Identity, error) {
	if data == "" {
		return nil, fmt.Errorf("X-Rh-Identity content cannot be an empty string")
	}
	bytes, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	identity := &identity.Identity{}
	if err = json.Unmarshal(bytes, identity); err != nil {
		return nil, err
	}
	// topLevelOrgIDFallback
	// See: https://github.com/RedHatInsights/identity/blob/main/identity.go#L164
	if identity.OrgID == "" && identity.Internal.OrgID != "" {
		identity.OrgID = identity.Internal.OrgID
	}
	return identity, nil
}

// EncodeIdentity serializes the data in a base64 of the
// json representation.
// data is the Identity struct to be encoded.
// Return empty if the process fails, or the base64
// representation of the identity.Identity provided.
func EncodeIdentity(data *identity.Identity) string {
	if data == nil {
		return ""
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return b64.StdEncoding.EncodeToString(bytes)
}
