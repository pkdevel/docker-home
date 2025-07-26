package task

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/pkdevel/docker-home/internal/pkg/docker"
	"github.com/pkdevel/docker-home/internal/pkg/model"
)

const (
	backoff = 3
)

type importer struct {
	docker     *docker.DockerClient
	containers *model.Containers
	endpoints  *model.Endpoints
}

func StartImporter() {
	slog.Info("starting importer")

	docker := docker.NewDockerClient()
	defer docker.Close()

	i := importer{
		docker:     docker,
		containers: model.GetContainers(),
		endpoints:  model.GetEndpoints(),
	}

	interval := time.Second * 1
	i.fetchAndSafe()
	go func() {
		for {
			err := i.run()
			if err != nil {
				if interval > time.Second*90 {
					slog.Error("unable to receive events, giving up...")
					return
				}
				slog.Error("unable to receive events, retrying...", "in", interval, "error", err)
				ticker := time.NewTimer(interval)
				for range ticker.C {
					interval = interval * backoff
					break
				}
			}
		}
	}()
}

func (i *importer) run() error {
	events, err := i.docker.Events()
	for {
		select {
		case event := <-events:
			handleEvent(event)
		case err := <-err:
			return err
		}
	}
}

func handleEvent(event events.Message) {
	switch event.Type {
	case "container":
		switch event.Action {
		case events.ActionStart:
			slog.Debug("container started", "name", event.Actor.Attributes["name"])
		case events.ActionDie:
			slog.Debug("container died", "name", event.Actor.Attributes["name"])
		}
	}
}

func (i *importer) fetchAndSafe() {
	hostname, err := i.docker.Hostname()
	if err != nil {
		slog.Error("unable to get hostname", "error", err)
		return
	}
	if len(hostname) == 0 {
		slog.Error("hostname not specified, using localhost")
		hostname = "localhost"
	}

	for _, container := range i.docker.ContainerList() {
		name := strings.TrimPrefix(container.Names[0], "/")
		data := []model.ContainerData{}
		links := []string{}
		for _, port := range container.Ports {
			if link := generateLink(hostname, port); link != nil {
				links = append(links, *link)
			}
			data = append(data, model.ContainerData{
				ID:          container.ID,
				Port:        port.PublicPort,
				PrivatePort: port.PrivatePort,
			})
		}
		slog.Debug("Importing", "container", name, "id", container.ID[:7], "links", links)

		err := i.containers.Save(&model.Container{Name: name, Data: data})
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if len(links) > 0 {
			err = i.endpoints.Save(&model.Endpoint{ID: name, Links: links})
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}
}

func generateLink(host string, port container.Port) *string {
	if port.Type != "tcp" {
		return nil
	}
	if port.PublicPort == 0 {
		return nil
	}
	if port.PrivatePort < 1024 && port.PrivatePort != 80 && port.PrivatePort != 443 {
		return nil
	}

	// TODO: test port for http and encryption
	// TODO: test for host network

	scheme := "http"
	if port.PrivatePort == 443 {
		scheme += "s"
	}
	result := fmt.Sprintf("%s://%s:%d", scheme, host, port.PublicPort)
	return &result
}
