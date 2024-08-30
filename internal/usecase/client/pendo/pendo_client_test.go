package pendo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	builder_pendo "github.com/podengo-project/idmsvc-backend/internal/test/builder/clients/pendo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL           = "http://localhost:8031"
	testAPIKey        = "test-api-key"
	testTrackEventKey = "test-track-event-key"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// PendoHash       kind       group      visitorId  fieldName => value
type PendoHash map[string]map[string]map[string]map[string]any

func helperPendoConfig() *config.Config {
	return &config.Config{
		Clients: config.Clients{
			PendoBaseURL:            baseURL,
			PendoAPIKey:             testAPIKey,
			PendoTrackEventKey:      testTrackEventKey,
			PendoRequestTimeoutSecs: 1,
		},
	}
}

func helperNewPendo(cfg *config.Config, fn RoundTripFunc) pendo.Pendo {
	client := newClient(cfg)
	client.Client.Transport = RoundTripFunc(fn)
	return client
}

func TestNewPendo(t *testing.T) {
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		newClient(nil)
	})

	assert.PanicsWithValue(t, "'PendoBaseURL' is empty", func() {
		newClient(&config.Config{
			Clients: config.Clients{
				PendoBaseURL:            "",
				PendoAPIKey:             "",
				PendoRequestTimeoutSecs: 0,
			},
		})
	})

	assert.PanicsWithValue(t, "'PendoAPIKey' is empty", func() {
		newClient(&config.Config{
			Clients: config.Clients{
				PendoBaseURL:            baseURL,
				PendoAPIKey:             "",
				PendoRequestTimeoutSecs: 0,
			},
		})
	})

	assert.PanicsWithValue(t, "'PendoTrackEventKey' is empty", func() {
		newClient(&config.Config{
			Clients: config.Clients{
				PendoBaseURL:            baseURL,
				PendoAPIKey:             testAPIKey,
				PendoTrackEventKey:      "",
				PendoRequestTimeoutSecs: 0,
			},
		})
	})

	client := newClient(&config.Config{
		Clients: config.Clients{
			PendoBaseURL:            baseURL,
			PendoAPIKey:             testAPIKey,
			PendoTrackEventKey:      testTrackEventKey,
			PendoRequestTimeoutSecs: 0,
		},
	})
	require.NotNil(t, client)
}

func TestGuardSetMetadata(t *testing.T) {
	cfg := helperPendoConfig()
	client := newClient(cfg)

	require.EqualError(t, client.guardSetMetadata("", "", nil), "'kind' is an empty string")
	require.EqualError(t, client.guardSetMetadata("mykind", "", nil), "'group' is an empty string")
	require.EqualError(t, client.guardSetMetadata("mykind", "mygroup", nil), "'metrics' is nil")
	pendoMetrics := make(pendo.SetMetadataRequest, 0, 1)
	require.NoError(t, client.guardSetMetadata("mykind", "mygroup", pendoMetrics), "an empty slice does not report error")
	pendoMetrics = append(pendoMetrics, pendo.SetMetadataDetailsRequest{})
	require.EqualError(t, client.guardSetMetadata("mykind", "mygroup", pendoMetrics), "'metrics[0].VisitorID' is an empty string")
	pendoMetrics[0].VisitorID = "my-test-visitor-id"

	// Success case
	require.NoError(t, client.guardSetMetadata("mykind", "mygroup", pendoMetrics))
}

func helperSetMetadataPrepareRequest(t *testing.T, req *http.Request, kind pendo.Kind, group pendo.Group) pendo.SetMetadataRequest {
	metrics := make(pendo.SetMetadataRequest, 0, 10)
	// Check request
	assert.Equal(t, baseURL+"/api/v1/metadata/"+string(kind)+"/"+string(group)+"/value", req.URL.String())
	reqBytes, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	require.NotNil(t, reqBytes)
	err = json.Unmarshal(reqBytes, &metrics)
	require.NoError(t, err)
	require.NotNil(t, metrics)
	return metrics
}

func helperSetMetadataPrepareResponse(t *testing.T, metrics pendo.SetMetadataRequest, kind pendo.Kind, group pendo.Group, store PendoHash) *pendo.SetMetadataResponse {
	var (
		ok           bool
		groups       map[string]map[string]map[string]any
		visitors     map[string]map[string]any
		storeMetrics map[string]any
	)
	respBuilder := builder_pendo.NewSetMetadataResponse().
		WithTotal(int64(len(metrics))).
		WithKind(kind)
	require.NotNil(t, respBuilder)
	for i := range metrics {
		if groups, ok = store[string(kind)]; !ok {
			respBuilder.AddMissing(metrics[i].VisitorID)
			respBuilder.IncFailed()
			continue
		}
		if visitors, ok = groups[string(group)]; !ok {
			respBuilder.AddMissing(metrics[i].VisitorID)
			respBuilder.IncFailed()
			continue
		}
		if storeMetrics, ok = visitors[metrics[i].VisitorID]; !ok {
			respBuilder.AddMissing(metrics[i].VisitorID)
			respBuilder.IncFailed()
			continue
		}
		respBuilder.IncUpdated()
		for k, v := range metrics[i].Values {
			storeMetrics[k] = v
		}
	}

	resp := respBuilder.Build()
	require.NotNil(t, resp)
	return resp
}

func TestSetMetadata(t *testing.T) {
	kind := pendo.KindAccount
	group := pendo.Group("custom")
	var store PendoHash = PendoHash{
		string(kind): {
			string(group) /* group */ : {
				"test-visitor-id" /* visitorId */ : {
					"test-field-1" /* fieldName */ : 4,                 /* value */
					"test-field-2" /* fieldName */ : "someStringValue", /* value */
				},
			},
		},
	}
	// https://hassansin.github.io/Unit-Testing-http-client-in-Go#2-by-replacing-httptransport
	cfg := helperPendoConfig()
	client := helperNewPendo(cfg, func(req *http.Request) *http.Response {
		metrics := helperSetMetadataPrepareRequest(t, req, kind, group)
		resp := helperSetMetadataPrepareResponse(t, metrics, kind, group, store)
		bytesResp, err := json.Marshal(resp)
		require.NoError(t, err)
		require.NotNil(t, bytesResp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(bytesResp)),
			Header:     make(http.Header),
		}
	})

	// Panic when context is nil
	assert.PanicsWithValue(t, "'ctx' is nil", func() {
		client.SetMetadata(nil, "", "", nil)
	})

	// Error when wrong argument
	ctx := app_context.CtxWithLog(context.TODO(), slog.Default())
	resp, err := client.SetMetadata(ctx, "", "", nil)
	require.EqualError(t, err, "bad arguments: 'kind' is an empty string")
	require.Nil(t, resp)

	// Success
	metrics := make(pendo.SetMetadataRequest, 0, 1)
	metrics = append(metrics, pendo.SetMetadataDetailsRequest{
		VisitorID: "test-visitor-id",
		Values: map[string]any{
			"test-field-1": 2,
			"test-field-2": "anotherString",
		},
	})
	metrics = append(metrics, pendo.SetMetadataDetailsRequest{
		VisitorID: "thisVisitorDoesNotExist",
	})

	// Call the data
	resp, err = client.SetMetadata(ctx, kind, group, metrics)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, int64(1), resp.Failed)
	assert.Equal(t, int64(1), resp.Updated)
	assert.Equal(t, []string{"thisVisitorDoesNotExist"}, resp.Missing)
}

func TestSetMetadataForceFailureOnDoingRequestToPendo(t *testing.T) {
	// This forces the path when c.Client.Do(req) returns nil at 'SetMetadata'
	kind := pendo.KindAccount
	group := pendo.Group("custom")
	// https://hassansin.github.io/Unit-Testing-http-client-in-Go#2-by-replacing-httptransport
	cfg := helperPendoConfig()
	client := helperNewPendo(cfg, func(req *http.Request) *http.Response {
		return nil
	})

	// Error path 1
	ctx := app_context.CtxWithLog(context.TODO(), slog.Default())
	metrics := make(pendo.SetMetadataRequest, 0, 1)
	metrics = append(metrics, pendo.SetMetadataDetailsRequest{
		VisitorID: "test-visitor-id",
		Values: map[string]any{
			"test-field-1": 2,
			"test-field-2": "anotherString",
		},
	})
	metrics = append(metrics, pendo.SetMetadataDetailsRequest{
		VisitorID: "thisVisitorDoesNotExist",
	})
	resp, err := client.SetMetadata(ctx, kind, group, metrics)
	require.EqualError(t, err, "Post \"http://localhost:8031/api/v1/metadata/account/custom/value\": http: RoundTripper implementation (pendo.RoundTripFunc) returned a nil *Response with a nil error")
	require.Nil(t, resp)
}

func TestSetMetadataErrorHttp(t *testing.T) {
	kind := pendo.KindAccount
	group := pendo.Group("custom")
	// https://hassansin.github.io/Unit-Testing-http-client-in-Go#2-by-replacing-httptransport
	cfg := helperPendoConfig()
	client := helperNewPendo(cfg, func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}
	})

	// Error path 1
	ctx := app_context.CtxWithLog(context.TODO(), slog.Default())
	metrics := make(pendo.SetMetadataRequest, 0, 1)
	metrics = append(metrics, pendo.SetMetadataDetailsRequest{
		VisitorID: "test-visitor-id",
		Values: map[string]any{
			"test-field-1": 2,
			"test-field-2": "anotherString",
		},
	})
	metrics = append(metrics, pendo.SetMetadataDetailsRequest{
		VisitorID: "thisVisitorDoesNotExist",
	})
	resp, err := client.SetMetadata(ctx, kind, group, metrics)
	require.EqualError(t, err, "unexpected StatusCode on SetMetadata response")
	require.Nil(t, resp)
}

func TestSetMetadataForceErrorUnmarshalling(t *testing.T) {
	kind := pendo.KindAccount
	group := pendo.Group("custom")
	// https://hassansin.github.io/Unit-Testing-http-client-in-Go#2-by-replacing-httptransport
	cfg := helperPendoConfig()
	client := helperNewPendo(cfg, func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("{")),
			Header:     make(http.Header),
		}
	})

	// Error path 1
	ctx := app_context.CtxWithLog(context.TODO(), slog.Default())
	metrics := make(pendo.SetMetadataRequest, 0, 1)
	resp, err := client.SetMetadata(ctx, kind, group, metrics)
	require.EqualError(t, err, "error parsing SetMetadata response: unexpected end of JSON input")
	require.Nil(t, resp)
}

func TestGuardSetTrack(t *testing.T) {
	cfg := helperPendoConfig()
	client := newClient(cfg)

	client.guardSendTrackEvent(nil)
	require.EqualError(t, client.guardSendTrackEvent(nil), "'track' is nil")
	track := pendo.TrackRequest{}
	require.EqualError(t, client.guardSendTrackEvent(&track), "'track.AccountID' is an empty string")
	track.AccountID = "my-account-id"
	require.EqualError(t, client.guardSendTrackEvent(&track), "'track.Type' is '', but 'track' was expected")
	track.Type = "track"
	require.EqualError(t, client.guardSendTrackEvent(&track), "'track.Event' is an empty string")
	track.Event = "guard-set-track-tested"
	require.EqualError(t, client.guardSendTrackEvent(&track), "'track.VisitorID' is an empty string")
	track.VisitorID = "my-visitor-id"
	require.EqualError(t, client.guardSendTrackEvent(&track), "'track.Timestamp' is invalid")
	track.Timestamp = time.Now().UTC().UnixMilli()
	require.NoError(t, client.guardSendTrackEvent(&track))
}

func TestSetTrackEvent(t *testing.T) {
	cfg := helperPendoConfig()
	client := helperNewPendo(cfg, func(req *http.Request) *http.Response {
		assert.Equal(t, baseURL+"/data/track", req.URL.String())
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	assert.PanicsWithValue(t, "'ctx' is nil", func() {
		client.SendTrackEvent(nil, nil)
	})
	ctx := app_context.CtxWithLog(context.TODO(), slog.Default())
	require.EqualError(t, client.SendTrackEvent(ctx, nil), "bad argument for SendTrackEvent: 'track' is nil")
	track := pendo.TrackRequest{
		AccountID: "my-account-id",
		Type:      "track",
		Event:     "my-event",
		VisitorID: "my-visitor-id",
		Timestamp: time.Now().UTC().UnixMilli(),
	}
	require.NoError(t, client.SendTrackEvent(ctx, &track))
}

func TestSetTrackEventErrorHttpStatus(t *testing.T) {
	cfg := helperPendoConfig()
	client := helperNewPendo(cfg, func(req *http.Request) *http.Response {
		assert.Equal(t, baseURL+"/data/track", req.URL.String())
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}
	})

	assert.PanicsWithValue(t, "'ctx' is nil", func() {
		client.SendTrackEvent(nil, nil)
	})
	ctx := app_context.CtxWithLog(context.TODO(), slog.Default())
	track := pendo.TrackRequest{
		AccountID: "my-account-id",
		Type:      "track",
		Event:     "my-event",
		VisitorID: "my-visitor-id",
		Timestamp: time.Now().UTC().UnixMilli(),
	}
	err := client.SendTrackEvent(ctx, &track)
	require.EqualError(t, err, "unexpected StatusCode on SendTrackEvent response")
}
