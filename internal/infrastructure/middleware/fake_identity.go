package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type FakeIdentityConfig struct {
	// Skipper function to skip for some request if necessary
	Skipper echo_middleware.Skipper
}

const headerXRhFakeIdentity = "X-Rh-Fake-Identity"

// FakeIdentityWithConfig middleware copy the x-rh-fake-identity to the
// x-rh-identity header when no skipper return true; it is intended
// to be called before the EnforceIdentity middleware.
func FakeIdentityWithConfig(config *FakeIdentityConfig) func(echo.HandlerFunc) echo.HandlerFunc {
	if config == nil {
		panic("'config' is nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}
			fakeIdentity := c.Request().Header[headerXRhFakeIdentity]
			if fakeIdentity != nil {
				c.Request().Header.Set(headerXRhIdentity, strings.Join(fakeIdentity, "; "))
				c.Request().Header.Del(headerXRhFakeIdentity)
			}
			return next(c)
		}
	}
}
