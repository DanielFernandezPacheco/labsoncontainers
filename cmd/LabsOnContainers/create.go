package main

import (
	"fmt"
	"log"
	"os"
	"bytes"

	"github.com/MarioRomanDono/LabsOnContainers/pkg/bindings"

	"gopkg.in/yaml.v3"
)

type LabEnviroments map[string]*bindings.LabEnviroment

func (l *LabEnviroments) Unmarshal(data []byte) error {
    err := yaml.NewDecoder(bytes.NewReader(data)).Decode(l)
    if err != nil {
        return err
    }
    for k, v := range *l {
        v.NombrePractica = k
    }
    return nil
}


/* type City struct {
    Name string
    AreaCode string `yaml:"area_code"`
    Landmarks []string
}

type State struct {
    Name string
    Cities []*City
    Rivers []string
    Presidents []string
}


type States map[string]*State

func (s *States) Unmarshal(data []byte) error {
    err := yaml.NewDecoder(bytes.NewReader(data)).Decode(s)
    if err != nil {
        return err
    }
    for k, v := range *s {
        v.Name= k
    }
    return nil
} */


// This function parses a YAML file and, if it has a correct format,
// it converts each element of the file and calls the LabsOnContainers API
func createLabEnviroment(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	/* // yamlMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &yamlMap)
    if err != nil {
        log.Fatal(err)
    }

	labEnv := bindings.LabEnviroment{} */

	var labEnvs LabEnviroments = map[string]*bindings.LabEnviroment{}
    err = labEnvs.Unmarshal([]byte(file))
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println(labEnvs)

	/* if yamlMap["nombre_practica"] == false {
		log.Fatal("Se debe especificar el campo nombre_practica en el fichero")
	}
	labEnv.NombrePractica = yamlMap["nombre_practica"].(string)

	for i := 0; i < len(yamlMap["perro"]); i++ {

	}

	fmt.Println(labEnv.NombrePractica) */
	
	/* for k, v := range labEnv.Contenedores {
        fmt.Printf("%v -> name: %v, presidents: %v, rivers: %v\n", k, v.N, v.Presidents, v.Rivers)
        if len(v.Cities) > 0 {
            fmt.Printf("cities:\n")
            for _, city := range v.Cities {
                fmt.Printf("%v\n", city)
            }
        }
    } */
}