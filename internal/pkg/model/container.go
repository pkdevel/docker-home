package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/pkdevel/docker-home/internal/pkg/persistence"
	bolt "go.etcd.io/bbolt"
)

type Containers struct{ db *bolt.DB }

func GetContainers() Containers { return Containers{persistence.Database()} }
func (c *Containers) Close()    { c.db.Close() }

type Container struct {
	ID        string        `json:"id"`
	UpdatedAt time.Time     `json:"updated_at"`
	Data      ContainerData `json:"data"`
}

type ContainerData struct {
	Name        string `json:"name"`
	Port        uint16 `json:"port"`
	PrivatePort uint16 `json:"private_port"`
}

func (c *Containers) Find(query string) []Container {
	result, err := c.find(query)
	if err != nil {
		log.Panic(err)
	}
	return result
}

func (c *Containers) find(query string) ([]Container, error) {
	tx, err := c.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var result []Container
	err = tx.Bucket([]byte("containers")).ForEach(func(k, v []byte) error {
		if query != "" {
			if !strings.Contains(strings.ToLower(string(k)), strings.ToLower(query)) {
				return nil
			}
		}

		var container Container
		err := json.Unmarshal(v, &container)
		if err == nil {
			result = append(result, container)
		}
		return err
	})

	return result, err
}

func (c *Containers) Save(item Container) {
	if err := c.save(item); err != nil {
		log.Panic(err)
	}
}

func (c *Containers) save(item Container) error {
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte("containers"))
	if err != nil {
		return err
	}

	item.UpdatedAt = time.Now()
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	if err = b.Put([]byte(item.ID), data); err == nil {
		err = tx.Commit()
	}
	return err
}
