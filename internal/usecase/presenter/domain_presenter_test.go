package presenter

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewTodoPresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewDomainPresenter()
	})
}

type mynewerror struct{}

func (e *mynewerror) Error() string {
	return "mynewerror"
}

func TestGet(t *testing.T) {
	testUuid := uuid.New()
	type TestCaseGiven struct {
		Input  *model.Domain
		Output *public.ReadDomainResponse
	}
	type TestCaseExpected struct {
		Err    error
		Output *public.ReadDomainResponse
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "error when 'in' is nil",
			Given: TestCaseGiven{
				Input: nil,
			},
			Expected: TestCaseExpected{
				Err:    fmt.Errorf("'domain' is nil"),
				Output: nil,
			},
		},
		{
			Name: "Success case",
			Given: TestCaseGiven{
				Input: &model.Domain{
					Model:                 gorm.Model{ID: 1},
					OrgId:                 "12345",
					DomainUuid:            testUuid,
					DomainName:            pointy.String("domain.example"),
					Type:                  pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						RealmName:    pointy.String("DOMAIN.EXAMPLE"),
						CaCerts:      []model.IpaCert{},
						Servers:      []model.IpaServer{},
						RealmDomains: pq.StringArray{"domain.example"},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Output: &public.ReadDomainResponse{
					AutoEnrollmentEnabled: true,
					DomainUuid:            testUuid.String(),
					DomainName:            "domain.example",
					Type:                  public.DomainType(model.DomainTypeString(model.DomainTypeIpa)),
					RhelIdm: &public.DomainIpa{
						RealmName:    "DOMAIN.EXAMPLE",
						CaCerts:      []public.DomainIpaCert{},
						Servers:      []public.DomainIpaServer{},
						RealmDomains: []string{"domain.example"},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := domainPresenter{}
		output, err := obj.Get(testCase.Given.Input)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			assert.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Nil(t, output)
		} else {
			assert.NoError(t, err)
			assert.Equal(t,
				testCase.Expected.Output.DomainUuid,
				output.DomainUuid)
			assert.Equal(t,
				testCase.Expected.Output.DomainName,
				output.DomainName)
			assert.Equal(t,
				testCase.Expected.Output.Type,
				output.Type)
			assert.Equal(t,
				testCase.Expected.Output.AutoEnrollmentEnabled,
				output.AutoEnrollmentEnabled)
			assert.Equal(t,
				testCase.Expected.Output.RhelIdm.RealmName,
				output.RhelIdm.RealmName)
			assert.Equal(t,
				testCase.Expected.Output.RhelIdm.CaCerts,
				output.RhelIdm.CaCerts)
			assert.Equal(t,
				testCase.Expected.Output.RhelIdm.Servers,
				output.RhelIdm.Servers)
		}
	}
}

func TestCreate(t *testing.T) {
	type TestCaseExpected struct {
		Response *public.CreateDomainResponse
		Err      error
	}
	type TestCase struct {
		Name     string
		Given    *model.Domain
		Expected TestCaseExpected
	}

	testCases := []TestCase{
		{
			Name:  "domain is nil",
			Given: nil,
			Expected: TestCaseExpected{
				Response: nil,
				Err:      fmt.Errorf("'domain' is nil"),
			},
		},
		{
			Name: "AutoEnrollmentEnabled is nil",
			Given: &model.Domain{
				AutoEnrollmentEnabled: nil,
			},
			Expected: TestCaseExpected{
				Response: nil,
				Err:      fmt.Errorf("'AutoEnrollmentEnabled' is nil"),
			},
		},
		{
			Name: "Type is nil",
			Given: &model.Domain{
				AutoEnrollmentEnabled: pointy.Bool(true),
				Type:                  nil,
			},
			Expected: TestCaseExpected{
				Response: nil,
				Err:      fmt.Errorf("'Type' is nil"),
			},
		},
		{
			Name: "IpaDomain is nil",
			Given: &model.Domain{
				AutoEnrollmentEnabled: pointy.Bool(true),
				DomainName:            pointy.String("domain.example"),
				Type:                  pointy.Uint(model.DomainTypeIpa),
				IpaDomain:             nil,
			},
			Expected: TestCaseExpected{
				Response: &public.Domain{
					DomainUuid:            "00000000-0000-0000-0000-000000000000",
					AutoEnrollmentEnabled: true,
					DomainName:            "domain.example",
					Type:                  public.DomainTypeRhelIdm,
					RhelIdm:               nil,
				},
				Err: nil,
			},
		},
		// {
		// 	Name: "RealmName is nil",
		// 	Given: &model.Domain{
		// 		AutoEnrollmentEnabled: pointy.Bool(true),
		// 		DomainName:            pointy.String("domain.example"),
		// 		Type:                  pointy.Uint(model.DomainTypeIpa),
		// 		IpaDomain: &model.Ipa{
		// 			RealmName: nil,
		// 		},
		// 	},
		// 	Expected: TestCaseExpected{
		// 		Response: nil,
		// 		Err:      fmt.Errorf("'RealmName' is nil"),
		// 	},
		// },
		// {
		// 	Name: "CaCerts is nil",
		// 	Given: &model.Domain{
		// 		AutoEnrollmentEnabled: pointy.Bool(true),
		// 		DomainName:            pointy.String("domain.example"),
		// 		Type:                  pointy.Uint(model.DomainTypeIpa),
		// 		IpaDomain: &model.Ipa{
		// 			RealmName: pointy.String("DOMAIN.EXAMPLE"),
		// 			CaCerts:   nil,
		// 		},
		// 	},
		// 	Expected: TestCaseExpected{
		// 		Response: nil,
		// 		Err:      fmt.Errorf("'CaCerts' is nil"),
		// 	},
		// },
		// {
		// 	Name: "Servers is nil",
		// 	Given: &model.Domain{
		// 		AutoEnrollmentEnabled: pointy.Bool(true),
		// 		DomainName:            pointy.String("domain.example"),
		// 		Type:                  pointy.Uint(model.DomainTypeIpa),
		// 		IpaDomain: &model.Ipa{
		// 			RealmName: pointy.String("DOMAIN.EXAMPLE"),
		// 			CaCerts:   []model.IpaCert{},
		// 			Servers:   nil,
		// 		},
		// 	},
		// 	Expected: TestCaseExpected{
		// 		Response: nil,
		// 		Err:      fmt.Errorf("'Servers' is nil"),
		// 	},
		// },
		{
			Name: "Success scenario",
			Given: &model.Domain{
				AutoEnrollmentEnabled: pointy.Bool(true),
				DomainName:            pointy.String("domain.example"),
				Type:                  pointy.Uint(model.DomainTypeIpa),
				// IpaDomain: &model.Ipa{
				// 	RealmName: pointy.String("DOMAIN.EXAMPLE"),
				// 	CaCerts:   []model.IpaCert{},
				// 	Servers:   []model.IpaServer{},
				// },
			},
			Expected: TestCaseExpected{
				Response: &public.CreateDomainResponse{
					DomainName:            "domain.example",
					AutoEnrollmentEnabled: true,
					Type:                  model.DomainTypeIpaString,
					DomainUuid:            "00000000-0000-0000-0000-000000000000",
					// RhelIdm: &public.DomainIpa{
					// 	RealmName: "DOMAIN.EXAMPLE",
					// 	CaCerts:   []public.DomainIpaCert{},
					// 	Servers:   []public.DomainIpaServer{},
					// },
				},
				Err: nil,
			},
		},
		// {
		// 	Name: "Success scenario with DomainName equals to nil",
		// 	Given: &model.Domain{
		// 		AutoEnrollmentEnabled: pointy.Bool(true),
		// 		DomainName:            nil,
		// 		Type:                  pointy.Uint(model.DomainTypeIpa),
		// 		IpaDomain: &model.Ipa{
		// 			RealmName: pointy.String("DOMAIN.EXAMPLE"),
		// 			CaCerts:   []model.IpaCert{},
		// 			Servers:   []model.IpaServer{},
		// 		},
		// 	},
		// 	Expected: TestCaseExpected{
		// 		Response: &public.CreateDomainResponse{
		// 			DomainName:            "",
		// 			AutoEnrollmentEnabled: true,
		// 			Type:                  model.DomainTypeIpaString,
		// 			DomainUuid:            "00000000-0000-0000-0000-000000000000",
		// 			RhelIdm: &public.DomainIpa{
		// 				RealmName: "DOMAIN.EXAMPLE",
		// 				CaCerts:   []public.DomainIpaCert{},
		// 				Servers:   []public.DomainIpaServer{},
		// 			},
		// 		},
		// 		Err: nil,
		// 	},
		// },
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		obj := domainPresenter{}
		response, err := obj.Create(testCase.Given)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			assert.EqualError(t, err, testCase.Expected.Err.Error())
			assert.Nil(t, response)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, testCase.Expected.Response.DomainName, response.DomainName)
			assert.Equal(t, testCase.Expected.Response.AutoEnrollmentEnabled, response.AutoEnrollmentEnabled)
			assert.Equal(t, testCase.Expected.Response.Type, response.Type)
			assert.Equal(t, testCase.Expected.Response.DomainUuid, response.DomainUuid)
			if testCase.Expected.Response.RhelIdm == nil {
				assert.Nil(t, response.RhelIdm)
			} else {
				require.NotNil(t, response.RhelIdm)
				assert.Equal(t, testCase.Expected.Response.RhelIdm.CaCerts, response.RhelIdm.CaCerts)
				assert.Equal(t, testCase.Expected.Response.RhelIdm.Servers, response.RhelIdm.Servers)
				assert.Equal(t, testCase.Expected.Response.RhelIdm.RealmName, response.RhelIdm.RealmName)
			}
		}
	}

}

func TestFillRhelmIdmCertsPanics(t *testing.T) {
	var err error
	p := &domainPresenter{}

	assert.Panics(t, func() {
		p.fillRhelIdmCerts(nil, nil)
	})

	domain := &model.Domain{}
	err = p.fillRhelIdmCerts(nil, domain)
	assert.NoError(t, err)

	domain.IpaDomain = &model.Ipa{}
	assert.Panics(t, func() {
		p.fillRhelIdmCerts(nil, domain)
	})

	output := &public.Domain{}
	assert.Panics(t, func() {
		p.fillRhelIdmCerts(output, domain)
	})

	domain.IpaDomain.CaCerts = []model.IpaCert{}
	assert.Panics(t, func() {
		p.fillRhelIdmCerts(output, domain)
	})

	// Minimal success
	output.RhelIdm = &public.DomainIpa{}
	err = p.fillRhelIdmCerts(output, domain)
	assert.NoError(t, err)
}

func TestFillRhelmIdmCerts(t *testing.T) {
	notValidBefore := time.Now()
	notValidAfter := notValidBefore.Add(time.Hour * 24)
	type TestCaseGiven struct {
		To   *public.Domain
		From *model.Domain
	}
	type TestCaseExpected struct {
		Err error
		To  *public.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "Full success copy",
			Given: TestCaseGiven{
				To: &public.Domain{
					Type:    public.DomainTypeRhelIdm,
					RhelIdm: &public.DomainIpa{},
				},
				From: &model.Domain{
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						CaCerts: []model.IpaCert{
							{
								Nickname:       "MYDOMAIN.EXAMPLE.IPA CA",
								Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								NotValidBefore: notValidBefore,
								NotValidAfter:  notValidAfter,
								SerialNumber:   "1",
								Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				To: &public.Domain{
					Type: public.DomainTypeRhelIdm,
					RhelIdm: &public.DomainIpa{
						CaCerts: []public.DomainIpaCert{
							{
								Nickname:       "MYDOMAIN.EXAMPLE.IPA CA",
								Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								NotValidBefore: notValidBefore,
								NotValidAfter:  notValidAfter,
								SerialNumber:   "1",
								Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
								Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		// Instantiate directly to access the private methods
		p := domainPresenter{}
		err := p.fillRhelIdmCerts(testCase.Given.To, testCase.Given.From)
		if testCase.Expected.Err != nil {
			assert.EqualError(t, err, testCase.Expected.Err.Error())
		} else {
			assert.NoError(t, err)
			require.NotNil(t, testCase.Expected.To)
			require.NotNil(t, testCase.Expected.To.RhelIdm)
			require.NotNil(t, testCase.Expected.To.RhelIdm.CaCerts)
			require.NotNil(t, testCase.Given.To)
			require.NotNil(t, testCase.Given.To.RhelIdm)
			require.NotNil(t, testCase.Given.To.RhelIdm.CaCerts)
			assert.Equal(t, testCase.Expected.To.RhelIdm.CaCerts, testCase.Given.To.RhelIdm.CaCerts)
		}
	}
}

func TestFillRhelIdmServers(t *testing.T) {
	type TestCaseGiven struct {
		To   *public.Domain
		From *model.Domain
	}
	type TestCaseExpected struct {
		Err error
		To  *public.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "Full success copy",
			Given: TestCaseGiven{
				To: &public.Domain{
					RhelIdm: &public.DomainIpa{},
				},
				From: &model.Domain{
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						Servers: []model.IpaServer{
							{
								FQDN:                "server1.mydomain.example",
								RHSMId:              "547ce70c-9eb5-4783-a619-086aa26f88e5",
								CaServer:            true,
								HCCEnrollmentServer: true,
								HCCUpdateServer:     true,
								PKInitServer:        true,
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				To: &public.Domain{
					Type: public.DomainTypeRhelIdm,
					RhelIdm: &public.DomainIpa{
						Servers: []public.DomainIpaServer{
							{
								Fqdn:                  "server1.mydomain.example",
								SubscriptionManagerId: "547ce70c-9eb5-4783-a619-086aa26f88e5",
								CaServer:              true,
								HccEnrollmentServer:   true,
								HccUpdateServer:       true,
								PkinitServer:          true,
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		// I instantiate directly because the public methods
		// are not part of the interface
		p := domainPresenter{}
		err := p.fillRhelIdmServers(testCase.Given.To, testCase.Given.From)
		if testCase.Expected.Err != nil {
			assert.EqualError(t, err, testCase.Expected.Err.Error())
		} else {
			assert.NoError(t, err)
			require.NotNil(t, testCase.Expected.To)
			require.NotNil(t, testCase.Expected.To.RhelIdm)
			require.NotNil(t, testCase.Expected.To.RhelIdm.Servers)
			require.NotNil(t, testCase.Given.To)
			require.NotNil(t, testCase.Given.To.RhelIdm)
			require.NotNil(t, testCase.Given.To.RhelIdm.Servers)
			assert.Equal(t, testCase.Expected.To.RhelIdm.Servers, testCase.Given.To.RhelIdm.Servers)
		}
	}
}

func TestRegisterAndUpdate(t *testing.T) {
	var (
		err    error
		output *public.RegisterDomainResponse
	)
	p := domainPresenter{}

	output, err = p.registerAndUpdate(nil)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'domain' is nil")

	domain := &model.Domain{}
	output, err = p.registerAndUpdate(domain)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'domain.Type' is invalid")

	domain.Type = pointy.Uint(model.DomainTypeIpa)
	output, err = p.registerAndUpdate(domain)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'domain.IpaDomain' is nil")

	domain.IpaDomain = &model.Ipa{}
	output, err = p.registerAndUpdate(domain)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'ipa.CaCerts' is nil")

	domain.IpaDomain.CaCerts = []model.IpaCert{}
	output, err = p.registerAndUpdate(domain)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'ipa.Servers' is nil")

	// Success case
	domain.IpaDomain.Servers = []model.IpaServer{}
	output, err = p.registerAndUpdate(domain)
	require.NotNil(t, output)
	assert.NoError(t, err)
	assert.Equal(t, false, output.AutoEnrollmentEnabled)
	assert.Equal(t, "", output.DomainName)
	assert.Equal(t, "", output.Title)
	assert.Equal(t, "", output.Description)
	assert.Equal(t, model.DomainTypeIpaString, string(output.Type))
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", output.DomainUuid)
	require.NotNil(t, output.RhelIdm)
	require.NotNil(t, output.RhelIdm.CaCerts)
	assert.Equal(t, 0, len(output.RhelIdm.CaCerts))
	require.NotNil(t, output.RhelIdm.Servers)
	assert.Equal(t, 0, len(output.RhelIdm.Servers))
	assert.Equal(t, 0, len(output.RhelIdm.RealmDomains))
	assert.Equal(t, "", output.RhelIdm.RealmName)
}

func TestRegister(t *testing.T) {
	testUUID := "ebac2444-e51b-11ed-a7f5-482ae3863d30"
	testDomainName := "mydomain.example"
	testModel := model.Domain{
		DomainUuid:            uuid.MustParse(testUUID),
		DomainName:            pointy.String(testDomainName),
		Title:                 pointy.String("My Example Domain Title"),
		Description:           pointy.String("My Example Domain Description"),
		OrgId:                 "12345",
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			RealmName:    pointy.String(testDomainName),
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []model.IpaCert{},
			Servers:      []model.IpaServer{},
		},
	}
	testExpected := public.Domain{
		DomainUuid:            testUUID,
		DomainName:            testDomainName,
		Title:                 "My Example Domain Title",
		Description:           "My Example Domain Description",
		AutoEnrollmentEnabled: true,
		Type:                  public.DomainTypeRhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    testDomainName,
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []public.DomainIpaCert{},
			Servers:      []public.DomainIpaServer{},
		},
	}

	p := domainPresenter{}

	domain, err := p.Register(nil)
	assert.EqualError(t, err, "'domain' is nil")
	assert.Nil(t, domain)

	domain, err = p.Register(&testModel)
	assert.NoError(t, err)
	assert.Equal(t, testExpected, *domain)
}

func TestUpdate(t *testing.T) {
	testUUID := "ebac2444-e51b-11ed-a7f5-482ae3863d30"
	testDomainName := "mydomain.example"
	testModel := model.Domain{
		DomainUuid:            uuid.MustParse(testUUID),
		DomainName:            pointy.String(testDomainName),
		Title:                 pointy.String("My Example Domain Title"),
		Description:           pointy.String("My Example Domain Description"),
		OrgId:                 "12345",
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
		IpaDomain: &model.Ipa{
			RealmName:    pointy.String(testDomainName),
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []model.IpaCert{},
			Servers:      []model.IpaServer{},
		},
	}
	testExpected := public.Domain{
		DomainUuid:            testUUID,
		DomainName:            testDomainName,
		Title:                 "My Example Domain Title",
		Description:           "My Example Domain Description",
		AutoEnrollmentEnabled: true,
		Type:                  public.DomainTypeRhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    testDomainName,
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []public.DomainIpaCert{},
			Servers:      []public.DomainIpaServer{},
		},
	}

	p := domainPresenter{}

	domain, err := p.Update(nil)
	assert.EqualError(t, err, "'domain' is nil")
	assert.Nil(t, domain)

	domain, err = p.Update(&testModel)
	assert.NoError(t, err, "")
	assert.Equal(t, testExpected, *domain)
}
