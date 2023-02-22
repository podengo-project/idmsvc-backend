package impl

import (
	"context"
	"sync"

	api_event "github.com/hmsidm/internal/api/event"
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/handler/impl"
	"github.com/hmsidm/internal/infrastructure/event"
	event_handler "github.com/hmsidm/internal/infrastructure/event/handler"
	"github.com/hmsidm/internal/infrastructure/service"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type kafkaConsumer struct {
	context   context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
	config    *config.Config

	db *gorm.DB
}

func NewKafkaConsumer(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, db *gorm.DB) service.ApplicationService {
	ctx, cancel := context.WithCancel(ctx)
	return &kafkaConsumer{
		context:   ctx,
		cancel:    cancel,
		waitGroup: wg,
		config:    cfg,

		db: db,
	}
}

func (s *kafkaConsumer) Start() error {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()

		// Create event router
		eventRouter := event_handler.NewRouter()
		eventRouter.Add(api_event.TopicTodoCreated, impl.NewTodoCreatedEventHandler(s.db))

		// Start service
		event.Start(s.context, &s.config.Kafka, eventRouter)
		log.Info().Msgf("kafkaConsumer stopped")
	}()
	return nil
}

func (s *kafkaConsumer) Stop() error {
	s.cancel()
	return nil
}
