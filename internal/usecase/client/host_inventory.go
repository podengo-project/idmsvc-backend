package client

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hmsidm/internal/config"
	interface_client "github.com/hmsidm/internal/interface/client"
)

type hostInventory struct {
	baseURL string
}

func NewHostInventory(cfg *config.Config) interface_client.HostInventory {
	return &hostInventory{
		baseURL: cfg.Clients.HostInventoryBaseUrl,
	}
}

func (c *hostInventory) ListHost(iden string) (hosts []interface_client.Host, err error) {
	// https://pkg.go.dev/net/http
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.baseURL+"/hosts", nil)
	req.Header.Add("X-Rh-Identity", iden)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	hosts = []interface_client.Host{}
	err = json.Unmarshal(body, &hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}
