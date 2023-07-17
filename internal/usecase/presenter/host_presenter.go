package presenter

import (
	"fmt"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/presenter"
)

type hostPresenter struct {
	cfg *config.Config
}

func NewHostPresenter(cfg *config.Config) presenter.HostPresenter {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	return &hostPresenter{cfg}
}

func (p *hostPresenter) HostConf(domain *model.Domain) (*public.HostConfResponse, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}
