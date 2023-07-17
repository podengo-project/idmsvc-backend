package presenter

import (
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
)

type HostPresenter interface {
	HostConf(domain *model.Domain) (*public.HostConfResponse, error)
}
