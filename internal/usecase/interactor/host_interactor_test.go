package interactor

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/openlyinc/pointy"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHostInteractor(t *testing.T) {
	var component interactor.HostInteractor
	assert.NotPanics(t, func() {
		component = NewHostInteractor()
	})
	assert.NotNil(t, component)
}

func TestHostConf(t *testing.T) {
	const (
		testCN         = "21258fc8-c755-11ed-afc4-482ae3863d30"
		testOrgID      = "12345"
		testFqdn       = "server.ipa.test"
		testDomainId   = "0851e1d6-003f-11ee-adf4-482ae3863d30"
		testDomainName = "ipa.test"
	)
	testInventoryId := uuid.MustParse("70d51c08-8831-4989-bca6-04d2c11a58e2")
	testDomainUUID := uuid.MustParse(testDomainId)
	testDomainType := api_public.DomainType(api_public.RhelIdm)

	xrhidSystem := identity.XRHID{
		Identity: identity.Identity{
			OrgID: testOrgID,
			Type:  "System",
			System: identity.System{
				CommonName: testCN,
				CertType:   "system",
			},
		},
	}
	xrhidUser := identity.XRHID{
		Identity: identity.Identity{
			OrgID: testOrgID,
			Type:  "User",
			User:  identity.User{},
			Internal: identity.Internal{
				OrgID: testOrgID,
			},
		},
	}

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
				InventoryId: testInventoryId,
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
				XRHID:       &xrhidUser,
				InventoryId: testInventoryId,
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
				XRHID:       &xrhidSystem,
				InventoryId: testInventoryId,
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
				XRHID:       &xrhidSystem,
				InventoryId: testInventoryId,
				Fqdn:        testFqdn,
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
				XRHID:       &xrhidSystem,
				InventoryId: testInventoryId,
				Fqdn:        testFqdn,
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
				XRHID:       &xrhidSystem,
				InventoryId: testInventoryId,
				Fqdn:        testFqdn,
				Params:      &api_public.HostConfParams{},
				Body:        &api_public.HostConf{},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &interactor.HostConfOptions{
					OrgId:       testOrgID,
					CommonName:  testCN,
					InventoryId: testInventoryId,
					Fqdn:        testFqdn,
					DomainId:    nil,
					DomainName:  nil,
					DomainType:  nil,
				},
			},
		},
		{
			Name: "valid call with params",
			Given: TestCaseGiven{
				XRHID:       &xrhidSystem,
				InventoryId: testInventoryId,
				Fqdn:        testFqdn,
				Params:      &api_public.HostConfParams{},
				Body: &api_public.HostConf{
					DomainId:   &testDomainUUID,
					DomainName: pointy.String(testDomainName),
					DomainType: &testDomainType,
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &interactor.HostConfOptions{
					OrgId:       testOrgID,
					CommonName:  testCN,
					InventoryId: testInventoryId,
					Fqdn:        testFqdn,
					DomainId:    &testDomainUUID,
					DomainName:  pointy.String(testDomainName),
					DomainType:  &testDomainType,
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
