package rbac_data

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRBACMap(t *testing.T) {
	var builder RBACMapBuilder

	require.NotPanics(t, func() {
		builder = NewRBACMapBuilder()
	})
	require.NotNil(t, builder)
}

func TestAdd(t *testing.T) {
	var builder RBACMapBuilder

	builder = NewRBACMapBuilder()
	builder.
		Add("/domains", http.MethodGet, NewRbacPermission("idmsvc", "domains", "read")).
		Add("/domains", http.MethodGet, NewRbacPermission("idmsvc", "domains", "order")).
		Build()
}
