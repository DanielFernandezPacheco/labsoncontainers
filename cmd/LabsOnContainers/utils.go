package main

import (
	"os/exec"
	"os"
	"syscall"
	"fmt"

	"github.com/marioromandono/labsoncontainers"
)

// createTerminalWindows creates terminal windows that make possible the
// interaction with the created containers. At the moment, this function
// only works with XFCE terminal, and it should be refactored in order to 
// support more terminal emulators.
func createTerminalWindows(containers []labsoncontainers.LabContainer) error {
	var foregroundContainersIds []string
	for _, container := range containers {
		if !container.Background {
			foregroundContainersIds = append(foregroundContainersIds, container.ID)
		}
	}

	args := []string{"-e"}
	for i := 0; i < len(foregroundContainersIds); i++ {
		if i > 0 {
			args = append(args, "--tab", "-e")
		}
		args = append(args, "ash -c 'docker container attach " + foregroundContainersIds[i] + "; exec ash'")
	}

	// GTK apps like xfce4-terminal won't run in setuid processes, so it is necessary
	// to create the process using the real UID and GID
	cmd := exec.Command("xfce4-terminal", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(os.Getuid()), 
			Gid: uint32(os.Getgid()),
		},
	}
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// printContainersInfo prints to the standard output the info of every container of 
// the lab enviroment, including its ID, name, image, background option and networks.
func printContainersInfo(containers []labsoncontainers.LabContainer) {
    for _, container := range containers {
        fmt.Println("ID:", container.ID)
        fmt.Println("Nombre:", container.Name)
        fmt.Println("Imagen:", container.Image)
		fmt.Println("Background:", container.Background)

        for _, network := range container.Networks {
            fmt.Println("Red:", network.Name, "IP:", network.IP)
        }
        fmt.Println()
    }
}