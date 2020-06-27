package main

import (
	"github.com/0ndi/go-privnote/internal/pkg/api"
	"github.com/0ndi/go-privnote/internal/pkg/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := api.NewRouter()
	server := api.NewServer(config.Conf.Addr, config.Conf.Port, config.Conf.WriteTimeout, config.Conf.ReadTimeout, router)

	log.Info("Start server..")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
	log.Info("Server stopped")
}
