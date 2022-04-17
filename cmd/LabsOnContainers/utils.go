package main

import (
	"os/exec"
	"os"
	"syscall"
)

// createTerminalWindows creates terminal windows that make possible the
// interaction with the created containers. At the moment, this function
// only works with XFCE terminal, and it should be refactored in order to 
// support more terminal emulators.
func createTerminalWindows(containersIds []string) error {
	args := []string{"-e"}
	for i := 0; i < len(containersIds); i++ {
		if i > 0 {
			args = append(args, "--tab", "-e")
		}
		args = append(args, "ash -c 'docker container attach " + containersIds[i] + "; exec ash'")
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