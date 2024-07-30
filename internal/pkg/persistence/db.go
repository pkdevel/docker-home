package persistence

import (
	"io/fs"
	"log"
	"log/slog"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

var instance *bolt.DB

func Database() *bolt.DB {
	Init()
	return instance
}

func Init() {
	if instance != nil {
		return
	}

	err := os.MkdirAll("data", fs.ModeDir|fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Opening database")
	instance, err = bolt.Open("data/bolt.db", fs.ModePerm, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
}
