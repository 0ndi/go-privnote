package api

import (
	"github.com/boltdb/bolt"
	"net/http"
)

type controller struct {
	db *bolt.DB
}

func newController(db *bolt.DB) *controller {
	return &controller{db}
}

func (c *controller) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Hello!`))
}
