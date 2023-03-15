package client

type HostInventory interface {
	ListHost(iden string) (hosts []Host, err error)
}

type Host struct {
	ID                    string `json:"id"`
	SubscriptionManagerId string `json:"subscription_manager_id"`
	FQDN                  string `json:"fqdn"`
}
