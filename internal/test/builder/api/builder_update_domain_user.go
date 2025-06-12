package api

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"go.openly.dev/pointy"
)

type UpdateDomainUserRequest interface {
	Build() *public.UpdateDomainUserRequest
	WithTitle(value *string) UpdateDomainUserRequest
	WithDescription(value *string) UpdateDomainUserRequest
	WithAutoEnrollmentEnabled(value *bool) UpdateDomainUserRequest
}

type updateDomainUserRequest public.UpdateDomainUserRequest

// NewUpdateDomainUserRequest instantiate a new generator for
// the user domain patch operation.
func NewUpdateDomainUserRequest() UpdateDomainUserRequest {
	letters := []rune("abcdefghijklmnopqrstuvwxyz 0123456789")
	return &updateDomainUserRequest{
		// TODO Enhance generator
		Title:                 pointy.String(helper.GenRandString(letters, 30)),
		AutoEnrollmentEnabled: pointy.Bool(helper.GenRandBool()),
		Description:           pointy.String(helper.GenRandParagraph(0)),
	}
}

func (b *updateDomainUserRequest) Build() *public.UpdateDomainUserJSONRequestBody {
	return (*public.UpdateDomainUserJSONRequestBody)(b)
}

func (b *updateDomainUserRequest) WithTitle(value *string) UpdateDomainUserRequest {
	b.Title = value
	return b
}

func (b *updateDomainUserRequest) WithDescription(value *string) UpdateDomainUserRequest {
	b.Description = value
	return b
}

func (b *updateDomainUserRequest) WithAutoEnrollmentEnabled(value *bool) UpdateDomainUserRequest {
	b.AutoEnrollmentEnabled = value
	return b
}
