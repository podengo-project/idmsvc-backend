package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hmsidm/internal/config"
	interface_client "github.com/hmsidm/internal/interface/client"
)

type hostInventory struct {
	baseURL string
}

// NewHostInventory initialize a new host inventory client.
// cfg is the reference to the configuration used by our service.
// Return an instance that accomplish HostInventory interface.
func NewHostInventory(cfg *config.Config) interface_client.HostInventory {
	return &hostInventory{
		baseURL: cfg.Clients.HostInventoryBaseUrl,
	}
}

// GetHostByCN Get the inventory host that match the cn for the
// given identity header.
// iden is the X-Rh-Identity header content as a base64 encoded string.
// cn is the Identity.System["cn"] field that is contained into the unmarshalled
// x-rh-identity header.
// Return the host matched and nil if the operation is successful, else it
// returns an empty host struct and an error with the details.
func (c *hostInventory) GetHostByCN(iden string, cn string) (
	interface_client.InventoryHost,
	error,
) {
	// https://pkg.go.dev/net/http
	client := &http.Client{}
	q := make(url.Values)
	q.Set("filter[system_profile][owner_id]", cn)
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/hosts?%s", c.baseURL, q.Encode()),
		nil)
	if err != nil {
		return interface_client.InventoryHost{}, err
	}
	req.Header.Add("X-Rh-Identity", iden)
	resp, err := client.Do(req)
	if err != nil {
		return interface_client.InventoryHost{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return interface_client.InventoryHost{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return interface_client.InventoryHost{}, fmt.Errorf("%s", resp.Status)
	}
	page := interface_client.InventoryHostPage{}
	err = json.Unmarshal(body, &page)
	if err != nil {
		return interface_client.InventoryHost{}, err
	}

	if page.Total != 1 {
		return interface_client.InventoryHost{},
			fmt.Errorf("Failed to look up 'cn=%s'", cn)
	}
	return page.Results[0], nil
}
