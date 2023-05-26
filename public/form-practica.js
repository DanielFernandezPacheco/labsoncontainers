var contador = 0;
// Obtener los parámetros de la URL
const urlParams = new URLSearchParams(window.location.search);

const file = urlParams.has('file') ? urlParams.get('file') : null;

actualizarTexto();

function actualizarTexto() {
    document.getElementById("tituloText").innerHTML = file ? "Editar práctica" : "Nueva práctica";
    document.getElementById("submitText").innerHTML = file ? "Actualizar configuración" : "Crear nueva configuración";

    if (file) {
        var inputNombre = document.getElementById("nombre-practica");
        inputNombre.value = file.replace(/\.yaml$/, "");
        inputNombre.readOnly = true;
        fillForm();
    }
}


function fillForm() {
    fetch(`/getConfig?file=${file}`)
        .then(response => response.json())
        .then(config => {
            // fill the form with the config data
            document.getElementById("num-contenedores").value = config.Contenedores.length;

            config.Contenedores.forEach((contenedor, i) => {
                nuevoContenedor();
                document.getElementById(`nombre-contenedor-${i}`).value = contenedor.Nombre;
                document.getElementById(`nombre-imagen-${i}`).value = contenedor.Imagen;
                actualizarNombreContenedor(i);
                contenedor.Redes.forEach((red, j) => {
                    nuevaRed(i);
                    document.getElementById(`nombre-red-${i}-${j}`).value = red.Nombre;
                    document.getElementById(`ip-red-${i}-${j}`).value = red.IP;
                });
            });
        })
        .catch(error => console.error(error));
}


function actualizarNombrePractica() {
    document.querySelector("#main-card-title").textContent = document.getElementById("nombre-practica").value;
}


function nuevaRed(indiceContenedor) {
    var padre = document.getElementById("campos-redes-" + indiceContenedor);
    var newNetwork = document.createElement("div");
    newNetwork.id = "red-" + indiceContenedor + "-" + padre.children.length
    newNetwork.className = "input-group mb-2";
    newNetwork.innerHTML = `
            <div class="input-group-prepend">
                <span class="input-group-text">Red</span>
            </div>
            <input required type="number" style="max-width:125px"
                class="form-control" placeholder="Nº de red" id="nombre-${newNetwork.id}" name="nombre-${newNetwork.id}" max="10">
            <input type="text" class="form-control" placeholder="IP (opcional)" id="ip-${newNetwork.id}" name="ip-${newNetwork.id}"  maxlength="20">

            <div class="input-group-append">
                <button class="btn btn-outline-danger" type="button"
                    onclick="eliminarRed('${newNetwork.id}', ${indiceContenedor})">Eliminar</button>
            </div>`;
    padre.appendChild(newNetwork);
    actualizarNumRedes(indiceContenedor);
}


function actualizarNumRedes(indiceContenedor) {
    var redes = document.getElementById("campos-redes-" + indiceContenedor).children;
    document.getElementById("num-redes-" + indiceContenedor).value = redes.length;
}


function actualizarRedes(indiceContenedor) {
    var hijos = document.getElementById("campos-redes-" + indiceContenedor).children;
    for (var i = 0; i < hijos.length; i++) {
        hijos[i].id = "red-" + indiceContenedor + "-" + i;
        var inputName = hijos[i].getElementsByTagName("input")[0];
        var inputIP = hijos[i].getElementsByTagName("input")[1];
        inputName.id = "nombre-" + hijos[i].id;
        inputName.name = "nombre-" + hijos[i].id;
        inputIP.id = "ip-" + hijos[i].id;
        inputIP.name = "ip-" + hijos[i].id;
        var removeBtn = hijos[i].getElementsByTagName("button")[0];
        removeBtn.onclick = eliminarRed.arguments(inputName.id, indiceContenedor);
    }
}


function eliminarRed(idRed, indiceContenedor) {
    document.getElementById(idRed).remove();
    actualizarRedes(indiceContenedor);
    actualizarNumRedes(indiceContenedor);
}


function nuevoContenedor() {
    var padre = document.getElementById("accordion");
    var newCard = document.createElement("div");
    newCard.className = "card"
    newCard.id = "card-" + contador
    newCard.innerHTML = `
        <div class="card-header d-flex">
            <button class="btn btn-link btn-block text-left flex-fill mr-6" type="button" data-toggle="collapse" data-target="#collapse-${contador}" aria-expanded="true" aria-controls="collapse-${contador}">
                Contenedor nuevo
            </button>
            <button type="button" class="btn btn-outline-danger flex-fill" id="eliminar-${contador}" onclick="eliminarContenedor(${contador})">Eliminar</button>
        
        </div>
        <div id="collapse-${contador}" class="collapse" data-parent="#accordion">
        <div class="card-body">
            <div class ="row">
                <div class='form-group col-6'>
                    <label for='nombre-contenedor-${contador}'>Nombre del contenedor:</label>
                    <input required class='form-control' type='text' id='nombre-contenedor-${contador}' name='nombre-contenedor-${contador}' onchange="actualizarNombreContenedor(${contador})"  maxlength="20">
                </div>
                <div class='form-group col-6'>
                    <label for='nombre-imagen-${contador}'>Nombre de la imagen:</label>
                    <input required class='form-control' type='text' id='nombre-imagen-${contador}' name='nombre-imagen-${contador}'  maxlength="30">
                </div>
            </div>
            <div class='form-group'>
                <button type="button" class="btn btn-dark btn-block" onclick="nuevaRed(${contador})">Añadir red</button>
            </div>
            
            <input id="num-redes-${contador}" name="num-redes-${contador}" type="hidden" value="0" max="10">
            <div id='campos-redes-${contador}'></div>
        </div>
        </div>`

    padre.appendChild(newCard);

    contador++;
    actualizarNumContenedores()
}


function actualizarNombreContenedor(indiceContenedor) {
    document.querySelector("button[data-target='#collapse-" + indiceContenedor + "']").textContent = document.getElementById("nombre-contenedor-" + indiceContenedor).value;
}


function actualizarNumContenedores() {
    document.getElementById("num-contenedores").value = contador;
}


function eliminarContenedor(indiceContenedor) {
    var node = document.getElementById("card-" + indiceContenedor);
    var parent = document.getElementById("accordion")
    parent.removeChild(node);
    contador--;
    actualizarContenedores(indiceContenedor);
    actualizarNumContenedores();
}


function actualizarContenedores(indiceContenedor) {
    const accordion = document.getElementById("accordion");

    const tarjetas = accordion.querySelectorAll('.card');

    for (let i = indiceContenedor; i < tarjetas.length; i++) {
        const tarjeta = tarjetas[i]
        tarjeta.id = `card-${i}`;

        const collapse = tarjetas[i].querySelector(`[data-target^="#collapse-"]`);
        collapse.dataset.target = `#collapse-${i}`;
        collapse.setAttribute(`aria-controls`, `collapse-${i}`);

        const contenido = tarjetas[i].querySelector(`[id^="collapse-"]`);
        contenido.id = `collapse-${i}`;
        contenido.dataset.parent = "#accordion";

        const nombreContenedor = tarjetas[i].querySelector(`[for^="nombre-contenedor-"]`);
        nombreContenedor.htmlFor = `nombre-contenedor-${i}`;

        const inputNombreContenedor = tarjetas[i].querySelector(`[id^="nombre-contenedor-"]`);
        inputNombreContenedor.id = `nombre-contenedor-${i}`;
        inputNombreContenedor.name = `nombre-contenedor-${i}`;
        inputNombreContenedor.onchange = actualizarNombreContenedor.bind(null, i);

        const inputNombreImagen = tarjetas[i].querySelector(`[id^="nombre-imagen-"]`);
        inputNombreImagen.id = `nombre-imagen-${i}`;
        inputNombreImagen.name = `nombre-imagen-${i}`;

        const divCamposRedes = tarjetas[i].querySelector(`[id^="campos-redes-"]`);
        divCamposRedes.id = `campos-redes-${i}`;
        divCamposRedes.name = `campos-redes-${i}`;

        const eliminarContenedorBtn = tarjetas[i].querySelector(`[onclick^="eliminarContenedor("]`);
        eliminarContenedorBtn.onclick = eliminarContenedor.bind(null, i);

        const nuevaRedBtn = tarjetas[i].querySelector(`[onclick^="nuevaRed("]`);
        nuevaRedBtn.onclick = nuevaRed.bind(null, i);

        actualizarRedes(i);
    }
}