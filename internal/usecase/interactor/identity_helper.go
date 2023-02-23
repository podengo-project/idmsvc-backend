package interactor

import (
	"fmt"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

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
