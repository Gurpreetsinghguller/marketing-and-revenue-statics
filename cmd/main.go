package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/config"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/logger"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/event"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/event/aggregator"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/event/eventservice"
	mqtt "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/event/mqttbroker"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	campaign_repo "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/campaign"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
	event_repo "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/event"
	metrics_repo "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/metrics"
	user_repo "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/user"
	analytics_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/usecase"
	auth_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/usecase"
	campaign_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/usecase"
	engagement_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/usecase"
	event_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/usecase"
	profile_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/usecase"
	v1 "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/router/v1"
	"github.com/sirupsen/logrus"
)

// TODO: Handle context cancellation and timeouts in handlers and usecases
// TODO: Improve error handling and logging with more context (e.g. request IDs, user IDs)
// TODO: Increase test coverage

// Log and Config are singletons that can be accessed throughout the application. We initialize them in init() function.
var log *logrus.Logger
var cfg *config.Config

func init() {
	log = logger.Get()
	log.Info("application initializing")

	cfg, cfgErr := config.Load(config.DefaultConfigPath)
	if cfgErr != nil {
		log.WithError(cfgErr).Warn("failed to load config file; using defaults")
		cfg = config.Default()
	}
	logger.Configure(cfg.Log.Level)
}

// The question is, Do we want to not start REST server if broker fails to start?
//
//	Or do we want to start REST server and log broker errors but keep REST functional?
func main() {

	// Initialize database
	storage := db.NewStorageMgr()
	defer storage.Close()

	/// -----------------------STORAGE will act as Strategy Pattern for Repositories-----------------------
	// We can easily swap out the storage implementation (e.g. switch from in-memory to SQL) without changing repository logic
	userRepo := user_repo.NewUserRepository(storage)
	campaignRepo := campaign_repo.NewCampaignRepository(storage)
	eventRepo := event_repo.NewEventRepository(storage)

	authUC := auth_usecase.NewAuthUseCase(userRepo)
	profileUC := profile_usecase.NewProfileUseCase(userRepo)
	campaignUC := campaign_usecase.NewCampaignUseCase(campaignRepo)
	eventUC := event_usecase.NewEventUseCase(eventRepo)
	analyticsUC := analytics_usecase.NewAnalyticsUseCase(campaignRepo, eventRepo)
	engagementUC := engagement_usecase.NewEngagementUseCase(eventRepo)

	metricsRepo := metrics_repo.NewMetricsRepository(storage)
	metricsAgg := aggregator.NewMetricsAggregator(metricsRepo)

	broker, err := InitializeBroker(cfg)
	if err != nil {
		log.WithError(err).Error("failed to initialize event broker")
		return
	}

	eventService := eventservice.New(
		broker,
		eventRepo,
		metricsAgg,
	)

	go func() {
		ctx := context.Background()
		if err := eventService.Start(ctx); err != nil {
			log.WithError(err).Error("failed to start event broker")
		}
	}()

	// Initialize router and setup routes
	router := v1.NewRouter()
	muxRouter := router.InitHTTPRoutes(
		authUC,
		profileUC,
		campaignUC,
		eventUC,
		analyticsUC,
		engagementUC,
	)

	// Add global middleware
	muxRouter.Use(middleware.LoggingMiddleware)
	muxRouter.Use(middleware.CORSMiddleware)

	// Start HTTP server
	port := strings.TrimSpace(cfg.Server.Port)
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.WithField("address", "http://localhost"+port).Info("server starting")
	log.WithField("url", "http://localhost"+port+"/api/v1/health").Info("health endpoint")
	log.WithField("url", "http://localhost"+port+"/api/v1/docs").Info("openapi endpoint")

	if err := http.ListenAndServe(port, muxRouter); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("server error")
	}

}

// This is Factory Pattern
func InitializeBroker(cfg *config.Config) (event.EventBroker, error) {
	// Initialize broker based on config
	var broker event.EventBroker
	switch strings.ToLower(cfg.Broker.Type) {
	case "mqtt":
		broker = mqtt.NewBroker(cfg.Broker.URL, cfg.Broker.ClientID, cfg.Broker.Topic)
	default:
		log.Warnf("unsupported broker type '%s'; no events will be consumed", cfg.Broker.Type)
		return nil, nil
	}
	return broker, nil
}
