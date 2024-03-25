package impl

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	cfg := helperConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, mockRbac := NewRbacMock(ctx, cfg)
	require.NoError(t, srv.Start())
	defer srv.Stop()
	require.NoError(t, mockRbac.WaitAddress(5*time.Second))

	mockRbac.SetPermissions([]string{
		"idmsvc:token:create",
		"idmsvc:domains:create",
		"idmsvc:domains:read",
		"idmsvc:domains:update",
		"idmsvc:domains:delete",
	})

	// Our mock support "/access" endpoint, but the real rbac
	// service require "/access/" endpoint.
	res, err := http.Get(mockRbac.GetBaseURL() + "/access/?application=idmsvc")
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	expected := "{\"meta\":{\"count\":5,\"limit\":100,\"offset\":0},\"links\":{\"first\":\"http://127.0.0.1:8020/api/rbac/v1/access/?application=idmsvc\\u0026limit=100\\u0026offset=0\",\"last\":\"http://127.0.0.1:8020/api/rbac/v1/access/?application=idmsvc\\u0026limit=100\\u0026offset=0\",\"next\":\"http://127.0.0.1:8020/api/rbac/v1/access/?application=idmsvc\\u0026limit=100\\u0026offset=0\",\"previous\":\"http://127.0.0.1:8020/api/rbac/v1/access/?application=idmsvc\\u0026limit=100\\u0026offset=0\"},\"data\":[{\"permission\":\"idmsvc:token:create\",\"resourceDefinitions\":[]},{\"permission\":\"idmsvc:domains:create\",\"resourceDefinitions\":[]},{\"permission\":\"idmsvc:domains:read\",\"resourceDefinitions\":[]},{\"permission\":\"idmsvc:domains:update\",\"resourceDefinitions\":[]},{\"permission\":\"idmsvc:domains:delete\",\"resourceDefinitions\":[]}]}\n"
	assert.Equal(t, expected, string(data))
}
