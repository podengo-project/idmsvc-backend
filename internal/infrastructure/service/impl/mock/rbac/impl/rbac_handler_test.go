package impl

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessHandler(t *testing.T) {
	cfg := helperConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, mockRbac := NewRbacMock(ctx, cfg)
	err := srv.Start()
	if err == nil {
		defer srv.Stop()
	}
	require.NoError(t, err)
	require.NoError(t, mockRbac.WaitAddress(5*time.Second))

	mockRbac.SetPermissions(Profiles[ProfileDomainAdmin])

	// Our mock support "/access" endpoint, but the real rbac
	// service require "/access/" endpoint.
	u, err := url.ParseRequestURI(mockRbac.GetBaseURL() + "/access/?application=idmsvc")
	require.NoError(t, err)
	require.NotNil(t, u)

	res, err := http.Get(u.String())
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		defer res.Body.Close()
	}
	require.NoError(t, err)
	require.NotNil(t, data)

	dataPage := &Page{}
	err = json.Unmarshal(data, dataPage)
	require.NoError(t, err)
	expected := &Page{
		Meta: map[string]any{
			"count":  float64(6),
			"limit":  float64(100),
			"offset": float64(0),
		},
		Links: map[string]string{
			"first": u.Path + "?application=idmsvc&limit=100&offset=0",
			"last":  u.Path + "?application=idmsvc&limit=100&offset=0",
		},
		Data: []Permission{
			{
				Permission:          "idmsvc:token:create",
				ResourceDefinitions: []any{},
			},
			{
				Permission:          "idmsvc:domains:create",
				ResourceDefinitions: []any{},
			},
			{
				Permission:          "idmsvc:domains:list",
				ResourceDefinitions: []any{},
			},
			{
				Permission:          "idmsvc:domains:update",
				ResourceDefinitions: []any{},
			},
			{
				Permission:          "idmsvc:domains:delete",
				ResourceDefinitions: []any{},
			},
			{
				Permission:          "idmsvc:domains:read",
				ResourceDefinitions: []any{},
			},
		},
	}
	assert.Equal(t, expected, dataPage)
}
