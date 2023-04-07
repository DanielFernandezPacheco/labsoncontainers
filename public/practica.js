var queryString = window.location.search;
var queryObject = new URLSearchParams(queryString);
var env = queryObject.get("env");

document.addEventListener("DOMContentLoaded", function() {
    var envName = document.getElementById("env-name");
    envName.innerHTML = env;
    var container = document.getElementById("config-container");
    container.appendChild(createLabControlItem(env + ".yaml"));
    fetch("http://localhost:8080/getContainers?env=" + env)
        .then(response => response.json())
        .then(data => generarDetalles(data));
  });


function MostrarContenedor(row, container) {
    const hostsCell = row.insertCell();
    const hostElement = document.createElement('div');
    
    hostElement.style.borderRadius = container.Networks.length > 1? '50%' : '5%';
    hostElement.style.display = 'inline-block';
    hostElement.style.marginRight = '10px';
    hostElement.style.textAlign = 'center';
    hostElement.style.cursor = 'pointer';
    hostElement.classList.add('alert', 'alert-success');
    hostElement.innerHTML = `<div><strong>${container.Name.replace(env + "_", '')}</strong><br>${container.Image} </div><div></div>`;
    hostElement.id = container.ID;
    container.Networks.forEach(network => {
        hostElement.innerHTML += `<div><strong>${network.Name.replace(env + "_", 'Red ')}</strong><br>${network.IP}</div>`;
    });

    hostElement.innerHTML += `</div>`;
    hostsCell.appendChild(hostElement);

    row.insertCell();

    // Agregar event listener
    hostElement.addEventListener('click', function() {
        if (hostElement.classList.contains('alert-success')) {
            if (confirm('¿Desea parar el contenedor?')) {
                // Cambiar clase a alert-danger
                hostElement.classList.remove('alert-success');
                hostElement.classList.add('alert-danger');
                fetch("http://localhost:8080/stopContainer?container=" + container.ID)
                .then(response => response.json())
                .then(data => console.log(data));

            }
        } else if (hostElement.classList.contains('alert-danger')) {
            if (confirm('¿Desea volver a iniciar el contenedor?')) {
                // Cambiar clase a alert-success
                hostElement.classList.remove('alert-danger');
                hostElement.classList.add('alert-success');
                fetch("http://localhost:8080/startContainer?" + env + "&container=" + container.ID)
                .then(response => response.json())
                .then(data => console.log(data));
            }
        }
    });
}

function generarDetalles(containers) {
    var networkMembers = {};


    containers.forEach(container => {
        //fetchImgOS(container);

        container.Networks.forEach(network => {
            if (!networkMembers[network.Name]) {
                networkMembers[network.Name] = {
                    hosts: [],
                    routers: []
                };
            }

            if (container.Networks.length > 1) {
                networkMembers[network.Name].routers.push(container);
            } else {
                networkMembers[network.Name].hosts.push(container);
            }
        });
    });

    let orderedNetworks = Object.entries(networkMembers);

    // Ordenar el arreglo de pares por nombre de red
    orderedNetworks.sort((a, b) => {
        if (a[0] < b[0]) return -1;
        if (a[0] > b[0]) return 1;
        return 0;
    });

    // Crear un nuevo objeto a partir de los pares ordenados
    const orderedNetworkData = {};
    for (const [networkName, networkData] of orderedNetworks) {
        orderedNetworkData[networkName] = networkData;
    }

    networkMembers = { ...orderedNetworkData };

    const table = document.createElement('table');
    table.classList.add('mt-3', 'ml-2', 'table-sm');
    table.style.fontSize = '8pt';
    const row = table.insertRow();
    const routersMostrados = [];

    for (let red in networkMembers) {
        const hosts = networkMembers[red].hosts;
        if (hosts.length > 0) {
            for (let host of hosts) {
                MostrarContenedor(row, host);
            }
        }

        const routers = networkMembers[red].routers;
        if (routers.length > 0) {
            for (let router of routers) {
                if (!routersMostrados.includes(router.ID)) {
                    routersMostrados.push(router.ID);
                    MostrarContenedor(row, router);
                }
            }
        }
    }

    document.getElementById("labContainer").appendChild(table);
}


