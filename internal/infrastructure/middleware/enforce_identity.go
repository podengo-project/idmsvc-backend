package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// The Identity header present into the public headers
const headerXRhIdentity = "X-Rh-Identity"

// FIXME Refactor to use the signature: func(c echo.Context) Error
//
//	so that the predicate has information about the http Request
//	context
type IdentityPredicate func(data *identity.XRHID) error

// identityConfig Represent the configuration for this middleware
// enforcement.
type identityConfig struct {
	// Skipper function to skip for some request if necessary
	skipper echo_middleware.Skipper
	// Map of predicates to be applied, all the predicates must
	// return true, if any of them fail, the enforcement will
	// return error for the request.
	predicates map[string]IdentityPredicate
}

var (
	systemEnforceRoutes = []string{
		"/api/hmsidm/v1/domains/:uuid/register",
		"/api/hmsidm/v1/domains/:uuid/update",
	}
	userEnforceRoutes = []string{
		"/api/hmsidm/v1/domains",
		"/api/hmsidm/v1/domains/:uuid",
	}
)

// NewIdentityConfig creates a new identityConfig for the
// EnforcementIdentity middleware.
// Return an identityConfig structure to configure the
// middleware.
func NewIdentityConfig() *identityConfig {
	return &identityConfig{
		predicates: map[string]IdentityPredicate{},
	}
}

// SetSkipper set a skipper function for the middleware.
// skipper is the function which check by using the current
// request context to check if the current request will be
// processed by this middleware.
// Return the identityConfig updated.
func (ic *identityConfig) SetSkipper(skipper echo_middleware.Skipper) *identityConfig {
	ic.skipper = skipper
	return ic
}

// AddPredicate add a predicate function to check the IdentityEnforcement,
// by allowing reuse the same middleware for different enforcements. We can
// add several functions, but if a key collide, the predicate will be overrided.
// key that will be associated to this predicate, it is used to report to the
// log which predicate failed.
// predicate is the check function to be added.
// Return the identityConfig updated.
func (ic *identityConfig) AddPredicate(key string, predicate IdentityPredicate) *identityConfig {
	if predicate != nil {
		ic.predicates[key] = predicate
	}
	return ic
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
		return fmt.Errorf("'Identity.Type' is not 'User'")
	}
	if !data.Identity.User.Active {
		return fmt.Errorf("'Identity.User.Active' is not true")
	}
	if data.Identity.User.UserID == "" {
		return fmt.Errorf("'Identity.User.UserID' cannot be empty")
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
	if data.Identity.System.CertType != "system" {
		return fmt.Errorf("'Identity.System.CertType' is not 'system'")
	}
	if data.Identity.System.CommonName == "" {
		return fmt.Errorf("'Identity.System.CommonName' is empty")
	}
	return nil
}

// SkipperUserPredicate applied when using EnforceUserPredicate.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func SkipperUserPredicate(ctx echo.Context) bool {
	route := ctx.Path()
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range userEnforceRoutes {
		if route == userEnforceRoutes[i] {
			return false
		}
	}
	return true
}

// SkipperSystemPredicate applied when using EnforceSystemPredicate.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func SkipperSystemPredicate(ctx echo.Context) bool {
	route := ctx.Path()
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range systemEnforceRoutes {
		if route == systemEnforceRoutes[i] {
			return false
		}
	}
	return true
}

// EnforceIdentityWithConfig instantiate a EnforceIdentity middleware
// for the configuration provided. This middleware depends on
// NewContext middleware. If the request pass the enforcement
// check, then the unmarshalled version of the identity is stored
// for the request context.
// config is the configuration with the skipper and predicates
// to be used for the middleware.
// Return an echo middleware function.
func EnforceIdentityWithConfig(config *identityConfig) func(echo.HandlerFunc) echo.HandlerFunc {
	if config == nil {
		panic("'config' cannot be nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.skipper != nil && config.skipper(c) {
				return next(c)
			}
			cc, ok := c.(DomainContextInterface)
			if !ok {
				c.Logger().Error("Expected a 'DomainContextInterface'")
				return echo.ErrInternalServerError
			}
			b64XRHID := cc.Request().Header.Get(headerXRhIdentity)
			if b64XRHID == "" {
				cc.Logger().Error("%s not present", headerXRhIdentity)
				return echo.ErrUnauthorized
			}
			stringXRHID, err := base64.StdEncoding.DecodeString(b64XRHID)
			if err != nil {
				cc.Logger().Error(err)
				return echo.ErrUnauthorized
			}
			xrhid := &identity.XRHID{}
			if err := json.Unmarshal([]byte(stringXRHID), xrhid); err != nil {
				cc.Logger().Error(err)
				return echo.ErrUnauthorized
			}

			// All the predicates should return true
			for _, predicate := range config.predicates {
				if err := predicate(xrhid); err != nil {
					if err != nil {
						cc.Logger().Error(err)
					}
					return echo.ErrUnauthorized
				}
			}
			cc.SetXRHID(xrhid)
			return next(c)
		}
	}
}
