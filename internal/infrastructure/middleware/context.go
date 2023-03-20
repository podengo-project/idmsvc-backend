package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// Represent the custom context for our backend service
type domainContext struct {
	echo.Context
	identity *identity.Identity
}

// Define the interface for our custom context.
type DomainContextInterface interface {
	echo.Context
	SetIdentity(iden *identity.Identity)
	Identity() *identity.Identity
}

// NewContext create our custom context
func NewContext(c echo.Context) DomainContextInterface {
	return &domainContext{
		Context:  c,
		identity: nil,
	}
}

// SetIdentity set the unmarshalled identity to the context
// so it can be retrieved without repeating all the operations
// to parse it.
// iden is a reference to the identity.Identity structure to
// store into the context. If it is nil, nothing is made.
// Return the DomainContext updated.
func (c *domainContext) SetIdentity(iden *identity.Identity) {
	if iden != nil {
		c.identity = iden
	}
}

// Identity retrieve the unserialized identity header from the request context
// Return the reference to the identity.Identity from the request context.
func (c *domainContext) Identity() *identity.Identity {
	return c.identity
}

// CreateContext is a middleware that create the IDM context
// used for the rest of the chain of actions into the request.
// Return the middleware that create the context.
func CreateContext() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &domainContext{Context: c}
			return next(cc)
		}
	}
}
