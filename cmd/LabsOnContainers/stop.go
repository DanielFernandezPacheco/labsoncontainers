package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"
)

// stopLabEnviroment stops the containers of the specified lab enviroment using LabsOnContainers API.
func stopLabEnviroment(labName string) {
	fmt.Printf("Deteniendo los contenedores y redes de %v...\n", labName)
	containersIds, err := labsoncontainers.GetEnviromentContainers(labName)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

	if len(containersIds) > 0 {
		err := labsoncontainers.StopEnviroment(labName)
   		if err != nil {
			fmt.Println(err)
			os.Exit(1)
   		}
		fmt.Println("Contenedores detenidos exitosamente")		
	} else {
		fmt.Println("No existen contenedores asociados a", labName)
	}
}