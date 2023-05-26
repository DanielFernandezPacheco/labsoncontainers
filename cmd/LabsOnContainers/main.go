// SPDX-FileCopyrightText: 2022 Mario Román Dono <mario.romandono@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/marioromandono/labsoncontainers"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

type Red struct {
	Nombre int
	IP     string `yaml:"ip,omitempty"`
}

type Contenedor struct {
	Nombre string `yaml:"nombre"`
	Imagen string `yaml:"imagen"`
	Redes  []Red  `yaml:"redes"`
}

type Practica struct {
	NombrePractica string       `yaml:"nombre_practica"`
	Contenedores   []Contenedor `yaml:"contenedores"`
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	// parses the form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form values", http.StatusBadRequest)
		return
	}
	nombrePractica := html.EscapeString(r.FormValue("nombre-practica"))
	numContenedores, _ := strconv.Atoi(r.FormValue("num-contenedores"))

	var p Practica
	p.NombrePractica = nombrePractica

	// iterates over the containers and networks to add them to the structure
	for i := 0; i < numContenedores; i++ {
		var c Contenedor
		c.Nombre = html.EscapeString(r.FormValue(fmt.Sprintf("nombre-contenedor-%d", i)))
		c.Imagen = html.EscapeString(r.FormValue(fmt.Sprintf("nombre-imagen-%d", i)))

		numRedes, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("num-redes-%d", i)))
		for j := 0; j < numRedes; j++ {
			var red Red
			red.Nombre, _ = html.EscapeString(strconv.Atoi(r.FormValue(fmt.Sprintf("nombre-red-%d-%d", i, j))))
			red.IP = html.EscapeString(r.FormValue(fmt.Sprintf("ip-red-%d-%d", i, j)))
			c.Redes = append(c.Redes, red)
		}

		p.Contenedores = append(p.Contenedores, c)
	}

	// converts the struct into yaml
	y, err := yaml.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile("/home/usuario/.labsoncontainers/recent_configs/"+nombrePractica+".yaml", y, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "http://localhost:8080/", http.StatusFound)
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	// get the name of the YAML file from the query string parameter "file"
	fileName := r.URL.Query().Get("file")
	// read the YAML file into a byte slice
	data, err := ioutil.ReadFile("/home/usuario/.labsoncontainers/recent_configs/" + fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse the YAML data into a struct
	var practica Practica
	err = yaml.Unmarshal(data, &practica)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse YAML file %s: %v", fileName, err), http.StatusInternalServerError)
		return
	}

	// encode the struct as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(practica)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FileInfo struct holds the name and ctime
type FileInfo struct {
	Name  string
	Ctime syscall.Timespec
}

// FileInfoList  holds a list of FileInfo
type FileInfoList []FileInfo

// Implement sort.Interface
func (f FileInfoList) Len() int { return len(f) }
func (f FileInfoList) Less(i, j int) bool {
	return f[i].Ctime.Nano() > f[j].Ctime.Nano()
}
func (f FileInfoList) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func listYAMLFiles(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var filesInfo FileInfoList
	for _, file := range files {
		if path.Ext(file.Name()) == ".yaml" {
			fileStat := syscall.Stat_t{}
			filePath := filepath.Join(directory, file.Name())
			syscall.Stat(filePath, &fileStat)
			filesInfo = append(filesInfo, FileInfo{file.Name(), fileStat.Ctim})
		}
	}
	sort.Sort(filesInfo)
	var yamlFiles []string
	for _, file := range filesInfo {
		yamlFiles = append(yamlFiles, file.Name)
	}
	return yamlFiles, nil
}

func main() {

	// Serve the files in the hidden directory
	fs := http.FileServer(http.Dir("/home/usuario/.labsoncontainers/"))

	http.Handle("/", fs)

	http.HandleFunc("/form", handleForm)
	http.HandleFunc("/getConfig", getConfig)

	http.HandleFunc("/yaml-files", func(w http.ResponseWriter, r *http.Request) {
		files, err := listYAMLFiles("/home/usuario/.labsoncontainers/recent_configs/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Query().Get("env")
		if createLabEnvironment("/home/usuario/.labsoncontainers/recent_configs/" + file) {
			http.Redirect(w, r, "http://localhost:8080/practica?env="+strings.TrimSuffix(file, ".yaml"), http.StatusOK)
			fmt.Fprintf(w, "El entorno del laboratorio ha sido creado exitosamente.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		env := r.URL.Query().Get("env")
		if startLabEnvironment(env) {
			http.Redirect(w, r, "http://localhost:8080/practica?env="+env, http.StatusOK)
			fmt.Fprintf(w, "El laboratorio se ha ejecutado correctamente.")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		env := r.URL.Query().Get("env")
		stopLabEnvironment(env)
		http.Redirect(w, r, "http://localhost:8080/"+env, http.StatusOK)
	})

	http.HandleFunc("/inspect", func(w http.ResponseWriter, r *http.Request) {
		env := r.URL.Query().Get("env")

		inspectLabEnvironment(env)
	})

	http.HandleFunc("/destroy", func(w http.ResponseWriter, r *http.Request) {
		env := r.URL.Query().Get("env")
		destroyLabEnvironment(env)
		http.Redirect(w, r, "http://localhost:8080/"+env, http.StatusOK)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		env := r.URL.Query().Get("env")
		file := env + ".yaml"

		destroyLabEnvironment(env)
		if createLabEnvironment("/home/usuario/.labsoncontainers/recent_configs/" + file) {
			http.Redirect(w, r, "http://localhost:8080/practica?env="+env, http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

	})

	http.HandleFunc("/getContainers", func(w http.ResponseWriter, r *http.Request) {
		env := r.URL.Query().Get("env")
		containers, err := labsoncontainers.GetEnvironmentContainers(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(containers)

	})

	http.HandleFunc("/stopContainer", func(w http.ResponseWriter, r *http.Request) {
		containerID := r.URL.Query().Get("container")
		g, ctx := errgroup.WithContext(context.Background())
		g.Go(func() error {
			return labsoncontainers.StopContainer(ctx, containerID)
		})
	})

	http.HandleFunc("/startContainer", func(w http.ResponseWriter, r *http.Request) {
		containerID := r.URL.Query().Get("container")
		env := r.URL.Query().Get("env")
		g, ctx := errgroup.WithContext(context.Background())
		g.Go(func() error {
			labsoncontainers.StartContainer(ctx, containerID)
			containers, err := labsoncontainers.GetEnvironmentContainers(env)
			if err != nil {
				return fmt.Errorf("error while attaching container %v: %v", containerID, err)
			}

			var container labsoncontainers.LabContainer
			for i := 0; i < len(containers); i++ {
				if containers[i].ID == containerID {
					container = containers[i]
					break
				}
			}

			createTerminalWindows([]labsoncontainers.LabContainer{container})
			return nil
		})

	})
	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
