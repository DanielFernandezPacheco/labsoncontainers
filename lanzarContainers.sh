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
    echo "Error en la creación de los contenedores"
    exit
}

# Primero eliminamos las redes si existen y luego se crean
for i in $(seq 1 $1); do
    docker network rm red${i} # Por si acaso existe. Es temporal, se tendría que modificar por algo más decente
    docker network create red${i}
done

# Creamos la cookie de X11 para usar aplicaciones con GUI
Cookiefile=~/containercookie
:> $Cookiefile
xauth -f $Cookiefile generate $DISPLAY . untrusted timeout 3600
Cookie="$(xauth -f $Cookiefile nlist $DISPLAY | sed -e 's/^..../ffff/')"
echo "$Cookie" | xauth -f "$Cookiefile" nmerge -

shift # Se va eliminando el primer argumento sucesivamente para utilizar siempre $1. Referencia: https://stackoverflow.com/questions/3575793/iterate-through-parameters-skipping-the-first
while [ ${#} -gt 0 ]; do

    # Se modifica el IFS para poder trabajar con pares separados por : en los bucles for. Referencia: https://stackoverflow.com/questions/918886/how-do-i-split-a-string-on-a-delimiter-in-bash
    OIFS=$IFS
    IFS=":"

    pos=1
    nombre=''

    for x in $1; do
        # Se almacena el nombre del contenedor
        if [ $pos -eq 1 ]; then
            nombre=$x
            docker container rm -f $nombre # Por si acaso existe. Es temporal, se tendría que modificar por algo más decente

            pos=$(( pos + 1 ))
        # Se lee la imagen deseada y se crea el contenedor
        elif [ $pos -eq 2 ]; then
            imagen=$x

            run_container || error

            docker network disconnect bridge $nombre

            # Se crea una nueva pestaña del terminal para poder interactuar con el contenedor y, cuando finalice, poder seguir usando la pestaña.
            # No es nada portable porque depende del programa de terminal usado, pero suponemos que este script solo se va a usar en la VM de Alpine
            xfce4-terminal --tab -e "ash -c 'docker container attach $nombre; exec ash'"
            
            pos=$(( pos + 1 ))
        # Se conecta el contenedor a las redes especificadas
        else
            # Volvemos a cambiar el IFS para trabajar ahora con pares separados por comas, que indican las redes a las que se debe conectar el contenedor
            OTHER_OIFS=$IFS
            IFS=","
            for red in $x; do
                docker network connect red${red} $nombre
            done
            pos=1
            IFS=OTHER_OIFS
        fi
    done
    IFS=$OIFS
    shift
done
