package hostconf_jwk

// JWK key state
type KeyState int

const (
	ValidKey KeyState = iota
	ExpiredKey
	InvalidKey
	RevokedKey
)
