package handler

import "net/http"

type MHandler interface {
	GetUrlPattern() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
