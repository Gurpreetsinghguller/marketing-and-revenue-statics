package v1

import (
	analytics_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/handler"
	auth_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/handler"
	campaign_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/handler"
	engagement_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/handler"
	event_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/handler"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/middleware"
	profile_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/handler"

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

// InitHTTPRoutes initializes all HTTP routes
func (r *Router) InitHTTPRoutes() *mux.Router {
	v1 := r.router.PathPrefix("/api/v1").Subrouter()

	// ============ Auth Routes ============
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", auth_handler.RegisterHandler).Methods("POST")
	auth.HandleFunc("/login", auth_handler.LoginHandler).Methods("POST")

	// ============ Profile Routes ============
	profile := v1.PathPrefix("/profile").Subrouter()
	profile.Use(middleware.AuthMiddleware)
	profile.HandleFunc("", profile_handler.GetProfileHandler).Methods("GET")
	profile.HandleFunc("", profile_handler.UpdateProfileHandler).Methods("PUT")

	// ============ Campaign Routes ============
	campaigns := v1.PathPrefix("/campaigns").Subrouter()

	// Public campaign preview (no auth required)
	campaigns.HandleFunc("/preview/{id}", campaign_handler.GetCampaignPreviewHandler).Methods("GET")

	// Authenticated campaign routes
	campaignsAuth := campaigns.PathPrefix("").Subrouter()
	campaignsAuth.Use(middleware.AuthMiddleware)
	campaignsAuth.HandleFunc("", campaign_handler.GetCampaignsHandler).Methods("GET")

	// Routes requiring Marketer or Admin role
	campaignsAuthRole := campaignsAuth.PathPrefix("").Subrouter()
	campaignsAuthRole.Use(middleware.RoleMiddleware("Marketer", "Admin"))
	campaignsAuthRole.HandleFunc("", campaign_handler.CreateCampaignHandler).Methods("POST")
	campaignsAuthRole.HandleFunc("/{id}", campaign_handler.UpdateCampaignHandler).Methods("PUT")
	campaignsAuthRole.HandleFunc("/{id}", campaign_handler.DeleteCampaignHandler).Methods("DELETE")

	campaignsAuth.HandleFunc("/search", campaign_handler.SearchCampaignsHandler).Methods("GET")

	// ============ Event Tracking Routes ============
	events := v1.PathPrefix("/events").Subrouter()
	events.Use(middleware.RateLimitMiddleware)
	events.HandleFunc("", event_handler.TrackEventHandler).Methods("POST")

	// Event logs (authenticated)
	eventsAuth := events.PathPrefix("").Subrouter()
	eventsAuth.Use(middleware.AuthMiddleware)
	eventsAuth.HandleFunc("", event_handler.GetEventsHandler).Methods("GET")

	// ============ Analytics & Reporting Routes ============
	analytics := v1.PathPrefix("/analytics").Subrouter()

	// Public stats (no auth required)
	analytics.HandleFunc("/public/stats", analytics_handler.GetPublicStatsHandler).Methods("GET")

	// Authenticated analytics routes
	analyticsAuth := analytics.PathPrefix("").Subrouter()
	analyticsAuth.Use(middleware.AuthMiddleware)
	analyticsAuth.HandleFunc("/reports", analytics_handler.GetAnalyticsReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/daily", analytics_handler.GetDailyReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/weekly", analytics_handler.GetWeeklyReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/monthly", analytics_handler.GetMonthlyReportHandler).Methods("GET")

	// ============ User Engagement & Behavioral Data Routes ============
	users := v1.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware)
	users.HandleFunc("/{user_id}/engagement", engagement_handler.GetUserEngagementHandler).Methods("GET")
	users.HandleFunc("/{user_id}/campaigns/{campaign_id}/engagement", engagement_handler.GetUserCampaignEngagementHandler).Methods("GET")

	// ============ Campaign Funnel Routes ============
	campaignFunnel := v1.PathPrefix("/campaigns").Subrouter()
	campaignFunnel.Use(middleware.AuthMiddleware)
	campaignFunnel.HandleFunc("/{campaign_id}/funnel", engagement_handler.GetCampaignFunnelHandler).Methods("GET")

	return r.router
}

// GetRouter returns the mux router
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
