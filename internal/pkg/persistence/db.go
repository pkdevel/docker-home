package persistence

import (
	"io/fs"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	version int = -1
)

var instance *bolt.DB

func database() *bolt.DB {
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

	slog.Info("opening database")
	instance, err = bolt.Open("data/bolt.db", fs.ModePerm, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	migrate(instance)
}

func Close() {
	err := instance.Close()
	if err != nil {
		log.Fatal(err)
	}
	instance = nil
}

func migrate(db *bolt.DB) {
	tx, err := db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	dbversion := 0
	if system := tx.Bucket([]byte("system")); system != nil {
		if version := system.Get([]byte("version")); version != nil {
			dbversion, _ = strconv.Atoi(string(version))
		}
	}

	if dbversion < 0 { // TODO: check for dev env
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			if err := tx.DeleteBucket(name); err == nil {
				slog.Info("dropped", "bucket", string(name))
			}
			return nil
		})
	}
	if dbversion < version {
		slog.Warn("TODO: migrating database")
	}

	system, err := tx.CreateBucketIfNotExists([]byte("system"))
	if err != nil {
		log.Fatal(err)
	}
	err = system.Put([]byte("version"), []byte(strconv.Itoa(version)))
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
