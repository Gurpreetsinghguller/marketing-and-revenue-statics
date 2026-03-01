package v1

import (
	auth_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/auth/handler"
	"github.com/gorilla/mux"
)

func (r *Router) registerAuthRoutes(v1 *mux.Router, authHandler *auth_handler.AuthHandler) {
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	auth.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
}
