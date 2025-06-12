package api

import (
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"go.openly.dev/pointy"
)

type Location interface {
	Build() public.Location
	WithName(value string) Location
	WithDescription(value string) Location
}

type location public.Location

func NewLocation() Location {
	return &location{}
}

func (b *location) Build() public.Location {
	return public.Location(*b)
}

func (b *location) WithName(value string) Location {
	b.Name = value
	return b
}

func (b *location) WithDescription(value string) Location {
	b.Description = pointy.String(value)
	return b
}
