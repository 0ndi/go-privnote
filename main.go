package main

import (
	"context"
	"fmt"
	"github.com/0ndi/go-privnote/internal/pkg/api"
	"github.com/0ndi/go-privnote/internal/pkg/config"
	"github.com/0ndi/go-privnote/internal/pkg/storage/file_storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

const (
	dbPath = "note.db"
)

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	db, err := file_storage.GetDB(dbPath)
	if err != nil {
		err := errors.Wrap(err, "GetDB")
		panic(err)
	}
	defer db.Close()

	storage, err := file_storage.NewStorage(db)
	if err != nil {
		err := errors.Wrap(err, "NewStorage")
		panic(err)
	}
	router := api.NewRouter(storage, fmt.Sprintf("%s:%d/", config.Conf.Host, config.Conf.Port))
	server := api.NewServer(config.Conf.Addr, config.Conf.Port, config.Conf.WriteTimeout, config.Conf.ReadTimeout, router)

	log.Infof("Start server at %s:%d", config.Conf.Addr, config.Conf.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			err := errors.Wrap(err, "ListenAndServe")
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Conf.WaitUntilShutdown)*time.Second)
	defer cancel()

	server.Shutdown(ctx)

	log.Info("Server stopped")
	os.Exit(0)

}
