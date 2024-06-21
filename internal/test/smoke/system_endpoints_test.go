package smoke

import (
	"net/http"
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	mock_rbac "github.com/podengo-project/idmsvc-backend/internal/infrastructure/service/impl/mock/rbac/impl"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteTokenCreate is the suite token for smoke tests at /api/idmsvc/v1/domains/token
type SuiteSystemEndpoints struct {
	SuiteBase
	token  *public.DomainRegTokenResponse
	domain *public.Domain
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
	hcdb := datastore.NewHostconfJwkDb(s.cfg)
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

func (s *SuiteSystemEndpoints) TestHostConfExecute() {
	t := s.T()
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileSuperAdmin])
	s.prepareDomainIpa(t)
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainNoPerms])

	t.Log("Calling SystemHostConfWithResponse")
	domainType := public.RhelIdm
	s.As(XRHIDSystem)
	res, err := s.HostConfWithResponse(
		s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String(),
		"client."+s.domain.DomainName,
		builder_api.NewHostConf().
			WithDomainName(pointy.String(s.domain.DomainName)).
			WithDomainType(&domainType).
			Build())
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func (s *SuiteSystemEndpoints) TestReadSigningKeys() {
	t := s.T()
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileSuperAdmin])
	s.prepareDomainIpa(t)
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainNoPerms])
	s.As(XRHIDSystem)
	res, err := s.ReadSigningKeysWithResponse()
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func (s *SuiteSystemEndpoints) TestSystemReadDomain() {
	t := s.T()
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileSuperAdmin])
	s.prepareDomainIpa(t)
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainNoPerms])
	s.As(XRHIDSystem)
	res, err := s.ReadDomainWithResponse(*s.domain.DomainId)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func (s *SuiteSystemEndpoints) TestSystemUpdateDomain() {
	t := s.T()
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileSuperAdmin])
	s.prepareDomainIpa(t)
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainNoPerms])
	subscriptionManagerID := s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String()
	domainID := s.domain.DomainId.String()

	s.As(XRHIDSystem)
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

	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainAdmin])

	// Create a token to register a domain
	s.As(XRHIDUser)
	s.token, err = s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, s.token)
	require.NotEqual(t, "", s.token.DomainToken)

	// Create the domains entry
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainNoPerms])
	s.As(XRHIDSystem)
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
