// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"
)

// destroyLabEnvironment destroys the specified lab environment using LabsOnContainers API.
func destroyLabEnvironment(labName string) {
	fmt.Printf("Eliminando los contenedores y redes de %v...\n", labName)
	containersIds, err := labsoncontainers.GetEnvironmentContainers(labName)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

	if len(containersIds) > 0 {
		err := labsoncontainers.DestroyEnvironment(labName)
   		if err != nil {
			fmt.Println(err)
			os.Exit(1)
   		}
		fmt.Println("Contenedores y redes eliminados exitosamente")		
	} else {
		fmt.Println("No existen contenedores asociados a", labName)
	}
}