package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/persistence/db"
	analytics_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/usecase"
	auth_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/usecase"
	campaign_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/usecase"
	engagement_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/usecase"
	event_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/usecase"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/middleware"
	profile_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/usecase"
	v1 "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/router/v1"
)

func main() {
	// Initialize database
	storage := db.NewStorageMgr()
	defer storage.Close()

	// Initialize usecases
	authUC := auth_usecase.NewAuthUseCase(storage)
	profileUC := profile_usecase.NewProfileUseCase(storage)
	campaignUC := campaign_usecase.NewCampaignUseCase(storage)
	eventUC := event_usecase.NewEventUseCase(storage)
	analyticsUC := analytics_usecase.NewAnalyticsUseCase(storage)
	engagementUC := engagement_usecase.NewEngagementUseCase(storage)

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
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("API Documentation available at: http://localhost:8080/api/v1")

	if err := http.ListenAndServe(port, muxRouter); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
