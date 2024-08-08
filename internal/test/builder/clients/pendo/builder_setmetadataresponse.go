package pendo

import (
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	pendo_api "github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

type SetMetadataResponse interface {
	Build() *pendo.SetMetadataResponse
	WithTotal(value int64) SetMetadataResponse
	WithUpdated(value int64) SetMetadataResponse
	IncUpdated() SetMetadataResponse
	WithFailed(value int64) SetMetadataResponse
	IncFailed() SetMetadataResponse
	AddMissing(value string) SetMetadataResponse
	WithKind(value pendo.Kind) SetMetadataResponse
}

type setMetadataResponse pendo_api.SetMetadataResponse

func NewSetMetadataResponse() SetMetadataResponse {
	return (*setMetadataResponse)(&pendo_api.SetMetadataResponse{})
}

func (b *setMetadataResponse) Build() *pendo.SetMetadataResponse {
	return (*pendo.SetMetadataResponse)(b)
}

func (b *setMetadataResponse) WithTotal(value int64) SetMetadataResponse {
	b.Total = value
	return b
}

func (b *setMetadataResponse) WithUpdated(value int64) SetMetadataResponse {
	b.Updated = value
	return b
}

func (b *setMetadataResponse) IncUpdated() SetMetadataResponse {
	b.Updated++
	return b
}

func (b *setMetadataResponse) WithFailed(value int64) SetMetadataResponse {
	b.Failed = value
	return b
}

func (b *setMetadataResponse) IncFailed() SetMetadataResponse {
	b.Failed++
	return b
}

func (b *setMetadataResponse) AddMissing(value string) SetMetadataResponse {
	b.Missing = append(b.Missing, value)
	return b
}

func (b *setMetadataResponse) WithKind(value pendo.Kind) SetMetadataResponse {
	b.Kind = value
	return b
}
