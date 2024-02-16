package hostconf_jwk

// JWK key state
type KeyState int

const (
	ValidKey KeyState = iota
	ExpiredKey
	InvalidKey
	RevokedKey
	EncryptionIdMismatch
	KeyDecryptionFailed
)

func KeyStateString(state KeyState) string {
	switch state {
	case ValidKey:
		return "valid"
	case ExpiredKey:
		return "expired"
	case InvalidKey:
		return "invalid"
	case RevokedKey:
		return "revoked"
	case EncryptionIdMismatch:
		return "encryptionIdMismatch"
	case KeyDecryptionFailed:
		return "keyDecryptionFailed"
	default:
		return "<unknown>"
	}
}
