package aggregator

import (
	"context"
	"sync"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/logger"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/domain"
	"github.com/sirupsen/logrus"
)

// MetricsAggregator computes metrics at event-time
type MetricsAggregator struct {
	metricsRepo domain.MetricsRepo
	mu          sync.RWMutex
	log         *logrus.Logger
	// In-memory buffers for windowing (e.g., hourly buckets)
	hourlyMetrics map[string]*domain.CampaignMetrics
	dailyMetrics  map[string]*domain.CampaignMetrics
}

// NewMetricsAggregator creates a new aggregator
func NewMetricsAggregator(metricsRepo domain.MetricsRepo) *MetricsAggregator {
	return &MetricsAggregator{
		metricsRepo:   metricsRepo,
		hourlyMetrics: make(map[string]*domain.CampaignMetrics),
		dailyMetrics:  make(map[string]*domain.CampaignMetrics),
		log:           logger.Get(),
	}
}

// ProcessEvent updates metrics when an event is received
func (m *MetricsAggregator) ProcessEvent(ctx context.Context, event *domain.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	timestamp, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
		// Try alternative timestamp format
		timestamp, err = time.Parse("2006-01-02T15:04:05Z07:00", event.Timestamp)
		if err != nil {
			m.log.WithError(err).WithField("timestamp", event.Timestamp).
				Error("failed to parse event timestamp")
			return err
		}
	}

	timestamp = timestamp.UTC()
	date := timestamp.Format("2006-01-02")
	hour := timestamp.Format("2006-01-02 15:00:00")

	// Update hourly metrics
	hourKey := event.CampaignID + "|" + hour
	if m.hourlyMetrics[hourKey] == nil {
		m.hourlyMetrics[hourKey] = &domain.CampaignMetrics{
			CampaignID:  event.CampaignID,
			Date:        date,
			Hour:        hour,
			LastUpdated: time.Now().UTC(),
		}
	}

	hourMetrics := m.hourlyMetrics[hourKey]

	// Update metrics based on event type
	switch event.EventType {
	case domain.EventType("impressions"):
		hourMetrics.Impressions++
	case domain.EventType("clicks"):
		hourMetrics.Clicks++
	case domain.EventType("conversions"):
		hourMetrics.Conversions++
		hourMetrics.TotalRevenue += event.Metadata.Amount
	}

	// Recalculate rates
	if hourMetrics.Impressions > 0 {
		hourMetrics.CTR = float64(hourMetrics.Clicks) / float64(hourMetrics.Impressions)
	}
	if hourMetrics.Conversions > 0 {
		hourMetrics.AvgRevenue = hourMetrics.TotalRevenue / float64(hourMetrics.Conversions)
	}
	if hourMetrics.Clicks > 0 {
		hourMetrics.ConversionRate = float64(hourMetrics.Conversions) / float64(hourMetrics.Clicks)
	}

	hourMetrics.LastUpdated = time.Now().UTC()

	// Also update daily metrics
	dateKey := event.CampaignID + "|" + date
	if m.dailyMetrics[dateKey] == nil {
		m.dailyMetrics[dateKey] = &domain.CampaignMetrics{
			CampaignID:  event.CampaignID,
			Date:        date,
			LastUpdated: time.Now().UTC(),
		}
	}

	dailyMetrics := m.dailyMetrics[dateKey]

	// Update daily metrics
	switch event.EventType {
	case domain.EventType("impressions"):
		dailyMetrics.Impressions++
	case domain.EventType("clicks"):
		dailyMetrics.Clicks++
	case domain.EventType("conversions"):
		dailyMetrics.Conversions++
		dailyMetrics.TotalRevenue += event.Metadata.Amount
	}

	// Recalculate daily rates
	if dailyMetrics.Impressions > 0 {
		dailyMetrics.CTR = float64(dailyMetrics.Clicks) / float64(dailyMetrics.Impressions)
	}
	if dailyMetrics.Conversions > 0 {
		dailyMetrics.AvgRevenue = dailyMetrics.TotalRevenue / float64(dailyMetrics.Conversions)
	}
	if dailyMetrics.Clicks > 0 {
		dailyMetrics.ConversionRate = float64(dailyMetrics.Conversions) / float64(dailyMetrics.Clicks)
	}

	dailyMetrics.LastUpdated = time.Now().UTC()

	// Persist hourly metrics to DB
	if err := m.metricsRepo.SaveMetrics(hourMetrics); err != nil {
		m.log.WithError(err).WithField("campaign_id", event.CampaignID).
			Error("failed to save hourly metrics")
		return err
	}

	return nil
}

// FlushMetrics periodically saves all buffered metrics
func (m *MetricsAggregator) FlushMetrics(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, metrics := range m.dailyMetrics {
		if err := m.metricsRepo.SaveMetrics(metrics); err != nil {
			m.log.WithError(err).WithField("campaign_id", metrics.CampaignID).
				Error("failed to flush daily metrics")
			return err
		}
	}

	m.log.Info("metrics flushed to database")
	return nil
}

// GetHourlyMetrics returns in-memory hourly metrics snapshot
func (m *MetricsAggregator) GetHourlyMetrics(campaignID, hour string) *domain.CampaignMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := campaignID + "|" + hour
	return m.hourlyMetrics[key]
}

// GetDailyMetrics returns in-memory daily metrics snapshot
func (m *MetricsAggregator) GetDailyMetrics(campaignID, date string) *domain.CampaignMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := campaignID + "|" + date
	return m.dailyMetrics[key]
}
