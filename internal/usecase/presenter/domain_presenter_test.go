package presenter

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/domain/model"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var cfg = config.Config{
	Application: config.Application{
		PaginationDefaultLimit: 10,
		PaginationMaxLimit:     100,
	},
}

func TestNewTodoPresenter(t *testing.T) {
	assert.Panics(t, func() {
		NewDomainPresenter(nil)
	})

	assert.NotPanics(t, func() {
		NewDomainPresenter(&cfg)
	})
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
					DomainType:            public.DomainDomainType(model.DomainTypeString(model.DomainTypeIpa)),
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
		obj := &domainPresenter{cfg: &cfg}
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
				testCase.Expected.Output.DomainType,
				output.DomainType)
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
			Name: "Type is nil",
			Given: &model.Domain{
				Type: nil,
			},
			Expected: TestCaseExpected{
				Response: nil,
				Err:      fmt.Errorf("'domain.Type' is nil"),
			},
		},
		{
			Name: "Domain Type is invalid",
			Given: &model.Domain{
				Type: pointy.Uint(model.DomainTypeUndefined),
			},
			Expected: TestCaseExpected{
				Response: nil,
				Err:      fmt.Errorf("'domain.Type' is invalid"),
			},
		},
		// {
		// 	Name: "AutoEnrollmentEnabled is nil",
		// 	Given: &model.Domain{
		// 		Type:                  pointy.Uint(model.DomainTypeIpa),
		// 		AutoEnrollmentEnabled: nil,
		// 	},
		// 	Expected: TestCaseExpected{
		// 		Response: nil,
		// 		Err:      fmt.Errorf("'AutoEnrollmentEnabled' is nil"),
		// 	},
		// },
		{
			Name: "IpaDomain is nil",
			Given: &model.Domain{
				AutoEnrollmentEnabled: pointy.Bool(true),
				DomainName:            pointy.String("domain.example"),
				Type:                  pointy.Uint(model.DomainTypeIpa),
				IpaDomain:             nil,
			},
			Expected: TestCaseExpected{
				Response: nil,
				Err:      fmt.Errorf("'domain.IpaDomain' is nil"),
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
				IpaDomain: &model.Ipa{
					RealmName: pointy.String("DOMAIN.EXAMPLE"),
					CaCerts:   []model.IpaCert{},
					Servers:   []model.IpaServer{},
				},
			},
			Expected: TestCaseExpected{
				Response: &public.CreateDomainResponse{
					DomainName:            "domain.example",
					AutoEnrollmentEnabled: true,
					DomainType:            model.DomainTypeIpaString,
					DomainUuid:            "00000000-0000-0000-0000-000000000000",
					RhelIdm: &public.DomainIpa{
						RealmName:    "DOMAIN.EXAMPLE",
						RealmDomains: []string{},
						CaCerts:      []public.DomainIpaCert{},
						Servers:      []public.DomainIpaServer{},
					},
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
		obj := &domainPresenter{cfg: &cfg}
		response, err := obj.Create(testCase.Given)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			assert.EqualError(t, err, testCase.Expected.Err.Error())
			assert.Nil(t, response)
		} else {
			assert.NoError(t, err)
			equalPresenterDomain(t, testCase.Expected.Response, response)
		}
	}

}

func TestFillRhelmIdmCertsPanics(t *testing.T) {
	var err error
	p := &domainPresenter{cfg: &cfg}

	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(nil, nil)
	})

	domain := &model.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(nil, domain)
	})

	domain.IpaDomain = &model.Ipa{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(nil, domain)
	})

	output := &public.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(output, domain)
	})

	domain.IpaDomain.CaCerts = []model.IpaCert{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(output, domain)
	})

	// Minimal success
	output.RhelIdm = &public.DomainIpa{}
	p.fillRhelIdmCerts(output, domain)
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
					DomainType: public.DomainDomainTypeRhelIdm,
					RhelIdm:    &public.DomainIpa{},
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
					DomainType: public.DomainDomainTypeRhelIdm,
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
		p := &domainPresenter{cfg: &cfg}
		if testCase.Expected.Err != nil {
			// assert.EqualError(t, err, testCase.Expected.Err.Error())
			assert.Panics(t, func() {
				p.fillRhelIdmCerts(testCase.Given.To, testCase.Given.From)
			})
		} else {
			// assert.NoError(t, err)
			assert.NotPanics(t, func() {
				p.fillRhelIdmCerts(testCase.Given.To, testCase.Given.From)
			})
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
		DomainType:            public.DomainDomainTypeRhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    testDomainName,
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []public.DomainIpaCert{},
			Servers:      []public.DomainIpaServer{},
		},
	}

	p := &domainPresenter{cfg: &cfg}

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
		DomainType:            public.DomainDomainTypeRhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    testDomainName,
			RealmDomains: pq.StringArray{testDomainName},
			CaCerts:      []public.DomainIpaCert{},
			Servers:      []public.DomainIpaServer{},
		},
	}

	p := &domainPresenter{cfg: &cfg}

	domain, err := p.Update(nil)
	assert.EqualError(t, err, "'domain' is nil")
	assert.Nil(t, domain)

	domain, err = p.Update(&testModel)
	assert.NoError(t, err, "")
	assert.Equal(t, testExpected, *domain)
}

func TestList(t *testing.T) {
	// https://consoledot.pages.redhat.com/docs/dev/developer-references/rest/pagination.html
	testOrgID := "12345"
	testUUID1 := "5427c3d6-eaa1-11ed-99da-482ae3863d30"
	testUUID2 := "5ae8e844-eaa1-11ed-8f71-482ae3863d30"
	prefix := "/api/hmsidm/v1"
	cfg := config.Config{
		Application: config.Application{
			PaginationDefaultLimit: 10,
			PaginationMaxLimit:     100,
		},
	}
	p := &domainPresenter{cfg: &cfg}

	// offset lower than 0
	count := int64(5)
	offset := -1
	limit := -1
	output, err := p.List(prefix, count, offset, limit, nil)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'offset' is lower than 0")

	// limit lower than 0
	offset = 5
	output, err = p.List(prefix, count, offset, limit, nil)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'limit' is lower than 0")

	// set default limit
	limit = 0
	output, err = p.List(prefix, count, offset, limit, nil)
	assert.NotNil(t, output)

	assert.Equal(t, count, output.Meta.Count)
	assert.Equal(t, p.cfg.Application.PaginationDefaultLimit, output.Meta.Limit)
	assert.Equal(t, offset, output.Meta.Offset)

	require.NotNil(t, output.Links.First)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationDefaultLimit), *output.Links.First)
	require.NotNil(t, output.Links.Previous)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationDefaultLimit), *output.Links.Previous)
	require.NotNil(t, output.Links.Next)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationDefaultLimit), *output.Links.Previous)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationDefaultLimit), *output.Links.Previous)

	// set max limit  paginationMaxLimit
	limit = p.cfg.Application.PaginationMaxLimit + 1
	output, err = p.List(prefix, count, offset, limit, nil)
	assert.NotNil(t, output)

	assert.Equal(t, count, output.Meta.Count)
	assert.Equal(t, p.cfg.Application.PaginationMaxLimit, output.Meta.Limit)
	assert.Equal(t, offset, output.Meta.Offset)

	require.NotNil(t, output.Links.First)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationMaxLimit), *output.Links.First)
	require.NotNil(t, output.Links.Previous)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationMaxLimit), *output.Links.Previous)
	require.NotNil(t, output.Links.Next)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationMaxLimit), *output.Links.Previous)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, p.buildPaginationLink(prefix, 0, p.cfg.Application.PaginationMaxLimit), *output.Links.Previous)

	// domain slice is nil return empty list
	count = int64(0)
	offset = 2
	limit = 2
	expected := public.ListDomainsResponseSchema{
		Meta: public.PaginationMeta{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		},
		Links: public.PaginationLinks{
			First:    pointy.String(p.buildPaginationLink(prefix, 0, limit)),
			Previous: pointy.String(p.buildPaginationLink(prefix, 0, limit)),
			Next:     pointy.String(p.buildPaginationLink(prefix, 0, limit)),
			Last:     pointy.String(p.buildPaginationLink(prefix, 0, limit)),
		},
		Data: []public.ListDomainsData{},
	}
	output, err = p.List(prefix, count, offset, limit, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, *output)

	// domain slice fill the list
	data := []model.Domain{
		{
			OrgId:                 testOrgID,
			DomainUuid:            uuid.MustParse(testUUID1),
			DomainName:            pointy.String("mydomain1.example"),
			AutoEnrollmentEnabled: pointy.Bool(true),
			Title:                 pointy.String("mydomain1 example title"),
			Description:           pointy.String("mydomain1.example located in Boston"),
			Type:                  pointy.Uint(model.DomainTypeIpa),
		},
		{
			OrgId:                 testOrgID,
			DomainUuid:            uuid.MustParse(testUUID2),
			DomainName:            pointy.String("mydomain2.example"),
			AutoEnrollmentEnabled: nil,
			Title:                 pointy.String("mydomain2 example title"),
			Description:           pointy.String("mydomain2.example located in Brno"),
			Type:                  pointy.Uint(model.DomainTypeIpa),
		},
	}
	count = int64(len(data))
	offset = 0
	limit = 2
	expected = public.ListDomainsResponseSchema{
		Meta: public.PaginationMeta{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		},
		Links: public.PaginationLinks{
			First:    pointy.String(p.buildPaginationLink(prefix, 0, limit)),
			Previous: pointy.String(p.buildPaginationLink(prefix, 0, limit)),
			Next:     pointy.String(p.buildPaginationLink(prefix, 0, limit)),
			Last:     pointy.String(p.buildPaginationLink(prefix, 0, limit)),
		},
		Data: []public.ListDomainsData{
			{
				AutoEnrollmentEnabled: true,
				DomainType:            public.ListDomainsDataDomainTypeRhelIdm,
				DomainName:            "mydomain1.example",
				DomainUuid:            testUUID1,
				Title:                 "mydomain1 example title",
				Description:           "mydomain1.example located in Boston",
			},
			{
				AutoEnrollmentEnabled: false,
				DomainType:            public.ListDomainsDataDomainTypeRhelIdm,
				DomainName:            "mydomain2.example",
				DomainUuid:            testUUID2,
				Title:                 "mydomain2 example title",
				Description:           "mydomain2.example located in Brno",
			},
		},
	}
	output, err = p.List(prefix, count, offset, limit, data)
	assert.NoError(t, err)
	assert.Equal(t, &expected, output)
}
