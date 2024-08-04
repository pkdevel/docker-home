package task

import (
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/pkdevel/docker-home/internal/pkg/docker"
	"github.com/pkdevel/docker-home/internal/pkg/model"
)

const interval = time.Second * 60

type importer struct {
	docker     *docker.DockerClient
	containers *model.Containers
	endpoints  *model.Endpoints
}

func StartImporter() {
	slog.Info("Starting importer", "interval", interval)

	docker := docker.NewDockerClient()
	defer docker.Close()

	i := importer{
		docker:     docker,
		containers: model.GetContainers(),
		endpoints:  model.GetEndpoints(),
	}

	i.fetchAndSafe()
	go i.run()
}

func (i *importer) run() {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		i.fetchAndSafe()
	}
}

func (i *importer) fetchAndSafe() {
	for _, container := range i.docker.List() {
		data := []model.ContainerData{}
		hostname := hostname(container.Host)
		links := []string{}
		for _, port := range container.Ports {
			if link := generateLink(hostname, port); link != nil {
				links = append(links, *link)
			}
			data = append(data, model.ContainerData{
				Name:        container.Name,
				Port:        port.Port,
				PrivatePort: port.PrivatePort,
			})
		}
		slog.Debug("Importing", "container", container.Name, "id", container.ID[:7], "links", links)

		err := i.containers.Save(&model.Container{ID: container.Name, Data: data})
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if len(links) > 0 {
			err = i.endpoints.Save(&model.Endpoint{ID: container.Name, Links: links})
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}
}

func hostname(uri string) string {
	host, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return host.Hostname()
}

func generateLink(host string, port docker.ContainerPort) *string {
	if len(host) == 0 {
		return nil
	}
	if port.Type != "tcp" {
		return nil
	}
	if port.Port == 0 {
		return nil
	}
	if port.PrivatePort < 1024 && port.PrivatePort != 80 && port.PrivatePort != 443 {
		return nil
	}

	// test port for http and encryption
	// test for host network

	scheme := "http"
	if port.PrivatePort == 443 {
		scheme += "s"
	}
	result := fmt.Sprintf("%s://%s:%d", scheme, host, port.Port)
	return &result
}
