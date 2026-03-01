package v1

import (
	"net/http"

	health_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/health/handler"
	"github.com/gorilla/mux"
)

func (r *Router) registerCoreRoutes(v1 *mux.Router, healthHandler *health_handler.HealthHandler) {
	v1.HandleFunc("/health", healthHandler.GetHealthHandler).Methods("GET")
	v1.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi.yaml")
	}).Methods("GET")
}
