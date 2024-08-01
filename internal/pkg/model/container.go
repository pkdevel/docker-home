package model

import (
	"log/slog"
	"time"

	"github.com/pkdevel/docker-home/internal/pkg/persistence"
)

var containers *Containers

func GetContainers() *Containers {
	if containers == nil {
		containers = &Containers{persistence.NewDAO[*Container]("containers")}
	}
	return containers
}

func (c *Containers) Find(query string) []*Container {
	return c.dao.Find(query)
}

func (c *Containers) Save(item *Container) {
	if err := c.dao.Save(item); err != nil {
		slog.Error(err.Error())
	}
}

type Containers struct {
	dao *persistence.DAO[*Container]
}

type Container struct {
	ID        string          `json:"id"`
	UpdatedAt time.Time       `json:"updated_at"`
	Data      []ContainerData `json:"data"`
}

type ContainerData struct {
	Name        string `json:"name"`
	Port        uint16 `json:"port"`
	PrivatePort uint16 `json:"private_port"`
}

func (c *Container) Identifier() []byte {
	return []byte(c.ID)
}

func (c *Container) UpdateTimestamp() {
	c.UpdatedAt = time.Now()
}
