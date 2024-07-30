package docker

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	api *client.Client
}

func NewDockerClient() *DockerClient {
	api, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &DockerClient{api}
}

func (c *DockerClient) Close() {
	c.api.Close()
}

func (c *DockerClient) List() []ContainerApp {
	result := []ContainerApp{}

	containers, err := c.api.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		slog.Error(err.Error())
		return result
	}
	slog.Debug(fmt.Sprintf("Found %d running container(s)", len(containers)))

	for _, ctr := range containers {
		if len(ctr.Ports) == 0 {
			slog.Debug("No open ports", "container", name(ctr))
			continue
		}

		if ctr.Ports[0].Type != "tcp" {
			slog.Debug("No tcp port", "container", name(ctr))
			continue
		}
		if ctr.Ports[0].PublicPort == 0 {
			continue
		}

		app := ContainerApp{
			ctr.ID,
			name(ctr),
			ctr.Ports[0].PublicPort,
			ctr.Ports[0].PrivatePort,
		}
		result = append(result, app)
	}
	return result
}

func name(ctr types.Container) string {
	return strings.TrimPrefix(ctr.Names[0], "/")
}

type ContainerApp struct {
	ID          string
	Name        string
	Port        uint16
	PrivatePort uint16
}
