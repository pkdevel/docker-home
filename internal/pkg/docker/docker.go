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

	result := make([]ContainerApp, len(containers))
	for i, ctr := range containers {
		name := strings.TrimPrefix(ctr.Names[0], "/")
		port := ctr.Ports[0].PublicPort
		result[i] = ContainerApp{name, port}
		log.Printf("Name: %s, Port: %v", name, port)
	}

	return result
}

type ContainerApp struct {
	Name string
	Port uint16
}
