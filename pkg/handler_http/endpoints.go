package handler_http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/shorten", wrapEndpoint(h.shorten)).Methods(http.MethodPost)
	router.HandleFunc("/{shortUrl}", wrapEndpoint(h.getFullURL)).Methods(http.MethodGet)
	router.HandleFunc("/update/{shortUrl}", wrapEndpoint(h.update)).Methods(http.MethodPut)
	router.HandleFunc("/{shortUrl}", wrapEndpoint(h.delete)).Methods(http.MethodDelete)
	return router
}
