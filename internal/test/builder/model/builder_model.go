package model

// FIXME Change the option pattern by a Builder pattern
//       this is, instead of using 'With...' that could
//       evoke name colisions, add the custom methods that
//       override the values to the ModelBuilder

import (
	"time"

	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"gorm.io/gorm"
)

// GormModel is the builder for the
type GormModel interface {
	Build() gorm.Model
	WithID(value uint) GormModel
	WithCreatedAt(value time.Time) GormModel
	WithUpdatedAt(value time.Time) GormModel
	WithDeletedAt(value gorm.DeletedAt) GormModel
}

// gormModel is the specific builder implementation
type gormModel struct {
	gorm.Model
}

// NewModel generate a gorm.Model with random information
// overrided by the customized options.
func NewModel() GormModel {
	genCreatedAt := builder_helper.GenPastNearTime(time.Hour * 24 * 10)
	genUpdatedAt := builder_helper.GenBetweenTimeUTC(genCreatedAt, time.Now())
	modelID := uint(builder_helper.GenRandNum(0, 2^63))

	return &gormModel{
		Model: gorm.Model{
			ID:        modelID,
			CreatedAt: genCreatedAt,
			UpdatedAt: genUpdatedAt,
		},
	}
}

func (b *gormModel) Build() gorm.Model {
	return b.Model
}

func (b *gormModel) WithID(value uint) GormModel {
	b.Model.ID = value
	return b
}

func (b *gormModel) WithCreatedAt(value time.Time) GormModel {
	b.CreatedAt = value
	return b
}

func (b *gormModel) WithUpdatedAt(value time.Time) GormModel {
	b.Model.UpdatedAt = value
	return b
}

func (b *gormModel) WithDeletedAt(value gorm.DeletedAt) GormModel {
	b.DeletedAt = value
	return b
}
