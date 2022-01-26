#!/bin/sh

# Copyright (C) 2021 Mario Román Dono

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.



# Uso: ./lanzarContainers.sh n_redes [<nombre:imagen:redes>]

# El contenedor se arranca en modo detached para poder saber cuándo termina de crearse y poder configurar las redes correctamente
# Se añade la capacidad NET_ADMIN para poder configurar las interfaces de red
# Es necesario pasar la variable de entorno DISPLAY para poder usar aplicaciones de interfaz gráfica
# Se monta la carpeta en la que se ejecuta el comando para compartir archivos con el contenedor
# Se monta el socket de X11 para poder usar aplicaciones de interfaz gráfica
# Se crea una nueva cookie de X11 para poder usar el servidor gráfico sin necesidad de xhost
# La opción init existe porque se van a ejecutar varios procesos dentro del contenedor

run_container() {
    docker container run \
        --name $nombre \
        --hostname $nombre \
        -d -it --cap-add=NET_ADMIN \
        --init \
        --env DISPLAY \
        --env XAUTHORITY=/cookie \
        --mount type=bind,source="$(pwd)",target=/mnt/shared \
        --mount type=bind,source=/tmp/.X11-unix,target=/tmp/.X11-unix \
        --mount type=bind,source="${Cookiefile}",target=/cookie \
	$imagen
}

error() {
    echo "$*" 1>&2
    exit 1
}

uso() {
    echo "Uso: $0 <-c fichero | -l nombre_practica | -p nombre_practica | -d nombre_practica | -h>"
    echo ""
    echo "-c: Crea el entorno de contenedores a partir del fichero YAML proporcionado (se destruyen los contenedores asociados al entorno del fichero)"
    echo "-l: Ejecuta todos los contenedores asociados al entorno proporcionado"
    echo "-p: Detiene todos los contenedores asociados al entorno proporcionado"
    echo "-d: Destruye todos los contenedores asociados al entorno proporcionado"
    echo "-h: Muestra este mensaje de ayuda"
}

destruir_entorno() {
    echo "hola" > /dev/null
}

crear_entorno() {
    nombre_practica=$(yq e '.nombre_practica' $fichero)
    if [ $nombre_practica == "null" ]; then
        error "Se debe especificar el campo nombre_practica en el fichero"
    fi

    destruir_entorno

    numero_maquinas=$(yq e '.maquinas | length' $fichero) 
    if [ $numero_maquinas == 0 ]; then
        error "Se deben especificar los contenedores deseados en el campo maquinas del fichero"
    fi

    # Primero se crean todos los contenedores
    i=0
    while [ $i -lt $numero_maquinas ]
    do
        nombre_VM=$(yq e '.maquinas['"$i"'].nombre' $fichero)
        if [ $nombre_VM == "null" ]; then
            destruir_entorno
            error "Se debe especificar el campo nombre en maquinas[$i]"
        fi

        # Se concatena el nombre de la practica con el nombre de la VM para poder hacer operaciones asociadas a entornos de prácticas
        nombre="${nombre_practica}_${nombre_VM}"

        imagen=$(yq e '.maquinas['"$i"'].imagen' $fichero)
        if [ $imagen == "null" ]; then
            destruir_entorno
            error "Se debe especificar el campo imagen en maquinas[$i]"
        fi
        
        # Se crea el contenedor y se le desconecta del adaptador bridge
        run_container || error "Error en la creación del contenedor ${nombre}"
        docker network disconnect bridge $nombre

        numero_redes=$(yq e '.maquinas['"$i"'].redes | length' $fichero) 
        if [ $numero_redes == 0 ]; then
            destruir_entorno
            error "Se debe especificar el campo redes en maquinas[$i]"
        fi

        j=0
        while [ $j -lt $numero_redes ]
        do
            red_fichero=$(yq e '.maquinas['"$i"'].redes['"$j"']' $fichero)
            red="red_${red_fichero}"

            # Solo se crea la red si no existe previamente
            docker network inspect $red > /dev/null 2>&1 || docker network create $red

            docker network connect $red $nombre

            j=$((j+1))
        done
        i=$((i+1))
    done

    # Si todos los contenedores se han creado exitosamente, se abren las terminales
    i=0
    while [ $i -lt $numero_maquinas ]
    do
        nombre_VM=$(yq e '.maquinas['"$i"'].nombre' $fichero)
        nombre="${nombre_practica}_${nombre_VM}"
        xfce4-terminal --tab -e "ash -c 'docker container attach $nombre; exec ash'"
        i=$((i+1))
    done
}

# Creamos la cookie de X11 para usar aplicaciones con GUI
Cookiefile=~/containercookie
:> $Cookiefile
xauth -f $Cookiefile generate $DISPLAY . untrusted timeout 3600
Cookie="$(xauth -f $Cookiefile nlist $DISPLAY | sed -e 's/^..../ffff/')"
echo "$Cookie" | xauth -f "$Cookiefile" nmerge -

# Controla que el comando se ejecute solo con un flag
opcion=false

while getopts ":c:l:d:p:h" o; do
    case "${o}" in
        c)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-d|-h>"
            fichero=${OPTARG}

            if [ ! -f $fichero ]; then
                error "El fichero $fichero no existe"
            fi
            comando="crear"
            opcion=true
            ;;
        l)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-d|-h>"
            nombre_practica=${OPTARG}

            if [ -z $nombre_practica ]; then
                error "Se debe proporcionar el valor nombre_practica"
            fi

            comando="lanzar"
            opcion=true
            ;;
        d)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-d|-h>"
            nombre_practica=${OPTARG}

            if [ -z $nombre_practica ]; then
                error "Se debe proporcionar el valor nombre_practica"
            fi

            comando="destruir"
            opcion=true
            ;;
        p)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-d|-h>"
            nombre_practica=${OPTARG}

            if [ -z $nombre_practica ]; then
                error "Se debe proporcionar el valor nombre_practica"
            fi

            comando="parar"
            opcion=true
            ;;
        h)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-d|-h>"

            uso
            exit 0
            ;;
        *)
            echo "Opción ${o} no válida"
            uso
            exit 1
            ;;
    esac
done
shift $((OPTIND-1))

case $comando in
    crear)
        crear_entorno && exit
        ;;
    lanzar)
        lanzar_entorno && exit
        ;;
    destruir)
        destruir_entorno && exit
        ;;
    parar)
        parar_entorno && exit
        ;;
    *)
        uso && exit 1
        ;;
esac