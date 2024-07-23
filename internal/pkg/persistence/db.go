package persistence

import (
	"io/fs"
	"log"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

var instance *bolt.DB

func Open() *bolt.DB {
	if instance == nil {
		err := os.MkdirAll("data", fs.ModeDir|fs.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		instance, err = bolt.Open("data/bolt.db", fs.ModePerm, &bolt.Options{Timeout: 5 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
	}
	return instance
}
