package bindings

type Network struct {
    Nombre string
    IP string
}

type Container struct {
    Nombre string
    Imagen string
    Redes []*Network
}

type LabEnviroment struct {
    NombrePractica string `yaml:"nombre_practica"`
    Contenedores []*Container
}