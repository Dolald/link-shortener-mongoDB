package handler_http

import (
	"net/http"
	"shortener/pkg/service"
)

type Service interface {
}

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) (any, error) {

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) (any, error) {

}

func (h *Handler) GetFullURL(w http.ResponseWriter, r *http.Request) (any, error) {

}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) (any, error) {

}

type EndpointHandler func(w http.ResponseWriter, r *http.Request) (any, error)

func WrapEndpoint(h EndpointHandler) (any, error) {

}
