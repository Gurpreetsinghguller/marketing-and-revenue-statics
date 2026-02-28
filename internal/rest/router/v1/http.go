package v1

import (
	"net/http"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	analytics_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/handler"
	analytics_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/usecase"
	auth_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/handler"
	auth_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/usecase"
	campaign_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/handler"
	campaign_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/usecase"
	engagement_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/handler"
	engagement_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/usecase"
	event_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/handler"
	event_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/usecase"
	health_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/health/handler"
	profile_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/handler"
	profile_usecase "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/usecase"

	"github.com/gorilla/mux"
)

// Router represents the HTTP router
type Router struct {
	router *mux.Router
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

// InitHTTPRoutes initializes all HTTP routes with handler instances
func (r *Router) InitHTTPRoutes(
	authUC *auth_usecase.AuthUseCase,
	profileUC *profile_usecase.ProfileUseCase,
	campaignUC *campaign_usecase.CampaignUseCase,
	eventUC *event_usecase.EventUseCase,
	analyticsUC *analytics_usecase.AnalyticsUseCase,
	engagementUC *engagement_usecase.EngagementUseCase,
) *mux.Router {
	// Initialize handlers with their usecases
	authHandler := auth_handler.NewAuthHandler(authUC)
	profileHandler := profile_handler.NewProfileHandler(profileUC)
	campaignHandler := campaign_handler.NewCampaignHandler(campaignUC)
	eventHandler := event_handler.NewEventHandler(eventUC)
	analyticsHandler := analytics_handler.NewAnalyticsHandler(analyticsUC)
	engagementHandler := engagement_handler.NewEngagementHandler(engagementUC)
	healthHandler := health_handler.NewHealthHandler()

	v1 := r.router.PathPrefix("/api/v1").Subrouter()

	// ============ Health Route ============
	v1.HandleFunc("/health", healthHandler.GetHealthHandler).Methods("GET")

	// ============ OpenAPI Spec Route ============
	v1.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi.yaml")
	}).Methods("GET")

	// ============ Auth Routes ============
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	auth.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")

	// ============ Profile Routes ============
	profile := v1.PathPrefix("/profile").Subrouter()
	profile.Use(middleware.AuthMiddleware)
	profile.HandleFunc("", profileHandler.GetProfileHandler).Methods("GET")
	profile.HandleFunc("", profileHandler.UpdateProfileHandler).Methods("PUT")

	// ============ Campaign Routes ============
	campaigns := v1.PathPrefix("/campaigns").Subrouter()

	// Public campaign preview (no auth required)
	campaigns.HandleFunc("/preview/{id}", campaignHandler.GetCampaignPreviewHandler).Methods("GET")

	// Authenticated campaign routes
	campaignsAuth := campaigns.PathPrefix("").Subrouter()
	campaignsAuth.Use(middleware.AuthMiddleware)
	campaignsAuth.HandleFunc("", campaignHandler.GetCampaignsHandler).Methods("GET")

	// Routes requiring Marketer or Admin role
	campaignsAuthRole := campaignsAuth.PathPrefix("").Subrouter()
	campaignsAuthRole.Use(middleware.RoleMiddleware("Marketer", "Admin"))
	campaignsAuthRole.HandleFunc("", campaignHandler.CreateCampaignHandler).Methods("POST")
	campaignsAuthRole.HandleFunc("/{id}", campaignHandler.UpdateCampaignHandler).Methods("PUT")
	campaignsAuthRole.HandleFunc("/{id}", campaignHandler.DeleteCampaignHandler).Methods("DELETE")

	campaignsAuth.HandleFunc("/search", campaignHandler.SearchCampaignsHandler).Methods("GET")

	// ============ Event Tracking Routes ============
	events := v1.PathPrefix("/events").Subrouter()
	events.Use(middleware.RateLimitMiddleware)
	events.HandleFunc("", eventHandler.TrackEventHandler).Methods("POST")

	// Event logs (authenticated)
	eventsAuth := events.PathPrefix("").Subrouter()
	eventsAuth.Use(middleware.AuthMiddleware)
	eventsAuth.HandleFunc("", eventHandler.GetEventsHandler).Methods("GET")

	// ============ Analytics & Reporting Routes ============
	analytics := v1.PathPrefix("/analytics").Subrouter()

	// Public stats (no auth required)
	analytics.HandleFunc("/public/stats", analyticsHandler.GetPublicStatsHandler).Methods("GET")

	// Authenticated analytics routes
	analyticsAuth := analytics.PathPrefix("").Subrouter()
	analyticsAuth.Use(middleware.AuthMiddleware)
	analyticsAuth.HandleFunc("/reports", analyticsHandler.GetAnalyticsReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/daily", analyticsHandler.GetDailyReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/weekly", analyticsHandler.GetWeeklyReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/monthly", analyticsHandler.GetMonthlyReportHandler).Methods("GET")

	// ============ User Engagement & Behavioral Data Routes ============
	users := v1.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware)
	users.HandleFunc("/{user_id}/engagement", engagementHandler.GetUserEngagementHandler).Methods("GET")
	users.HandleFunc("/{user_id}/campaigns/{campaign_id}/engagement", engagementHandler.GetUserCampaignEngagementHandler).Methods("GET")

	// ============ Campaign Funnel Routes ============
	campaignFunnel := v1.PathPrefix("/campaigns").Subrouter()
	campaignFunnel.Use(middleware.AuthMiddleware)
	campaignFunnel.HandleFunc("/{campaign_id}/funnel", engagementHandler.GetCampaignFunnelHandler).Methods("GET")

	return r.router
}

// GetRouter returns the mux router
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
