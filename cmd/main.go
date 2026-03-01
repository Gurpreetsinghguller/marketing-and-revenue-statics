package main

import (
	"net/http"
	"strings"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/config"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/logger"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	campaign_repo "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/campaign"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
	event_repo "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/event"
	user_repo "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/user"
	analytics_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/usecase"
	auth_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/usecase"
	campaign_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/usecase"
	engagement_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/usecase"
	event_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/usecase"
	profile_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/usecase"
	v1 "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/router/v1"
)

// TODO: Handle context cancellation and timeouts in handlers and usecases
// TODO: Improve error handling and logging with more context (e.g. request IDs, user IDs)
func main() {
	cfg, cfgErr := config.Load(config.DefaultConfigPath)
	if cfgErr != nil {
		cfg = config.Default()
	}

	logger.Configure(cfg.Log.Level)
	log := logger.Get()

	if cfgErr != nil {
		log.WithError(cfgErr).Warn("failed to load config file; using defaults")
	}

	// Initialize database
	storage := db.NewStorageMgr()
	defer storage.Close()

	// Initialize usecases
	userRepo := user_repo.NewUserRepository(storage)
	campaignRepo := campaign_repo.NewCampaignRepository(storage)
	eventRepo := event_repo.NewEventRepository(storage)

	authUC := auth_usecase.NewAuthUseCase(userRepo)
	profileUC := profile_usecase.NewProfileUseCase(userRepo)
	campaignUC := campaign_usecase.NewCampaignUseCase(campaignRepo)
	eventUC := event_usecase.NewEventUseCase(eventRepo)
	analyticsUC := analytics_usecase.NewAnalyticsUseCase(campaignRepo, eventRepo)
	engagementUC := engagement_usecase.NewEngagementUseCase(eventRepo)

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
