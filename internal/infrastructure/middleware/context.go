package middleware

import (
	"github.com/labstack/echo/v4"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
)

// Represent the custom context for our backend service
type domainContext struct {
	echo.Context
	xrhid *identity.XRHID
}

// Define the interface for our custom context.
type DomainContextInterface interface {
	echo.Context
	SetXRHID(iden *identity.XRHID)
	XRHID() *identity.XRHID
}

// NewContext create our custom context
// Return an initialized
func NewContext(c echo.Context) DomainContextInterface {
	return &domainContext{
		Context: c,
		xrhid:   nil,
	}
}

func orgFallback(data *identity.XRHID) {
	// See: https://github.com/RedHatInsights/identity/blob/main/identity.go#L164
	if data != nil && data.Identity.OrgID == "" && data.Identity.Internal.OrgID != "" {
		data.Identity.OrgID = data.Identity.Internal.OrgID
	}
}

// SetXRHID set the unmarshalled identity to the context
// so it can be retrieved without repeating all the operations
// to parse it.
// iden is a reference to the identity.Identity structure to
// store into the context. If it is nil, nothing is made.
// Return the DomainContext updated.
func (c *domainContext) SetXRHID(xrhid *identity.XRHID) {
	if xrhid != nil {
		orgFallback(xrhid)
		c.xrhid = xrhid
	}
}

// XRHID retrieve the unserialized identity header from the request context
// Return the reference to the identity.Identity from the request context.
func (c *domainContext) XRHID() *identity.XRHID {
	return c.xrhid
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
