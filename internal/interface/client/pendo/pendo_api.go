package pendo

type SetMetadataDetailsRequest struct {
	VisitorID string         `json:"visitorId"`
	Values    map[string]any `json:"values"`
}

type SetMetadataRequest []SetMetadataDetailsRequest

type SetMetadataResponse struct {
	Total   int64    `json:"total"`
	Updated int64    `json:"updated"`
	Failed  int64    `json:"failed"`
	Missing []string `json:"missing"`
	Kind    Kind     `json:"visitor"`
}

type TrackRequest struct {
	Type       string         `json:"type"`
	Event      string         `json:"event"`
	VisitorID  string         `json:"visitorId"`
	AccountID  string         `json:"accountId"`
	Timestamp  int64          `json:"timestamp"`
	Properties map[string]any `json:"properties,omitempty"`
	Context    map[string]any `json:"context"`
	IP         string         `json:"ip"`
	UserAgent  string         `json:"userAgent"`
	URL        string         `json:"url,omitempty"`
	Title      string         `json:"title,omitempty"`
}
