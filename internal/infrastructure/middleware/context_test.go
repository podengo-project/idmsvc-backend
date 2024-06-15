package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := NewContext(e.NewContext(req, rec))
	assert.NotNil(t, ctx)
	_, ok := ctx.(DomainContextInterface)
	assert.Equal(t, true, ok)
}

func TestSetIdentityIdentity(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := NewContext(e.NewContext(req, rec))
	xrhid := &identity.XRHID{
		Identity: identity.Identity{
			OrgID:         "11111",
			Type:          "System",
			AccountNumber: "123",
			System: &identity.System{
				CommonName: "a6355a76-c6e8-11ed-8aed-482ae3863d30",
				CertType:   "system",
			},
		},
	}
	ctx.SetXRHID(xrhid)
	assert.Equal(t, xrhid, ctx.XRHID())
}

func TestOrgFallback(t *testing.T) {
	var xrhid *identity.XRHID

	// Test nil argument
	xrhid = nil
	assert.NotPanics(t, func() {
		orgFallback(xrhid)
	})

	// Test with Identity.OrgID filled
	xrhid = &identity.XRHID{
		Identity: identity.Identity{
			OrgID: "11111",
			Internal: identity.Internal{
				OrgID: "22222",
			},
		},
	}
	orgFallback(xrhid)
	assert.Equal(t, "11111", xrhid.Identity.OrgID)

	// Test with Identity.OrgID empty
	xrhid = &identity.XRHID{
		Identity: identity.Identity{
			OrgID: "",
			Internal: identity.Internal{
				OrgID: "22222",
			},
		},
	}
	orgFallback(xrhid)
	assert.Equal(t, "22222", xrhid.Identity.OrgID)
}
