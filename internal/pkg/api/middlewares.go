package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	livnessProbeUrl = "/ping"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == livnessProbeUrl {
			next.ServeHTTP(w, r)
		} else {
			t := time.Now()
			next.ServeHTTP(w, r)
			log.WithFields(log.Fields{
				"spent":     time.Since(t).Seconds(),
				"url":       r.RequestURI,
			}).Info()
		}
	})
}
