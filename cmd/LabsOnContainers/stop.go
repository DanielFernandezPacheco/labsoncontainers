// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"

	"github.com/marioromandono/labsoncontainers"

	"os/user"
	"time"
	"errors"
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

		res_output := "[" + formatedTime + "] [+] Operación: Detención exitosa del entorno || [+] Nombre del entorno: " + labName + " (Usuario: " + currentUserUsername + ")\n"


		//Escritura en el fichero
        	if _, err := f.Write([]byte(res_output)); err != nil {
                	fmt.Printf("error writing to log file: %v\n", err)
                	os.Exit(1)
	        }


		if err := f.Close(); err != nil {
                	fmt.Printf("error closing log file: %v\n", err)
                	os.Exit(1)
       	 	}


	} else {
		fmt.Println("No existen contenedores asociados a", labName)


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

                res_output := "[" + formatedTime + "] [+] Operación: Detención fallida del entorno (No existen contenedores asociados al entorno) || [+] Nombre del entorno: " + labName + " (Usuario: " + currentUserUsername + ")\n"

		
		//Escritura en el fichero
                if _, err := f.Write([]byte(res_output)); err != nil {
                        fmt.Printf("error writing to log file: %v\n", err)
                        os.Exit(1)
                }

                if err := f.Close(); err != nil {
                        fmt.Printf("error closing log file: %v\n", err)
                        os.Exit(1)
                }
	}
}
