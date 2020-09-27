package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(storage NoteStorage, host string) *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	ctrl := newController(storage, host)
	r.HandleFunc("/", ctrl.Home).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/n/{note_id}", ctrl.GetNote).Methods(http.MethodGet)
	return r
}
