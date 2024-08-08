package pendo

import (
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	pendo_api "github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

type SetMetadataRequest interface {
	Add(item pendo.SetMetadataDetailsRequest) SetMetadataRequest
	Build() *pendo.SetMetadataRequest
}

type setMetadataRequest pendo_api.SetMetadataRequest

func NewSetMetadataRequest() SetMetadataRequest {
	return (*setMetadataRequest)(&pendo_api.SetMetadataRequest{})
}

func (b *setMetadataRequest) Build() *pendo.SetMetadataRequest {
	return (*pendo.SetMetadataRequest)(b)
}

func (b *setMetadataRequest) Add(item pendo.SetMetadataDetailsRequest) SetMetadataRequest {
	*b = append(*b, item)
	return b
}
