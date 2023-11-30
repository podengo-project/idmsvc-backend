package api

import (
	"github.com/openlyinc/pointy"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
)

type UpdateDomainUserJSONRequestBody interface {
	Build() *public.UpdateDomainUserJSONRequestBody
	WithTitle(value *string) UpdateDomainUserJSONRequestBody
	WithDescription(value *string) UpdateDomainUserJSONRequestBody
	WithAutoEnrollmentEnabled(value *bool) UpdateDomainUserJSONRequestBody
}

type updateDomainUserJSONRequestBody public.UpdateDomainUserJSONRequestBody

// NewUpdateDomainUserJSONRequestBody instantiate a new generator for
// the user domain patch operation.
func NewUpdateDomainUserJSONRequestBody() UpdateDomainUserJSONRequestBody {
	letters := []rune("abcdefghijklmnopqrstuvwxyz 0123456789")
	return &updateDomainUserJSONRequestBody{
		// TODO Enhance generator
		Title:                 pointy.String(helper.GenRandString(letters, 30)),
		AutoEnrollmentEnabled: pointy.Bool(helper.GenRandBool()),
		Description:           pointy.String(helper.GenRandParagraph(0)),
	}
}

func (b *updateDomainUserJSONRequestBody) Build() *public.UpdateDomainUserJSONRequestBody {
	return (*public.UpdateDomainUserJSONRequestBody)(b)
}

func (b *updateDomainUserJSONRequestBody) WithTitle(value *string) UpdateDomainUserJSONRequestBody {
	b.Title = value
	return b
}

func (b *updateDomainUserJSONRequestBody) WithDescription(value *string) UpdateDomainUserJSONRequestBody {
	b.Description = value
	return b
}

func (b *updateDomainUserJSONRequestBody) WithAutoEnrollmentEnabled(value *bool) UpdateDomainUserJSONRequestBody {
	b.AutoEnrollmentEnabled = value
	return b
}
