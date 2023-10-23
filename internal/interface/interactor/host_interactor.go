package interactor

import (
	"github.com/google/uuid"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	api_public "github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type HostConfOptions struct {
	OrgId       string
	CommonName  public.SubscriptionManagerId
	InventoryId public.HostId
	Fqdn        public.Fqdn
	DomainId    *uuid.UUID
	DomainName  *string
	DomainType  *api_public.DomainType
}

type HostInteractor interface {
	HostConf(xrhid *identity.XRHID, inventoryId api_public.HostId, fqdn string, params *api_public.HostConfParams, body *api_public.HostConf) (*HostConfOptions, error)
}
