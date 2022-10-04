// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"
)

// startLabEnvironment starts the containers of the specified lab environment using LabsOnContainers API.
func startLabEnvironment(labName string) {
	fmt.Printf("Lanzando de nuevo los contenedores y redes de %v...\n", labName)
	containers, err := labsoncontainers.GetEnvironmentContainers(labName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(containers) > 0 {
		err := labsoncontainers.StartEnvironment(labName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// It is necessary to call again to GetEnvironmentContainers because we can only know
		// the containers IPs after the restart
		containers, err = labsoncontainers.GetEnvironmentContainers(labName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		printContainersInfo(containers)

		err = createTerminalWindows(containers)
		if err != nil {
			fmt.Printf("error while creating terminal windows: %v\n", err)
			stopErr := labsoncontainers.StopEnvironment(labName)
			if stopErr != nil {
				fmt.Printf("error on labsoncontainers: %v\n", err)
			}
			os.Exit(1)
		}

		fmt.Println("Contenedores lanzados exitosamente")
	} else {
		fmt.Println("No existen contenedores asociados a", labName)
	}
}
