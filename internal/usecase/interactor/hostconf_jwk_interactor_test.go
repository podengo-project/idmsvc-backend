package interactor

import (
	"testing"

	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	internal_errors "github.com/podengo-project/idmsvc-backend/internal/errors"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"go.openly.dev/pointy"
)

func TestNewHostconfJwkInteractor(t *testing.T) {
	var component interactor.HostconfJwkInteractor
	assert.NotPanics(t, func() {
		component = NewHostconfJwkInteractor()
	})
	assert.NotNil(t, component)
}

func TestGetSigningKeys(t *testing.T) {
	var (
		orgID string
		err   error
	)
	xrhid := test.SystemXRHID
	params := &api_public.GetSigningKeysParams{
		XRhInsightsRequestId: pointy.String("requestid"),
	}

	i := NewHostconfJwkInteractor()

	orgID, err = i.GetSigningKeys(nil, nil)
	assert.Equal(t, orgID, "")
	assert.EqualError(t, err, internal_errors.NilArgError("xrhid").Error())

	orgID, err = i.GetSigningKeys(&xrhid, nil)
	assert.Equal(t, orgID, "")
	assert.EqualError(t, err, internal_errors.NilArgError("params").Error())

	orgID, err = i.GetSigningKeys(&xrhid, params)
	assert.Equal(t, orgID, xrhid.Identity.OrgID)
	assert.Nil(t, err)
}
