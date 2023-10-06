package token

// JWK key state
type KeyState int

const (
	ValidKey KeyState = iota
	ExpiredKey
	InvalidKey
	RevokedKey
)

// Token
const (
	TokenIssuer      = "idmsvc/v1"
	AudJoinHost      = "join host"
	ClaimOrgId       = "rhorg"
	ClaimDomainId    = "rhdomid"
	ClaimFqdn        = "rhfqdn"
	ClaimInventoryId = "rhinvid"
)
