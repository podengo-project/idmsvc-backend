package presenter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterIpaCaCerts(t *testing.T) {
	var (
		err error
	)
	p := &domainPresenter{}
	notValidBefore := time.Now()
	notValidAfter := notValidBefore.Add(time.Hour * 24)
	ipa := model.Domain{
		IpaDomain: &model.Ipa{
			CaCerts: []model.IpaCert{
				{
					Nickname:       "MYDOMAIN.EXAMPLE.IPA CA",
					Issuer:         "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
					NotValidBefore: notValidBefore.UTC(),
					NotValidAfter:  notValidAfter.UTC(),
					SerialNumber:   "1",
					Subject:        "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE.COM",
					Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				},
			},
		},
	}
	ipaNilCerts := model.Domain{
		IpaDomain: &model.Ipa{},
	}
	ipaNil := model.Domain{
		IpaDomain: nil,
	}
	output := public.DomainResponse{}

	assert.Panics(t, func() {
		err = p.registerIpaCaCerts(nil, nil)
	})

	assert.Panics(t, func() {
		err = p.registerIpaCaCerts(&ipaNil, nil)
	})

	err = p.registerIpaCaCerts(&ipaNilCerts, nil)
	assert.EqualError(t, err, "'ipa.CaCerts' is nil")

	assert.Panics(t, func() {
		err = p.registerIpaCaCerts(&ipa, nil)
	})

	err = p.registerIpaCaCerts(&ipa, &output)
	assert.NoError(t, err)
	assert.Equal(t, ipa.IpaDomain.CaCerts[0].Nickname, output.Ipa.CaCerts[0].Nickname)
	assert.Equal(t, ipa.IpaDomain.CaCerts[0].Issuer, output.Ipa.CaCerts[0].Issuer)
	assert.Equal(t, ipa.IpaDomain.CaCerts[0].NotValidAfter, output.Ipa.CaCerts[0].NotValidAfter)
	assert.Equal(t, ipa.IpaDomain.CaCerts[0].NotValidBefore, output.Ipa.CaCerts[0].NotValidBefore)
	assert.Equal(t, ipa.IpaDomain.CaCerts[0].SerialNumber, output.Ipa.CaCerts[0].SerialNumber)
	assert.Equal(t, ipa.IpaDomain.CaCerts[0].Subject, output.Ipa.CaCerts[0].Subject)
	assert.Equal(t, ipa.IpaDomain.CaCerts[0].Pem, output.Ipa.CaCerts[0].Pem)
}

func TestRegisterIpaServers(t *testing.T) {
	var (
		err error
	)
	p := domainPresenter{}
	assert.Panics(t, func() {
		err = p.registerIpaServers(nil, nil)
	})

	domain := &model.Domain{}
	err = p.registerIpaServers(domain, nil)
	assert.EqualError(t, err, "'IpaDomain' is nil")

	domain.IpaDomain = &model.Ipa{}
	err = p.registerIpaServers(domain, nil)
	assert.EqualError(t, err, "'output' is nil")

	output := public.RegisterDomainResponse{}
	err = p.registerIpaServers(domain, &output)
	assert.EqualError(t, err, "'ipa.Servers' is nil")

	domain.IpaDomain.Servers = append(domain.IpaDomain.Servers, model.IpaServer{})
	err = p.registerIpaServers(domain, &output)
	assert.NoError(t, err)
}

func TestGuardsRegisterIpa(t *testing.T) {
	var (
		err error
	)
	p := domainPresenter{}
	assert.Panics(t, func() {
		err = p.registerIpa(nil, nil)
	})

	domain := &model.Domain{}
	domain.Type = pointy.Uint(model.DomainTypeIpa)
	err = p.registerIpa(domain, nil)
	assert.EqualError(t, err, "'domain.IpaDomain' is nil")

	domain.IpaDomain = &model.Ipa{}
	assert.Panics(t, func() {
		err = p.registerIpa(domain, nil)
	})

	output := &public.RegisterDomainResponse{}
	err = p.registerIpa(domain, output)
	assert.EqualError(t, err, "'ipa.CaCerts' is nil")

	domain.IpaDomain.CaCerts = []model.IpaCert{}
	err = p.registerIpa(domain, output)
	assert.EqualError(t, err, "'ipa.Servers' is nil")
}

func TestRegisterIpa(t *testing.T) {
	tokenExpiration := &time.Time{}
	*tokenExpiration = time.Now()
	type TestCaseExpected struct {
		Domain *public.DomainResponse
		Err    error
	}
	type TestCase struct {
		Name     string
		Given    *model.Domain
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "Success minimal rhel-idm content",
			Given: &model.Domain{
				Type: pointy.Uint(model.DomainTypeIpa),
				IpaDomain: &model.Ipa{
					RealmName:       pointy.String(""),
					CaCerts:         []model.IpaCert{},
					Servers:         []model.IpaServer{},
					RealmDomains:    pq.StringArray{},
					Token:           pointy.String("71ad4978-c768-11ed-ad69-482ae3863d30"),
					TokenExpiration: tokenExpiration,
				},
			},
			Expected: TestCaseExpected{
				Domain: &public.DomainResponse{
					Type: model.DomainTypeIpaString,
					Ipa: public.DomainResponseIpa{
						RealmName:    "",
						CaCerts:      []public.DomainIpaCert{},
						Servers:      []public.DomainIpaServer{},
						RealmDomains: []string{},
					},
				},
				Err: nil,
			},
		},
		{
			Name: "Success full rhel-idm content",
			Given: &model.Domain{
				Type: pointy.Uint(model.DomainTypeIpa),
				IpaDomain: &model.Ipa{
					RealmName:       pointy.String("MYDOMAIN.EXAMPLE"),
					Token:           pointy.String("71ad4978-c768-11ed-ad69-482ae3863d30"),
					TokenExpiration: tokenExpiration,
					RealmDomains:    pq.StringArray{"mydomain.example"},
					CaCerts: []model.IpaCert{
						{
							Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
							Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
							Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
							SerialNumber: "1",
							Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
						},
					},
					Servers: []model.IpaServer{
						{
							FQDN:                "server1.mydomain.example",
							RHSMId:              "c4a80438-c768-11ed-a60e-482ae3863d30",
							PKInitServer:        true,
							CaServer:            true,
							HCCEnrollmentServer: true,
							HCCUpdateServer:     true,
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Domain: &public.DomainResponse{
					Ipa: public.DomainResponseIpa{
						RealmName:       "MYDOMAIN.EXAMPLE",
						Token:           nil,
						TokenExpiration: nil,
						RealmDomains:    pq.StringArray{"mydomain.example"},
						CaCerts: []public.DomainIpaCert{
							{
								Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
								Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								SerialNumber: "1",
								Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
						Servers: []public.DomainIpaServer{
							{
								Fqdn:                  "server1.mydomain.example",
								SubscriptionManagerId: "c4a80438-c768-11ed-a60e-482ae3863d30",
								PkinitServer:          true,
								CaServer:              true,
								HccEnrollmentServer:   true,
								HccUpdateServer:       true,
							},
						},
					},
				},
				Err: nil,
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		p := NewDomainPresenter()
		ipa, err := p.Register(testCase.Given)
		if testCase.Expected.Err != nil {
			assert.EqualError(t, err, testCase.Expected.Err.Error())
			assert.Nil(t, ipa)
		} else {
			assert.NoError(t, err)
			require.NotNil(t, ipa)
			assert.Equal(t, testCase.Expected.Domain.Ipa.RealmName, ipa.Ipa.RealmName)
			require.Equal(t, len(testCase.Expected.Domain.Ipa.RealmDomains), len(ipa.Ipa.RealmDomains))
			for i := range ipa.Ipa.RealmDomains {
				assert.Equal(t, testCase.Expected.Domain.Ipa.RealmDomains[i], ipa.Ipa.RealmDomains[i])
			}
			require.Equal(t, len(testCase.Expected.Domain.Ipa.CaCerts), len(ipa.Ipa.CaCerts))
			for i := range ipa.Ipa.CaCerts {
				assert.Equal(t, testCase.Expected.Domain.Ipa.CaCerts[i], ipa.Ipa.CaCerts[i])
			}
			require.Equal(t, len(testCase.Expected.Domain.Ipa.Servers), len(ipa.Ipa.Servers))
			for i := range ipa.Ipa.Servers {
				assert.Equal(t, testCase.Expected.Domain.Ipa.Servers[i], ipa.Ipa.Servers[i])
			}
		}
	}
}

func TestRegisterFillDomainData(t *testing.T) {
	p := domainPresenter{}
	assert.Panics(t, func() {
		p.registerFillDomainData(nil, nil)
	})

	domain := &model.Domain{}
	assert.Panics(t, func() {
		p.registerFillDomainData(domain, nil)
	})

	output := public.RegisterDomainResponse{}
	domainUUID := "6d9575f2-de94-11ed-af6e-482ae3863d30"
	domain.DomainUuid = uuid.MustParse(domainUUID)
	domain.AutoEnrollmentEnabled = pointy.Bool(true)
	domain.DomainName = pointy.String("mydomain.example")
	domain.Title = pointy.String("My Domain Example")
	domain.Description = pointy.String("My Domain Example Description")
	p.registerFillDomainData(domain, &output)
	assert.Equal(t, domainUUID, output.DomainUuid)
	assert.Equal(t, true, output.AutoEnrollmentEnabled)
	assert.Equal(t, "mydomain.example", output.DomainName)
	assert.Equal(t, "My Domain Example", output.Title)
	assert.Equal(t, "My Domain Example Description", output.Description)
}
