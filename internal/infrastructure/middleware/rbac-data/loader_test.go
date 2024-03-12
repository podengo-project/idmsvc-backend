package rbac_data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBACMapLoadErrors(t *testing.T) {
	require.PanicsWithError(t, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `???` into rbac_data.RbacFile", func() {
		_, _ = RBACMapLoad([]byte(`???`))
	}, "Panic on wrong yaml format")

	require.PanicsWithError(t, "it was expected rbac map config data version '1.0', but '0.0' was found", func() {
		_, _ = RBACMapLoad([]byte(`version: "0.0"`))
	}, "Panic on wrong version")
}

func TestRBACMapLoadSuccess(t *testing.T) {
	var (
		rbacMap RBACMap
		prefix  string
		example = `---
version: "1.0"
prefix: "/api/idmsvc/v1"
data:
  "/domains/token":
    POST: "idmsvc:token:create"
  "/domains":
    POST: "idmsvc:domains:create"
    GET: "idmsvc:domains:list"
  "/domains/:uuid":
    GET: "idmsvc:domains:read"
    PUT: "idmsvc:domains:update"
    PATCH: "idmsvc:domains:update"
    DELETE: "idmsvc:domains:delete"
  "/host-conf/:inventory_id/:fqdn":
    POST: "idmsvc:host_conf:execute"
  "/signing_keys":
    GET: "idmsvc:signing_keys:execute"
`
	)

	require.NotPanics(t, func() {
		prefix, rbacMap = RBACMapLoad([]byte(example))
	})
	assert.NotEqual(t, RBACMap{}, rbacMap)
	assert.Equal(t, "/api/idmsvc/v1", prefix)
}
