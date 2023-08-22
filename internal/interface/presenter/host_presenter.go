package presenter

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

type HostPresenter interface {
	HostConf(domain *model.Domain) (*public.HostConfResponse, error)
}
