package handler_http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/shorten", h.shorten).Methods(http.MethodPost)
	// router.HandleFunc("/{shortUrl}", h.getFullURL).Methods(http.MethodGet)
	// router.HandleFunc("/update/{shortUrl}", h.update).Methods(http.MethodPut)
	// router.HandleFunc("/{shortUrl}", h.delete).Methods(http.MethodDelete)
	router.HandleFunc("/ping", h.ping).Methods(http.MethodGet)
	return router
}
