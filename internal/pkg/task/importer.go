package task

import (
	"log"
	"time"

	"github.com/pkdevel/docker-home/internal/pkg/docker"
	"github.com/pkdevel/docker-home/internal/pkg/model"
)

func StartImporter() {
	log.Println("Starting importer")

	ticker := time.NewTicker(60 * time.Second)
	go func() {
		m := model.GetContainers()
		defer m.Close()

		importContainers(m)

		for range ticker.C {
			importContainers(m)
		}
	}()
}

func importContainers(m model.Containers) {
	containers := docker.List()
	for _, container := range containers {
		m.Save(model.Container{
			Name: container.Name,
			Data: model.ContainerData{Port: container.Port, PrivatePort: container.PrivatePort},
		})
		log.Printf("Importing %s (%s)", container.Name, container.ID)
	}
}
