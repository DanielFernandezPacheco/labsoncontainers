// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package labsoncontainers

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/docker/docker/client"
)

// StopEnvironment stops all containers of the provided lab environment.
func StopEnvironment(labName string) error {
	containers, err := GetEnvironmentContainers(labName)
	if err != nil {
		return fmt.Errorf("error while stopping environment: %w", err)
	}

	err = stopContainers(containers)
	if err != nil {
		return fmt.Errorf("error while stopping environment: %w", err)
	}

	return nil
}

// stopsContainers concurrently stops (using errgroups) all the specified containers.
func stopContainers(containers []LabContainer) error {
	g, ctx := errgroup.WithContext(context.Background())

	for _, container := range containers {
		id := container.ID
		g.Go(func() error {
			return StopContainer(ctx, id)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error while stopping containers: %w", err)
	}

	return nil
}

// StopContainer stops the specified container.
func StopContainer(errGroupCtx context.Context, containerID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return fmt.Errorf("error while destroying container %v: %v", containerID, err)
	}

	err = cli.ContainerStop(ctx, containerID, nil)

	if err != nil {
		return fmt.Errorf("error while destroying container %v: %v", containerID, err)
	}

	return nil
}
