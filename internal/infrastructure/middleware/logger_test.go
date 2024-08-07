package middleware

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSource(t *testing.T) {
	require.Equal(t, "", getSource(""))
	require.Equal(t, "client1", getSource("client1, proxy1, proxy2"))
	require.Equal(t, "129.78.138.66", getSource("129.78.138.66, 129.78.64.103"))
}
