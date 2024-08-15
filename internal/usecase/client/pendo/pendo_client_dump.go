package pendo

import (
	"context"
	"fmt"
	"log/slog"

	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

type pendoClientDump struct{}

// NewClientDump
func NewClientDump() pendo.Pendo {
	return newClientDump()
}

// SetMetadata launch a request to pendo service to register some metric
func (c *pendoClientDump) SetMetadata(ctx context.Context, kind pendo.Kind, group pendo.Group, metrics pendo.SetMetadataRequest) (*pendo.SetMetadataResponse, error) {
	logger := app_context.LogFromCtx(ctx)
	logger.Debug("pendo SetMetadata called",
		slog.String("kind", string(kind)),
		slog.String("group", string(group)),
	)
	metaDataResponse := &pendo.SetMetadataResponse{
		Total:   int64(len(metrics)),
		Updated: int64(len(metrics)),
		Failed:  0,
		Kind:    kind,
	}
	logger.Debug("pendo SetMetadata returning",
		slog.Int64("total", metaDataResponse.Total),
		slog.Int64("failed", metaDataResponse.Failed),
		slog.Int64("updated", metaDataResponse.Updated),
		slog.Any("missing", metaDataResponse.Missing),
		slog.String("kind", string(metaDataResponse.Kind)),
	)
	return metaDataResponse, nil
}

// SendTrackEvent
func (c *pendoClientDump) SendTrackEvent(ctx context.Context, track *pendo.TrackRequest) error {
	logger := app_context.LogFromCtx(ctx)
	logger.Debug(fmt.Sprintf("pendo SendTrackEvent called: %s", track.Event))
	return nil
}

//
// ----- Private methods ------
//

func newClientDump() *pendoClientDump {
	return &pendoClientDump{}
}
