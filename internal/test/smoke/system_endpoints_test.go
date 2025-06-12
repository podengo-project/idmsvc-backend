package smoke

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	mock_pendo "github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client/pendo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.openly.dev/pointy"
)

const (
	pendoHostConfSuccess = "idmsvc-host-conf-success"
	pendoHostConfFailure = "idmsvc-host-conf-failure"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteSystemEndpoints struct {
	SuiteBase
	token  *public.DomainRegTokenResponse
	domain *public.Domain
}

func (s *SuiteSystemEndpoints) SetupTest() {
	s.PendoClient = mock_pendo.NewPendo(s.T())
	s.SuiteBase.SetupTest()
}

func (s *SuiteSystemEndpoints) prepareDomainIpaCreate(t *testing.T) {
	s.As(XRHIDUser)
	token, err := s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotEqual(t, "", token.DomainToken)
	s.token = token
}

func (s *SuiteSystemEndpoints) prepareDomainIpa(t *testing.T) {
	var err error

	// Add key to the database
	t.Log("Adding key")
	hcdb := datastore.NewHostconfJwkDb(s.Config, slog.Default())
	hcdb.ListKeys()
	hcdb.Purge()
	hcdb.Refresh()
	hcdb.ListKeys()

	// Create a token to register a domain
	t.Log("Creating token")
	s.As(XRHIDUser)
	s.token, err = s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, s.token)
	require.NotEqual(t, "", s.token.DomainToken)

	// This operation set AutoEnrollmentEnabled = False whatever
	// is the value we indicate here; we have to PATCH in a second
	// operation
	t.Log("Registering a domain")
	domain := "test.example"
	s.As(XRHIDSystem)
	s.domain, err = s.RegisterIpaDomain(s.token.DomainToken,
		builder_api.NewDomain(domain).
			WithDomainID(&s.token.DomainId).
			WithRhelIdm(builder_api.NewRhelIdmDomain(domain).
				WithServers([]public.DomainIpaServer{}).
				AddServer(builder_api.NewDomainIpaServer("1."+domain).
					WithHccUpdateServer(true).
					WithHccEnrollmentServer(true).
					WithSubscriptionManagerId(s.systemXRHID.Identity.System.CommonName).
					Build(),
				).Build(),
			).Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, s.domain)

	// Enable auto-join for the domain
	t.Log("Enabling auto-enrollment")
	s.As(XRHIDUser)
	s.domain, err = s.PatchDomain(
		s.domain.DomainId.String(),
		builder_api.NewUpdateDomainUserRequest().
			WithTitle(pointy.String(s.domain.DomainName)).
			WithDescription(nil).
			WithAutoEnrollmentEnabled(pointy.Bool(true)).
			Build())
	require.NoError(t, err)
	require.NotNil(t, s.domain)
}

func (s *SuiteSystemEndpoints) TestHostConfExecuteSuccess() {
	// Given
	t := s.T()
	s.As(RBACSuperAdmin)
	s.prepareDomainIpa(t)
	domainType := public.RhelIdm
	s.As(XRHIDSystem, RBACNoPermis)

	mockPendo, ok := s.PendoClient.(*mock_pendo.Pendo)
	require.True(t, ok)
	mockPendo.On("SendTrackEvent", mock.Anything, mock.MatchedBy(func(r *pendo.TrackRequest) bool {
		return (r.Event == pendoHostConfSuccess &&
			r.AccountID == s.systemXRHID.Identity.OrgID &&
			r.VisitorID == s.systemXRHID.Identity.System.CommonName)
	})).Return(nil)

	// When
	res, err := s.HostConfWithResponse(
		s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String(),
		"client."+s.domain.DomainName,
		builder_api.NewHostConf().
			WithDomainName(pointy.String(s.domain.DomainName)).
			WithDomainType(&domainType).
			Build())

	// Then
	require.NoError(t, err)
	require.NotNil(t, res)
	err = res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	mockPendo.AssertExpectations(t)
}

func (s *SuiteSystemEndpoints) TestHostConfExecuteFailure() {
	// Given
	t := s.T()
	s.As(RBACSuperAdmin)
	s.prepareDomainIpa(t)
	domainType := public.RhelIdm
	s.As(XRHIDSystem, RBACNoPermis)

	mockPendo, ok := s.PendoClient.(*mock_pendo.Pendo)
	require.True(t, ok)
	mockPendo.On("SendTrackEvent", mock.Anything, mock.MatchedBy(func(r *pendo.TrackRequest) bool {
		return (r.Event == pendoHostConfFailure &&
			r.AccountID == s.systemXRHID.Identity.OrgID &&
			r.VisitorID == s.systemXRHID.Identity.System.CommonName)
	})).Return(nil)

	// When
	res, err := s.HostConfWithResponse(
		s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String(),
		"client."+s.domain.DomainName,
		builder_api.NewHostConf().
			WithDomainName(pointy.String("invaliddomain.test")).
			WithDomainType(&domainType).
			Build())

	// Then
	require.NoError(t, err)
	require.NotNil(t, res)
	err = res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	mockPendo.AssertExpectations(t)
}

func (s *SuiteSystemEndpoints) TestInvalidRouteCauses404() {
	// Given
	t := s.T()
	s.As(RBACSuperAdmin)
	s.prepareDomainIpa(t)
	s.As(XRHIDSystem, RBACNoPermis)

	// When
	inventoryID := s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String()
	hdr := http.Header{}
	url := s.DefaultPublicBaseURL() + "/host-conf/" + inventoryID // MISSING HOSTNAME
	method := http.MethodPost
	s.addRequestID(&hdr, "test_system_host_conf")
	body := ""
	res, err := s.DoRequest(
		method,
		url,
		hdr,
		body,
	)

	// Then
	require.NoError(t, err)
	require.NotNil(t, res)
	err = res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func (s *SuiteSystemEndpoints) TestReadSigningKeys() {
	t := s.T()
	s.As(RBACSuperAdmin)
	s.prepareDomainIpa(t)
	s.As(RBACNoPermis, XRHIDSystem)
	res, err := s.ReadSigningKeysWithResponse()
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func (s *SuiteSystemEndpoints) TestSystemReadDomain() {
	t := s.T()
	s.As(RBACSuperAdmin)
	s.prepareDomainIpa(t)
	s.As(RBACNoPermis, XRHIDSystem)
	res, err := s.ReadDomainWithResponse(*s.domain.DomainId)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func (s *SuiteSystemEndpoints) TestSystemUpdateDomain() {
	t := s.T()
	s.As(RBACSuperAdmin)
	s.prepareDomainIpa(t)
	subscriptionManagerID := s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String()
	domainID := s.domain.DomainId.String()

	s.As(RBACNoPermis, XRHIDSystem)
	res, err := s.UpdateDomainWithResponse(
		domainID,
		builder_api.NewUpdateDomainAgent("test.example").
			WithHCCUpdate(true).
			WithDomainRhelIdm(*builder_api.NewRhelIdmDomain("test.example").
				WithServers([]public.DomainIpaServer{}).
				AddServer(
					builder_api.NewDomainIpaServer("1.test.example").
						WithHccUpdateServer(true).
						WithSubscriptionManagerId(subscriptionManagerID).
						Build(),
				).Build(),
			).WithSubscriptionManagerID(subscriptionManagerID).
			Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func (s *SuiteSystemEndpoints) TestSystemCreateDomain() {
	t := s.T()
	var err error

	// Create a token to register a domain
	s.As(RBACAdmin, XRHIDUser)
	s.token, err = s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, s.token)
	require.NotEqual(t, "", s.token.DomainToken)

	// Create the domains entry
	s.As(RBACNoPermis, XRHIDSystem)
	s.domain, err = s.RegisterIpaDomain(s.token.DomainToken,
		builder_api.NewDomain("test.example").
			WithDomainID(&s.token.DomainId).
			WithRhelIdm(builder_api.NewRhelIdmDomain("test.example").
				WithServers([]public.DomainIpaServer{}).
				AddServer(builder_api.NewDomainIpaServer("1.test.example").
					WithHccUpdateServer(true).
					WithSubscriptionManagerId(s.systemXRHID.Identity.System.CommonName).
					Build(),
				).Build(),
			).Build(),
	)
	require.NoError(t, err)
	require.NotNil(t, s.domain)
}
