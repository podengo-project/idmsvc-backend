package handler

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/api/private"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
)

type Application interface {
	public.ServerInterface
	private.ServerInterface
	metrics.ServerInterface
}
