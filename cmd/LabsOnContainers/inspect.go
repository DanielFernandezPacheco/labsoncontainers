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

// inspectLabEnvironment prints to the standard output the result of running inspect
// on the lab environment containers, using LabsOnContainers API.
func inspectLabEnvironment(labName string) {
	fmt.Printf("Información de los contenedores de %v:\n", labName)
	containersIds, err := labsoncontainers.GetEnvironmentContainers(labName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(containersIds) > 0 {
		inspectMap, err := labsoncontainers.InspectEnvironment(labName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for container, inspect := range inspectMap {
			fmt.Println()
			fmt.Printf("%v:\n", container)
			fmt.Println(string(inspect))
		}


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

		res_output := "[" + formatedTime + "] [+] Operación: Inspección exitosa del entorno || [+] Nombre del entorno: " + labName + " (Usuario: " + currentUserUsername + ")\n"


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

                res_output := "[" + formatedTime + "] [+] Operación: Inspección fallida del entorno (No existen contenedores asociados al entorno) || [+] Nombre del entorno: " + labName + " (Usuario: " + currentUserUsername + ")\n"


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
