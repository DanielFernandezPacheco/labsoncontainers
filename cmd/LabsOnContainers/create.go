package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"

	"gopkg.in/yaml.v3"
)

// createLabEnviroment parses a YAML file and, if it has a correct format,
// it converts each element of the file and calls the LabsOnContainers API.
func createLabEnviroment(filePath string) {
    fmt.Println("Creando el entorno de laboratorio...")

    file, err := os.ReadFile(filePath)
    if err != nil {
        fmt.Printf("error while opening file: %v\n", err)
        os.Exit(1)
    }

    var labEnv labsoncontainers.LabEnviroment

    err = yaml.Unmarshal(file, &labEnv)
    if err != nil {
        fmt.Printf("error while parsing yaml file: %v\n", err)
        os.Exit(1)
    }

    containersIds, err := labsoncontainers.CreateEnviroment(&labEnv)
    if err != nil {
        fmt.Printf("error on labsoncontainers: %v\n", err)
        os.Exit(1)
    }

    var foregroundContainersIds []string
    for _, container := range labEnv.Containers {
        if !container.Background {
            foregroundContainersIds = append(foregroundContainersIds, containersIds[container.Name])
        }
    }    
    err = createTerminalWindows(foregroundContainersIds) 
    if err != nil {
        fmt.Printf("error while creating terminal windows: %v\n", err)
        destroyErr := labsoncontainers.DestroyEnviroment(labEnv.LabName)
        if destroyErr != nil {
            fmt.Printf("error on labsoncontainers: %v\n", err)
        }
        os.Exit(1)
    }

    fmt.Println("Entorno creado exitosamente")
}