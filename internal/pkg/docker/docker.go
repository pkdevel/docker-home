package docker

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	api *client.Client
}

func NewDockerClient() *DockerClient {
	api, err := client.NewClientWithOpts(
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
		log.Println(err)
		return result
	}
	log.Printf("Found %d running container(s)", len(containers))

	for _, ctr := range containers {
		if len(ctr.Ports) == 0 {
			continue
		}
		if ctr.Ports[0].Type != "tcp" {
			continue
		}
		if ctr.Ports[0].PublicPort == 0 {
			continue
		}

		app := ContainerApp{
			ctr.ID,
			strings.TrimPrefix(ctr.Names[0], "/"),
			ctr.Ports[0].PublicPort,
			ctr.Ports[0].PrivatePort,
		}
		result = append(result, app)
	}
	return result
}

type ContainerApp struct {
	ID          string
	Name        string
	Port        uint16
	PrivatePort uint16
}
