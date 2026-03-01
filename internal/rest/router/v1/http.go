package v1

import (
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

	r.registerCoreRoutes(v1, healthHandler)
	r.registerAuthRoutes(v1, authHandler)
	r.registerProfileRoutes(v1, profileHandler)
	r.registerCampaignRoutes(v1, campaignHandler)
	r.registerEventRoutes(v1, eventHandler)
	r.registerAnalyticsRoutes(v1, analyticsHandler)
	r.registerEngagementRoutes(v1, engagementHandler)

	return r.router
}

// GetRouter returns the mux router
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
