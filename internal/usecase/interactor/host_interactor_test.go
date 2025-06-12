package interactor

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
)

func TestNewHostInteractor(t *testing.T) {
	var component interactor.HostInteractor
	assert.NotPanics(t, func() {
		component = NewHostInteractor()
	})
	assert.NotNil(t, component)
}

func TestHostConf(t *testing.T) {
	type TestCaseGiven struct {
		XRHID       *identity.XRHID
		InventoryId uuid.UUID
		Fqdn        string
		Params      *api_public.HostConfParams
		Body        *api_public.HostConf
	}
	type TestCaseExpected struct {
		Err error
		Out *interactor.HostConfOptions
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "nil 'xrhid'",
			Given: TestCaseGiven{
				XRHID:       nil,
				InventoryId: test.Client1.InventoryUUID,
				Fqdn:        "",
				Params:      nil,
				Body:        &api_public.HostConf{},
			},
			Expected: TestCaseExpected{
				Err: internal_errors.NilArgError("xrhid"),
				Out: nil,
			},
		},
		{
			Name: "'xrhid' wrong type",
			Given: TestCaseGiven{
				XRHID:       &test.UserXRHID,
				InventoryId: test.Client1.InventoryUUID,
				Fqdn:        "",
				Params:      nil,
				Body:        &api_public.HostConf{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("invalid 'xrhid' type 'User'"),
				Out: nil,
			},
		},
		{
			Name: "empty fqdn",
			Given: TestCaseGiven{
				XRHID:       &test.Client1XRHID,
				InventoryId: test.Client1.InventoryUUID,
				Fqdn:        "",
				Params:      nil,
				Body:        &api_public.HostConf{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'fqdn' is empty"),
				Out: nil,
			},
		},
		{
			Name: "'params' is nil",
			Given: TestCaseGiven{
				XRHID:       &test.Client1XRHID,
				InventoryId: test.Client1.InventoryUUID,
				Fqdn:        test.Client1.Fqdn,
				Params:      nil,
				Body:        &api_public.HostConf{},
			},
			Expected: TestCaseExpected{
				Err: internal_errors.NilArgError("params"),
				Out: nil,
			},
		},
		{
			Name: "nil for the 'body'",
			Given: TestCaseGiven{
				XRHID:       &test.Client1XRHID,
				InventoryId: test.Client1.InventoryUUID,
				Fqdn:        test.Client1.Fqdn,
				Params:      &api_public.HostConfParams{},
				Body:        nil,
			},
			Expected: TestCaseExpected{
				Err: internal_errors.NilArgError("body"),
				Out: nil,
			},
		},
		{
			Name: "valid call without params",
			Given: TestCaseGiven{
				XRHID:       &test.Client1XRHID,
				InventoryId: test.Client1.InventoryUUID,
				Fqdn:        test.Client1.Fqdn,
				Params:      &api_public.HostConfParams{},
				Body:        &api_public.HostConf{},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &interactor.HostConfOptions{
					OrgId:       test.OrgId,
					CommonName:  test.Client1.CertUUID,
					InventoryId: test.Client1.InventoryUUID,
					Fqdn:        test.Client1.Fqdn,
					DomainId:    nil,
					DomainName:  nil,
					DomainType:  nil,
				},
			},
		},
		{
			Name: "valid call with params",
			Given: TestCaseGiven{
				XRHID:       &test.Client1XRHID,
				InventoryId: test.Client1.InventoryUUID,
				Fqdn:        test.Client1.Fqdn,
				Params:      &api_public.HostConfParams{},
				Body: &api_public.HostConf{
					DomainId:   &test.DomainUUID,
					DomainName: pointy.String(test.DomainName),
					DomainType: &test.DomainType,
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &interactor.HostConfOptions{
					OrgId:       test.OrgId,
					CommonName:  test.Client1.CertUUID,
					InventoryId: test.Client1.InventoryUUID,
					Fqdn:        test.Client1.Fqdn,
					DomainId:    &test.DomainUUID,
					DomainName:  pointy.String(test.DomainName),
					DomainType:  &test.DomainType,
				},
			},
		},
	}

	component := NewHostInteractor()
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		options, err := component.HostConf(testCase.Given.XRHID, testCase.Given.InventoryId, testCase.Given.Fqdn, testCase.Given.Params, testCase.Given.Body)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			require.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Nil(t, options)
		} else {
			assert.NoError(t, err)
			require.NotNil(t, testCase.Expected.Out)
			require.Equal(t, testCase.Expected.Out.OrgId, options.OrgId)
			require.Equal(t, testCase.Expected.Out.CommonName, options.CommonName)
			require.Equal(t, testCase.Expected.Out.InventoryId, options.InventoryId)
			require.Equal(t, testCase.Expected.Out.Fqdn, options.Fqdn)
			require.Equal(t, testCase.Expected.Out.DomainId, options.DomainId)
			require.Equal(t, testCase.Expected.Out.DomainName, options.DomainName)
			require.Equal(t, testCase.Expected.Out.DomainType, options.DomainType)
		}
	}
}
