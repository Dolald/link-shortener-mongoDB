package handler_http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/todos", h.Shorten).Methods(http.MethodGet)
	return router
}
