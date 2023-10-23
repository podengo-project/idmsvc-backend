package middleware

import "github.com/labstack/echo/v4"

// DefaultNooperation is a default instance for Nooperation middleware
var DefaultNooperation = Nooperation()

// Nooperation is a middleware that do nothing. This is useful to decouple
// middleware initialisation from middleware wiring, so if some middleware
// is option based on some configuration, we only have to assign this
// middleware instead of add middlewares in a conditional way.
func Nooperation() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
