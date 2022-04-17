package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"
)

// destroyLabEnviroment destroys the specified lab enviroment using LabsOnContainers API.
func destroyLabEnviroment(labName string) {
	fmt.Printf("Eliminando los contenedores y redes de %v...\n", labName)
	containersIds, err := labsoncontainers.GetEnviromentContainers(labName)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

	if len(containersIds) > 0 {
		err := labsoncontainers.DestroyEnviroment(labName)
   		if err != nil {
			fmt.Println(err)
			os.Exit(1)
   		}
		fmt.Println("Contenedores y redes eliminados exitosamente")		
	} else {
		fmt.Println("No existen contenedores asociados a", labName)
	}

}