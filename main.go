package main

import (
	"github.com/0ndi/go-privnote/internal/pkg/api"
	"github.com/0ndi/go-privnote/internal/pkg/config"
	"github.com/0ndi/go-privnote/internal/pkg/note"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

const (
	dbPath = "note.db"
)

func main() {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(note.BucketName))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	router := api.NewRouter(db)
	server := api.NewServer(config.Conf.Addr, config.Conf.Port, config.Conf.WriteTimeout, config.Conf.ReadTimeout, router)

	log.Info("Start server..")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
	log.Info("Server stopped")
}
