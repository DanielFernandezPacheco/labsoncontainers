# LabsOnContainers

[![Go Reference](https://pkg.go.dev/badge/github.com/marioromandono/labsoncontainers.svg)](https://pkg.go.dev/github.com/marioromandono/labsoncontainers) 
[![Go Report Card](https://goreportcard.com/badge/github.com/marioromandono/labsoncontainers)](https://goreportcard.com/report/github.com/marioromandono/labsoncontainers)
 
**LabsOnContainers** is a tool written in Go that allows the automatic deployment and configuration of containerized lab environments, using Docker. These lab environments are designed to practice Computer Science exercises: monitor network traffic, create VPNs, perform mock cyberattacks...

This tool was created as part of my Final Degree Project for my Bachelor’s Degree in Computer Science Engineering. You can read the whole paper (***Implementación de un laboratorio usando contenedores***) in the [E-Prints Complutense Repository](https://eprints.ucm.es/id/eprint/74396/).
 
The created containers have the following features:
 
- Allow interaction with the container `(-it)`
- Possibility of connecting the containers to one or more networks in order to recreate full network infrastructures.
- Management of network settings `(--cap-add=NET_ADMIN)`
- Running init system `(--init)`
- Execution of GUI apps is possible using X11 untrusted cookies (https://github.com/mviereck/x11docker/wiki/X-authentication-with-cookies-and-xhost-("No-protocol-specified"-error)#untrusted-cookie-for-container-applications)
- Home directory is bind mounted on all the lab containers on `/mnt/shared`
 
This project is divided in two packages:

- Package `labsoncontainers`: This package is conceived as an API that provides the necessary functions for the management of the lab environments, encapsulating the calls to the Docker Engine API and the X Window System. This package is contained in the .go files of the root directory of this repository. The documentation of this package is available in the [Go Packages website](https://pkg.go.dev/github.com/marioromandono/labsoncontainers).
- Package `main`: This package contains a CLI that uses the functions of the `labsoncontainers` package. This package is contained in the `cmd/LabsOnContainers` directory.

## Usage

To use LabsOnContainers, just run `labsoncontainers [-option] [file.yaml | labName]`. The available options are:

- `-c file.yaml`: Creates a new lab environment from the infrastructure specified in a YAML file. The YAML file must follow this structure:
```
nombre_practica: ejemplo
contenedores:
- nombre: vm1
    imagen: rys
    redes:
    - nombre: 1
    ip: 192.168.1.10/24
    background: true
- nombre: router
    imagen: rys
    redes:
    - nombre: 1
    - nombre: 2
    ip: 192.168.2.10/24
- nombre: vm2
    imagen: rys
    redes:
    - nombre: 2
```
As you can see, every container can be connected to one or more networks. The IP field is optional: if it is not provided, Docker will assign an IP automatically. The background field is also optional: if it is true, the terminal window corresponding to the container will not be created. This is useful if no interaction is needed with the container, for example using a web server.

- `-i labName`: Prints information of every container of the lab environment. This is similar to running `docker container inspect`.
- `-p labName`: Stops every container of the lab environment.
- `-l labName`: Restarts every container of the lab environment.
- `-r labName`: Destroys every container and networks of the lab environment.

In these four options, `labName` must be equal to the `nombre_practica` specified in the YAML file when creating the lab environment. 

## Requirements

- Docker Engine
- X Window System with the security extension enabled. It is also needed to have installed `xauth` in your system.
- Currently, the terminal windows creation process is heavily attached to the operative system used in my Final Degree Project, Alpine Linux with the XFCE desktop environment. Therefore, it is needed to have installed `xfce4-terminal`.
- If you want to compile the package by yourself, you need to have installed Go in your system (you can also use `make` to simplify the process).

This tool was designed with the idea of restricting the users —students in Computer Science faculty labs— from creating containers by themselves, which could lead to exploiting vulnerabilities and taking advantage of the underlying system. That way, LabsOnContainers would be the only available tool for creating and managing containers. In order to accomplish this objective, LabsOnContainers is designed to be used as a `setuid` executable, owned by `root`, so the users could create container-based labs environments as the executable would have enough privileges to call the Docker Engine API, but not allowing the students to use Docker directly. However, if you don't need these restrictions, theoretically you could run LabsOnContainers without setting `setuid`, choosing between adding your user to the `docker` group or running LabsOnContainers with `sudo`. **Keep in mind that these two options have not been tested and there could be additional failures.**

Independently of the option you've chosen, it is also necessary to have installed `sudo` and add this line to the sudoers file (replace `youruser` by your actual user): `youruser ALL = (root) NOPASSWD: /usr/bin/docker *attach*`. This is because, during the creation of LabsOnContainers, I struggled with the creation of the terminal windows for interacting with the containers, due to GTK applications like `xfce4-terminal` won't run in setuid processes. This was the easiest solution at the moment, but it is not at all convenient. If you know an equivalent solution that does not require adding the user to the sudoers file, please open an issue or pull request.

## Downloading and/or compiling LabsOnContainers

You can download LabsOnContainers from the releases page. Alternatively, you can compile the executable yourself.

The simplest way to compile LabsOnContainers is to run `make` in the root directory. If `make` is not installed, alternatively you can run `go build cmd/LabsOnContainers/*`. Setting `CGO_ENABLED=0` before the `go build` command can be useful in order to get a statically linked executable instead of a dynamically linked executable.

After downloading or compiling the executable, you have to set its permissions as described in [Requirements](#requirements). `sudo make install` is available too, which will copy `labsoncontainers` executable to `/usr/local/bin` and enable the `setuid` bit. 

## License and contributions

The whole project is under the GNU General Public License version 3. See [COPYING](../COPYING).

You can create your own tools using the `labsoncontainers` package as long as you follow the terms of the GPL. Also, all contributions are welcomed! Please open an issue or create a pull request if you have any problem or want to improve LabsOnContainers.

## TODO

- Improve the terminal windows creation process, with the objective of not having to require the installation of `xfce4-terminal` nor having to add the current user to the sudoers file.
