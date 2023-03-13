package interactor

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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
					XRhIdentity: EncodeIdentity(
						&identity.Identity{
							OrgID: "12345",
							Internal: identity.Internal{
								OrgID: "12345",
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
					DomainType:            pointy.Uint(model.DomainTypeIpa),
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
					XRhIdentity: EncodeIdentity(
						&identity.Identity{
							OrgID: "12345",
							Internal: identity.Internal{
								OrgID: "12345",
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
								Fqdn:                "server1.domain.example",
								CaServer:            true,
								HccEnrollmentServer: true,
								PkinitServer:        true,
								RhsmId:              rhsmId,
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
					DomainType:            pointy.Uint(model.DomainTypeIpa),
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
			assert.Equal(t, *testCase.Expected.Out.DomainType, *data.DomainType)
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
		Fqdn:                "server1.mydomain.example",
		RhsmId:              "001a2314-bf37-11ed-a42a-482ae3863d30",
		PkinitServer:        true,
		CaServer:            true,
		HccEnrollmentServer: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "server1.mydomain.example", to.FQDN)
	assert.Equal(t, "001a2314-bf37-11ed-a42a-482ae3863d30", to.RHSMId)
	assert.Equal(t, true, to.PKInitServer)
	assert.Equal(t, true, to.CaServer)
	assert.Equal(t, true, to.HCCEnrollmentServer)
}
