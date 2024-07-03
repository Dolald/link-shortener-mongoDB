package handler_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/internal/service"

	"github.com/gorilla/mux"
)

type Handler interface {
	InitRoutes() *mux.Router
}

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) Handler {
	return &handler{
		service: service,
	}
}

type Request struct {
	Url     string `json:"url"`
	TtlDays int    `json:"ttlDays"`
}

func (h *handler) shorten(w http.ResponseWriter, r *http.Request) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("shorten/ReadAll: %w", err)
		return
	}

	// put URL json into our json struct
	var requestJson Request
	if err := json.Unmarshal(bodyData, &requestJson); err != nil {
		fmt.Errorf("shorten/Unmarshal: %w", err)
		return
	}

	// check URL
	if _, err := url.ParseRequestURI(requestJson.Url); err != nil {
		fmt.Errorf("shorten/ParseRequestURI: %w", err)
		return
	}

	h.service.Shorten(r.Context(), requestJson.Url, requestJson.TtlDays)
}

func (h *handler) getFullURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl, ok := vars["shortUrl"]
	if !ok {
		fmt.Errorf("missing shortUrl path parameter")
		return
	}

	fullUrl, err := h.service.GetFullURL(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// redirect our client to a full URL with a short URL
	http.Redirect(w, r, fullUrl, http.StatusFound)
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl, ok := vars["shortUrl"]
	if !ok {
		fmt.Errorf("missing shortUrl path parameter")
		return
	}

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("shorten/ReadAll: %w", err)
		return
	}

	var URLRequest Request
	if err := json.Unmarshal(bodyData, &URLRequest); err != nil {
		fmt.Errorf("shorten/Unmarshal: %w", err)
		return
	}

	if _, err := url.ParseRequestURI(URLRequest.Url); err != nil {
		fmt.Errorf("shorten/ParseRequestURI: %w", err)
		return
	}

	_, err = h.service.Update(r.Context(), shortUrl, URLRequest.Url, URLRequest.TtlDays)
	if err != nil {
		fmt.Errorf("shorten/Update: %w", err)
		return
	}
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	_, ok := vars["shortUrl"]
	if !ok {
		fmt.Errorf("missing shortUrl path parameter")
		return
	}

}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
