package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api_header "github.com/hmsidm/internal/api/header"
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/infrastructure/middleware"
	"github.com/hmsidm/internal/interface/client"
	test_client "github.com/hmsidm/internal/test/client"
	"github.com/labstack/echo/v4"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockInventoryServer(
	t *testing.T,
	iden *identity.Identity,
	status int,
	body string,
) (
	*echo.Echo,
	echo.Context,
	*httptest.ResponseRecorder,
) {
	requestHeaders := map[string]string{
		"X-Rh-Identity": api_header.EncodeIdentity(iden),
	}
	return test_client.NewHandlerTester(t,
		http.MethodGet,
		"/api/inventory/v1/hosts",
		"/api/inventory/v1/hosts",
		nil,
		"",
		requestHeaders,
		status,
		body,
		nil,
		middleware.CreateContext(),
		middleware.EnforceIdentityWithConfig(middleware.NewIdentityConfig()),
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
			HostInventoryBaseUrl: "http://localhost:8010/api/inventory/v1",
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
		fqdn  = "server.hmsidm-dev.test"
	)
	cfg := config.Config{}
	// e, ctx, rec := newMockInventoryServer(t,
	iden := identity.Identity{
		OrgID: orgId,
		Type:  "System",
		System: identity.System{
			CommonName: cn,
			CertType:   "system",
		},
	}
	e, _, rec := newMockInventoryServer(t,
		&iden,
		http.StatusOK,
		helperBodySuccess(id, orgId, fqdn, cn),
	)
	defer e.Shutdown(context.Background())

	// TODO Call the echo handler using the context
	// sqlMock, db, err := test.NewSqlMock(nil)
	// _, db, err := test.NewSqlMock(nil)
	// require.NoError(t, err)
	// app := impl.NewHandler(&cfg, db, nil, NewHostInventory(&cfg))
	// h := public.ServerInterfaceWrapper{
	// 	Handler: app,
	// }

	// err = h.RegisterIpaDomain(ctx)
	cfg.Clients.HostInventoryBaseUrl = fmt.Sprintf("http://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	cli := NewHostInventory(&cfg)
	host, err := cli.GetHostByCN(api_header.EncodeIdentity(&iden), cn)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, client.InventoryHost{
		ID:                    id,
		SubscriptionManagerId: cn,
		FQDN:                  fqdn,
	}, host)
}

func TestGetHostByCNErrors(t *testing.T) {
	const (
		id    = "93bb346a-4297-4952-9ec4-f53b3a5006c2"
		cn    = "c91e72f6-c518-11ed-bd88-482ae3863d30"
		orgId = "11111"
		fqdn  = "server.hmsidm-dev.test"
	)
	var (
		// rec *httptest.ResponseRecorder
		e *echo.Echo
	)
	cfg := config.Config{}
	iden := identity.Identity{
		OrgID: orgId,
		Type:  "System",
		System: identity.System{
			CommonName: cn,
			CertType:   "system",
		},
	}

	e, _, _ = newMockInventoryServer(t,
		&iden,
		http.StatusBadRequest,
		"",
	)
	defer e.Shutdown(context.Background())

	// Failure because a wrong base url
	cfg.Clients.HostInventoryBaseUrl = fmt.Sprintf("httpvf://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	cli := NewHostInventory(&cfg)
	host, err := cli.GetHostByCN(api_header.EncodeIdentity(&iden), cn)
	require.EqualError(t,
		err,
		fmt.Sprintf("Get \"httpvf://localhost:%s/api/inventory/v1/hosts?filter%%5Bsystem_profile%%5D%%5Bowner_id%%5D=c91e72f6-c518-11ed-bd88-482ae3863d30\": unsupported protocol scheme \"httpvf\"",
			readPort(
				e.Listener.
					Addr().
					String(),
			),
		),
	)
	assert.Equal(t, client.InventoryHost{}, host)

	// Error unmarshalling the body response
	// for GET /api/inventory/v1/hosts
	e, _, _ = newMockInventoryServer(t,
		&iden,
		http.StatusBadRequest,
		"{",
	)
	defer e.Shutdown(context.Background())
	// Error unmarshalling response
	cfg.Clients.HostInventoryBaseUrl = fmt.Sprintf("http://localhost:%s/api/inventory/v1", readPort(e.Listener.Addr().String()))
	t.Logf("Listening for: %s", cfg.Clients.HostInventoryBaseUrl)
	cli = NewHostInventory(&cfg)
	host, err = cli.GetHostByCN(api_header.EncodeIdentity(&iden), cn)
	require.EqualError(t, err, "400 Bad Request")
	assert.Equal(t, client.InventoryHost{}, host)

}
