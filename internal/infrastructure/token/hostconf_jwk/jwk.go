package hostconf_jwk

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

const KeyCurve = jwa.P256

// Generate a private key with additional properties
// alg: based on key type (ES256 for P256)
// exp: expiration time (Unix timestamp)
// kid: base64 SHA-256 thumbprint
// use: "sig"
func GeneratePrivateJWK(expiration time.Time) (key jwk.Key, err error) {
	var (
		crv elliptic.Curve
		alg jwa.SignatureAlgorithm
		raw *ecdsa.PrivateKey
	)

	switch KeyCurve {
	case jwa.P256:
		crv = elliptic.P256()
		alg = jwa.ES256
	default:
		return nil, fmt.Errorf("Unsupported JWK curve %s", KeyCurve)
	}

	raw, err = ecdsa.GenerateKey(crv, rand.Reader)
	if err != nil {
		return nil, err
	}

	key, err = jwk.FromRaw(raw)
	if err != nil {
		return nil, err
	}

	// set kid to truncated SHA-256 thumbprint (RFC 7638)
	tp, err := key.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	if err = key.Set(jwk.KeyIDKey, base64.RawURLEncoding.EncodeToString(tp)[:8]); err != nil {
		return nil, err
	}

	// P-256 key is used for signing with ES256
	if err = key.Set(jwk.KeyUsageKey, jwk.ForSignature); err != nil {
		return nil, err
	}
	if err = key.Set(jwk.AlgorithmKey, alg); err != nil {
		return nil, err
	}

	// non-standard but common expiration for key
	if err = key.Set("exp", expiration.Unix()); err != nil {
		return nil, err
	}

	return key, nil
}

// Get public key of a JWK
func GetPublicJWK(key jwk.Key) (jwk.Key, error) {
	return key.PublicKey()
}

// Parse and validate a single JWK
func ParseJWK(src []byte) (key jwk.Key, state KeyState, err error) {
	key, err = jwk.ParseKey(src)
	if err != nil {
		return nil, InvalidKey, err
	}
	state, err = checkJWK(key)
	if state == ValidKey {
		return key, state, nil
	} else {
		return nil, state, err
	}
}

// Parse and validate a JWKSet
func ParseJWKSet(src []byte) (rs jwk.Set, err error) {
	s, err := jwk.Parse(src)
	if err != nil {
		return nil, err
	}
	rs = jwk.NewSet()
	for i := 0; i < s.Len(); i++ {
		key, _ := s.Key(i)
		state, err := checkJWK((key))
		switch state {
		case ValidKey:
			if err = rs.AddKey(key); err != nil {
				return nil, err
			}
		case ExpiredKey:
			// skip expired key
			continue
		case InvalidKey:
			return nil, fmt.Errorf("Invalid key %d: %v", i, err)
		}
	}
	return rs, nil
}

// Verify a JWK and check that it matches our requirements
func checkJWK(key jwk.Key) (KeyState, error) {
	if key.KeyType() != jwa.EC {
		return InvalidKey, fmt.Errorf("Invalid key type %s", key.KeyType())
	}

	// Crv
	switch raw := key.(type) {
	case jwk.ECDSAPrivateKey:
		if raw.Crv() != jwa.P256 {
			return InvalidKey, fmt.Errorf("Invalid curve %s", raw.Crv().String())
		}
	case jwk.ECDSAPublicKey:
		if raw.Crv() != jwa.P256 {
			return InvalidKey, fmt.Errorf("Invalid curve %s", raw.Crv().String())
		}
	default:
		return InvalidKey, fmt.Errorf("Invalid key")
	}

	if key.KeyID() == "" {
		return InvalidKey, fmt.Errorf("KeyID is empty")
	}
	if key.KeyUsage() != jwk.ForSignature.String() {
		return InvalidKey, fmt.Errorf("Invalid key usage %s", key.KeyUsage())
	}
	if key.Algorithm() != jwa.ES256 {
		return InvalidKey, fmt.Errorf("Invalid key alg %s", key.Algorithm().String())
	}

	expif, ok := key.Get("exp")
	if !ok {
		return InvalidKey, fmt.Errorf("Missing or invalid 'exp'")
	}
	exp, ok := expif.(int64)
	if !ok {
		return InvalidKey, fmt.Errorf("Missing or invalid 'exp'")
	}
	if exp <= time.Now().Unix() {
		return ExpiredKey, fmt.Errorf("Key has expired")
	}

	return ValidKey, nil
}

func init() {
	var exp int64
	jwk.RegisterCustomField("exp", exp)
}
