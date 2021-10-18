#!/bin/sh

# Uso: ./lanzarContainers.sh n_redes [<maquina:redes>]

# Problemas: 
# 1. Como la creación del contenedor y la conexión a redes no se esperan entre sí, puede generar problemas. Por ejemplo, 
# en esta versión del script los contenedores no se desconectan de bridge porque todavía no se han terminado de crear.
# Si se utilizara una imagen que no estuviera ya pulleada, ni siquiera se conectaría a las redes especificadas debido a la demora.
#
# 2. Se podría modificar la sintaxis de tal forma que se especificara el contenedor que se quiere usar. 
# Por ejemplo, busybox:1 crearía un contenedor de Busybox conectado a la red1.
# 
# 3. Los nombres de las máquinas son del estilo "vm<numero>". Se podría modificar


# Primero eliminamos las redes si existen y luego se crean
for i in $(seq 1 $1); do
    docker network rm red${i} # Por si acaso existe. Es temporal, se tendría que modificar por algo más decente
    docker network create red${i}
done

shift # Se va eliminando el primer argumento sucesivamente para utilizar siempre $1. Referencia: https://stackoverflow.com/questions/3575793/iterate-through-parameters-skipping-the-first
while [ ${#} -gt 0 ]; do

    # Se modifica el IFS para poder trabajar con pares separados por : en los bucles for. Referencia: https://stackoverflow.com/questions/918886/how-do-i-split-a-string-on-a-delimiter-in-bash
    OIFS=$IFS
    IFS=":"

    pos=1
    maquina=''

    for x in $1; do
        # Si se trata del primer elemento, se crea el contenedor. Si es el segundo, se trata de la parte de redes
        if [ $pos -eq 1 ]; then
            maquina=$x # Se almacena el número de contenedor para luego usarlo a la hora de conectarlo a las redes
            docker container rm -f vm${maquina} # Por si acaso existe. Es temporal, se tendría que modificar por algo más decente

            # Se crea una nueva pestaña del terminal para poder interactuar con el contenedor y, cuando finalice, poder seguir usando la pestaña.
            # No es nada portable porque depende del programa de terminal usado, pero suponemos que este script solo se va a usar en la VM de Alpine
            xfce4-terminal --tab -e "ash -c 'docker container run --name vm${maquina} --hostname vm${maquina} -it --cap-add=NET_ADMIN busybox; exec ash'"

            docker network disconnect bridge vm${maquina}
            pos=$(( pos + 1 ))
        else
            # Volvemos a cambiar el IFS para trabajar ahora con pares separados por comas, que indican las redes a las que se debe conectar el contenedor
            OTHER_OIFS=$IFS
            IFS=","
            for red in $x; do
                docker network connect red${red} vm${maquina}
            done
            pos=1
            IFS=OTHER_OIFS
        fi
    done
    IFS=$OIFS
    shift
done
