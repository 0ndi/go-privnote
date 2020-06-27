package api

import (
	"net/http"
)

type controller struct {}

func newController() *controller {
	return &controller{}
}

func (c *controller) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Hello!`))
}
