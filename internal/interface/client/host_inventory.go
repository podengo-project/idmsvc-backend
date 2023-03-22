package client

// HostInventory represent the client to reach
// out the host inventory service and abstract
// the necessary operations.
type HostInventory interface {
	GetHostByCN(iden string, requestId string, cn string) (InventoryHost, error)
}

// InventoryHost only cover the necessary information
// used by idm-domains-backend service from the host
// inventory service when requesting a filtered /hosts
// request.
type InventoryHost struct {
	ID                    string `json:"id"`
	SubscriptionManagerId string `json:"subscription_manager_id"`
	FQDN                  string `json:"fqdn"`
}

// InventoryHostPage represent a paginated list of results
// from the host inventory for the GET /hosts endpoint.
type InventoryHostPage struct {
	Total   int             `json:"total"`
	Count   int             `json:"count"`
	Page    int             `json:"page"`
	PerPage int             `json:"per_page"`
	Results []InventoryHost `json:"results"`
}
