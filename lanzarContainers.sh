#!/bin/sh

# Uso: ./lanzarContainers.sh n_redes [<maquina:redes>]

for i in $(seq 1 $1); do
    # Habría que comprobar si la red ya está creada
    docker network create red${i}
done

for i in "${@:2}" do
    tmp=$(echo $i | tr ":" "\n") 
    maquina=$tmp[0]
    redes=$(echo $tmp[1] | tr "," "\n")

    xfce4-terminal --tab -e "ash -c 'docker container run --name vm${maquina} --network red${redes[0]} -it busybox; exec ash'"
    
    if [${#redes[@]} -ne 1] then
        for ((i=0; i<${#redes[@]}; i++)); do
            docker network connect red${i} vm${maquina}
        done
    fi
done


# xfce4-terminal --tab -e "ash -c 'docker container run -it busybox; exec ash'"
