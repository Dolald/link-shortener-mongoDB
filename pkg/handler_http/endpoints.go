package handler_http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/shorten", WrapEndpoint(h.Shorten)).Methods(http.MethodPost)
	router.HandleFunc("/:shortUrl", WrapEndpoint(h.GetFullURL)).Methods(http.MethodGet)
	return router
}
