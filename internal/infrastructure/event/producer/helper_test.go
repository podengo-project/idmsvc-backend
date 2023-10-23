package producer

import (
	"github.com/podengo-project/idmsvc-backend/internal/config"
)

func helperGetKafkaConfig() *config.Kafka {
	return &config.Kafka{
		Bootstrap: struct{ Servers string }{
			Servers: "localhost:9092",
		},
		Request: struct {
			Timeout  struct{ Ms int }
			Required struct{ Acks int }
		}{
			Required: struct{ Acks int }{
				Acks: -1,
			},
		},
		Message: struct {
			Send struct{ Max struct{ Retries int } }
		}{
			Send: struct{ Max struct{ Retries int } }{
				Max: struct{ Retries int }{
					Retries: 3,
				},
			},
		},
		Retry: struct{ Backoff struct{ Ms int } }{
			Backoff: struct{ Ms int }{
				Ms: 300,
			},
		},
	}
}
