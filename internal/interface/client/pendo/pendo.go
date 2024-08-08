package pendo

import "context"

type Kind string
type Group string
type ID string
type FieldName string

const (
	KindVisitor Kind = "visitor"
	KindAccount Kind = "account"
)

type Pendo interface {
	SetMetadata(ctx context.Context, kind Kind, group Group, metrics SetMetadataRequest) (*SetMetadataResponse, error)
	SendTrackEvent(ctx context.Context, track *TrackRequest) error
}
