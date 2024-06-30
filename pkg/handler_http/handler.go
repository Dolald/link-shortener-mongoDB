package handler_http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/pkg"
	"shortener/pkg/service"

	"github.com/gorilla/mux"
)

type Service interface {
	Shorten(ctx context.Context, url string, ttlDays int) (*pkg.ShortURL, error)
	GetFullURL(ctx context.Context, shortURL string) (string, error)
}

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type URLRequest struct {
	Url     string `json:"url"`
	TtlDays int    `json:"ttlDays"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) (any, error) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var requestJson URLRequest

	if err := json.Unmarshal(bodyData, &requestJson); err != nil {
		return nil, err
	}

	if _, err := url.ParseRequestURI(requestJson.Url); err != nil {
		return nil, err
	}

	return h.service.Shorten(r.Context(), requestJson.Url, requestJson.TtlDays)
}

// func (h *Handler) Update(w http.ResponseWriter, r *http.Request) (any, error) {

// }

// func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) (any, error) {

// }

func (h *Handler) GetFullURL(w http.ResponseWriter, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	val, ok := vars["shortUrl"]
	if !ok {
		return nil, fmt.Errorf("missing shortUrl path parameter")
	}
	return h.service.GetFullURL(r.Context(), val)
}

// func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) (any, error) {

// }

type EndpointHandler func(w http.ResponseWriter, r *http.Request) (any, error)

func WrapEndpoint(e EndpointHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := e(w, r)
		if err != nil {
			fmt.Println(err.Error())
			pkg.WriteErrorCode(w, err)
			return
		}

		data, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err.Error())
			pkg.WriteErrorCode(w, err)
			return
		}

		// define our content type
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
