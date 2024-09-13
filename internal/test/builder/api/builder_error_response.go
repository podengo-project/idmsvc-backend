package api

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
)

type ErrorResponse interface {
	Build() *public.ErrorResponse
	Add(err public.ErrorInfo) ErrorResponse
}

type errorResponse public.ErrorResponse

func NewErrorResponse() ErrorResponse {
	return &errorResponse{}
}

func (b *errorResponse) Build() *public.ErrorResponse {
	return (*public.ErrorResponse)(b)
}

func (b *errorResponse) Add(err public.ErrorInfo) ErrorResponse {
	if b.Errors == nil {
		errors := make([]public.ErrorInfo, 0, 8)
		b.Errors = &errors
	}
	*b.Errors = append(*b.Errors, err)
	return b
}
