package router

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMajorVersion(t *testing.T) {
	assert.Equal(t, "", getMajorVersion(""))
	assert.Equal(t, "1", getMajorVersion("1.0"))
	assert.Equal(t, "1", getMajorVersion("1.0.3"))
	assert.Equal(t, "1", getMajorVersion("1."))
	assert.Equal(t, "a", getMajorVersion("a.b.c"))
}

func TestCheckRouterConfig(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    RouterConfig
		Expected error
	}
	testCases := []TestCase{
		{
			Name: "PublicPath is empty",
			Given: RouterConfig{
				PublicPath: "",
			},
			Expected: fmt.Errorf("PublicPath cannot be empty"),
		},
		{
			Name: "PrivatePath is empty",
			Given: RouterConfig{
				PublicPath:  "/api/hmsidm/v1",
				PrivatePath: "",
			},
			Expected: fmt.Errorf("PrivatePath cannot be empty"),
		},
		{
			Name: "PublicPath and PrivatePath are equal",
			Given: RouterConfig{
				PublicPath:  "/api/hmsidm/v1",
				PrivatePath: "/api/hmsidm/v1",
				Version:     "",
			},
			Expected: fmt.Errorf("PublicPath and PrivatePath cannot be equal"),
		},
		{
			Name: "Version is empty",
			Given: RouterConfig{
				PublicPath:  "/api/hmsidm/v1",
				PrivatePath: "/private",
				Version:     "",
			},
			Expected: fmt.Errorf("Version cannot be empty"),
		},
		{
			Name: "Success scenario",
			Given: RouterConfig{
				PublicPath:  "/api/hmsidm/v1",
				PrivatePath: "/private",
				Version:     "1.0",
			},
			Expected: nil,
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		err := checkRouterConfig(testCase.Given)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
