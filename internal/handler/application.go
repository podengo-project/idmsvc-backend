package handler

import (
	"github.com/hmsidm/internal/api/metrics"
	"github.com/hmsidm/internal/api/private"
	"github.com/hmsidm/internal/api/public"
)

type Application interface {
	public.ServerInterface
	private.ServerInterface
	metrics.ServerInterface
}
