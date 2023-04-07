// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"

	"gopkg.in/yaml.v3"

	"errors"
	"os/user"
	"time"
)

// createLabEnvironment parses a YAML file and, if it has a correct format,
// it converts each element of the file and calls the LabsOnContainers API.
func createLabEnvironment(filePath string) bool {
	fmt.Println("Creando el entorno de laboratorio...")

	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error while opening file: %v\n", err)
		os.Exit(1)
	}

	var labEnv labsoncontainers.LabEnvironment

	err = yaml.Unmarshal(file, &labEnv)
	if err != nil {
		fmt.Printf("error while parsing yaml file: %v\n", err)
		os.Exit(1)
	}

	_, err = labsoncontainers.CreateEnvironment(&labEnv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	containers, err := labsoncontainers.GetEnvironmentContainers(labEnv.LabName)
	if err != nil {
		fmt.Println(err)
		destroyErr := labsoncontainers.DestroyEnvironment(labEnv.LabName)
		if destroyErr != nil {
			fmt.Printf("error on labsoncontainers: %v\n", err)
		}
		os.Exit(1)
	}

	printContainersInfo(containers)

	err = createTerminalWindows(containers)
	if err != nil {
		fmt.Printf("error while creating terminal windows: %v\n", err)
		destroyErr := labsoncontainers.DestroyEnvironment(labEnv.LabName)
		if destroyErr != nil {
			fmt.Printf("error on labsoncontainers: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("Entorno creado exitosamente")

	//Comprobación de la existencia del directorio que almacenará el fichero de logs (en caso de no existir, se creará)
	dir_path := "/var/log/labsoncontainers"

	if _, err := os.Stat(dir_path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir_path, 0755)
		if err != nil {
			fmt.Printf("error creating log directory: %v\n", err)
			os.Exit(1)
		}
	}

	//Creación o apertura del fichero que almacenará los logs
	f, err := os.OpenFile("/var/log/labsoncontainers/operations.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		fmt.Printf("error creating/opening log file: %v\n", err)
		os.Exit(1)
	}

	//Obtención del nombre del usuario
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("error geting current user name: %v\n", err)
		os.Exit(1)
	}

	currentUserUsername := currentUser.Username

	//Obtención de la fecha actual
	currentTime := time.Now()
	formatedTime := currentTime.Format("2006-01-02 15:04:05")

	res_output := "[" + formatedTime + "] [+] Operación: Creación del entorno || [+] Nombre del entorno: " + labEnv.LabName + " (Usuario: " + currentUserUsername + ")\n"

	//Escritura en el fichero
	if _, err := f.Write([]byte(res_output)); err != nil {
		fmt.Printf("error writing to log file: %v\n", err)
		os.Exit(1)
	}

	//Cierre del fichero
	if err := f.Close(); err != nil {
		fmt.Printf("error closing log file: %v\n", err)
		os.Exit(1)
	}

	return true
}
