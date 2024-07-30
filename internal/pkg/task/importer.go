package task

import (
	"log/slog"
	"time"

	"github.com/pkdevel/docker-home/internal/pkg/docker"
	"github.com/pkdevel/docker-home/internal/pkg/model"
)

func StartImporter() {
	slog.Info("Starting importer")

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
		go func() {
			containers.Save(model.Container{
				ID: container.Name,
				Data: model.ContainerData{
					Name:        container.Name,
					Port:        container.Port,
					PrivatePort: container.PrivatePort,
				},
			})
			slog.Debug("Imported", container.Name, container.ID[:7])
		}()
	}
}
