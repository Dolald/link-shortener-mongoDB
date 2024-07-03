package pkg

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrCollision  = errors.New("failed to create short link due to collision")
)

func WriteErrorCode(rw http.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		rw.WriteHeader(http.StatusNotFound)
	case ErrBadRequest:
		rw.WriteHeader(http.StatusBadRequest)
	default:
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
