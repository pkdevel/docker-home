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
	slog.Debug(fmt.Sprintf("Found %d container(s)", len(containers)))

	for _, ctr := range containers {
		app := ContainerApp{
			ctr.ID,
			name(ctr),
			c.api.DaemonHost(),
			[]ContainerPort{},
		}

		for _, port := range ctr.Ports {
			if port.PublicPort == 0 {
				continue
			}

			app.Ports = append(app.Ports, ContainerPort{
				port.Type,
				port.PublicPort,
				port.PrivatePort,
			})
		}
		result = append(result, app)
	}
	return result
}

func name(ctr types.Container) string {
	return strings.TrimPrefix(ctr.Names[0], "/")
}

type ContainerApp struct {
	ID    string
	Name  string
	Host  string
	Ports []ContainerPort
}

type ContainerPort struct {
	Type        string
	Port        uint16
	PrivatePort uint16
}
