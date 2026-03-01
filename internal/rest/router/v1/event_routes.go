package v1

import (
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/middleware"
	event_handler "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/rest/event/handler"
	"github.com/gorilla/mux"
)

func (r *Router) registerEventRoutes(v1 *mux.Router, eventHandler *event_handler.EventHandler) {
	events := v1.PathPrefix("/events").Subrouter()
	events.Use(middleware.RateLimitMiddleware)
	events.HandleFunc("", eventHandler.TrackEventHandler).Methods("POST")

	eventsAuth := events.PathPrefix("").Subrouter()
	eventsAuth.Use(middleware.AuthMiddleware)
	eventsAuth.HandleFunc("", eventHandler.GetEventsHandler).Methods("GET")
}
