package v1

import (
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	engagement_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/engagement/handler"
	"github.com/gorilla/mux"
)

func (r *Router) registerEngagementRoutes(v1 *mux.Router, engagementHandler *engagement_handler.EngagementHandler) {
	users := v1.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware)
	users.HandleFunc("/{user_id}/engagement", engagementHandler.GetUserEngagementHandler).Methods("GET")
	users.HandleFunc("/{user_id}/campaigns/{campaign_id}/engagement", engagementHandler.GetUserCampaignEngagementHandler).Methods("GET")

	campaignFunnel := v1.PathPrefix("/campaigns").Subrouter()
	campaignFunnel.Use(middleware.AuthMiddleware)
	campaignFunnel.HandleFunc("/{campaign_id}/funnel", engagementHandler.GetCampaignFunnelHandler).Methods("GET")
}
