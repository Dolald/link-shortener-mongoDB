package handler_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"shortener/pkg"
	"shortener/pkg/service"

	"github.com/gorilla/mux"
)

// type Service interface {
// 	Shorten(ctx context.Context, url string, ttlDays int) (*pkg.ShortURL, error)
// 	GetFullURL(ctx context.Context, shortURL string) (string, error)
// 	Update(ctx context.Context, id, url, string, ttl int) (*pkg.ShortURL, error)
// }

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

func (h *handler) shorten(w http.ResponseWriter, r *http.Request) (any, error) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// put URL json into our json struct
	var requestJson Request
	if err := json.Unmarshal(bodyData, &requestJson); err != nil {
		return nil, err
	}

	// check URL
	if _, err := url.ParseRequestURI(requestJson.Url); err != nil {
		return nil, err
	}

	return h.service.Shorten(r.Context(), requestJson.Url, requestJson.TtlDays)
}

func (h *handler) getFullURL(w http.ResponseWriter, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	shortUrl, ok := vars["shortUrl"]
	if !ok {
		return nil, fmt.Errorf("missing shortUrl path parameter")
	}

	fullUrl, err := h.service.GetFullURL(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil, err
	}

	// redirect our client to a full URL with a short URL
	http.Redirect(w, r, fullUrl, http.StatusFound)

	return nil, nil
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	shortUrl, ok := vars["shortUrl"]
	if !ok {
		return nil, fmt.Errorf("missing shortUrl path parameter")
	}

	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var URLRequest Request
	if err := json.Unmarshal(bodyData, &URLRequest); err != nil {
		return nil, err
	}

	if _, err := url.ParseRequestURI(URLRequest.Url); err != nil {
		return nil, err
	}

	return h.service.Update(r.Context(), shortUrl, URLRequest.Url, URLRequest.TtlDays)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) (any, error) {
	vars := mux.Vars(r)

	_, ok := vars["shortUrl"]
	if !ok {
		return nil, fmt.Errorf("missing shortUrl path parameter")
	}

	return nil, nil
}

// func (h *Handler) ping(w http.ResponseWriter, r *http.Request) (any, error) {

// }

type EndpointHandler func(w http.ResponseWriter, r *http.Request) (any, error)

func wrapEndpoint(e EndpointHandler) http.HandlerFunc {
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
		w.WriteHeader(http.StatusOK) // in my opinion it doesn't need ?
		w.Write(data)
	}
}
