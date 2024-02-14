package presenter

import (
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHostconfJwkPresenter(t *testing.T) {
	assert.Panics(t, func() {
		NewHostconfJwkPresenter(nil)
	})
	cfg := test.GetTestConfig()
	assert.NotPanics(t, func() {
		NewHostconfJwkPresenter(cfg)
	})
}

func TestPublicSigningKeys(t *testing.T) {
	type TestCaseGiven struct {
		Keys   []string
		Output *public.SigningKeysResponse
	}
	type TestCaseExpected struct {
		Err    error
		Output *public.SigningKeysResponse
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "error when 'keys' is nil",
			Given: TestCaseGiven{
				Keys: nil,
			},
			Expected: TestCaseExpected{
				Err:    internal_errors.NilArgError("keys"),
				Output: nil,
			},
		},
		{
			Name: "Success",
			Given: TestCaseGiven{
				Keys: []string{"key1", "key2"},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Output: &public.SigningKeysResponse{
					Keys: []string{"key1", "key2"},
				},
			},
		},
	}
	cfg := test.GetTestConfig()
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := NewHostconfJwkPresenter(cfg)
		output, err := obj.PublicSigningKeys(testCase.Given.Keys)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Nil(t, output)
		} else {
			assert.NoError(t, err)
			assert.Equal(t,
				testCase.Expected.Output.Keys,
				output.Keys)
		}
	}
}
