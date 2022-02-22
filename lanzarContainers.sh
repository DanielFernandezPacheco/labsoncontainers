#!/bin/sh

# Copyright (C) 2021-2022 Mario Román Dono

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



# El contenedor se arranca en modo detached para poder saber cuándo termina de crearse y poder configurar las redes correctamente
# Se añade la capacidad NET_ADMIN para poder configurar las interfaces de red
# Es necesario pasar la variable de entorno DISPLAY para poder usar aplicaciones de interfaz gráfica
# Se monta la carpeta en la que se ejecuta el comando para compartir archivos con el contenedor
# Se monta el socket de X11 para poder usar aplicaciones de interfaz gráfica
# Se crea una nueva cookie de X11 para poder usar el servidor gráfico sin necesidad de xhost
# La opción init existe porque se van a ejecutar varios procesos dentro del contenedor

run_container() {
    docker container run \
        --name "$nombre" \
        --hostname "$nombre" \
        -d -it --cap-add=NET_ADMIN \
        --init \
        --env DISPLAY \
        --env XAUTHORITY=/cookie \
        --mount type=bind,source="$(pwd)",target=/mnt/shared \
        --mount type=bind,source=/tmp/.X11-unix,target=/tmp/.X11-unix \
        --mount type=bind,source="${Cookiefile}",target=/cookie \
        --label background="$background" \
	"$imagen"
}

# Para cada contenedor, se crea una cookie para acceder al servidor X
# Fuente: https://github.com/mviereck/x11docker/wiki/X-authentication-with-cookies-and-xhost-("No-protocol-specified"-error)#untrusted-cookie-for-container-applications
crear_cookie() {
    if [ ! -d ~/.cookies ]; then
        mkdir ~/.cookies
    fi
    Cookiefile=~/.cookies/"$nombre"_cookie
    :> "$Cookiefile"
    xauth -f "$Cookiefile" generate "$DISPLAY" . untrusted timeout 3600
    Cookie="$(xauth -f "$Cookiefile" nlist "$DISPLAY" | sed -e 's/^..../ffff/')"
    echo "$Cookie" | xauth -f "$Cookiefile" nmerge -
}

# Fuente: http://mywiki.wooledge.org/BashFAQ/050
crear_terminal() {
    primer_nombre=$1
    shift
    first=1
    for nombre; do
        if [ "$first" = 1 ]; then set --; first=0; fi
        set -- "$@" --tab -e "ash -c 'docker container attach $nombre; exec ash'"
    done
    xfce4-terminal -e "ash -c 'docker container attach $primer_nombre; exec ash'" "$@"
}

error() {
    echo "$*" 1>&2
    exit 1
}

uso() {
    echo "Uso: $0 <-c fichero | -l nombre_practica | -i nombre_practica | -p nombre_practica | -r nombre_practica | -h>"
    echo ""
    echo "-c: Crea el entorno de contenedores a partir del fichero YAML proporcionado (se destruyen los contenedores asociados al entorno del fichero)"
    echo "-l: Ejecuta todos los contenedores asociados al entorno proporcionado"
    echo "-i: Muestra la información de todos los contenedores asociados al entorno proporcionado"
    echo "-p: Detiene todos los contenedores asociados al entorno proporcionado"
    echo "-r: Destruye todos los contenedores asociados al entorno proporcionado"
    echo "-h: Muestra este mensaje de ayuda"
}

abortar_creacion() {
    destruir_entorno > /dev/null 2>&1
    error "$1"
}

mostrar_informacion_contenedor() {
    _nombre="$1"
    _imagen=$(docker inspect --format='{{.Config.Image}}' "$_nombre")
    _background=$(docker inspect --format='{{.Config.Labels.background}}' "$_nombre")

    echo "Nombre del contenedor: $_nombre"
    echo "Imagen: $_imagen"
    echo "Background: $_background"
    for _red in $(docker inspect --format='{{range $key, $value := .NetworkSettings.Networks}}{{ println $key}}{{end}}' $_nombre | sed '/^$/d'); do
        _ip=$(docker container inspect -f '{{ (index .NetworkSettings.Networks "'"$_red"'").IPAddress }}' $_nombre)
        echo "Red: $_red IP: $_ip"
    done
    echo ""
}

destruir_entorno() {
    if [ "$(docker ps -aq -f name="$nombre_practica")" ]; then
        echo "Eliminando los contenedores y redes de $nombre_practica..."
        docker container rm -f $(docker ps -a --filter name="$nombre_practica" --format '{{.Names}}') > /dev/null || error "No se ha podido borrar el entorno $nombre_practica"
        if [ "$(docker network ls -q --filter name="$nombre_practica")" ]; then
            docker network rm $(docker network ls -q --filter name="$nombre_practica") > /dev/null || error "No se ha podido borrar el entorno $nombre_practica"
        fi
        echo "Contenedores y redes eliminados exitosamente"
    else
        echo "No existen contenedores asociados a $nombre_practica"
    fi    
}

lanzar_entorno() {
    if [ "$(docker ps -aq -f name="$nombre_practica")" ]; then
        echo "Lanzando de nuevo los contenedores de $nombre_practica..."
        listaContenedoresTerminal=""
        for nombre in $(docker ps -a --filter name="$nombre_practica" --format '{{.Names}}'); do
            docker container start "$nombre" > /dev/null || error "No se ha podido lanzar el contenedor $nombre"

            if [ "$(docker inspect -f '{{.Config.Labels.background}}' "$nombre" )" = "false" ]; then
                if [ -z "$listaContenedoresTerminal" ]; then
                    listaContenedoresTerminal="$nombre"
                else
                    listaContenedoresTerminal="$listaContenedoresTerminal $nombre"
                fi
            fi

            mostrar_informacion_contenedor "$nombre"
        done

        # Se abren las terminales
        crear_terminal $listaContenedoresTerminal

        echo "Contenedores lanzados exitosamente"
    else
        echo "No existen contenedores asociados a $nombre_practica"
    fi    
}

parar_entorno() {
    if [ "$(docker ps -aq -f name="$nombre_practica")" ]; then
        echo "Deteniendo los contenedores de $nombre_practica..."
        docker container stop $(docker ps -a --filter name="$nombre_practica" --format '{{.Names}}') > /dev/null || error "No se ha podido parar el entorno $nombre_practica"
        echo "Contenedores detenidos exitosamente"
    else
        echo "No existen contenedores asociados a $nombre_practica"
    fi    
}

inspeccionar_entorno() {
    if [ "$(docker ps -aq -f name="$nombre_practica")" ]; then
        echo "Información de los contenedores de $nombre_practica"
        docker container inspect $(docker ps -a --filter name="$nombre_practica" --format '{{.Names}}') || error "No se ha podido parar el entorno $nombre_practica"
    else
        echo "No existen contenedores asociados a $nombre_practica"
    fi    
}

crear_entorno() {
    nombre_practica=$(yq e '.nombre_practica' "$fichero")
    if [ "$nombre_practica" = "null" ]; then
        error "Se debe especificar el campo nombre_practica en el fichero"
    fi

    # Si ya existen contenedores asociados a la práctica, se destruyen previamente
    if [ "$(docker ps -aq -f name="$nombre_practica")" ]; then
        destruir_entorno
        echo "" # Se imprime esta línea vacía para dejar más bonita la salida
    fi

    numero_contenedores=$(yq e '.contenedores | length' "$fichero") 
    if [ "$numero_contenedores" = 0 ]; then
        error "Se deben especificar los contenedores deseados en el campo contenedores del fichero"
    fi

    # Primero se crean todos los contenedores
    i=0
    listaContenedoresTerminal=""
    while [ $i -lt "$numero_contenedores" ]
    do
        nombre_contenedor=$(yq e '.contenedores['"$i"'].nombre' "$fichero")
        if [ "$nombre_contenedor" = "null" ]; then
            abortar_creacion "Se debe especificar el campo nombre en contenedores[$i]"
        fi

        # Se concatena el nombre de la practica con el nombre del contenedor para poder hacer operaciones asociadas a entornos de prácticas
        nombre="${nombre_practica}_${nombre_contenedor}"

        imagen=$(yq e '.contenedores['"$i"'].imagen' "$fichero")
        if [ "$imagen" = "null" ]; then
            abortar_creacion "Se debe especificar el campo imagen en contenedores[$i]"
        fi

        background=$(yq e '.contenedores['"$i"'].background' "$fichero")
        if [ "$background" != true ]; then
            background=false
            if [ -z "$listaContenedoresTerminal" ]; then
                listaContenedoresTerminal="$nombre"
            else
                listaContenedoresTerminal="$listaContenedoresTerminal $nombre"
            fi
        fi

        # Se crea la cookie de X11 para usar aplicaciones con GUI
        crear_cookie || abortar_creacion "Error al crear la cookie para $nombre"
        
        # Se crea el contenedor y se le desconecta del adaptador bridge
        run_container > /dev/null || abortar_creacion "Error en la creación del contenedor $nombre"
        docker network disconnect bridge "$nombre" || abortar_creacion "Error en la creación del contenedor $nombre"

        numero_redes=$(yq e '.contenedores['"$i"'].redes | length' "$fichero") 
        if [ "$numero_redes" = 0 ]; then
            abortar_creacion "Se debe especificar el campo redes en contenedores[$i]"
        fi

        j=0
        while [ $j -lt "$numero_redes" ]
        do
            red_fichero=$(yq e '.contenedores['"$i"'].redes['"$j"'].nombre' "$fichero")
            if [ "$red_fichero" = "null" ]; then
                abortar_creacion "Se debe especificar el campo nombre en contenedores[$i].redes[$j]"
            fi
            red="${nombre_practica}_red_${red_fichero}"

            ip=$(yq e '.contenedores['"$i"'].redes['"$j"'].ip' "$fichero")
            if [ "$ip" = "null" ]; then
                subnet_option=""
                ip_option=""
            else
                # Comprueba que la IP introducida esté en formato CIDR válido
                if [ "$(ipcalc -s -n "$ip" )" ]; then
                    direccion_red="$(ipcalc -n "$ip" | sed 's/NETWORK=//')"
                    prefijo="$(ipcalc -p "$ip" | sed 's/PREFIX=//')"
                    subnet_option="--subnet=$direccion_red/$prefijo"

                    ip=$(echo "$ip" | sed 's/\/[0-9]*//')
                    ip_option="--ip $ip"
                else
                    abortar_creacion "Debe proporcionar una IP válida en .contenedores[$i].redes[$j]"
                fi
            fi
            
            # Solo se crea la red si no existe previamente
            docker network inspect "$red" > /dev/null 2>&1 || docker network create $subnet_option "$red" > /dev/null || \
              abortar_creacion "No se ha podido crear la red $red"
            
            docker network connect $ip_option "$red" "$nombre" || abortar_creacion "No se ha podido conectar el contenedor $nombre"

            j=$((j+1))
        done

        mostrar_informacion_contenedor "$nombre"

        i=$((i+1))
    done

    # Si todos los contenedores se han creado exitosamente, se abren las terminales
    if [ -n "$listaContenedoresTerminal" ]; then
        crear_terminal $listaContenedoresTerminal
    fi

    echo "Contenedores creados exitosamente"
}

if [ ! "$(which yq)" ]; then
    error "yq debe estar instalado en el sistema para utilizar este script"
fi

if [ ! "$(which xauth)" ]; then
    error "xauth debe estar instalado en el sistema para utilizar este script"
fi

if [ ! "$(which docker)" ]; then
    error "docker debe estar instalado en el sistema para utilizar este script"
fi

# Controla que el comando se ejecute solo con un flag
opcion=false

while getopts ":c:l:r:p:i:h" o; do
    case "${o}" in
        c)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-r|-p|-i|-h>"
            fichero=${OPTARG}

            if [ ! -f "$fichero" ]; then
                error "El fichero $fichero no existe"
            fi
            comando="crear"
            opcion=true
            ;;
        l)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-r|-p|-i|-h>"
            nombre_practica=${OPTARG}

            if [ -z "$nombre_practica" ]; then
                error "Se debe proporcionar el valor nombre_practica"
            fi

            comando="lanzar"
            opcion=true
            ;;
        r)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-r|-p|-i|-h>"
            nombre_practica=${OPTARG}

            if [ -z "$nombre_practica" ]; then
                error "Se debe proporcionar el valor nombre_practica"
            fi

            comando="destruir"
            opcion=true
            ;;
        p)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-r|-p|-i|-h>"
            nombre_practica=${OPTARG}

            if [ -z "$nombre_practica" ]; then
                error "Se debe proporcionar el valor nombre_practica"
            fi

            comando="parar"
            opcion=true
            ;;
        i)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-r|-p|-i|-h|>"
            nombre_practica=${OPTARG}

            if [ -z "$nombre_practica" ]; then
                error "Se debe proporcionar el valor nombre_practica"
            fi

            comando="inspeccionar"
            opcion=true
            ;;
        h)
            $opcion && error "Solo se puede especificar una opción <-c|-l|-r|-p|-i|-h>"

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

case "$comando" in
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
    inspeccionar)
        inspeccionar_entorno && exit
        ;;
    *)
        uso && exit 1
        ;;
esac