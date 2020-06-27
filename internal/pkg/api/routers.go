package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	ctrl := newController()
	r.HandleFunc("/", ctrl.Hello).Methods(http.MethodGet)
	return r
}
