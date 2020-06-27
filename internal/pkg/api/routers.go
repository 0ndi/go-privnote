package api

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(db *bolt.DB) *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	ctrl := newController(db)
	r.HandleFunc("/", ctrl.Hello).Methods(http.MethodGet)
	return r
}
