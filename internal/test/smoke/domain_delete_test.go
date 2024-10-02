package smoke

import (
	"fmt"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_api "github.com/podengo-project/idmsvc-backend/internal/test/builder/api"
	"github.com/stretchr/testify/require"
)

// SuiteDeleteDomain is the suite to validate the smoke test when read domain endpoint at GET /api/idmsvc/v1/domains/:domain_id
type SuiteDeleteDomain struct {
	SuiteBase
}

func (s *SuiteDeleteDomain) TestDeleteDomain() {
	t := s.T()
	var (
		token  *public.DomainRegTokenResponse
		domain *public.Domain
		err    error
	)
	xrhids := []XRHIDProfile{XRHIDUser, XRHIDServiceAccount}

	// Execute the test cases
	for i, xrhid := range xrhids {
		// GIVEN
		s.As(RBACAdmin, XRHIDUser)
		token, err = s.CreateToken()
		require.NoError(t, err)
		require.NotNil(t, token)

		domainName := fmt.Sprintf("domain%d.test", i)
		domainRequest := builder_api.NewDomain(domainName).Build()
		setFirstAsUpdateServer(domainRequest)
		setFirstServerRHSMId(s.T(), domainRequest, s.systemXRHID)

		s.As(RBACAdmin, XRHIDSystem)
		domain, err = s.RegisterIpaDomain(token.DomainToken, domainRequest)
		require.NoError(t, err)
		require.NotNil(t, domain)
		require.NotNil(t, domain.DomainId)

		// WHEN
		s.As(RBACAdmin, xrhid)
		err = s.DeleteDomain(*domain.DomainId)

		// THEN
		require.NoError(t, err)
	}
}
