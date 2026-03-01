package v1

import (
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	profile_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/profile/handler"
	"github.com/gorilla/mux"
)

func (r *Router) registerProfileRoutes(v1 *mux.Router, profileHandler *profile_handler.ProfileHandler) {
	profile := v1.PathPrefix("/profile").Subrouter()
	profile.Use(middleware.AuthMiddleware)
	profile.HandleFunc("", profileHandler.GetProfileHandler).Methods("GET")
	profile.HandleFunc("", profileHandler.UpdateProfileHandler).Methods("PUT")
}
