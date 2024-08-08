package pendo

import (
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	pendo_api "github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

type SetMetadataDetailsRequest interface {
	Build() *pendo.SetMetadataDetailsRequest
	SetVisitorID(value string) SetMetadataDetailsRequest
	AddValue(fieldName string, value any) SetMetadataDetailsRequest
}

type setMetadataDetailsRequest pendo_api.SetMetadataDetailsRequest

func NewSetMetadataDetailsRequest() SetMetadataDetailsRequest {
	return (*setMetadataDetailsRequest)(&pendo_api.SetMetadataDetailsRequest{})
}

func (b *setMetadataDetailsRequest) Build() *pendo.SetMetadataDetailsRequest {
	return (*pendo.SetMetadataDetailsRequest)(b)
}

func (b *setMetadataDetailsRequest) SetVisitorID(value string) SetMetadataDetailsRequest {
	b.VisitorID = value
	return b
}

func (b *setMetadataDetailsRequest) AddValue(fieldName string, value any) SetMetadataDetailsRequest {
	b.Values[fieldName] = value
	return b
}
