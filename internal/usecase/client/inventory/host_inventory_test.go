package inventory

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	api_header "github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	client_inventory "github.com/podengo-project/idmsvc-backend/internal/interface/client/inventory"
	test_client "github.com/podengo-project/idmsvc-backend/internal/test/client"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockInventoryServer(
	t *testing.T,
	xrhid *identity.XRHID,
	status int,
	body string,
	requestHeaders map[string]string,
) (
	*echo.Echo,
	echo.Context,
	*httptest.ResponseRecorder,
) {
	rh := map[string]string{
		"X-Rh-Identity": api_header.EncodeXRHID(xrhid),
	}
	return test_client.NewHandlerTester(t,
		http.MethodGet,
		"/api/inventory/v1/hosts",
		"/api/inventory/v1/hosts",
		nil,
		"",
		rh,
		status,
		body,
		nil,
		middleware.CreateContext(),
		middleware.EnforceIdentityWithConfig(
			&middleware.IdentityConfig{}),
	)
}

func helperBodySuccess(
	id string,
	orgId string,
	fqdn string,
	subscriptionManagerId string,
) string {
	displayName := fqdn
	return fmt.Sprintf(`{
		"total": 1,
		"count": 1,
		"page": 1,
		"per_page": 50,
		"results": [
		  {
			"insights_id": "ca5b5f40-d97a-425e-94ca-0f099295677d",
			"subscription_manager_id": "%s",
			"satellite_id": null,
			"bios_uuid": "3c89d5a1-7287-4856-a20e-b85337a1771a",
			"ip_addresses": [
			  "10.0.169.222"
			],
			"fqdn": "%s",
			"mac_addresses": [
			  "fa:16:3e:7b:5e:3a",
			  "00:00:00:00:00:00"
			],
			"provider_id": null,
			"provider_type": null,
			"id": "%s",
			"account": "11474377",
			"org_id": "%s",
			"display_name": "%s",
			"ansible_host": null,
			"facts": [],
			"reporter": "cloud-connector",
			"per_reporter_staleness": {
			  "cloud-connector": {
				"check_in_succeeded": true,
				"stale_timestamp": "2023-03-16T10:54:40+00:00",
				"last_check_in": "2023-03-15T08:54:40.613575+00:00"
			  },
			  "puptoo": {
				"check_in_succeeded": true,
				"stale_timestamp": "2023-03-16T13:54:24.432198+00:00",
				"last_check_in": "2023-03-15T08:54:24.714327+00:00"
			  }
			},
			"stale_timestamp": "2023-03-16T10:54:40+00:00",
			"stale_warning_timestamp": "2023-03-23T10:54:40+00:00",
			"culled_timestamp": "2023-03-30T10:54:40+00:00",
			"created": "2023-03-15T08:54:24.740431+00:00",
			"updated": "2023-03-15T08:55:14.950975+00:00"
		  }
		]
	  }`, subscriptionManagerId, fqdn, id, orgId, displayName)
}

func readPort(addr string) string {
	items := strings.Split(addr, ":")
	if len(items) < 2 {
		return ""
	}
	return items[len(items)-1]
}

func helperBodyEmpty() string {
	return fmt.Sprintf(`{
		"total": 0,
		"count": 0,
		"page": 0,
		"per_page": 50,
		"results": []
	  }`)
}

func TestNewHostInventory(t *testing.T) {
	cfg := &config.Config{
		Clients: config.Clients{
			InventoryBaseURL: "http://localhost:8010/api/inventory/v1",
		},
	}
	result := NewHostInventory(cfg)
	assert.NotNil(t, result)
}

func TestGetHostByCN(t *testing.T) {
	const (
		id    = "93bb346a-4297-4952-9ec4-f53b3a5006c2"
		cn    = "c91e72f6-c518-11ed-bd88-482ae3863d30"
		orgId = "11111"
		fqdn  = "server.idmsvc-dev.test"
	)
	cfg := config.Config{}

	// See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/identities/cert.json
	xrhid := identity.XRHID{
		Identity: identity.Identity{
			OrgID: orgId,
			Type:  "System",
			System: identity.System{
				CommonName: cn,
				CertType:   "system",
			},
		},
	}
	e, _, rec := newMockInventoryServer(t,
		&xrhid,
		http.StatusOK,
		helperBodySuccess(id, orgId, fqdn, cn),
		nil,
	)
	defer e.Shutdown(context.Background())

	cfg.Clients.InventoryBaseURL = fmt.Sprintf("http://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	cli := NewHostInventory(&cfg)
	host, err := cli.GetHostByCN(api_header.EncodeXRHID(&xrhid), "test", cn)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, client_inventory.InventoryHost{
		ID:                    id,
		SubscriptionManagerId: cn,
		FQDN:                  fqdn,
	}, host)
}

func TestGetHostByCNErrors(t *testing.T) {
	const (
		id      = "93bb346a-4297-4952-9ec4-f53b3a5006c2"
		cn      = "c91e72f6-c518-11ed-bd88-482ae3863d30"
		cnWrong = "f2f45bc2-c897-11ed-a09f-482ae3863d30"
		orgId   = "11111"
		fqdn    = "server.idmsvc-dev.test"
	)

	// rec *httptest.ResponseRecorder
	var e *echo.Echo
	cfg := config.Config{}
	xrhid := identity.XRHID{
		Identity: identity.Identity{
			OrgID: orgId,
			Type:  "System",
			System: identity.System{
				CommonName: cn,
				CertType:   "system",
			},
		},
	}

	e, _, _ = newMockInventoryServer(t,
		&xrhid,
		http.StatusBadRequest,
		"",
		nil,
	)
	defer e.Shutdown(context.Background())

	// Failure because a wrong base url
	cfg.Clients.InventoryBaseURL = fmt.Sprintf("lhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	cli := NewHostInventory(&cfg)
	host, err := cli.GetHostByCN(api_header.EncodeXRHID(&xrhid), "test", cn)
	require.EqualError(t,
		err,
		fmt.Sprintf("Get \"lhost:%s/api/inventory/v1/hosts?filter%%5Bsystem_profile%%5D%%5Bowner_id%%5D=c91e72f6-c518-11ed-bd88-482ae3863d30\": unsupported protocol scheme \"lhost\"",
			readPort(
				e.Listener.
					Addr().
					String(),
			),
		),
	)
	assert.Equal(t, client_inventory.InventoryHost{}, host)

	// Failure request because wrong Location header parsing
	// (forcing req.Do operation to fail)
	cfg.Clients.InventoryBaseURL = fmt.Sprintf("http://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	cli = NewHostInventory(&cfg)
	host, err = cli.GetHostByCN(api_header.EncodeXRHID(&xrhid), "test", cn)
	require.EqualError(t,
		err,
		"400 Bad Request",
	)
	assert.Equal(t, client_inventory.InventoryHost{}, host)

	// Error unmarshalling the body response
	e, _, _ = newMockInventoryServer(t,
		&xrhid,
		http.StatusOK,
		"{",
		nil,
	)
	defer e.Shutdown(context.Background())
	// Error unmarshalling response
	cfg.Clients.InventoryBaseURL = fmt.Sprintf("http://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	t.Logf("Listening for: %s", cfg.Clients.InventoryBaseURL)
	cli = NewHostInventory(&cfg)
	host, err = cli.GetHostByCN(api_header.EncodeXRHID(&xrhid), "test", cn)
	require.EqualError(t, err, "unexpected end of JSON input")
	assert.Equal(t, client_inventory.InventoryHost{}, host)

	// Force 'Failed to look up'
	e, _, _ = newMockInventoryServer(t,
		&xrhid,
		http.StatusOK,
		helperBodyEmpty(),
		nil,
	)
	defer e.Shutdown(context.Background())
	// Error unmarshalling response
	cfg.Clients.InventoryBaseURL = fmt.Sprintf("http://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	t.Logf("Listening for: %s", cfg.Clients.InventoryBaseURL)
	cli = NewHostInventory(&cfg)
	host, err = cli.GetHostByCN(api_header.EncodeXRHID(&xrhid), "test", cn)
	require.EqualError(t, err, fmt.Sprintf("Failed to look up 'cn=%s'", cn))
	assert.Equal(t, client_inventory.InventoryHost{}, host)

	// Force 'Look up does not match'
	e, _, _ = newMockInventoryServer(t,
		&xrhid,
		http.StatusOK,
		helperBodySuccess(id, orgId, fqdn, cnWrong),
		nil,
	)
	defer e.Shutdown(context.Background())
	// Error unmarshalling response
	cfg.Clients.InventoryBaseURL = fmt.Sprintf("http://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	t.Logf("Listening for: %s", cfg.Clients.InventoryBaseURL)
	cli = NewHostInventory(&cfg)
	host, err = cli.GetHostByCN(api_header.EncodeXRHID(&xrhid), "test", cn)
	require.EqualError(t, err, fmt.Sprintf("Looked up 'cn=%s' does not match", cn))
	assert.Equal(t, client_inventory.InventoryHost{}, host)
}
