package interactor

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/token"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestRegisterIpa(t *testing.T) {
	const (
		cn          = "21258fc8-c755-11ed-afc4-482ae3863d30"
		orgID       = "12345"
		description = "My Example Domain Description"
	)
	tok, _, err := token.NewDomainRegistrationToken(
		[]byte("TODO secret"),
		string(api_public.RhelIdm),
		orgID,
		time.Hour,
	)
	assert.NoError(t, err)
	var (
		rhsmID      = uuid.MustParse("cf26cd96-c75d-11ed-ae20-482ae3863d30")
		domainID    = token.TokenDomainId(tok)
		requestID   = pointy.String("TW9uIE1hciAyMCAyMDo1Mzoz")
		testTitle   = pointy.String("My Domain Example")
		xrhidSystem = identity.XRHID{
			Identity: identity.Identity{
				OrgID: orgID,
				Type:  "System",
				System: identity.System{
					CommonName: cn,
					CertType:   "system",
				},
			},
		}
		clientVersionParsed = &header.XRHIDMVersion{
			IPAHCCVersion:      "0.7",
			IPAVersion:         "4.10.0-8.el9_1",
			OSReleaseID:        "rhel",
			OSReleaseVersionID: "8",
		}
		clientVersion         = header.EncodeXRHIDMVersion(clientVersionParsed)
		paramsNoClientVersion = &api_public.RegisterDomainParams{
			XRhInsightsRequestId:    requestID,
			XRhIdmRegistrationToken: string(tok),
		}
		params = &api_public.RegisterDomainParams{
			XRhInsightsRequestId:    requestID,
			XRhIdmRegistrationToken: string(tok),
			XRhIdmVersion:           clientVersion,
		}
		NotBefore = time.Now().UTC()
		NotAfter  = NotBefore.Add(24 * time.Hour)
	)
	type TestCaseGiven struct {
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
				XRHID:  nil,
				UUID:   uuid.Nil,
				Params: nil,
				Body:   nil,
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         fmt.Errorf("'xrhid' is nil"),
			},
		},
		{
			Name: "No X-Rh-IDM-Version",
			Given: TestCaseGiven{
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: paramsNoClientVersion,
				Body:   &api_public.Domain{},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         fmt.Errorf("'X-Rh-Idm-Version' is invalid"),
			},
		},
		{
			Name: "Type is invalid",
			Given: TestCaseGiven{
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
				// Error:         fmt.Errorf("Unsupported domain_type='somethingwrong'"),
				Error: fmt.Errorf("Signature mismatch"),
			},
		},
		{
			Name: "Empty slices",
			Given: TestCaseGiven{
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					Title:       testTitle,
					Description: pointy.String(description),
					DomainName:  "mydomain.example",
					DomainType:  api_public.DomainType(api_public.RhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName: "",
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 testTitle,
					Description:           pointy.String(description),
					AutoEnrollmentEnabled: nil,
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
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					Title:       testTitle,
					Description: pointy.String(description),
					DomainName:  "mydomain.example",
					DomainType:  api_public.RhelIdm,
					RhelIdm: &api_public.DomainIpa{
						RealmName: "MYDOMAIN.EXAMPLE",
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 testTitle,
					Description:           pointy.String(description),
					AutoEnrollmentEnabled: nil,
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
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					Title:       testTitle,
					Description: pointy.String(description),
					DomainName:  "mydomain.example",
					DomainType:  api_public.RhelIdm,
					RhelIdm: &api_public.DomainIpa{
						RealmName:    "MYDOMAIN.EXAMPLE",
						RealmDomains: []string{"server.domain.example"},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         orgID,
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 testTitle,
					Description:           pointy.String(description),
					AutoEnrollmentEnabled: nil,
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
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					Title:       testTitle,
					Description: pointy.String(description),
					DomainName:  "mydomain.example",
					DomainType:  api_public.RhelIdm,
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
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 testTitle,
					Description:           pointy.String(description),
					AutoEnrollmentEnabled: nil,
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
				XRHID:  &xrhidSystem,
				UUID:   domainID,
				Params: params,
				Body: &api_public.Domain{
					Title:       testTitle,
					Description: pointy.String(description),
					DomainName:  "mydomain.example",
					DomainType:  api_public.RhelIdm,
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
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					OrgId:                 orgID,
					DomainUuid:            domainID,
					DomainName:            pointy.String("mydomain.example"),
					Title:                 testTitle,
					Description:           pointy.String(description),
					AutoEnrollmentEnabled: nil,
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
		t.Log(testCase.Name)
		i := NewDomainInteractor()
		orgID, clientVersion, output, err := i.Register(testCase.Given.XRHID, testCase.Given.Params, testCase.Given.Body)
		if testCase.Expected.Error != nil {
			assert.EqualError(t, err, testCase.Expected.Error.Error())
			assert.Equal(t, testCase.Expected.OrgId, orgID)
			assert.Equal(t, testCase.Expected.Output, output)
			assert.Equal(t, testCase.Expected.ClientVersion, clientVersion)
		} else {
			require.NoError(t, err)
			assert.Equal(t, testCase.Expected.OrgId, orgID)
			assert.Equal(t, testCase.Expected.Output, output)
			assert.Equal(t, testCase.Expected.ClientVersion, clientVersion)
		}
	}
}

func TestUpdateAgent(t *testing.T) {
	const testOrgID = "12345"
	testID := uuid.MustParse("658700b8-005b-11ee-9e09-482ae3863d30")
	testTitle := pointy.String("My Example Domain Title")
	testDescription := "My Example Domain Description"
	testXRHID := identity.XRHID{
		Identity: identity.Identity{
			OrgID: testOrgID,
			Type:  "user",
			User:  identity.User{},
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
	testWrongTypeBody := api_public.Domain{
		AutoEnrollmentEnabled: pointy.Bool(true),
		Title:                 testTitle,
		Description:           pointy.String(testDescription),
		DomainName:            "mydomain.example",
		DomainId:              &testID,
		DomainType:            "aninvalidtype",
	}
	testBody := api_public.Domain{
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
	i := NewDomainInteractor()

	// Get an error in guards
	orgID, xrhidmVersion, domain, err := i.UpdateAgent(nil, uuid.Nil, nil, nil)
	assert.EqualError(t, err, "'xrhid' is nil")
	assert.Equal(t, "", orgID)
	assert.Nil(t, xrhidmVersion)
	assert.Nil(t, domain)

	// Error retrieving ipa-hcc version information
	orgID, xrhidmVersion, domain, err = i.UpdateAgent(&testXRHID, testID, &testBadParams, &testBody)
	assert.EqualError(t, err, "'X-Rh-Idm-Version' is invalid")
	assert.Equal(t, "", orgID)
	assert.Nil(t, xrhidmVersion)
	assert.Nil(t, domain)

	// Error because of wrongtype
	orgID, xrhidmVersion, domain, err = i.UpdateAgent(&testXRHID, testID, &testParams, &testWrongTypeBody)
	assert.EqualError(t, err, "Unsupported domain_type='aninvalidtype'")
	assert.Equal(t, "", orgID)
	assert.Nil(t, xrhidmVersion)
	assert.Nil(t, domain)

	// Success result
	orgID, xrhidmVersion, domain, err = i.UpdateAgent(&testXRHID, testID, &testParams, &testBody)
	assert.NoError(t, err)
	assert.Equal(t, testOrgID, orgID)
	assert.Equal(t, testID, *testBody.DomainId)
	require.NotNil(t, xrhidmVersion)
	require.NotNil(t, domain)
}

func assertListEqualError(t *testing.T, err error, msg string, orgID string, offset int, limit int) {
	assert.EqualError(t, err, msg)
	assert.Equal(t, "", orgID)
	assert.Equal(t, -1, offset)
	assert.Equal(t, -1, limit)
}

func TestList(t *testing.T) {
	i := NewDomainInteractor()
	testOrgID := "12345"

	// xrhid is nil
	orgID, offset, limit, err := i.List(nil, nil)
	assertListEqualError(t, err, "'xrhid' is nil", orgID, offset, limit)

	// params is nil
	xrhid := identity.XRHID{
		Identity: identity.Identity{
			OrgID: testOrgID,
		},
	}
	orgID, offset, limit, err = i.List(&xrhid, nil)
	assertListEqualError(t, err, "'params' is nil", orgID, offset, limit)

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

// --------- Private methods -----------

func TestGuardRegister(t *testing.T) {
	var err error

	i := domainInteractor{}

	err = i.guardRegister(nil, nil, nil)
	assert.EqualError(t, err, "'xrhid' is nil")

	xrhid := &identity.XRHID{}
	err = i.guardRegister(xrhid, nil, nil)
	assert.EqualError(t, err, "'params' is nil")

	params := &api_public.RegisterDomainParams{}
	err = i.guardRegister(xrhid, params, nil)
	assert.EqualError(t, err, "'body' is nil")

	body := &public.Domain{}
	err = i.guardRegister(xrhid, params, body)
	assert.NoError(t, err)
}

func TestGuardUpdate(t *testing.T) {
	var err error

	i := domainInteractor{}

	err = i.guardUpdate(nil, uuid.Nil, nil)
	assert.EqualError(t, err, "'xrhid' is nil")

	xrhid := &identity.XRHID{}
	err = i.guardUpdate(xrhid, uuid.Nil, nil)
	assert.EqualError(t, err, "'UUID' is invalid")

	UUID := uuid.MustParse("b0264600-005c-11ee-ba48-482ae3863d30")
	err = i.guardUpdate(xrhid, UUID, nil)
	assert.EqualError(t, err, "'body' is nil")

	body := &public.Domain{}
	err = i.guardUpdate(xrhid, UUID, body)
	assert.NoError(t, err)
}

func TestCommonRegisterUpdate(t *testing.T) {
	testOrgID := "12345"
	testID := uuid.MustParse("c95c6e74-005c-11ee-82b5-482ae3863d30")
	testTitle := pointy.String("My Example Domain Title")
	i := domainInteractor{}
	assert.Panics(t, func() {
		i.commonRegisterUpdate("", uuid.Nil, nil)
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

	domain, err := i.commonRegisterUpdate(testOrgID, testID, &testWrongTypeBody)
	assert.EqualError(t, err, "Unsupported domain_type='wrongtype'")
	assert.Nil(t, domain)

	// Success case
	domain, err = i.commonRegisterUpdate(testOrgID, testID, &testBody)
	assert.NoError(t, err)
	assert.NotNil(t, domain)
}

func TestGetByID(t *testing.T) {
	i := NewDomainInteractor()

	orgID, err := i.GetByID(nil, nil)
	assert.EqualError(t, err, "'xrhid' is nil")
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
	assert.EqualError(t, err, "'params' is nil")
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
		i.registerOrUpdateRhelIdmLocations(item.Given, ipa)
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
				User:  identity.User{},
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
				Err: fmt.Errorf("'xrhid' is nil"),
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
				Err: fmt.Errorf("'params' is nil"),
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
				Err: fmt.Errorf("'body' is nil"),
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
