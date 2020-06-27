package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func NewServer(addr string, port, readTimeout, writeTimeout int, router *mux.Router) *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", addr, port),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		Handler:      router,
	}
	return server
}
