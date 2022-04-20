// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package labsoncontainers

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// StartEnviroment starts all the containers of the provided lab enviroment.
func StartEnviroment(labName string) error {
	containers, err := GetEnviromentContainers(labName)
	if err != nil {
		return fmt.Errorf("error while starting enviroment: %w", err)
	}

	err = startContainers(containers)
	if err != nil {
		return fmt.Errorf("error while starting enviroment: %w", err)
	}

	return nil
}

// stopsContainers concurrently starts (using errgroups) all the specified containers.
func startContainers(containers []LabContainer) error {
	g, ctx := errgroup.WithContext(context.Background())

	for _, container := range containers {
		id := container.ID
		g.Go(func() error {
			return startContainer(ctx, id)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error while starting containers: %w", err)
	}

	return nil
}

// startContainer starts the specified container.
func startContainer(errGroupCtx context.Context, containerID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error while starting container %v: %v", containerID, err)
	}

	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("error while starting container %v: %v", containerID, err)
	}

	return nil
}