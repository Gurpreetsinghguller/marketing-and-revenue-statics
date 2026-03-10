package eventservice

import (
	"context"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/logger"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/event"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/event/aggregator"
	"github.com/sirupsen/logrus"
)

// EventService orchestrates event processing from broker and persistence
type EventService struct {
	broker     event.EventBroker
	eventRepo  domain.EventRepo
	metricsAgg *aggregator.MetricsAggregator
	log        *logrus.Logger
}

// NewEventService creates a new event service
func New(
	broker event.EventBroker,
	eventRepo domain.EventRepo,
	metricsAgg *aggregator.MetricsAggregator,
) *EventService {
	return &EventService{
		broker:     broker,
		eventRepo:  eventRepo,
		metricsAgg: metricsAgg,
		log:        logger.Get(),
	}
}

// Start begins consuming events and processing them
func (s *EventService) Start(ctx context.Context) error {
	return s.broker.Start(ctx, s.handleEvent)
}

// handleEvent processes each event received from broker
func (s *EventService) handleEvent(ctx context.Context, event *event.Event) error {
	// 1. Persist event to DB
	// we have to adapt the Broker Event to our domain Event Model

	if err := s.eventRepo.Create(s.adaptEventToDomain(event)); err != nil {
		s.log.WithError(err).
			WithField("campaign_id", event.CampaignID).
			WithField("event_type", event.EventType).
			Error("failed to persist event to database")
		return err
	}

	// 2. Update real-time metrics (event-time computation)
	if err := s.metricsAgg.ProcessEvent(ctx, s.adaptEventToDomain(event)); err != nil {
		// Log but don't fail - event is already persisted
		s.log.WithError(err).
			WithField("campaign_id", event.CampaignID).
			Warn("failed to update metrics for event")
		// Return nil so that subsequent events are still processed
	}

	return nil
}

// Close gracefully shuts down the event service
func (s *EventService) Close() error {
	// Flush any remaining metrics
	ctx := context.Background()
	if err := s.metricsAgg.FlushMetrics(ctx); err != nil {
		s.log.WithError(err).Warn("failed to flush metrics on shutdown")
	}

	// Close broker connection
	if err := s.broker.Close(); err != nil {
		s.log.WithError(err).Warn("failed to close broker")
		return err
	}

	return nil
}

// Health checks the health of the event service
func (s *EventService) Health(ctx context.Context) error {
	return s.broker.Health(ctx)
}

func (s *EventService) adaptEventToDomain(event *event.Event) *domain.Event {
	return &domain.Event{
		ID: event.ID}
}

func (s *EventService) PublishEvent(ctx context.Context, event *domain.Event) error {
	return nil
}
