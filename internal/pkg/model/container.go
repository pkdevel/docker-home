package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/pkdevel/docker-home/internal/pkg/persistence"
)

type Containers struct {
	db *sql.DB
}

func GetContainers() Containers {
	db := persistence.Open()
	query := `CREATE TABLE IF NOT EXISTS containers (
    id TEXT PRIMARY KEY UNIQUE NOT NULL,
    data BLOB NOT NULL,
    updated_at TEXT NOT NULL
  )`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return Containers{db}
}

func (c *Containers) Close() {
	c.db.Close()
}

type Container struct {
	updatedAt time.Time     `db:"updated_at"`
	Name      string        `db:"id"`
	Data      ContainerData `db:"data"`
}

type ContainerData struct {
	Port        uint16 `json:"port"`
	PrivatePort uint16 `json:"private_port"`
}

func (s *ContainerData) Scan(src any) error {
	switch t := src.(type) {
	case []byte:
		return json.Unmarshal(t, &s)
	default:
		return errors.New("invalid type")
	}
}

func (c *Containers) Find() []Container {
	var result []Container
	rows, err := c.db.Query("SELECT * FROM containers")
	if err != nil {
		log.Panic(err)
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var container Container

		if err := rows.Scan(
			&container.Name,
			&container.Data,
			&container.updatedAt,
		); err != nil {
			log.Panic(err)
			return result
		}

		result = append(result, container)
	}

	if err := rows.Err(); err != nil {
		log.Panic(err)
	}
	return result
}

func (c *Containers) Save(item Container) {
	data, err := json.Marshal(item.Data)
	if err != nil {
		log.Panic(err)
	}
	query := `INSERT OR REPLACE
    INTO containers (id, data, updated_at)
    VALUES (?, ?, ?)
  `
	item.updatedAt = time.Now()
	_, err = c.db.Exec(query, item.Name, data, item.updatedAt)
	if err != nil {
		log.Panic(err)
	}
}
