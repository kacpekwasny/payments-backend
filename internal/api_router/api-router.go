package api_router

import (
	"github.com/gorilla/mux"
	"github.com/kacpekwasny/payments-backend/internal/wrappers"
)

func NewRouter() (*mux.Router, *mux.Router) {
	// logger takes all traffic loggs it, passes it to rtr and the rtr then responds

	rtr := mux.NewRouter().StrictSlash(true)
	// rtr.HandleFunc("/login", handleGetLogin).Methods("GET")

	// All requests are first handled by logger which then relays them to rtr.
	// logger loggsdata from http.Request
	logAndRelay := wrappers.LogRequests(rtr.ServeHTTP)
	logger := mux.NewRouter()
	logger.PathPrefix("/").HandlerFunc(logAndRelay)
	return rtr, logger
}
