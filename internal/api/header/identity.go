package header

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

// https://github.com/RedHatInsights/identity-schemas/blob/main/3scale/schema.json#L13
// https://github.com/RedHatInsights/identity-schemas/blob/main/turnpike/schema.json#L7
const (
	identityTypeUser           = "User"
	identityTypeServiceAccount = "ServiceAccount"
	identityTypeSystem         = "System"
	identityTypeX509           = "X509"
	identityTypeAssociate      = "Associate"
	identityTypeUnknown        = "unknown"
)

// DecodeXRHID from a base64 representation and return
// the unmarshaled value.
// data is the base64 x-rh-identity header representation.
// Return the identity unmarshalled on success and nil, else
// nil and an error.
func DecodeXRHID(data string) (*identity.XRHID, error) {
	if data == "" {
		return nil, fmt.Errorf("'" + HeaderXRHID + "' is an empty string")
	}
	bytes, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	xrhid := &identity.XRHID{}
	if err = json.Unmarshal(bytes, xrhid); err != nil {
		return nil, err
	}
	// topLevelOrgIDFallback
	// See: https://github.com/RedHatInsights/identity/blob/main/identity.go#L164
	if xrhid.Identity.OrgID == "" && xrhid.Identity.Internal.OrgID != "" {
		xrhid.Identity.OrgID = xrhid.Identity.Internal.OrgID
	}
	return xrhid, nil
}

// EncodeXRHID serializes the data in a base64 of the
// json representation.
// data is the Identity struct to be encoded.
// Return empty if the process fails, or the base64
// representation of the identity.Identity provided.
func EncodeXRHID(data *identity.XRHID) string {
	if data == nil {
		return ""
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return b64.StdEncoding.EncodeToString(bytes)
}

// GetPrincipal returns the principal of the identity.
// A principal is, based on type, User.UserID, ServiceAccount.ClientId,
// System.CommonName, X509.SubjectDN, Associate.RHatUUID.
func GetPrincipal(xrhid *identity.XRHID) string {
	principal := ""
	switch xrhid.Identity.Type {
	case identityTypeUser:
		principal = xrhid.Identity.User.UserID
	case identityTypeServiceAccount:
		principal = xrhid.Identity.ServiceAccount.ClientId
	case identityTypeSystem:
		principal = xrhid.Identity.System.CommonName
	case identityTypeX509:
		principal = xrhid.Identity.X509.SubjectDN
	case identityTypeAssociate:
		principal = xrhid.Identity.Associate.RHatUUID
	default:
		principal = identityTypeUnknown
	}
	return principal
}
