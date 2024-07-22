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
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	log.Printf("Found %d containers", len(containers))

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
			strings.TrimPrefix(ctr.Names[0], "/"),
			ctr.Ports[0].PublicPort,
		}
		result = append(result, app)
		log.Printf("Name: %s, Port: %v", app.Name, app.Port)
	}
	return result
}

type ContainerApp struct {
	Name string
	Port uint16
}
