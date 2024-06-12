package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	slog "golang.org/x/exp/slog"
)

// FIXME Refactor to use the signature: func(c echo.Context) Error
//
//	so that the predicate has information about the http Request
//	context
type IdentityPredicate func(data *identity.XRHID) error

// IdentityPredicateEntry represents a predicate in the chain of
// responsibility established.
type IdentityPredicateEntry struct {
	Name      string
	Predicate IdentityPredicate
}

// IdentityConfig Represent the configuration for this middleware
// enforcement.
type IdentityConfig struct {
	// Skipper function to skip for some request if necessary
	Skipper echo_middleware.Skipper
	// Map of predicates to be applied, all the predicates must
	// return true, if any of them fail, the enforcement will
	// return error for the request.
	Predicates []IdentityPredicateEntry
}

// IdentityAlwaysTrue is a predicate that always return nil
// so everything was ok.
// data is the reference to the identity.Identity data.
// Return nil on success or an error with additional
// information about the predicate failure.
func IdentityAlwaysTrue(data *identity.XRHID) error {
	return nil
}

// EnforceUserPredicate is a predicate that enforce identity
// is a user and some additional checks for a user identity.
// data is the XRHID to enforce.
// Return nil if the enforce is passed, else details about the
// enforce process.
func EnforceUserPredicate(data *identity.XRHID) error {
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/basic.json
	if data == nil {
		return fmt.Errorf("'data' cannot be nil")
	}
	if data.Identity.Type != "User" {
		return fmt.Errorf("'Identity.Type=%s' is not 'User'", data.Identity.Type)
	}
	if data.Identity.User == nil {
		return fmt.Errorf("'Identity.User' is nil")
	}
	if !data.Identity.User.Active {
		return fmt.Errorf("'Identity.User.Active' is not true")
	}
	if data.Identity.User.Username == "" {
		return fmt.Errorf("'Identity.User.Username' cannot be empty")
	}
	return nil
}

// EnforceSystemPredicate is a predicate that enforce identity
// is a system and some additional checks for a user identity.
// data is the XRHID to enforce.
// Return nil if the enforce is passed, else details about the
// enforce process.
func EnforceSystemPredicate(data *identity.XRHID) error {
	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/cert.json
	if data == nil {
		return fmt.Errorf("'data' cannot be nil")
	}
	if data.Identity.Type != "System" {
		return fmt.Errorf("'Identity.Type' must be 'System'")
	}
	if data.Identity.System == nil {
		return fmt.Errorf("'Identity.System' is nil")
	}
	if data.Identity.System.CertType != "system" {
		return fmt.Errorf("'Identity.System.CertType' is not 'system'")
	}
	if data.Identity.System.CommonName == "" {
		return fmt.Errorf("'Identity.System.CommonName' is empty")
	}
	return nil
}

// NewEnforceOr allow to create new predicates by composing a
// logical OR with existing predicates.
func NewEnforceOr(predicates ...IdentityPredicate) IdentityPredicate {
	return func(data *identity.XRHID) error {
		var firsterr error
		for i := range predicates {
			if err := predicates[i](data); err == nil {
				return nil
			} else {
				if firsterr == nil {
					firsterr = err
				}
			}
		}
		return firsterr
	}
}

// EnforceIdentityWithConfig instantiate a EnforceIdentity middleware
// for the configuration provided. This middleware depends on
// NewContext middleware. If the request pass the enforcement
// check, then the unmarshalled version of the identity is stored
// for the request context.
// config is the configuration with the skipper and predicates
// to be used for the middleware.
// Return an echo middleware function.
func EnforceIdentityWithConfig(config *IdentityConfig) func(echo.HandlerFunc) echo.HandlerFunc {
	if config == nil {
		panic("'config' is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				xrhid *identity.XRHID
				err   error
			)
			if config.Skipper != nil && config.Skipper(c) {
				slog.DebugContext(c.Request().Context(), "Skipping EnforceIdentity middleware")
				return next(c)
			}
			cc, ok := c.(DomainContextInterface)
			if !ok {
				slog.ErrorContext(c.Request().Context(), "'DomainContextInterface' is expected")
				return echo.ErrInternalServerError
			}
			xrhidRaw := cc.Request().Header.Get(header.HeaderXRHID)
			if xrhid, err = decodeXRHID(xrhidRaw); err != nil {
				slog.ErrorContext(c.Request().Context(), err.Error())
				return echo.ErrBadRequest
			}

			// The predicate must return no error, otherwise
			// the request is not authorised.
			for _, entry := range config.Predicates {
				key := entry.Name
				predicate := entry.Predicate
				if err = predicate(xrhid); err != nil {
					slog.ErrorContext(
						c.Request().Context(),
						fmt.Sprintf("'%s' IdentityPredicate failed: %s", key, err.Error()),
					)
					return echo.ErrUnauthorized
				}
			}

			// Set the unserialized Identity into the request context
			cc.SetXRHID(xrhid)
			return next(c)
		}
	}
}

func decodeXRHID(b64XRHID string) (*identity.XRHID, error) {
	if b64XRHID == "" {
		return nil, fmt.Errorf("%s not present", header.HeaderXRHID)
	}
	stringXRHID, err := base64.StdEncoding.DecodeString(b64XRHID)
	if err != nil {
		return nil, err
	}
	xrhid := &identity.XRHID{}
	if err := json.Unmarshal([]byte(stringXRHID), xrhid); err != nil {
		return nil, err
	}
	return xrhid, nil
}
