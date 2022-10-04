// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"
)

// stopLabEnvironment stops the containers of the specified lab environment using LabsOnContainers API.
func stopLabEnvironment(labName string) {
	fmt.Printf("Deteniendo los contenedores y redes de %v...\n", labName)
	containersIds, err := labsoncontainers.GetEnvironmentContainers(labName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(containersIds) > 0 {
		err := labsoncontainers.StopEnvironment(labName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Contenedores detenidos exitosamente")
	} else {
		fmt.Println("No existen contenedores asociados a", labName)
	}
}
