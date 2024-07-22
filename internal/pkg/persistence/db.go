package persistence

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

func Open() *sql.DB {
	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	dbName := "file:./data/local.db"
	db, err := sql.Open("libsql", dbName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
