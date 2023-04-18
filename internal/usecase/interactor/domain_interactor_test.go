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
					DomainName:            "domain.example",
					DomainType:            api_public.CreateDomainDomainTypeIpa,
					Ipa: api_public.CreateDomainIpa{
						RealmName:    "DOMAIN.EXAMPLE",
						CaCerts:      []api_public.CreateDomainIpaCert{},
						Servers:      &[]api_public.CreateDomainIpaServer{},
						RealmDomains: []string{},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &model.Domain{
					OrgId:                 "12345",
					DomainName:            pointy.String("domain.example"),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String("DOMAIN.EXAMPLE"),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
						RealmDomains: pq.StringArray{},
					},
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
					DomainName:            "domain.example",
					DomainType:            api_public.CreateDomainDomainTypeIpa,
					Ipa: api_public.CreateDomainIpa{
						RealmName: "DOMAIN.EXAMPLE",
						CaCerts: []api_public.CreateDomainIpaCert{
							{
								Nickname:       pointy.String("DOMAIN.EXAMPLE IPA CA"),
								Issuer:         pointy.String("CN=Certificate Authority,O=DOMAIN.EXAMPLE"),
								Subject:        pointy.String("CN=Certificate Authority,O=DOMAIN.EXAMPLE"),
								SerialNumber:   pointy.String("1"),
								NotValidAfter:  &notValidAfter,
								NotValidBefore: &notValidBefore,
								Pem:            pointy.String("-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n"),
							},
						},
						Servers: &[]api_public.CreateDomainIpaServer{
							{
								Fqdn:                  "server1.domain.example",
								CaServer:              true,
								HccEnrollmentServer:   true,
								HccUpdateServer:       true,
								PkinitServer:          true,
								SubscriptionManagerId: rhsmId,
							},
						},
						RealmDomains: []string{
							"server1.domain.example",
							"server2.domain.example",
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &model.Domain{
					OrgId:                 "12345",
					DomainName:            pointy.String("domain.example"),
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
			assert.Equal(t, *testCase.Expected.Out.DomainName, *data.DomainName)
			assert.Equal(t, *testCase.Expected.Out.Type, *data.Type)
			assert.Equal(t,
				*testCase.Expected.Out.IpaDomain.RealmName,
				*data.IpaDomain.RealmName)
			assert.Equal(t,
				testCase.Expected.Out.IpaDomain.CaCerts,
				data.IpaDomain.CaCerts)
			assert.Equal(t,
				testCase.Expected.Out.IpaDomain.Servers,
				data.IpaDomain.Servers)
			assert.Equal(t,
				testCase.Expected.Out.IpaDomain.RealmDomains,
				data.IpaDomain.RealmDomains)
		}
	}
}

func TestHelperDomainTypeToUint(t *testing.T) {
	var (
		result uint
	)

	result = helperDomainTypeToUint("")
	assert.Equal(t, model.DomainTypeUndefined, result)

	result = helperDomainTypeToUint(public.CreateDomainDomainTypeIpa)
	assert.Equal(t, model.DomainTypeIpa, result)
}

func TestFillCert(t *testing.T) {
	notValidBefore := time.Now()
	notValidAfter := time.Now().Add(time.Hour * 24)

	type TestCaseGiven struct {
		To   *model.IpaCert
		From *api_public.CreateDomainIpaCert
	}
	type TestCaseExpected struct {
		Err error
		To  model.IpaCert
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}

	testCases := []TestCase{
		{
			Name: "'to' cannot be nil",
			Given: TestCaseGiven{
				To:   nil,
				From: nil,
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'to' cannot be nil"),
				To:  model.IpaCert{},
			},
		},
		{
			Name: "'from' cannot be nil",
			Given: TestCaseGiven{
				To:   &model.IpaCert{},
				From: nil,
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'from' cannot be nil"),
				To:  model.IpaCert{},
			},
		},

		{
			Name: "Fill all the fields",
			Given: TestCaseGiven{
				To: &model.IpaCert{},
				From: &api_public.CreateDomainIpaCert{
					Nickname:       pointy.String("Nickname"),
					Issuer:         pointy.String("Issuer"),
					Subject:        pointy.String("Subject"),
					NotValidBefore: &notValidBefore,
					NotValidAfter:  &notValidAfter,
					SerialNumber:   pointy.String("1"),
					Pem:            pointy.String("-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n"),
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				To: model.IpaCert{
					Nickname:       "Nickname",
					Issuer:         "Issuer",
					Subject:        "Subject",
					NotValidBefore: notValidBefore,
					NotValidAfter:  notValidAfter,
					SerialNumber:   "1",
					Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				},
			},
		},
		{
			Name: "Fill empty all the fields",
			Given: TestCaseGiven{
				To: &model.IpaCert{
					Nickname:       "Nickname",
					Issuer:         "Issuer",
					Subject:        "Subject",
					NotValidBefore: notValidBefore,
					NotValidAfter:  notValidAfter,
					SerialNumber:   "1",
					Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				},
				From: &api_public.CreateDomainIpaCert{
					Nickname:       nil,
					Issuer:         nil,
					Subject:        nil,
					NotValidBefore: nil,
					NotValidAfter:  nil,
					SerialNumber:   nil,
					Pem:            nil,
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				To: model.IpaCert{
					Nickname:       "",
					Issuer:         "",
					Subject:        "",
					NotValidBefore: time.Time{},
					NotValidAfter:  time.Time{},
					SerialNumber:   "",
					Pem:            "",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := domainInteractor{}
		err := obj.FillCert(testCase.Given.To, testCase.Given.From)
		if testCase.Expected.Err != nil {
			assert.EqualError(t, err, testCase.Expected.Err.Error())
		} else {
			assert.Equal(t, testCase.Expected.To.Nickname, testCase.Given.To.Nickname)
			assert.Equal(t, testCase.Expected.To.Issuer, testCase.Given.To.Issuer)
			assert.Equal(t, testCase.Expected.To.Subject, testCase.Given.To.Subject)
			assert.Equal(t, testCase.Expected.To.NotValidBefore, testCase.Given.To.NotValidBefore)
			assert.Equal(t, testCase.Expected.To.NotValidAfter, testCase.Given.To.NotValidAfter)
			assert.Equal(t, testCase.Expected.To.SerialNumber, testCase.Given.To.SerialNumber)
			assert.Equal(t, testCase.Expected.To.Pem, testCase.Given.To.Pem)
		}
	}
}

func TestFillServer(t *testing.T) {
	var (
		err error
	)
	obj := domainInteractor{}

	err = obj.FillServer(nil, nil)
	assert.EqualError(t, err, "'to' cannot be nil")

	err = obj.FillServer(&model.IpaServer{}, nil)
	assert.EqualError(t, err, "'from' cannot be nil")

	to := model.IpaServer{}
	err = obj.FillServer(&to, &api_public.CreateDomainIpaServer{
		Fqdn:                  "server1.mydomain.example",
		SubscriptionManagerId: "001a2314-bf37-11ed-a42a-482ae3863d30",
		PkinitServer:          true,
		CaServer:              true,
		HccEnrollmentServer:   true,
		HccUpdateServer:       true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "server1.mydomain.example", to.FQDN)
	assert.Equal(t, "001a2314-bf37-11ed-a42a-482ae3863d30", to.RHSMId)
	assert.Equal(t, true, to.PKInitServer)
	assert.Equal(t, true, to.CaServer)
	assert.Equal(t, true, to.HCCEnrollmentServer)
	assert.Equal(t, true, to.HCCUpdateServer)
}

func TestRegisterIpa(t *testing.T) {
	const (
		cn        = "21258fc8-c755-11ed-afc4-482ae3863d30"
		requestID = "TW9uIE1hciAyMCAyMDo1Mzoz"
		token     = "3fa8caf6-c759-11ed-99dd-482ae3863d30"
		rhsmId    = "cf26cd96-c75d-11ed-ae20-482ae3863d30"
	)
	var (
		xrhidSystem = identity.XRHID{
			Identity: identity.Identity{
				Type: "System",
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
			XRhIDMRegistrationToken: token,
		}
		params = &api_public.RegisterDomainParams{
			XRhIdentity:             xrhidSystemBase64,
			XRhInsightsRequestId:    requestID,
			XRhIDMRegistrationToken: token,
			XRhIdmVersion:           "eyJpcGEtaGNjIjogIjAuNyIsICJpcGEiOiAiNC4xMC4wLTguZWw5XzEifQo=",
		}
		notValidBefore = time.Now().UTC()
		notValidAfter  = notValidBefore.Add(24 * time.Hour)
	)
	type TestCaseGiven struct {
		XRHID  *identity.XRHID
		Params *api_public.RegisterDomainParams
		Body   *public.RegisterDomain
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
				Body:   &api_public.RegisterDomain{},
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
				Body: &api_public.RegisterDomain{
					Type: "somethingwrong",
					RhelIdm: api_public.RegisterDomainIpa{
						RealmDomains: nil,
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: nil,
				Output:        nil,
				Error:         fmt.Errorf("Type='somethingwrong' is invalid"),
			},
		},
		{
			Name: "Empty slices",
			Given: TestCaseGiven{
				XRHID:  &xrhidSystem,
				Params: params,
				Body: &api_public.RegisterDomain{
					Type: api_public.RhelIdm,
					RhelIdm: api_public.RegisterDomainIpa{
						RealmDomains: nil,
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
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
				Body: &api_public.RegisterDomain{
					Type: api_public.RhelIdm,
					RhelIdm: api_public.RegisterDomainIpa{
						RealmDomains: []string{"server.domain.example"},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
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
				Body: &api_public.RegisterDomain{
					Type: api_public.RhelIdm,
					RhelIdm: api_public.RegisterDomainIpa{
						CaCerts: []api_public.CreateDomainIpaCert{
							{
								Nickname:       pointy.String("MYDOMAIN.EXAMPLE IPA CA"),
								SerialNumber:   pointy.String("1"),
								Issuer:         pointy.String("CN=Certificate Authority,O=MYDOMAIN.EXAMPLE"),
								Subject:        pointy.String("CN=Certificate Authority,O=MYDOMAIN.EXAMPLE"),
								NotValidBefore: &notValidBefore,
								NotValidAfter:  &notValidAfter,
								Pem:            pointy.String("-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n"),
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				OrgId:         "",
				ClientVersion: clientVersionParsed,
				Output: &model.Domain{
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
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
				Body: &api_public.RegisterDomain{
					Type: api_public.RhelIdm,
					RhelIdm: api_public.RegisterDomainIpa{
						Servers: []api_public.CreateDomainIpaServer{
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
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						CaCerts: []model.IpaCert{},
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
