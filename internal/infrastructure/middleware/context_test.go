package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/redhatinsights/platform-go-middlewares/identity"
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
	iden := &identity.Identity{
		OrgID:         "11111",
		Type:          "System",
		AccountNumber: "123",
		System: identity.System{
			CommonName: "a6355a76-c6e8-11ed-8aed-482ae3863d30",
			CertType:   "system",
		},
	}
	ctx.SetIdentity(iden)
	assert.Equal(t, iden, ctx.Identity())
}
