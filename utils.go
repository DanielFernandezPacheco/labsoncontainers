// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package labsoncontainers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// GetEnviromentContainers gets all containers info (including non-running containers) that belong to the provided lab enviroment.
func GetEnviromentContainers(labName string) ([]LabContainer, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error while getting enviroment containers: %v", err)
	}

	dockerContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", labName),
		),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting enviroment containers: %v", err)
	}

	labContainers := make([]LabContainer, len(dockerContainers))

	i := 0
	for _, container := range dockerContainers {
		labContainers[i].Name = container.Names[0][1:] // Leading '/' must be removed
		labContainers[i].Image = container.Image
		labContainers[i].ID = container.ID

		labContainers[i].Background, err = strconv.ParseBool(container.Labels["background"])
		if err != nil {
			return nil, fmt.Errorf("error while getting enviroment containers: %v", err)
		}

		labContainers[i].Networks = make([]LabNetwork, len(container.NetworkSettings.Networks))
		j := 0
		for name, network := range container.NetworkSettings.Networks {
			labContainers[i].Networks[j].Name = name
			labContainers[i].Networks[j].IP = network.IPAddress
			j++
		}
		i++
	}
	return labContainers, nil
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