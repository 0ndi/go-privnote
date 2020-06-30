package api

import (
	"fmt"
	"github.com/0ndi/go-privnote/internal/pkg/note"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

const (
	currentHost = "http://127.0.0.1:8080"
)

type HomeData struct {
	Url   string
	Error string
}

type GetNoteData struct {
	Text  string
	Error string
}

type controller struct {
	db *bolt.DB
}

func newController(db *bolt.DB) *controller {
	return &controller{db}
}

func (c *controller) Home(w http.ResponseWriter, r *http.Request) {
	var data HomeData
	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.PostFormValue("text")

		n := note.NewNote(text)
		noteID, err := n.Save(c.db)
		if err != nil {
			log.Error(err)
			data.Error = err.Error()
		} else {
			data.Url = fmt.Sprintf("%s/n/%s", currentHost, noteID)
		}
	}
	t, err := template.New("index.html").ParseFiles("web/templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if err := t.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (c *controller) GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := note.GetNote(c.db, vars["note_id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Error(err)
		return
	}

	if n.Data == "" {
		c.Home(w, r)
		return
	}
	w.Write([]byte(n.Data))
}
