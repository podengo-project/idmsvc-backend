package interactor

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/header"
	"github.com/hmsidm/internal/api/public"
	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
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

func TestCreate(t *testing.T) {
	notValidBefore := time.Now()
	notValidAfter := time.Now().Add(time.Hour * 24)
	rhsmId := uuid.New().String()
	type TestCaseGiven struct {
		Params *api_public.CreateDomainParams
		Body   *api_public.CreateDomain
	}
	type TestCaseExpected struct {
		Err   error
		OrgId string
		Out   *model.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "nil for the 'params'",
			Given: TestCaseGiven{
				Params: nil,
				Body:   &api_public.CreateDomain{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'params' cannot be nil"),
				Out: nil,
			},
		},
		{
			Name: "nil for the 'body'",
			Given: TestCaseGiven{
				Params: &api_public.CreateDomainParams{},
				Body:   nil,
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'body' cannot be nil"),
				Out: nil,
			},
		},
		{
			Name: "nil for the returned Model",
			Given: TestCaseGiven{
				Params: &api_public.CreateDomainParams{},
				Body:   &api_public.CreateDomain{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("X-Rh-Identity content cannot be an empty string"),
				Out: nil,
			},
		},
		{
			Name: "success case",
			Given: TestCaseGiven{
				Params: &api_public.CreateDomainParams{
					XRhIdentity: header.EncodeXRHID(
						&identity.XRHID{
							Identity: identity.Identity{
								OrgID: "12345",
								Internal: identity.Internal{
									OrgID: "12345",
								},
							},
						},
					),
				},
				Body: &api_public.CreateDomain{
					AutoEnrollmentEnabled: true,
					Type:                  api_public.CreateDomainTypeRhelIdm,
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &model.Domain{
					OrgId:                 "12345",
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
				},
			},
		},
		{
			Name: "success case - not empty ca list",
			Given: TestCaseGiven{
				Params: &api_public.CreateDomainParams{
					XRhIdentity: header.EncodeXRHID(
						&identity.XRHID{
							Identity: identity.Identity{
								OrgID: "12345",
								Internal: identity.Internal{
									OrgID: "12345",
								},
							},
						},
					),
				},
				Body: &api_public.CreateDomain{
					AutoEnrollmentEnabled: true,
					Type:                  api_public.CreateDomainTypeRhelIdm,
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &model.Domain{
					OrgId:                 "12345",
					DomainName:            nil,
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName: pointy.String("DOMAIN.EXAMPLE"),
						CaCerts: []model.IpaCert{
							{
								Nickname:       "DOMAIN.EXAMPLE IPA CA",
								Issuer:         "CN=Certificate Authority,O=DOMAIN.EXAMPLE",
								Subject:        "CN=Certificate Authority,O=DOMAIN.EXAMPLE",
								SerialNumber:   "1",
								NotValidAfter:  notValidAfter,
								NotValidBefore: notValidBefore,
								Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
						Servers: []model.IpaServer{
							{
								FQDN:                "server1.domain.example",
								CaServer:            true,
								HCCEnrollmentServer: true,
								HCCUpdateServer:     true,
								PKInitServer:        true,
								RHSMId:              rhsmId,
							},
						},
						RealmDomains: pq.StringArray{"server1.domain.example", "server2.domain.example"},
					},
				},
			},
		},
	}

	component := NewDomainInteractor()
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		orgId, data, err := component.Create(testCase.Given.Params, testCase.Given.Body)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			require.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Equal(t, "", orgId)
			assert.Nil(t, data)
		} else {
			assert.NoError(t, err)
			require.NotNil(t, testCase.Expected.Out)

			assert.Equal(t, *testCase.Expected.Out.AutoEnrollmentEnabled, *data.AutoEnrollmentEnabled)
			assert.Equal(t, testCase.Expected.Out.OrgId, data.OrgId)
			assert.Nil(t, testCase.Expected.Out.DomainName)
			assert.Equal(t, *testCase.Expected.Out.Type, *data.Type)
		}
	}
}

func TestHelperDomainTypeToUint(t *testing.T) {
	var (
		result uint
	)

	result = helperDomainTypeToUint("")
	assert.Equal(t, model.DomainTypeUndefined, result)

	result = helperDomainTypeToUint(public.DomainTypeRhelIdm)
	assert.Equal(t, model.DomainTypeIpa, result)
}

func TestRegisterIpa(t *testing.T) {
	const (
		cn        = "21258fc8-c755-11ed-afc4-482ae3863d30"
		requestID = "TW9uIE1hciAyMCAyMDo1Mzoz"
		token     = "3fa8caf6-c759-11ed-99dd-482ae3863d30"
		rhsmId    = "cf26cd96-c75d-11ed-ae20-482ae3863d30"
		orgID     = "12345"
	)
	var (
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
			IPAHCCVersion: "0.7",
			IPAVersion:    "4.10.0-8.el9_1",
		}
		xrhidSystemBase64     = header.EncodeXRHID(&xrhidSystem)
		paramsNoClientVersion = &api_public.RegisterDomainParams{
			XRhIdentity:             xrhidSystemBase64,
			XRhInsightsRequestId:    requestID,
			XRhIdmRegistrationToken: token,
		}
		params = &api_public.RegisterDomainParams{
			XRhIdentity:             xrhidSystemBase64,
			XRhInsightsRequestId:    requestID,
			XRhIdmRegistrationToken: token,
			XRhIdmVersion:           "eyJpcGEtaGNjIjogIjAuNyIsICJpcGEiOiAiNC4xMC4wLTguZWw5XzEifQo=",
		}
		notValidBefore = time.Now().UTC()
		notValidAfter  = notValidBefore.Add(24 * time.Hour)
	)
	type TestCaseGiven struct {
		XRHID  *identity.XRHID
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
			Name: "xrhid is nil",
			Given: TestCaseGiven{
				XRHID:  nil,
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
			Name: "param is nil",
			Given: TestCaseGiven{
				XRHID:  &xrhidSystem,
				Params: nil,
				Body:   nil,
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         fmt.Errorf("'params' is nil"),
			},
		},
		{
			Name: "body is nil",
			Given: TestCaseGiven{
				XRHID:  &xrhidSystem,
				Params: params,
				Body:   nil,
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         fmt.Errorf("'body' is nil"),
			},
		},
		{
			Name: "No X-Rh-IDM-Version",
			Given: TestCaseGiven{
				XRHID:  &xrhidSystem,
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
				Params: params,
				Body: &api_public.Domain{
					Type: "somethingwrong",
					RhelIdm: &api_public.DomainIpa{
						RealmDomains: nil,
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         fmt.Errorf("'Type=somethingwrong' is invalid"),
			},
		},
		{
			Name: "Empty slices",
			Given: TestCaseGiven{
				XRHID:  &xrhidSystem,
				Params: params,
				Body: &api_public.Domain{
					Title:       "My Example Domain",
					Description: "My Example Domain Description",
					DomainName:  "mydomain.example",
					Type:        api_public.DomainType(api_public.DomainTypeRhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName: "",
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("My Example Domain"),
					Description:           pointy.String("My Example Domain Description"),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String(""),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
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
				Params: params,
				Body: &api_public.Domain{
					Title:       "My Example Domain",
					Description: "My Example Domain Description",
					DomainName:  "mydomain.example",
					Type:        api_public.DomainType(api_public.DomainTypeRhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName: "MYDOMAIN.EXAMPLE",
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("My Example Domain"),
					Description:           pointy.String("My Example Domain Description"),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
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
				Params: params,
				Body: &api_public.Domain{
					Title:       "My Example Domain",
					Description: "My Example Domain Description",
					DomainName:  "mydomain.example",
					Type:        api_public.DomainType(api_public.DomainTypeRhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName:    "MYDOMAIN.EXAMPLE",
						RealmDomains: []string{"server.domain.example"},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("My Example Domain"),
					Description:           pointy.String("My Example Domain Description"),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
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
				Params: params,
				Body: &api_public.Domain{
					Title:       "My Example Domain",
					Description: "My Example Domain Description",
					DomainName:  "mydomain.example",
					Type:        api_public.DomainType(api_public.DomainTypeRhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName: "MYDOMAIN.EXAMPLE",
						CaCerts: []api_public.DomainIpaCert{
							{
								Nickname:       "MYDOMAIN.EXAMPLE IPA CA",
								SerialNumber:   "1",
								Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								NotValidBefore: notValidBefore,
								NotValidAfter:  notValidAfter,
								Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("My Example Domain"),
					Description:           pointy.String("My Example Domain Description"),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName: pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts: []model.IpaCert{
							{
								Nickname:       "MYDOMAIN.EXAMPLE IPA CA",
								SerialNumber:   "1",
								Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								NotValidBefore: notValidBefore,
								NotValidAfter:  notValidAfter,
								Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
						Servers:      []model.IpaServer{},
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
				Params: params,
				Body: &api_public.Domain{
					Title:       "My Example Domain",
					Description: "My Example Domain Description",
					DomainName:  "mydomain.example",
					Type:        api_public.DomainType(api_public.DomainTypeRhelIdm),
					RhelIdm: &api_public.DomainIpa{
						RealmName: "MYDOMAIN.EXAMPLE",
						Servers: []api_public.DomainIpaServer{
							{
								Fqdn:                  "server.mydomain.example",
								SubscriptionManagerId: rhsmId,
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
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					DomainName:            pointy.String("mydomain.example"),
					Title:                 pointy.String("My Example Domain"),
					Description:           pointy.String("My Example Domain Description"),
					AutoEnrollmentEnabled: pointy.Bool(false),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						RealmName: pointy.String("MYDOMAIN.EXAMPLE"),
						CaCerts:   []model.IpaCert{},
						Servers: []model.IpaServer{
							{
								FQDN:                "server.mydomain.example",
								RHSMId:              rhsmId,
								CaServer:            true,
								HCCEnrollmentServer: true,
								HCCUpdateServer:     true,
								PKInitServer:        true,
							},
						},
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
		orgId, clientVersion, output, err := i.Register(testCase.Given.XRHID, testCase.Given.Params, testCase.Given.Body)
		if testCase.Expected.Error != nil {
			assert.EqualError(t, err, testCase.Expected.Error.Error())
			assert.Equal(t, testCase.Expected.OrgId, orgId)
			assert.Equal(t, testCase.Expected.Output, output)
			assert.Equal(t, testCase.Expected.ClientVersion, clientVersion)
		} else {
			require.NoError(t, err)
			assert.Equal(t, testCase.Expected.OrgId, orgId)
			assert.Equal(t, testCase.Expected.Output, output)
			assert.Equal(t, testCase.Expected.ClientVersion, clientVersion)
		}
	}
}
