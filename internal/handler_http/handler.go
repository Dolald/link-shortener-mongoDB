package handler_http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"shortener/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	InitRoutes() *mux.Router
}

type handler struct {
	service service.Service
	logger  *logrus.Logger
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
		h.logger.Errorf("ReadAll failed: %w", err)
		return
	}

	// put URL json into our json struct
	var requestJson Request
	if err := json.Unmarshal(bodyData, &requestJson); err != nil {
		h.logger.Errorf("Unmarshal failed: %w", err)
		return
	}

	// check URL
	if _, err := url.ParseRequestURI(requestJson.Url); err != nil {
		h.logger.Errorf("ParseRequestURI failed: %w", err)
		return
	}

	shortURL, err := h.service.Shorten(r.Context(), requestJson.Url, requestJson.TtlDays)
	if err != nil {
		h.logger.Errorf("Shorten failed: %w", err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(shortURL)
}

func (h *handler) getFullURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl, ok := vars["shortUrl"]
	if !ok {
		h.logger.Error("missing shortUrl path parameter")
		return
	}

	fullUrl, err := h.service.GetFullURL(r.Context(), shortUrl)
	if err != nil {
		h.logger.Errorf("GetFullURL failed: %w", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	// redirect our client to a full URL with a short URL
	http.Redirect(w, r, fullUrl.URL, http.StatusFound)
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl, ok := vars["shortUrl"]
	if !ok {
		h.logger.Error("missing shortUrl path parameter")
		return
	}

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("ReadAll failed: %w", err)
		return
	}

	var URLRequest Request
	if err := json.Unmarshal(bodyData, &URLRequest); err != nil {
		h.logger.Errorf("Unmarshal failed: %w", err)
		return
	}

	if _, err := url.ParseRequestURI(URLRequest.Url); err != nil {
		h.logger.Errorf("ParseRequestURI failed: %w", err)
		return
	}

	err = h.service.Update(r.Context(), shortUrl, URLRequest.Url, URLRequest.TtlDays)
	if err != nil {
		h.logger.Errorf("Update failed: %w", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl, ok := vars["shortUrl"]
	if !ok {
		h.logger.Error("missing shortUrl path parameter")
		return
	}

	if err := h.service.Delete(r.Context(), shortUrl); err != nil {
		h.logger.Errorf("Delete failed: %w", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
