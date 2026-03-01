package v1

import (
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/constant"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	campaign_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/campaigns/handler"
	"github.com/gorilla/mux"
)

func (r *Router) registerCampaignRoutes(v1 *mux.Router, campaignHandler *campaign_handler.CampaignHandler) {
	campaigns := v1.PathPrefix("/campaigns").Subrouter()
	// campaigns.HandleFunc("", campaignHandler.ListPublicCampaigns).Methods("GET")

	campaigns.HandleFunc("/preview/{id}", campaignHandler.GetCampaignPreviewHandler).Methods("GET")

	campaignsAuth := campaigns.PathPrefix("").Subrouter()
	campaignsAuth.Use(middleware.AuthMiddleware)
	campaignsAuth.HandleFunc("", campaignHandler.GetCampaignsHandler).Methods("GET")
	campaignsAuth.HandleFunc("/search", campaignHandler.SearchCampaignsHandler).Methods("GET")

	campaignsAuthRole := campaignsAuth.PathPrefix("").Subrouter()
	campaignsAuthRole.Use(middleware.RoleMiddleware(constant.RoleMarketer, constant.RoleAdmin))
	campaignsAuthRole.HandleFunc("", campaignHandler.CreateCampaignHandler).Methods("POST")
	campaignsAuthRole.HandleFunc("/{id}", campaignHandler.UpdateCampaignHandler).Methods("PUT")
	campaignsAuthRole.HandleFunc("/{id}", campaignHandler.DeleteCampaignHandler).Methods("DELETE")
	campaignsAuthRole.HandleFunc("/{id}/status", campaignHandler.PatchCampaignStatusHandler).Methods("PATCH")
	campaignsAuthRole.HandleFunc("/{id}/end", campaignHandler.EndCampaignHandler).Methods("PATCH")
}
