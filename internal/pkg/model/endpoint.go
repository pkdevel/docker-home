package model

import (
	"time"

	"github.com/pkdevel/docker-home/internal/pkg/persistence"
)

var endpoints *Endpoints

func GetEndpoints() *Endpoints {
	if endpoints == nil {
		endpoints = &Endpoints{persistence.NewDAO[*Endpoint]("endpoints")}
	}
	return endpoints
}

func (e *Endpoints) Find(query string) []*Endpoint {
	return e.dao.Find(query)
}

func (e *Endpoints) Save(item *Endpoint) error {
	return e.dao.Save(item)
}

type Endpoints struct {
	dao *persistence.DAO[*Endpoint]
}

type Endpoint struct {
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	Links     []string  `json:"links"`
}

func (e *Endpoint) Identifier() []byte {
	return []byte(e.ID)
}

func (e *Endpoint) UpdateTimestamp() {
	e.UpdatedAt = time.Now()
}
