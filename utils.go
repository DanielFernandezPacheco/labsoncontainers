package labsoncontainers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// GetEnviromentContainers gets all containers IDs (including non-running containers) that belong to the provided lab enviroment.
func GetEnviromentContainers(labName string) ([]string, error) {
	var containersIds []string

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error while getting enviroment containers: %v", err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet: true,
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", labName),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting enviroment containers: %v", err)
	}

	for _, container := range containers {
		containersIds = append(containersIds, container.ID)
	}

	return containersIds, nil
}

// GetEnviromentNetworks gets all networks IDs that belong to the provided lab enviroment.
func GetEnviromentNetworks(labName string) ([]string, error) {
	var networksIds []string

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error while getting enviroment containers: %v", err)
	}

	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", labName),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting enviroment containers: %v", err)
	}

	for _, network := range networks {
		networksIds = append(networksIds, network.ID)
	}

	return networksIds, nil
}

// boolPointer returns a bool pointer.
func boolPointer(b bool) *bool {
	return &b
}