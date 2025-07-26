package docker

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
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

func (c *DockerClient) Hostname() (string, error) {
	host, err := url.Parse(c.api.DaemonHost())
	if err != nil {
		return "", err
	}
	return host.Hostname(), nil
}

func (c *DockerClient) Events() (<-chan events.Message, <-chan error) {
	return c.api.Events(context.Background(), events.ListOptions{})
}

func (c *DockerClient) ContainerList() []container.Summary {
	result, err := c.api.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		slog.Error("unable to get containers", "error", err)
	} else {
		slog.Debug(fmt.Sprintf("found %d container(s)", len(result)))
	}
	return result
}
