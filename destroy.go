package labsoncontainers

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DestroyEnviroment removes all containers (including running containers),
// networks and the X11 cookie directory of the provided lab enviroment.
func DestroyEnviroment(labName string) error {
	err := destroyCookieDir(labName)
	if err != nil {
		return fmt.Errorf("error while destroying enviroment: %w", err)
	}

	containersIds, err := GetEnviromentContainers(labName)
	if err != nil {
		return fmt.Errorf("error while destroying enviroment: %w", err)
	}

	err = destroyContainers(containersIds)
	if err != nil {
		return fmt.Errorf("error while destroying enviroment: %w", err)
	}

	networksIds, err := GetEnviromentNetworks(labName)
	if err != nil {
		return fmt.Errorf("error while destroying enviroment: %w", err)
	}

	err = destroyNetworks(networksIds)
	if err != nil {
		return fmt.Errorf("error while destroying enviroment: %w", err)
	}

	return nil
}

// destroyCookieDir removes the X11 cookie directory of the provided lab enviroment.
// It does not return an error if the directory does not exist.
func destroyCookieDir(labName string) error {
	labCookieDir := cookieDir + labName + "/"

	err := os.RemoveAll(labCookieDir)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	} else {
		return fmt.Errorf("error while deleting %v: %v", labCookieDir, err)
	}
}

// destroyContainers concurrently removes (using errgroups) all the specified containers.
func destroyContainers(containersIds []string) error {
	g, ctx := errgroup.WithContext(context.Background())

	for _, id := range containersIds {
		id := id
		g.Go(func() error {
			return destroyContainer(ctx, id)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error while destroying containers: %w", err)
	}

	return nil
}

// destroyContainer removes the specified container. It uses Docker force option.
func destroyContainer(errGroupCtx context.Context, containerID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error while destroying container %v: %v", containerID, err)
	}

	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("error while destroying container %v: %v", containerID, err)
	}

	return nil
}

// destroyNetworks concurrently removes (using errgroups) all the specified networks.
func destroyNetworks(networksIds []string) error {
	g, ctx := errgroup.WithContext(context.Background())

	for _, id := range networksIds {
		id := id
		g.Go(func() error {
			return destroyNetwork(ctx, id)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error while destroying networks: %w", err)
	}

	return nil
}

// destroyNetwork removes the specified network.
func destroyNetwork(errGroupCtx context.Context, networkID string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error while destroying network %v: %v", networkID, err)
	}

	err = cli.NetworkRemove(ctx, networkID)
	if err != nil {
		return fmt.Errorf("error while destroying network %v: %v", networkID, err)
	}

	return nil
}