package api

import (
	"fmt"
	"github.com/0ndi/go-privnote/internal/pkg/note"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	currentHost = "http://127.0.0.1:8080"
)

type controller struct {
	db *bolt.DB
}

func newController(db *bolt.DB) *controller {
	return &controller{db}
}

func (c *controller) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte(`Hello!`))
		return
	}
	r.ParseForm()
	text := r.PostFormValue("text")
	n := note.NewNote(text)
	noteID, err := n.Save(c.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("%s/n/%s", currentHost, noteID)))
}

func (c *controller) GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := note.GetNote(c.db, vars["note_id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(n.Data))
}
