package pendo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

type pendoClient struct {
	Config *config.Config
	Client *http.Client
}

// NewClient
func NewClient(cfg *config.Config) (client pendo.Pendo) {
	defer func() {
		if recover() != nil {
			slog.Default().Warn("Using pendoClientDump as fallback")
			client = newClientDump()
		}
	}()
	client = newClient(cfg)
	return client
}

// SetMetadata launch a request to pendo service to register some metric
func (c *pendoClient) SetMetadata(ctx context.Context, kind pendo.Kind, group pendo.Group, metrics pendo.SetMetadataRequest) (*pendo.SetMetadataResponse, error) {
	logger := app_context.LogFromCtx(ctx)
	if err := c.guardSetMetadata(kind, group, metrics); err != nil {
		logger.Error("wrong arguments: " + err.Error())
		return nil, fmt.Errorf("bad arguments: %w", err)
	}

	// https://github.com/RedHatInsights/cloud-connector/blob/master/internal/pendo_transmitter/pendo_reporter.go#L51
	reqBody, err := json.Marshal(metrics)
	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("error unserializing SetMetaDataResponse: %w", err)
	}

	// Prepare the request
	url := c.Config.Clients.PendoBaseURL + "/api/v1/metadata/" + url.PathEscape(string(kind)) + "/" + url.PathEscape(string(group)) + "/value"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("error making SetMetaData request: %w", err)
	}

	// Add headers to the request
	req.Header.Set("content-type", "application/json")
	req.Header.Set(header.HeaderXPendoIntegrationKey, c.Config.Clients.PendoAPIKey)

	// Launch request
	resp, err := c.Client.Do(req)
	if err != nil {
		logger.Error("doing request to pendo service")
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		logger.Error("expected StatusCode=" + http.StatusText(http.StatusOK) + " but received StatusCode=" + http.StatusText(resp.StatusCode))
		return nil, fmt.Errorf("unexpected StatusCode on SetMetadata response")
	}

	// Maybe we don't need the body, but adding general approach
	respBytes, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("error reading response body for SetMetadata: %w", err)
	}
	metaDataResponse := &pendo.SetMetadataResponse{}
	err = json.Unmarshal(respBytes, metaDataResponse)
	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("error parsing SetMetadata response: %w", err)
	}

	logger.Debug("pendo SetMetadata successful",
		slog.String("kind", string(metaDataResponse.Kind)),
		slog.Int64("total", metaDataResponse.Total),
		slog.Int64("failed", metaDataResponse.Failed),
		slog.Int64("updated", metaDataResponse.Updated),
		slog.Any("missing", metaDataResponse.Missing),
	)
	return metaDataResponse, nil
}

// SendTrackEvent
// See: https://engageapi.pendo.io/?bash%23#e45be48e-e01f-4f0a-acaa-73ef6851c4ac
func (c *pendoClient) SendTrackEvent(ctx context.Context, track *pendo.TrackRequest) error {
	logger := app_context.LogFromCtx(ctx)
	if err := c.guardSendTrackEvent(track); err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("bad argument for SendTrackEvent: %w", err)
	}

	// https://github.com/RedHatInsights/cloud-connector/blob/master/internal/pendo_transmitter/pendo_reporter.go#L51
	reqBody, err := json.Marshal(track)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("error generating request body for SendTrackEvent")
	}

	// Prepare the request
	url := c.Config.Clients.PendoBaseURL + "/data/track"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("error preparing request for SendTrackEvent")
	}

	// Add headers to the request
	req.Header.Set("content-type", "application/json")
	req.Header.Set(header.HeaderXPendoIntegrationKey, c.Config.Clients.PendoTrackEventKey)

	// Launch request
	resp, err := c.Client.Do(req)
	if err != nil {
		logger.Error("doing request to pendo service")
		return err
	}
	defer resp.Body.Close() // nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		logger.Error("expected StatusCode=" + http.StatusText(http.StatusOK) + " but received StatusCode=" + http.StatusText(resp.StatusCode))
		return fmt.Errorf("unexpected StatusCode on SendTrackEvent response")
	}

	logger.Debug("pendo SendTrackEvent successful")
	return nil
}

//
// ----- Private methods ------
//

func newClient(cfg *config.Config) *pendoClient {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if cfg.Clients.PendoBaseURL == "" {
		panic("'PendoBaseURL' is empty")
	}
	if cfg.Clients.PendoAPIKey == "" {
		panic("'PendoAPIKey' is empty")
	}
	if cfg.Clients.PendoTrackEventKey == "" {
		panic("'PendoTrackEventKey' is empty")
	}
	if cfg.Clients.PendoRequestTimeoutSecs == 0 {
		cfg.Clients.PendoRequestTimeoutSecs = 3
	}
	client := &http.Client{
		Timeout: time.Duration(cfg.Clients.PendoRequestTimeoutSecs) * time.Second,
	}
	return &pendoClient{
		Config: cfg,
		Client: client,
	}
}

func (c *pendoClient) guardSetMetadata(kind pendo.Kind, group pendo.Group, metrics pendo.SetMetadataRequest) error {
	if kind == "" {
		return fmt.Errorf("'kind' is an empty string")
	}
	if group == "" {
		return fmt.Errorf("'group' is an empty string")
	}
	if metrics == nil {
		return fmt.Errorf("'metrics' is nil")
	}
	for i := range metrics {
		if metrics[i].VisitorID == "" {
			return fmt.Errorf("'metrics[%d].VisitorID' is an empty string", i)
		}
	}
	// TODO check 'kind' and 'group' contains valid characters
	//      we need additional information to do that
	return nil
}

func (c *pendoClient) guardSendTrackEvent(track *pendo.TrackRequest) error {
	if track == nil {
		return fmt.Errorf("'track' is nil")
	}
	if track.AccountID == "" {
		return fmt.Errorf("'track.AccountID' is an empty string")
	}
	if track.Type != "track" {
		return fmt.Errorf("'track.Type' is '%s', but '%s' was expected", track.Type, "track")
	}
	if track.Event == "" {
		return fmt.Errorf("'track.Event' is an empty string")
	}
	if track.VisitorID == "" {
		return fmt.Errorf("'track.VisitorID' is an empty string")
	}
	if track.Timestamp <= 0 {
		return fmt.Errorf("'track.Timestamp' is invalid")
	}
	return nil
}
