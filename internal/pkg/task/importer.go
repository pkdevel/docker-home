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
		containers := model.GetContainers()
		defer containers.Close()

		docker := docker.NewDockerClient()
		defer docker.Close()

		importContainers(docker, containers)
		for range ticker.C {
			importContainers(docker, containers)
		}
	}()
}

func importContainers(docker *docker.DockerClient, containers model.Containers) {
	for _, container := range docker.List() {
		containers.Save(model.Container{
			Name: container.Name,
			Data: model.ContainerData{Port: container.Port, PrivatePort: container.PrivatePort},
		})
		log.Printf("Importing %s (%s)", container.Name, container.ID)
	}
}
