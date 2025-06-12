package model

import (
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"go.openly.dev/pointy"
	"gorm.io/gorm"
)

type IpaLocation interface {
	Build() model.IpaLocation
	WithIpaID(value uint) IpaLocation
	WithModel(value gorm.Model) IpaLocation
	WithName(value string) IpaLocation
	WithDescription(value string) IpaLocation
}

type ipaLocation struct {
	IpaLocation model.IpaLocation
}

func NewIpaLocation(gormModel gorm.Model) IpaLocation {
	return &ipaLocation{
		IpaLocation: model.IpaLocation{
			Model:       gormModel,
			Name:        helper.GenRandLocationLabel(),
			Description: pointy.String(helper.GenRandLocationDescription()),
		},
	}
}

func (b *ipaLocation) Build() model.IpaLocation {
	return b.IpaLocation
}

func (b *ipaLocation) WithIpaID(value uint) IpaLocation {
	b.IpaLocation.IpaID = value
	return b
}

func (b *ipaLocation) WithModel(value gorm.Model) IpaLocation {
	b.IpaLocation.Model = value
	return b
}

func (b *ipaLocation) WithName(value string) IpaLocation {
	b.IpaLocation.Name = value
	return b
}

func (b *ipaLocation) WithDescription(value string) IpaLocation {
	b.IpaLocation.Description = pointy.String(value)
	return b
}
