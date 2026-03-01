package v1

import (
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	analytics_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/analytics/handler"
	"github.com/gorilla/mux"
)

func (r *Router) registerAnalyticsRoutes(v1 *mux.Router, analyticsHandler *analytics_handler.AnalyticsHandler) {
	analytics := v1.PathPrefix("/analytics").Subrouter()

	analytics.HandleFunc("/public/stats", analyticsHandler.GetPublicStatsHandler).Methods("GET")

	analyticsAuth := analytics.PathPrefix("").Subrouter()
	analyticsAuth.Use(middleware.AuthMiddleware)
	analyticsAuth.HandleFunc("/reports", analyticsHandler.GetAnalyticsReportHandler).Methods("GET")
	// TODO: these 3 APIS can be optimized to a single API with a query parameter for the time range (daily, weekly, monthly)
	analyticsAuth.HandleFunc("/reports/daily", analyticsHandler.GetDailyReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/weekly", analyticsHandler.GetWeeklyReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/reports/monthly", analyticsHandler.GetMonthlyReportHandler).Methods("GET")
	analyticsAuth.HandleFunc("/campaigns/{campaign_id}/stats", analyticsHandler.GetCampaignStatsHandler).Methods("GET")
}
