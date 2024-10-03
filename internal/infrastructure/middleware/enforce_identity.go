package middleware

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
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

// EnforceUserPredicate is a predicate that enforce identity
// is a user and some additional checks for a user identity.
// data is the XRHID to enforce.
// Return nil if the enforce is passed, else details about the
// enforce process.
func EnforceUserPredicate(data *identity.XRHID) error {
	// See: https://github.com/RedHatInsights/identity-schemas/blob/main/3scale/identities/jwt.json
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
	// See: https://github.com/RedHatInsights/identity-schemas/blob/main/3scale/identities/cert.json
	if data == nil {
		return fmt.Errorf("'data' cannot be nil")
	}
	if data.Identity.Type != "System" {
		return fmt.Errorf("'Identity.Type' must be 'System'")
	}
	if data.Identity.AuthType != "cert-auth" {
		return fmt.Errorf("'Identity.AuthType' is not 'cert-auth'")
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

// EnforceServiceAccountPredicate is a predicate that check fields for
// ServiceAccount identities.
// Return nil if the enforce is passed, else details about the
// enforce process.
func EnforceServiceAccountPredicate(data *identity.XRHID) error {
	// See: https://github.com/RedHatInsights/identity-schemas/blob/main/3scale/identities/service-account.json
	if data == nil {
		return fmt.Errorf("'data' cannot be nil")
	}
	if data.Identity.Type != "ServiceAccount" {
		return fmt.Errorf("'Identity.Type' must be 'ServiceAccount'")
	}
	if data.Identity.AuthType != "jwt-auth" {
		return fmt.Errorf("'Identity.AuthType' is not 'jwt-auth'")
	}
	if data.Identity.ServiceAccount == nil {
		return fmt.Errorf("'Identity.ServiceAccount' is nil")
	}
	if data.Identity.ServiceAccount.ClientId == "" {
		return fmt.Errorf("'Identity.ServiceAccount.ClientId' is empty")
	}
	if data.Identity.ServiceAccount.Username == "" {
		return fmt.Errorf("'Identity.ServiceAccount.Username' is empty")
	}
	return nil
}

// NewEnforceOr allow to create new predicates by composing a
// logical OR with existing predicates.
func NewEnforceOr(predicates ...IdentityPredicate) IdentityPredicate {
	return func(data *identity.XRHID) error {
		var allErrors error
		for i := range predicates {
			if err := predicates[i](data); err == nil {
				return nil
			} else {
				allErrors = errors.Join(allErrors, err)
			}
		}
		return allErrors
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
				err error
				cc  DomainContextInterface
				ok  bool
			)
			ctx := c.Request().Context()
			logger := app_context.LogFromCtx(ctx)
			if config.Skipper != nil && config.Skipper(c) {
				logger.Debug("Skipping EnforceIdentity middleware")
				return next(c)
			}
			if cc, ok = c.(DomainContextInterface); !ok {
				logger.Error("'DomainContextInterface' is expected")
				return echo.ErrInternalServerError
			}

			xrhid := cc.XRHID()

			// The predicate must return no error, otherwise
			// the request is not authorised.
			for _, entry := range config.Predicates {
				key := entry.Name
				predicate := entry.Predicate
				if err = predicate(xrhid); err != nil {
					logger.Error(fmt.Sprintf("'%s' IdentityPredicate failed: %s", key, err.Error()))
					return echo.ErrUnauthorized
				}
			}

			return next(c)
		}
	}
}

type ParseXRHIDMiddlewareConfig struct {
	// Skipper function to skip for some request if necessary
	Skipper echo_middleware.Skipper
}

// Parse the X-RH-Identity header and set it into the request context.
// This must be called AFTER the "Fake Identity" middleware (if used),
// but BEFORE the EnforceIdentity middlewares.
func ParseXRHIDMiddlewareWithConfig(config *ParseXRHIDMiddlewareConfig) func(echo.HandlerFunc) echo.HandlerFunc {
	if config == nil {
		panic("'config' is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				xrhid *identity.XRHID
				err   error
				cc    DomainContextInterface
				ok    bool
			)

			ctx := c.Request().Context()
			logger := app_context.LogFromCtx(ctx)

			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}

			if cc, ok = c.(DomainContextInterface); !ok {
				logger.Error("'DomainContextInterface' is expected")
				return echo.ErrInternalServerError
			}

			xrhidRaw := cc.Request().Header.Get(header.HeaderXRHID)
			if xrhidRaw == "" {
				return echo.ErrUnauthorized
			}
			if xrhid, err = header.DecodeXRHID(xrhidRaw); err != nil {
				logger.Error(err.Error())
				return echo.ErrBadRequest
			}

			// Aggregate additional information to the logs
			principal := header.GetPrincipal(xrhid)
			logger = logger.With(
				slog.String("org_id", xrhid.Identity.OrgID),
				slog.String("identity_type", xrhid.Identity.Type),
				slog.String("identity_principal", principal),
			)
			ctx = app_context.CtxWithLog(ctx, logger)
			c.SetRequest(c.Request().Clone(ctx))

			// Set the unserialized Identity into the request context
			cc.SetXRHID(xrhid)
			return next(c)
		}
	}
}
