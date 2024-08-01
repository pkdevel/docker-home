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

	slog.Info("Opening database")
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

	bucket, err := tx.CreateBucketIfNotExists([]byte("system"))
	if err != nil {
		log.Fatal(err)
	}

	dbversion := bucket.Get([]byte("version"))
	if dbversion != nil {
		current, err := strconv.Atoi(string(dbversion))
		if err != nil {
			log.Fatal(err)
		}
		if current < 0 { // TODO: check for dev env
			drop(tx, "containers")
			drop(tx, "endpoints")
		}
		if current < version {
			slog.Info("TODO: Migrating database")
		}
	}

	err = bucket.Put([]byte("version"), []byte(strconv.Itoa(version)))
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func drop(tx *bolt.Tx, bucket string) {
	if err := tx.DeleteBucket([]byte(bucket)); err == nil {
		slog.Info("Dropped", "bucket", bucket)
	}
}
