// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package labsoncontainers

import (
	"context"
	"fmt"
	"encoding/json"

	"golang.org/x/sync/errgroup"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// InspectEnvironment returns low-level information of all the containers of the provided lab environment.
// On success, it returns a map of the containers' names as keys and their information as values.
func InspectEnvironment(labName string) (map[string][]byte, error) {
	containers, err := GetEnvironmentContainers(labName)
	if err != nil {
		return nil, fmt.Errorf("error while inspecting environment: %w", err)
	}

	inspectMap, err := inspectContainers(containers)
	if err != nil {
		return nil, fmt.Errorf("error while inspecting environment: %w", err)
	}

	inspectJSONMap := make(map[string][]byte, len(inspectMap))

	for container, inspectInfo := range inspectMap {
		inspectJSON, err := json.MarshalIndent(inspectInfo, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("error while inspecting environment: %w", err)
		}
		inspectJSONMap[container] = inspectJSON
	}	

	return inspectJSONMap, nil
}

// inspectContainers concurrently inspects (using errgroups) all the specified containers.
func inspectContainers(containers []LabContainer) (map[string]types.ContainerJSON, error) {
	g, ctx := errgroup.WithContext(context.Background())

	inspectMap := make(map[string]types.ContainerJSON, len(containers))

	for _, container := range containers {
		name, id := container.Name, container.ID
		g.Go(func() error {
			inspectInfo, err := inspectContainer(ctx, id)
			if err == nil {
				inspectMap[name] = inspectInfo
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("error while inspecting containers: %w", err)
	}

	return inspectMap, nil
}

// inspectContainer inspects the specified container and returns the required information.
func inspectContainer(errGroupCtx context.Context, containerID string) (types.ContainerJSON, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("error while inspecting container %v: %v", containerID, err)
	}

	inspectJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("error while inspecting container %v: %v", containerID, err)
	}

	return inspectJSON, nil
}