package smoke

import (
	"net/http"
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/datastore"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	mock_pendo "github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client/pendo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SuiteMiscelanea struct {
	SuiteBase
	token  *public.DomainRegTokenResponse
	domain *public.Domain
}

func (s *SuiteMiscelanea) SetupTest() {
	s.PendoClient = mock_pendo.NewPendo(s.T())
	s.SuiteBase.SetupTest()
}

func (s *SuiteMiscelanea) prepareDomainIpaCreate(t *testing.T) {
	s.As(XRHIDUser)
	token, err := s.CreateToken()
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotEqual(t, "", token.DomainToken)
	s.token = token
}

func (s *SuiteMiscelanea) prepareDomainIpa(t *testing.T) {
	var err error

	// Add key to the database
	t.Log("Adding key")
	hcdb := datastore.NewHostconfJwkDb(s.Config)
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

// TestBugHmsXXXXNotFound check a condition when a path
// that is not found for /api/idmsvc/domains/v1/host-conf/:inventory_id
// where the API should return with 404 status code.
func (s *SuiteMiscelanea) TestBugHmsXXXXNotFound() {
	// Given
	t := s.T()
	s.As(RBACSuperAdmin)
	s.prepareDomainIpa(t)
	domainType := public.RhelIdm
	s.As(XRHIDSystem, RBACNoPermis)

	mockPendo, ok := s.PendoClient.(*mock_pendo.Pendo)
	require.True(t, ok)

	// When
	hdr := http.Header{}
	inventoryID := s.domain.RhelIdm.Servers[0].SubscriptionManagerId.String()
	fqdn := ""
	hostconf := builder_api.NewHostConf().
		WithDomainName(pointy.String(s.domain.DomainName)).
		WithDomainType(&domainType).
		Build()
	url := s.DefaultPublicBaseURL() + "/host-conf/" + inventoryID + "/" + fqdn
	method := http.MethodPost
	s.addRequestID(&hdr, "test_miscelanea_bug_hms_xxxx_not_found")
	body := hostconf
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
	mockPendo.AssertExpectations(t)
}

func TestSuiteMiscelanea(t *testing.T) {
	suite.Run(t, new(SuiteMiscelanea))
}
