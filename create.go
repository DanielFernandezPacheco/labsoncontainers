// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package labsoncontainers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// CreateEnvironment creates a lab environment using Docker Engine SDK and LabEnvironment type.
//
// First, it destroy any other lab environment with the same name using DestroyEnvironment, then
// it creates all the desired networks and finally, it creates all the containers. On success, it
// returns a map of the created containers names and their IDs.
//
// Note that, during container creation, it will not be pulled any image, so the desired images will have
// to be previously built or pulled.
func CreateEnvironment(labEnv *LabEnvironment) (map[string]string, error) {
	if labEnv.LabName == "" {
		return nil, fmt.Errorf("error while creating environment: lab name cannot be empty")
	}

	err := DestroyEnvironment(labEnv.LabName)
	if err != nil {
		return nil, fmt.Errorf("error while creating environment: %w", err)
	}

	networks, err := parseNetworks(labEnv)
	if err != nil {
		return nil, fmt.Errorf("error while creating environment: %w", err)
	}

	err = createNetworks(networks)
	if err != nil {
		createErr := fmt.Errorf("error while creating environment: %w", err)
		destroyErr := DestroyEnvironment(labEnv.LabName)
		if destroyErr != nil {
			createErr = fmt.Errorf("%w, %v", createErr, destroyErr)
		}
		return nil, createErr
	}

	containerIds, err := createContainers(labEnv.Containers, labEnv.LabName)
	if err != nil {
		createErr := fmt.Errorf("error while creating environment: %w", err)
		destroyErr := DestroyEnvironment(labEnv.LabName)
		if destroyErr != nil {
			createErr = fmt.Errorf("%w, %v", createErr, destroyErr)
		}
		return nil, createErr
	}

	return containerIds, nil
}

// parseNetworks traverses the LabEnvironment struct and returns a slice of LabNetwork. In case that two or more networks
// have the same name but their IP address does not belong to the same subnet, it will use the first network address.
func parseNetworks(labEnv *LabEnvironment) ([]LabNetwork, error) {
	uniqueNetworks := make(map[string]string)

	for _, container := range labEnv.Containers {
		for _, network := range container.Networks {
			networkFullName := labEnv.LabName + "_" + network.Name

			if uniqueNetworkIP, ok := uniqueNetworks[networkFullName]; !ok || uniqueNetworkIP == "" {
				var ip string
				if network.IP != "" {
					_, ipv4Net, err := net.ParseCIDR(network.IP)
					if err != nil {
						return nil, fmt.Errorf("error while parsing network addresses: %v", err)
					}
					ip = ipv4Net.String()
				}
				uniqueNetworks[networkFullName] = ip
			}
		}
	}

	var labNetworkSlice []LabNetwork

	for name, ip := range uniqueNetworks {
		network := LabNetwork{
			Name: name,
			IP: ip,
		}
		labNetworkSlice = append(labNetworkSlice, network)
	}

	return labNetworkSlice, nil
}

// createNetworks concurrently creates (using errgroups) all the specified networks.
func createNetworks(networks []LabNetwork) error {
	g, ctx := errgroup.WithContext(context.Background())

	for _, network := range networks {
		network := network
		g.Go(func() error {
			return createNetwork(ctx, &network)
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error while creating networks: %w", err)
	}

	return nil
}

// createNetwork creates the specified network. If address is not empty, it will be used 
// as the network subnet address.
func createNetwork(errGroupCtx context.Context, labNetwork * LabNetwork) error {
	var IPAM *network.IPAM

	if labNetwork.IP != "" {
		IPAM = &network.IPAM{
			Config: []network.IPAMConfig{
				{Subnet: labNetwork.IP},
			},
		}
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error while creating network %v: %v", labNetwork.Name, err)
	}

	_, err = cli.NetworkCreate(ctx, labNetwork.Name, types.NetworkCreate{
		CheckDuplicate: true,
		IPAM:           IPAM,
	})
	if err != nil {
		return fmt.Errorf("error while creating network %v: %v", labNetwork.Name, err)
	}

	return nil
}

// createX11Cookie creates an untrusted X11 cookie in order to allow the container to run graphical apps.
// This function is based on the commands specified on: https://github.com/mviereck/x11docker/wiki/X-authentication-with-cookies-and-xhost-("No-protocol-specified"-error)#untrusted-cookie-for-container-applications
func createX11Cookie(containerName string, labName string) (string, error) {
	labCookieDir := cookieDir + labName + "/"

	err := os.MkdirAll(labCookieDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return "", fmt.Errorf("%v", err)
	}
	if err := os.Chown(labCookieDir, os.Getuid(), os.Getgid()); err != nil {
		return "", fmt.Errorf("%v", err)
	}

	cookiePath := labCookieDir + containerName + "_cookie"
	if _, err := os.Create(cookiePath); err != nil {
		return "", fmt.Errorf("%v", err)
	}
	if err := os.Chown(cookiePath, os.Getuid(), os.Getgid()); err != nil {
		return "", fmt.Errorf("%v", err)
	}  

	DISPLAY := os.Getenv("DISPLAY")
	xauthCmd := exec.Command("xauth", "-f", cookiePath, "generate", DISPLAY, ".", "untrusted", "timeout", "3600")
	if err := xauthCmd.Run(); err != nil {
		return "", fmt.Errorf("%v", err)
	}
	if err := os.Chown(cookiePath, os.Getuid(), os.Getgid()); err != nil {
		return "", fmt.Errorf("%v", err)
	}
	
	xauthCmd = exec.Command("xauth", "-f", cookiePath, "nlist", DISPLAY)
	stdout, err := xauthCmd.Output()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}	

	r, err := regexp.Compile("^....")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	cookie := r.ReplaceAll(stdout, []byte("ffff"))

	
	xauthCmd = exec.Command("xauth", "-f", cookiePath, "nmerge", "-")
	xauthCmd.Stdin = bytes.NewReader(cookie)
	if err := xauthCmd.Run(); err != nil {
		return "", fmt.Errorf("%v", err)
	}
	if err := os.Chown(cookiePath, os.Getuid(), os.Getgid()); err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return cookiePath, nil
}

// createContainers concurrently creates (using errgroups) all the specified containers.
func createContainers(containers []LabContainer, labName string) (map[string]string, error) {
	containerIds := make(map[string]string)

	g, ctx := errgroup.WithContext(context.Background())

	for _, container := range containers {
		container := container
		g.Go(func() error {
			containerId, err := createContainer(ctx, &container, labName)
			if err == nil {
				containerIds[container.Name] = containerId
				return nil
			} else {
				return fmt.Errorf("error while creating container %v: %w", container.Name, err)
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("error while creating containers: %w", err)
	}

	return containerIds, nil
}

// createContainer creates the container. This function is equivalent to run:
//
// docker container run --name (labName_labContainer.Name) --hostname (labName_labContainer.Name) -d -it --cap-add=NET_ADMIN \
// --init --env DISPLAY --env XAUTHORITY=cookiePath --mount type=bind,source=~,target=/mnt/shared \
// --mount type=bind,source=/tmp/.X11-unix,target=/tmp/.X11-unix \ --mount type=bind,source=cookiePath,target=cookiePath \
// --label background=labContainer.Background labContainer.Image
//
// If the container is successfully created, it is disconnected from the bridge network and connected to the specified networks.
func createContainer(errGroupCtx context.Context, labContainer *LabContainer, labName string) (string, error) {
	if labContainer.Name == "" {
		return "", fmt.Errorf("container name cannot be empty")
	}
	containerFullName := labName + "_" + labContainer.Name

	if labContainer.Image == "" {
		return "", fmt.Errorf("container image cannot be empty")
	}
	imageName := labContainer.Image

	cookiePath, err := createX11Cookie(containerFullName, labName)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	displayEnvVar := "DISPLAY=" + os.Getenv("DISPLAY")
	xauthorityEnvVar := "XAUTHORITY=" + cookiePath
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Hostname:  containerFullName,
		Tty:       true,
		OpenStdin: true,
		Env: []string{
			displayEnvVar,
			xauthorityEnvVar,
		},
		Labels: map[string]string{
			"background": strconv.FormatBool(labContainer.Background),
		},
		Image: imageName,
	}, &container.HostConfig{
		Init: boolPointer(true),
		CapAdd: []string{
			"NET_ADMIN",
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: homeDir,
				Target: "/mnt/shared",
			},
			{
				Type:   mount.TypeBind,
				Source: "/tmp/.X11-unix",
				Target: "/tmp/.X11-unix",
			},
			{
				Type:   mount.TypeBind,
				Source: cookiePath,
				Target: cookiePath,
			},
		},
	}, nil, nil, containerFullName)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	containerID := resp.ID

	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("%v", err)
	}

	if err := cli.NetworkDisconnect(ctx, "bridge", containerID, false); err != nil {
		return "", fmt.Errorf("%v", err)
	}

	if err := connectToNetworks(labContainer.Networks, containerID, labName, labContainer.Name); err != nil {
		return "", fmt.Errorf("%w", err)
	}	

	return containerID, nil
}

// connectToNetworks concurrently connects (using errgroups) the container to the specified networks.
func connectToNetworks(networks []LabNetwork, containerID string, labName string, containerName string) error {
	g, ctx := errgroup.WithContext(context.Background())

	for _, network := range networks {
		network := network
		g.Go(func() error {
			return connectToNetwork(ctx, containerID, &network, labName)
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error while connecting container %v to networks: %w", containerName, err)
	}
	return nil
}

// connectToNetwork connects the container to the specified lab network. 
// If the lab network has an IP address, the container will use that address (keep in mind that this function can return
// an error if the IP address is not included in the network subnet).
func connectToNetwork(errGroupCtx context.Context, containerID string, labNetwork *LabNetwork, labName string) error {
	var endpointSettings *network.EndpointSettings
	if labNetwork.IP != "" {
		ipv4Addr, _, err := net.ParseCIDR(labNetwork.IP)
		if err != nil {
			return fmt.Errorf("error while connecting to network %v: %v", labNetwork.Name, err)
		}
		endpointSettings = &network.EndpointSettings{
			IPAMConfig: &network.EndpointIPAMConfig{
				IPv4Address: ipv4Addr.String(),
			},
		}
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error while connecting to network %v: %v", labNetwork.Name, err)
	}

	networkFullName := labName + "_" + labNetwork.Name
	err = cli.NetworkConnect(ctx, networkFullName, containerID, endpointSettings)
	if err != nil {
		return fmt.Errorf("error while connecting to network %v: %v", labNetwork.Name, err)
	}
	return nil
}