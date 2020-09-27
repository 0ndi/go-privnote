package api

import (
	"fmt"
	"github.com/0ndi/go-privnote/internal/pkg/note"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type HomeData struct {
	Url   string
	Error string
	Data  string
}

type GetNoteData struct {
	Text  string
	Error string
}

type NoteStorage interface {
	Save(n *note.Note) (string, error)
	GetNote(key string) (*note.Note, error)
}

type controller struct {
	host    string
	storage NoteStorage
}

func newController(storage NoteStorage, host string) *controller {
	return &controller{host, storage}
}

func (c *controller) Home(w http.ResponseWriter, r *http.Request) {
	var data HomeData
	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.PostFormValue("text")
		if len(text) > 0 {
			n := note.NewNote(text)
			noteID, err := c.storage.Save(n)
			if err != nil {
				log.Error(err)
				data.Error = err.Error()
			} else {
				data.Url = fmt.Sprintf("%s/n/%s", c.host, noteID)
			}
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
	n, err := c.storage.GetNote(vars["note_id"])
	if err != nil {
		log.Error(err)
	}

	if n.Data == "" {
		c.Home(w, r)
		return
	}
	t, err := template.New("index.html").ParseFiles("web/templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	if err := t.Execute(w, HomeData{Data: n.Data}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
