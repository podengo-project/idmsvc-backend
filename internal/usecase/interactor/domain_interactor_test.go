package interactor

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token/domain_token"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	identity "github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
)

func TestNewTodoInteractor(t *testing.T) {
	var component interactor.DomainInteractor
	assert.NotPanics(t, func() {
		component = NewDomainInteractor()
	})
	assert.NotNil(t, component)
}

func TestHelperDomainTypeToUint(t *testing.T) {
	var (
		result uint
	)

	result = helperDomainTypeToUint("")
	assert.Equal(t, model.DomainTypeUndefined, result)

	result = helperDomainTypeToUint(public.RhelIdm)
	assert.Equal(t, model.DomainTypeIpa, result)
}

func createFakeSystemIdentity(orgID string) identity.XRHID {
	return identity.XRHID{
		Identity: identity.Identity{
			OrgID: orgID,
			Type:  "System",
			System: &identity.System{
				CommonName: "21258fc8-c755-11ed-afc4-482ae3863d30",
				CertType:   "system",
			},
		},
	}
}

func createFakeXRHIDMVersion() header.XRHIDMVersion {
	return header.XRHIDMVersion{
		IPAHCCVersion:      "0.7",
		IPAVersion:         "4.10.0-8.el9_1",
		OSReleaseID:        "rhel",
		OSReleaseVersionID: "8",
	}
}

func createFakeRegisterDomainsParams(token domain_token.DomainRegistrationToken, idmVersion header.XRHIDMVersion) *api_public.RegisterDomainParams {
	return &api_public.RegisterDomainParams{
		XRhInsightsRequestId:    pointy.String("TW9uIE1hciAyMCAyMDo1Mzoz"),
		XRhIdmRegistrationToken: string(token),
		XRhIdmVersion:           header.EncodeXRHIDMVersion(&idmVersion),
	}
}

func TestRegisterIpa(t *testing.T) {
	const (
		orgID = "12345"
	)
	secret := []byte("token secret")
	tok, _, err := domain_token.NewDomainRegistrationToken(
		secret,
		string(api_public.RhelIdm),
		orgID,
		time.Hour,
	)
	assert.NoError(t, err)
	var (
		rhsmID                = uuid.MustParse("cf26cd96-c75d-11ed-ae20-482ae3863d30")
		domainID              = domain_token.TokenDomainId(tok)
		requestID             = pointy.String("TW9uIE1hciAyMCAyMDo1Mzoz")
		xrhidSystem           = createFakeSystemIdentity(orgID)
		clientVersionParsed   = createFakeXRHIDMVersion()
		clientVersion         = header.EncodeXRHIDMVersion(&clientVersionParsed)
		paramsNoClientVersion = &api_public.RegisterDomainParams{
			XRhInsightsRequestId:    requestID,
			XRhIdmRegistrationToken: string(tok),
		}
		params = &api_public.RegisterDomainParams{
			XRhInsightsRequestId:    requestID,
			XRhIdmRegistrationToken: string(tok),
			XRhIdmVersion:           clientVersion,
		}
		NotBefore          = time.Now().UTC()
		NotAfter           = NotBefore.Add(24 * time.Hour)
		ESignatureMismatch = echo.NewHTTPError(
			http.StatusUnauthorized,
			"Domain registration token is invalid: Signature mismatch",
		)
	)
	type TestCaseGiven struct {
		Secret []byte
		XRHID  *identity.XRHID
		UUID   uuid.UUID
		Params *api_public.RegisterDomainParams
		Body   *public.Domain
	}
	type TestCaseExpected struct {
		OrgId         string
		ClientVersion *header.XRHIDMVersion
		Output        *model.Domain
		Error         error
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "fail guards with xrhid is nil",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  nil,
				UUID:   uuid.Nil,
				Params: nil,
				Body:   nil,
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         internal_errors.NilArgError("xrhid"),
			},
		},
		{
			Name: "No X-Rh-IDM-Version",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: paramsNoClientVersion,
				Body:   &api_public.Domain{},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         fmt.Errorf("'" + header.HeaderXRHIDMVersion + "' is invalid"),
			},
		},
		{
			Name: "Type is invalid",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					DomainType: "somethingwrong",
					RhelIdm: &api_public.DomainIpa{
						RealmDomains: nil,
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         ESignatureMismatch,
				// "Unsupported domain_type='somethingwrong'" is a signature mismatch
			},
		},
		{
			Name: "Secret is invalid",
			Given: TestCaseGiven{
				Secret: []byte("invalid"),
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					DomainName: "mydomain.example",
					DomainType: api_public.DomainType(api_public.RhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName: "",
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         ESignatureMismatch,
			},
		},
		{
			Name: "Empty slices",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					DomainName: "mydomain.example",
					DomainType: api_public.DomainType(api_public.RhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName: "",
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: &clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("mydomain.example"),
					Description:           pointy.String(""),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String(""),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
						Locations:    []model.IpaLocation{},
						RealmDomains: pq.StringArray{},
					},
				},
				Error: nil,
			},
		},
		{
			Name: "Empty slices and RealmName filled",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					DomainName: "mydomain.example",
					DomainType: api_public.RhelIdm,
					RhelIdm: &api_public.DomainIpa{
						RealmName: "MYDOMAIN.EXAMPLE",
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: &clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("mydomain.example"),
					Description:           pointy.String(""),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
						Locations:    []model.IpaLocation{},
						RealmDomains: pq.StringArray{},
					},
				},
				Error: nil,
			},
		},
		{
			Name: "RealmDomains with some content",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					DomainName: "mydomain.example",
					DomainType: api_public.RhelIdm,
					RhelIdm: &api_public.DomainIpa{
						RealmName:    "MYDOMAIN.EXAMPLE",
						RealmDomains: []string{"server.domain.example"},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: &clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("mydomain.example"),
					Description:           pointy.String(""),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
						Locations:    []model.IpaLocation{},
						RealmDomains: pq.StringArray{"server.domain.example"},
					},
				},
				Error: nil,
			},
		},
		{
			Name: "CaCerts with some content",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					DomainName: "mydomain.example",
					DomainType: api_public.RhelIdm,
					RhelIdm: &api_public.DomainIpa{
						RealmName: "MYDOMAIN.EXAMPLE",
						CaCerts: []api_public.Certificate{
							{
								Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
								SerialNumber: "1",
								Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								NotBefore:    NotBefore,
								NotAfter:     NotAfter,
								Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: &clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("mydomain.example"),
					Description:           pointy.String(""),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName: pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts: []model.IpaCert{
							{
								Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
								SerialNumber: "1",
								Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								NotBefore:    NotBefore,
								NotAfter:     NotAfter,
								Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
						Servers:      []model.IpaServer{},
						Locations:    []model.IpaLocation{},
						RealmDomains: pq.StringArray{},
					},
				},
				Error: nil,
			},
		},
		{
			Name: "Servers as some content",
			Given: TestCaseGiven{
				Secret: secret,
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					DomainName: "mydomain.example",
					DomainType: api_public.RhelIdm,
					RhelIdm: &api_public.DomainIpa{
						RealmName: "MYDOMAIN.EXAMPLE",
						Servers: []api_public.DomainIpaServer{
							{
								Fqdn:                  "server.mydomain.example",
								SubscriptionManagerId: &rhsmID,
								Location:              pointy.String("europe"),
								CaServer:              true,
								PkinitServer:          true,
								HccEnrollmentServer:   true,
								HccUpdateServer:       true,
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: &clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("mydomain.example"),
					Description:           pointy.String(""),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName: pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts:   []model.IpaCert{},
						Servers: []model.IpaServer{
							{
								FQDN:                "server.mydomain.example",
								RHSMId:              pointy.String(rhsmID.String()),
								Location:            pointy.String("europe"),
								CaServer:            true,
								HCCEnrollmentServer: true,
								HCCUpdateServer:     true,
								PKInitServer:        true,
							},
						},
						Locations:    []model.IpaLocation{},
						RealmDomains: pq.StringArray{},
					},
				},
				Error: nil,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			i := domainInteractor{}
			orgID, clientVersion, output, tErr := i.Register(
				testCase.Given.Secret,
				testCase.Given.XRHID,
				testCase.Given.Params,
				testCase.Given.Body,
			)
			if testCase.Expected.Error != nil {
				require.EqualError(t, tErr, testCase.Expected.Error.Error())
				assert.Equal(t, testCase.Expected.OrgId, orgID)
				assert.Equal(t, testCase.Expected.Output, output)
				assert.Equal(t, testCase.Expected.ClientVersion, clientVersion)
			} else {
				require.NoError(t, tErr)
				assert.Equal(t, testCase.Expected.OrgId, orgID)
				assert.Equal(t, testCase.Expected.Output, output)
				assert.Equal(t, testCase.Expected.ClientVersion, clientVersion)
			}
		})
	}
}

//nolint:gocognit // cognitive complexity 13 of func `TestRegisterDefaultValues` is high (> 10)
func TestRegisterDefaultValues(t *testing.T) {
	// nolint:govet  // disable for "fieldalignment: struct with 64 pointer bytes could be 56"
	type TestCase struct {
		Name                          string
		Title                         *string
		Description                   *string
		AutoEnrollmentEnabled         *bool
		ExpectedTitle                 *string
		ExpectedDescription           *string
		ExpectedAutoEnrollmentEnabled *bool
	}

	// Prep
	const (
		orgID = "12345"
	)
	secret := []byte("token secret")
	tok, _, err := domain_token.NewDomainRegistrationToken(
		secret,
		string(api_public.RhelIdm),
		orgID,
		time.Hour,
	)
	require.NoError(t, err)

	var (
		xrhidSystem         = createFakeSystemIdentity(orgID)
		clientVersionParsed = createFakeXRHIDMVersion()
		rhsmID              = uuid.MustParse(xrhidSystem.Identity.System.CommonName)
		params              = createFakeRegisterDomainsParams(tok, clientVersionParsed)
	)

	domainRequest := &api_public.Domain{
		DomainName: "mydomain.example",
		DomainType: api_public.RhelIdm,
		RhelIdm: &api_public.DomainIpa{
			RealmName: "MYDOMAIN.EXAMPLE",
			Servers: []api_public.DomainIpaServer{
				{
					Fqdn:                  "server.mydomain.example",
					SubscriptionManagerId: &rhsmID,
					Location:              pointy.String("europe"),
					CaServer:              true,
					PkinitServer:          true,
					HccEnrollmentServer:   true,
					HccUpdateServer:       true,
				},
			},
		},
	}

	// Test Cases
	testCases := []TestCase{
		{
			Name:                          "Nil/undefined values",
			Title:                         nil,
			Description:                   nil,
			AutoEnrollmentEnabled:         nil,
			ExpectedTitle:                 pointy.String("mydomain.example"),
			ExpectedDescription:           pointy.String(""),
			ExpectedAutoEnrollmentEnabled: pointy.Bool(false),
		},
		{
			Name:                          "Empty values",
			Title:                         pointy.String(""),
			Description:                   pointy.String(""),
			AutoEnrollmentEnabled:         pointy.Bool(false),
			ExpectedTitle:                 pointy.String("mydomain.example"),
			ExpectedDescription:           pointy.String(""),
			ExpectedAutoEnrollmentEnabled: pointy.Bool(false),
		},
		{
			Name:                          "Filled values",
			Title:                         pointy.String("My Domain Title"),
			Description:                   pointy.String("My Domain Description"),
			AutoEnrollmentEnabled:         pointy.Bool(true),
			ExpectedTitle:                 pointy.String("My Domain Title"),
			ExpectedDescription:           pointy.String("My Domain Description"),
			ExpectedAutoEnrollmentEnabled: pointy.Bool(true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			i := domainInteractor{}
			// Given
			domainRequest.Title = testCase.Title
			domainRequest.Description = testCase.Description
			domainRequest.AutoEnrollmentEnabled = testCase.AutoEnrollmentEnabled

			// When
			_, _, output, tErr := i.Register(
				secret,
				&xrhidSystem,
				params,
				domainRequest,
			)

			// Then
			require.NoError(t, tErr)
			if testCase.ExpectedTitle != nil {
				assert.Equal(t, *testCase.ExpectedTitle, *output.Title)
			} else {
				assert.Nil(t, output.Title)
			}
			if testCase.ExpectedDescription != nil {
				assert.Equal(t, *testCase.ExpectedDescription, *output.Description)
			} else {
				assert.Nil(t, output.Description)
			}
			if testCase.ExpectedAutoEnrollmentEnabled != nil {
				assert.Equal(t, *testCase.ExpectedAutoEnrollmentEnabled, *output.AutoEnrollmentEnabled)
			} else {
				assert.Nil(t, output.AutoEnrollmentEnabled)
			}
		})
	}
}

func TestUpdateAgent(t *testing.T) {
	const testOrgID = "12345"
	testID := uuid.MustParse("658700b8-005b-11ee-9e09-482ae3863d30")
	testXRHID := identity.XRHID{
		Identity: identity.Identity{
			OrgID: testOrgID,
			Type:  "user",
			User:  &identity.User{},
			Internal: identity.Internal{
				OrgID: testOrgID,
			},
		},
	}
	testXRHIDMVersion := header.XRHIDMVersion{
		IPAHCCVersion:      "",
		IPAVersion:         "",
		OSReleaseID:        "rhel",
		OSReleaseVersionID: "8",
	}
	testParams := api_public.UpdateDomainAgentParams{
		XRhInsightsRequestId: pointy.String("put_update_test"),
		XRhIdmVersion:        header.EncodeXRHIDMVersion(&testXRHIDMVersion),
	}
	testBadParams := api_public.UpdateDomainAgentParams{
		XRhInsightsRequestId: pointy.String("put_update_test"),
		XRhIdmVersion:        "{",
	}
	testWrongTypeBody := api_public.UpdateDomainAgentRequest{
		DomainName: "mydomain.example",
		DomainType: "aninvalidtype",
	}
	testBody := api_public.UpdateDomainAgentRequest{
		DomainName: "mydomain.example",
		DomainType: api_public.RhelIdm,
		RhelIdm: api_public.DomainIpa{
			RealmName:    "mydomain.example",
			RealmDomains: []string{"mydomain.example"},
			CaCerts:      []api_public.Certificate{},
			Servers:      []api_public.DomainIpaServer{},
		},
	}
	i := NewDomainInteractor()

	// Get an error in guards
	orgID, xrhidmVersion, domain, err := i.UpdateAgent(nil, uuid.Nil, nil, nil)
	require.EqualError(t, err, "code=500, message='xrhid' cannot be nil")
	assert.Equal(t, "", orgID)
	assert.Nil(t, xrhidmVersion)
	assert.Nil(t, domain)

	// Get an error with nil param
	orgID, xrhidmVersion, domain, err = i.UpdateAgent(&testXRHID, testID, nil, &testBody)
	require.EqualError(t, err, "code=500, message='params' cannot be nil")
	assert.Equal(t, "", orgID)
	assert.Nil(t, xrhidmVersion)
	assert.Nil(t, domain)

	// Error retrieving ipa-hcc version information
	orgID, xrhidmVersion, domain, err = i.UpdateAgent(&testXRHID, testID, &testBadParams, &testBody)
	require.EqualError(t, err, "'"+header.HeaderXRHIDMVersion+"' is invalid")
	assert.Equal(t, "", orgID)
	assert.Nil(t, xrhidmVersion)
	assert.Nil(t, domain)

	// Error because of wrongtype
	orgID, xrhidmVersion, domain, err = i.UpdateAgent(&testXRHID, testID, &testParams, &testWrongTypeBody)
	require.EqualError(t, err, "Unsupported domain_type='aninvalidtype'")
	assert.Equal(t, "", orgID)
	assert.Nil(t, xrhidmVersion)
	assert.Nil(t, domain)

	// Success result
	orgID, xrhidmVersion, domain, err = i.UpdateAgent(&testXRHID, testID, &testParams, &testBody)
	assert.NoError(t, err)
	assert.Equal(t, testOrgID, orgID)
	require.NotNil(t, xrhidmVersion)
	require.NotNil(t, domain)
}

func assertListEqualError(t *testing.T, err error, msg string, orgID string, offset int, limit int) {
	assert.EqualError(t, err, msg)
	assert.Equal(t, "", orgID)
	assert.Equal(t, -1, offset)
	assert.Equal(t, -1, limit)
}

func TestUpdateUser(t *testing.T) {
	i := NewDomainInteractor()

	// given mock XRHID and header parasm
	testXRHID := test.UserXRHID
	testParams := api_public.UpdateDomainUserParams{
		XRhInsightsRequestId: pointy.String("some-request-id"),
	}

	testOrgID := test.OrgId
	testUUID := test.DomainUUID
	testAutoEnrollment := pointy.Bool(true)
	testDescription := "My Example Domain Description"

	t.Run("empty body", func(t *testing.T) {
		// when body is nil
		orgID, domain, err := i.UpdateUser(&testXRHID, testUUID, &testParams, nil)

		// return error
		require.EqualError(t, err, "code=500, message='body' cannot be nil")
		assert.Equal(t, "", orgID)
		assert.Nil(t, domain)
	})

	t.Run("empty title", func(t *testing.T) {
		// given
		testEmptyTitle := ""
		body := public.UpdateDomainUserRequest{
			AutoEnrollmentEnabled: testAutoEnrollment,
			Title:                 pointy.String(testEmptyTitle),
			Description:           pointy.String(testDescription),
		}

		// when title is empty
		orgID, domain, err := i.UpdateUser(&testXRHID, testUUID, &testParams, &body)

		// return bad-request error
		require.EqualError(t, err, "code=400, message='title' cannot be empty")
		assert.Equal(t, "", orgID)
		assert.Nil(t, domain)
	})

	t.Run("valid request", func(t *testing.T) {
		// given
		testTitle := "My Example Domain Title"
		testBody := public.UpdateDomainUserRequest{
			Title: pointy.String(testTitle),
		}

		// success
		orgID, domain, err := i.UpdateUser(&testXRHID, testUUID, &testParams, &testBody)
		require.NoError(t, err)
		assert.Equal(t, testOrgID, orgID)
		require.NotNil(t, domain)

		// assert org and uuid is set
		assert.Equal(t, testUUID, domain.DomainUuid)
		assert.Equal(t, testOrgID, domain.OrgId)

		// changed field is set
		assert.Equal(t, testTitle, *domain.Title)

		// undefined field are not set
		assert.Nil(t, domain.AutoEnrollmentEnabled)
		assert.Nil(t, domain.Description)

		// other model.Domain fields are not st
		assert.Nil(t, domain.DomainName)
		assert.Nil(t, domain.Type)
		assert.Nil(t, domain.IpaDomain)
	})
}

func TestList(t *testing.T) {
	i := NewDomainInteractor()
	testOrgID := "12345"

	// xrhid is nil
	orgID, offset, limit, err := i.List(nil, nil)
	assertListEqualError(t, err, "code=500, message='xrhid' cannot be nil", orgID, offset, limit)

	// params is nil
	xrhid := identity.XRHID{
		Identity: identity.Identity{
			OrgID: testOrgID,
		},
	}
	orgID, offset, limit, err = i.List(&xrhid, nil)
	assertListEqualError(t, err, "code=500, message='params' cannot be nil", orgID, offset, limit)

	// params.Offset is nil
	params := api_public.ListDomainsParams{}
	orgID, offset, limit, err = i.List(&xrhid, &params)
	assert.NoError(t, err)
	assert.Equal(t, testOrgID, orgID)
	assert.Equal(t, 0, offset)
	assert.Equal(t, 10, limit)

	// params.Offset is not nil
	params.Offset = pointy.Int(20)
	orgID, offset, limit, err = i.List(&xrhid, &params)
	assert.NoError(t, err)
	assert.Equal(t, testOrgID, orgID)
	assert.Equal(t, 20, offset)
	assert.Equal(t, 10, limit)

	// params.Limit is not nil
	params.Limit = pointy.Int(30)
	orgID, offset, limit, err = i.List(&xrhid, &params)
	assert.NoError(t, err)
	assert.Equal(t, testOrgID, orgID)
	assert.Equal(t, 20, offset)
	assert.Equal(t, 30, limit)
}

func TestGuardRegister(t *testing.T) {
	var err error

	i := domainInteractor{}

	err = i.guardRegister(nil, nil, nil)
	assert.EqualError(t, err, "code=500, message='xrhid' cannot be nil")

	xrhid := &identity.XRHID{}
	err = i.guardRegister(xrhid, nil, nil)
	assert.EqualError(t, err, "code=500, message='params' cannot be nil")

	params := &api_public.RegisterDomainParams{}
	err = i.guardRegister(xrhid, params, nil)
	assert.EqualError(t, err, "code=500, message='body' cannot be nil")

	body := &public.Domain{}
	err = i.guardRegister(xrhid, params, body)
	assert.NoError(t, err)
}

func TestGuardXrhidUUID(t *testing.T) {
	var err error

	i := domainInteractor{}

	err = i.guardXrhidUUID(nil, uuid.Nil)
	assert.EqualError(t, err, "code=500, message='xrhid' cannot be nil")

	xrhid := &identity.XRHID{}
	err = i.guardXrhidUUID(xrhid, uuid.Nil)
	assert.EqualError(t, err, "'UUID' is invalid")

	UUID := uuid.MustParse("b0264600-005c-11ee-ba48-482ae3863d30")
	err = i.guardXrhidUUID(xrhid, UUID)
	assert.NoError(t, err)
}

func TestGuardUserUpdate(t *testing.T) {
	i := domainInteractor{}

	body := &public.UpdateDomainUserRequest{}
	err := i.guardUserUpdate(body)
	assert.NoError(t, err)

	body.Title = pointy.String("")
	err = i.guardUserUpdate(body)
	require.EqualError(t, err, "code=400, message='title' cannot be empty")

	body.Title = pointy.String("Some title")
	err = i.guardUserUpdate(body)
	require.NoError(t, err)
}

func TestTranslateDomain(t *testing.T) {
	testOrgID := "12345"
	testID := uuid.MustParse("c95c6e74-005c-11ee-82b5-482ae3863d30")
	testTitle := pointy.String("My Example Domain Title")
	i := domainInteractor{}
	assert.Panics(t, func() {
		//nolint:errcheck,gosec
		i.translateDomain("", uuid.Nil, nil)
	})

	testDescription := "My Example Domain Description"
	testBody := public.Domain{
		AutoEnrollmentEnabled: pointy.Bool(true),
		Title:                 testTitle,
		Description:           pointy.String(testDescription),
		DomainName:            "mydomain.example",
		DomainId:              &testID,
		DomainType:            api_public.RhelIdm,
		RhelIdm: &api_public.DomainIpa{
			RealmName:    "mydomain.example",
			RealmDomains: []string{"mydomain.example"},
			CaCerts:      []api_public.Certificate{},
			Servers:      []api_public.DomainIpaServer{},
		},
	}
	testWrongTypeBody := public.Domain{
		AutoEnrollmentEnabled: pointy.Bool(true),
		Title:                 testTitle,
		Description:           pointy.String(testDescription),
		DomainName:            "mydomain.example",
		DomainId:              &testID,
		DomainType:            "wrongtype",
		RhelIdm: &api_public.DomainIpa{
			RealmName:    "mydomain.example",
			RealmDomains: []string{"mydomain.example"},
			CaCerts:      []api_public.Certificate{},
			Servers:      []api_public.DomainIpaServer{},
		},
	}

	domain, err := i.translateDomain(testOrgID, testID, &testWrongTypeBody)
	assert.EqualError(t, err, "Unsupported domain_type='wrongtype'")
	assert.Nil(t, domain)

	// Success case
	domain, err = i.translateDomain(testOrgID, testID, &testBody)
	assert.NoError(t, err)
	assert.NotNil(t, domain)
}

func TestTranslateUpdateDomainAgent(t *testing.T) {
	// Given
	i := domainInteractor{}
	testOrgID := "12345"
	testUUID := uuid.MustParse("c95c6e74-005c-11ee-82b5-482ae3863d30")
	agentRequest := public.UpdateDomainAgentRequest{
		DomainName: "mydomain.example",
		DomainType: api_public.RhelIdm,
		RhelIdm: public.DomainIpa{
			RealmName:    "mydomain.example",
			RealmDomains: []string{"mydomain.example"},
			CaCerts:      []api_public.Certificate{},
			Servers:      []public.DomainIpaServer{},
		},
	}

	t.Run("Correct input", func(t *testing.T) {
		// When translateUpdateDomainAgent is called with expected values
		domain, err := i.translateUpdateDomainAgent(testOrgID, testUUID, &agentRequest)

		// Then it returns the expected domain and no error
		require.NoError(t, err)
		require.NotNil(t, domain)

		// Fields are set as expected
		assert.Equal(t, testOrgID, domain.OrgId)
		assert.Equal(t, testUUID, domain.DomainUuid)
		assert.Equal(t, "mydomain.example", *domain.DomainName)

		// The domain type matches input (RhelIdm) and respected domain field is set
		assert.Equal(t, model.DomainTypeIpa, *domain.Type)
		assert.NotNil(t, domain.IpaDomain)
		assert.Equal(t, "mydomain.example", *domain.IpaDomain.RealmName)
	})

	t.Run("Invalid domain type", func(t *testing.T) {
		// Given an invalid domain type
		agentRequest.DomainType = "invalid"

		// When translateUpdateDomainAgent is called
		domain, err := i.translateUpdateDomainAgent(testOrgID, testUUID, &agentRequest)

		// Then it returns an error
		require.EqualError(t, err, "Unsupported domain_type='invalid'")
		assert.Nil(t, domain)
	})
}

func TestGetByID(t *testing.T) {
	i := NewDomainInteractor()

	orgID, err := i.GetByID(nil, nil)
	assert.EqualError(t, err, "code=500, message='xrhid' cannot be nil")
	assert.Equal(t, "", orgID)

	testOrgID := "12345"
	xrhid := identity.XRHID{
		Identity: identity.Identity{
			OrgID: testOrgID,
			Internal: identity.Internal{
				OrgID: testOrgID,
			},
		},
	}
	orgID, err = i.GetByID(&xrhid, nil)
	assert.EqualError(t, err, "code=500, message='params' cannot be nil")
	assert.Equal(t, "", orgID)

	testRequestID := "getByID"
	params := public.ReadDomainParams{
		XRhInsightsRequestId: &testRequestID,
	}
	orgID, err = i.GetByID(&xrhid, &params)
	assert.NoError(t, err)
	assert.Equal(t, testOrgID, orgID)
}

func TestRegisterOrUpdateRhelIdmLocations(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    *api_public.Domain
		Expected *model.Ipa
	}
	testCases := []TestCase{
		{
			Name: "nil Locations",
			Given: &api_public.Domain{
				RhelIdm: &api_public.DomainIpa{
					Locations: nil,
				},
			},
			Expected: &model.Ipa{
				Locations: []model.IpaLocation{},
			},
		},
		{
			Name: "Empty Locations slice",
			Given: &api_public.Domain{
				RhelIdm: &api_public.DomainIpa{
					Locations: []api_public.Location{},
				},
			},
			Expected: &model.Ipa{
				Locations: []model.IpaLocation{},
			},
		},
		{
			Name: "Location without description",
			Given: &api_public.Domain{
				RhelIdm: &api_public.DomainIpa{
					Locations: []api_public.Location{
						{
							Name:        "boston",
							Description: nil,
						},
					},
				},
			},
			Expected: &model.Ipa{
				Locations: []model.IpaLocation{
					{
						Name:        "boston",
						Description: nil,
					},
				},
			},
		},
		{
			Name: "Location with description",
			Given: &api_public.Domain{
				RhelIdm: &api_public.DomainIpa{
					Locations: []api_public.Location{
						{
							Name:        "boston",
							Description: pointy.String("Boston data center"),
						},
					},
				},
			},
			Expected: &model.Ipa{
				Locations: []model.IpaLocation{
					{
						Name:        "boston",
						Description: pointy.String("Boston data center"),
					},
				},
			},
		},
	}
	i := domainInteractor{}
	for _, item := range testCases {
		t.Log(item.Name)
		ipa := &model.Ipa{}
		i.translateIdmLocations(item.Given.RhelIdm, ipa)
		assert.Equal(t, item.Expected, ipa)
	}
}

func TestCreateDomainToken(t *testing.T) {
	const (
		testOrgID = "12345"
	)
	var (
		xrhidUser = identity.XRHID{
			Identity: identity.Identity{
				OrgID: testOrgID,
				Type:  "User",
				User:  &identity.User{},
				Internal: identity.Internal{
					OrgID: testOrgID,
				},
			},
		}
	)

	type TestCaseGiven struct {
		XRHID  *identity.XRHID
		Params *api_public.CreateDomainTokenParams
		Body   *api_public.DomainRegTokenRequest
	}
	type TestCaseExpected struct {
		OrgID      string
		DomainType public.DomainType
		Err        error
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
				XRHID:  nil,
				Params: nil,
				Body:   &api_public.DomainRegTokenRequest{},
			},
			Expected: TestCaseExpected{
				Err: internal_errors.NilArgError("xrhid"),
			},
		},
		{
			Name: "nil 'params'",
			Given: TestCaseGiven{
				XRHID:  &xrhidUser,
				Params: nil,
				Body:   &api_public.DomainRegTokenRequest{},
			},
			Expected: TestCaseExpected{
				Err: internal_errors.NilArgError("params"),
			},
		},
		{
			Name: "nil 'body'",
			Given: TestCaseGiven{
				XRHID:  &xrhidUser,
				Params: &api_public.CreateDomainTokenParams{},
				Body:   nil,
			},
			Expected: TestCaseExpected{
				Err: internal_errors.NilArgError("body"),
			},
		},
		{
			Name: "wrong domain type",
			Given: TestCaseGiven{
				XRHID:  &xrhidUser,
				Params: &api_public.CreateDomainTokenParams{},
				Body: &api_public.DomainRegTokenRequest{
					DomainType: api_public.DomainType("invalid"),
				},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("Unsupported domain_type='invalid'"),
			},
		},
		{
			Name: "success case",
			Given: TestCaseGiven{
				XRHID:  &xrhidUser,
				Params: &api_public.CreateDomainTokenParams{},
				Body: &api_public.DomainRegTokenRequest{
					DomainType: api_public.RhelIdm,
				},
			},
			Expected: TestCaseExpected{
				OrgID:      testOrgID,
				DomainType: api_public.RhelIdm,
				Err:        nil,
			},
		},
	}

	component := NewDomainInteractor()
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		orgID, domainType, err := component.CreateDomainToken(
			testCase.Given.XRHID,
			testCase.Given.Params,
			testCase.Given.Body,
		)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			require.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Equal(t, "", orgID)
			assert.Equal(t, public.DomainType(""), domainType)
		} else {
			assert.Equal(t, testCase.Expected.OrgID, orgID)
			assert.Equal(t, testCase.Expected.DomainType, domainType)
			assert.NoError(t, err)
		}
	}
}

func TestDelete(t *testing.T) {
	i := NewDomainInteractor()

	xrhidUser := test.UserXRHID
	testID := test.DomainUUID
	params := api_public.DeleteDomainParams{}

	// Guard xrhid is nil
	orgID, UUID, err := i.Delete(nil, testID, &params)
	assert.Equal(t, "", orgID)
	assert.Equal(t, uuid.UUID{}, UUID)
	assert.EqualError(t, err, "code=500, message='xrhid' cannot be nil")

	// Guard params is nil
	orgID, UUID, err = i.Delete(&xrhidUser, testID, nil)
	assert.Equal(t, "", orgID)
	assert.Equal(t, uuid.UUID{}, UUID)
	assert.EqualError(t, err, "'params' cannot be nil")

	// Success result
	orgID, UUID, err = i.Delete(&xrhidUser, testID, &params)
	assert.Equal(t, xrhidUser.Identity.OrgID, orgID)
	assert.Equal(t, testID, UUID)
	assert.NoError(t, err)
}
