package middleware

import (
	"encoding/base64"
	"encoding/json"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

const headerXRhIdentity = "X-Rh-Identity"

// FIXME Refactor to use the signature: func(c echo.Context) Error
//       so that the predicate has information about the http Request
//       context
type Predicate func(data *identity.Identity) error

type IdentityConfig struct {
	skipper    echo_middleware.Skipper
	predicates map[string]Predicate
}

func NewIdentityConfig(skipper echo_middleware.Skipper) *IdentityConfig {
	return &IdentityConfig{
		skipper:    skipper,
		predicates: map[string]Predicate{},
	}
}

func (ic *IdentityConfig) Add(key string, predicate Predicate) *IdentityConfig {
	if predicate != nil {
		ic.predicates[key] = predicate
	}
	return ic
}

func IdentityAlwaysTrue(data *identity.Identity) error {
	return nil
}

func EnforceIdentityWithConfig(config *IdentityConfig) func(echo.HandlerFunc) echo.HandlerFunc {
	if config == nil {
		panic("config cannot be nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.skipper != nil && config.skipper(c) {
				return next(c)
			}
			b64Identity := c.Request().Header.Get(headerXRhIdentity)
			if b64Identity == "" {
				c.Logger().Error("%s not present", headerXRhIdentity)
				return echo.ErrUnauthorized
			}
			stringIdentity, err := base64.StdEncoding.DecodeString(b64Identity)
			if err != nil {
				c.Logger().Error(err)
				return echo.ErrUnauthorized
			}
			var data identity.Identity
			if err := json.Unmarshal([]byte(stringIdentity), &data); err != nil {
				c.Logger().Error(err)
				return echo.ErrUnauthorized
			}

			for _, predicate := range config.predicates {
				if err := predicate(&data); err != nil {
					if err != nil {
						c.Logger().Error(err)
					}
					return echo.ErrUnauthorized
				}
			}
			return next(c)
		}
	}
}
