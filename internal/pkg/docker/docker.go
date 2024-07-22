package docker

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func List() []ContainerApp {
	apiClient, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d running container(s)", len(containers))

	result := []ContainerApp{}
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
