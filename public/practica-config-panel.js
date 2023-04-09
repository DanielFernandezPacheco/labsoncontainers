function createLabButton(env, url, color, svg, description, reload = false) {
    var button = document.createElement("button");
    button.title = description;
    button.className = "btn btn-" + color + " d-flex w-100";

    button.innerHTML = svg;
    button.onclick = function () {
        fetch("http://localhost:8080/" + url + "?env=" + env)
            .then(response => {
                if (reload && response.ok) {
                    window.location.replace("practica.html?env=" + env.replace(".yaml", ""));
                }
            })
            .catch(error => {
                console.error(error);
            });
    };

    return button;
}

function createLabControlItem(file) {
    var item = document.createElement("div");
    var iconSize = "10pt";
    item.className = "input-group mb-2 d-flex flex-column border-1";

    var nameWithoutExtension = file.substring(0, file.indexOf('.'));
    // create input text
    var inputText = document.createElement("input");
    inputText.type = "text";
    inputText.className = "form-control border-0 rounded-0";
    inputText.readOnly = true;
    inputText.value = file;
    inputText.style.height = "30px";
    inputText.style.backgroundColor = "rgb(249, 249, 249)";

    // create input group append
    var btnGroupDiv = document.createElement("div");
    btnGroupDiv.className = "input-group-append d-flex flex-fill";

    var inputTextDiv = document.createElement("div");
    inputTextDiv.className = "d-flex flex-fill";

    // create button
    var createButton = createLabButton(file, "create", "light",
        `<svg id="svg" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="${iconSize}" height="${iconSize}" viewBox="0, 0, 400,400"><g id="svgg"><path id="path0" d="M188.396 1.519 C 179.017 4.986,172.186 11.185,167.746 20.258 L 165.234 25.391 165.023 95.101 L 164.811 164.811 95.101 165.023 L 25.391 165.234 19.922 167.906 C -2.004 178.620,-7.043 207.925,10.094 225.062 C 12.766 227.735,16.569 230.455,19.922 232.094 L 25.391 234.766 95.101 234.977 L 164.811 235.189 165.023 304.899 L 165.234 374.609 167.906 380.078 C 169.545 383.431,172.265 387.234,174.938 389.906 C 192.075 407.043,221.380 402.004,232.094 380.078 L 234.766 374.609 234.977 304.899 L 235.189 235.189 304.899 234.977 L 374.609 234.766 380.078 232.094 C 383.431 230.455,387.234 227.735,389.906 225.062 C 407.043 207.925,402.004 178.620,380.078 167.906 L 374.609 165.234 304.899 165.023 L 235.189 164.811 234.977 95.101 L 234.766 25.391 232.094 19.922 C 224.320 4.012,204.376 -4.387,188.396 1.519 " stroke="none" fill="#07a" fill-rule="evenodd"></path></g></svg>`
        , "Crea y ejecuta el entorno", true);

    // exec button
    var executeButton = createLabButton(nameWithoutExtension, "execute", "light",
        `<svg id="svg" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="${iconSize}" height="${iconSize}" viewBox="0, 0, 400,400"><g id="svgg"><path id="path0" d="M45.703 1.049 C 31.847 5.649,23.172 17.089,20.349 34.484 C 19.233 41.366,19.233 358.634,20.349 365.516 C 25.373 396.477,51.261 408.735,80.605 394.049 C 85.205 391.746,348.424 240.867,355.914 236.239 C 388.455 216.133,388.455 183.867,355.914 163.761 C 349.251 159.644,85.731 8.536,81.050 6.148 C 70.568 0.801,53.757 -1.625,45.703 1.049 " stroke="none" fill="#015e36" fill-rule="evenodd"></path></g></svg>`
        , "Ejecuta el entorno guardado", true);


    // stop button
    var stopButton = createLabButton(nameWithoutExtension, "stop", "light", `
    <svg id="svg" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="${iconSize}" height="${iconSize}" viewBox="0, 0, 400,400"><g id="svgg"><path id="path0" d="M38.095 14.017 C 26.775 17.636,17.071 27.593,13.739 39.009 C 11.642 46.193,11.873 354.876,13.980 361.648 C 19.802 380.355,32.269 387.500,59.086 387.500 L 74.919 387.500 75.172 328.320 C 75.419 270.395,75.459 269.037,77.083 264.259 C 81.184 252.188,90.624 243.010,102.734 239.319 C 110.072 237.083,291.047 237.108,297.656 239.346 C 309.831 243.469,318.888 252.402,322.917 264.259 C 324.541 269.037,324.581 270.395,324.828 328.320 L 325.081 387.500 340.987 387.500 C 367.924 387.500,380.708 380.015,386.261 360.991 C 388.052 354.858,388.165 97.267,386.379 93.907 C 385.025 91.358,308.141 14.589,305.715 13.363 C 304.653 12.826,298.546 12.500,289.541 12.500 L 275.072 12.500 274.805 64.648 C 274.607 103.235,274.287 117.544,273.576 119.670 C 271.146 126.930,264.002 133.923,256.768 136.121 C 252.067 137.550,147.936 137.551,143.236 136.122 C 135.869 133.883,128.898 127.062,126.424 119.670 C 125.713 117.544,125.393 103.235,125.195 64.648 L 124.928 12.500 83.753 12.542 C 46.038 12.580,42.201 12.704,38.095 14.017 M150.000 62.500 L 150.000 112.500 200.000 112.500 L 250.000 112.500 250.000 62.500 L 250.000 12.500 200.000 12.500 L 150.000 12.500 150.000 62.500 M105.657 264.058 C 99.653 267.719,100.006 263.657,100.003 329.102 L 100.000 387.500 200.000 387.500 L 300.000 387.500 299.997 329.102 C 299.994 263.657,300.347 267.719,294.343 264.058 C 290.342 261.619,109.658 261.619,105.657 264.058 " stroke="none" fill="#424242" fill-rule="evenodd"></path></g></svg>`
        , "Cierra y guarda el estado del entorno");

    // info button
    var infoButton = createLabButton(nameWithoutExtension, "inspect", "light",
        `<svg id="svg" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="${iconSize}" height="${iconSize}" viewBox="0, 0, 400,400"><g id="svgg"><path id="path0" d="M225.000 1.250 C 190.212 10.535,176.547 50.906,199.585 76.335 C 230.420 110.371,288.615 90.855,288.447 46.534 C 288.329 15.549,256.511 -7.161,225.000 1.250 M180.078 135.642 C 166.375 137.714,148.871 142.881,127.734 151.094 L 116.797 155.343 113.992 166.930 C 112.450 173.302,111.285 178.640,111.404 178.792 C 111.523 178.944,115.422 177.834,120.068 176.326 C 143.687 168.658,162.276 170.581,167.719 181.257 C 173.987 193.550,171.946 206.481,153.293 272.656 C 147.298 293.926,141.573 315.547,140.570 320.703 C 130.524 372.372,151.845 399.418,202.734 399.561 C 225.644 399.625,228.647 398.933,263.760 385.494 L 277.911 380.078 280.737 368.409 C 284.027 354.826,284.213 355.867,278.860 357.908 C 254.550 367.177,230.380 364.795,225.296 352.630 C 220.120 340.240,222.311 326.973,239.792 264.844 C 253.053 217.714,252.281 220.668,253.586 212.109 C 260.077 169.532,247.672 145.635,214.844 137.475 C 208.338 135.858,186.326 134.697,180.078 135.642 " stroke="none" fill="#102b8f" fill-rule="evenodd"></path></g></svg>`
        , "Muestra la información de todos los contenedores por consola");

    // destroy button
    var destroyButton = createLabButton(nameWithoutExtension, "destroy", "light",
        `<svg id="svg" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="${iconSize}" height="${iconSize}" viewBox="0, 0, 400,400"><g id="svgg"><path id="path0" d="M68.359 1.053 C 60.621 3.612,58.129 5.721,31.885 31.925 C 0.156 63.606,0.391 63.262,0.391 78.125 C 0.391 93.782,-3.110 89.357,56.470 149.023 L 107.372 200.000 56.470 250.977 C -3.110 310.643,0.391 306.218,0.391 321.875 C 0.391 336.730,0.157 336.387,31.885 368.115 C 63.613 399.843,63.270 399.609,78.125 399.609 C 93.782 399.609,89.357 403.110,149.023 343.530 L 200.000 292.628 250.977 343.530 C 310.643 403.110,306.218 399.609,321.875 399.609 C 336.730 399.609,336.387 399.843,368.115 368.115 C 399.843 336.387,399.609 336.730,399.609 321.875 C 399.609 306.218,403.110 310.643,343.530 250.977 L 292.628 200.000 343.530 149.023 C 403.110 89.357,399.609 93.782,399.609 78.125 C 399.609 63.270,399.843 63.613,368.115 31.885 C 336.387 0.157,336.730 0.391,321.875 0.391 C 306.218 0.391,310.643 -3.110,250.977 56.470 L 200.000 107.372 149.023 56.470 C 90.104 -2.363,93.787 0.608,79.297 0.226 C 74.373 0.096,70.329 0.402,68.359 1.053 " stroke="none" fill="#8f1023" fill-rule="evenodd"></path></g></svg>`
        , "Destruye el entorno");

    // reset button
    var resetButton = createLabButton(nameWithoutExtension, "reset", "light",
        `<svg id="svg" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="${iconSize}" height="${iconSize}" viewBox="0, 0, 400,400"><g id="svgg"><path id="path0" d="M61.507 33.805 C 57.458 35.125,53.428 38.530,51.543 42.226 C 50.003 45.243,49.970 46.484,50.178 92.515 L 50.391 139.710 52.734 143.037 C 54.023 144.867,56.628 147.182,58.521 148.182 L 61.965 150.000 108.357 150.000 C 160.528 150.000,158.024 150.272,163.156 144.046 C 169.841 135.935,168.234 127.452,158.049 117.096 L 152.196 111.145 157.934 108.773 C 244.103 73.157,329.719 161.186,288.625 243.146 C 247.708 324.752,128.875 315.268,100.414 228.125 C 91.118 199.662,62.669 190.630,43.132 209.940 C 29.437 223.474,30.173 241.186,45.733 272.541 C 104.175 390.309,267.519 399.920,339.830 289.844 C 390.001 213.470,368.049 111.121,290.563 60.149 C 237.316 25.121,165.386 23.990,110.191 57.311 L 102.803 61.771 89.878 48.932 C 74.020 33.179,70.002 31.037,61.507 33.805 " stroke="none" fill="#3a0f52" fill-rule="evenodd"></path></g></svg>`
        , "Resetea y vuelve a lanzar el entorno", true);

    // edit button
    var editButton = createLabButton(nameWithoutExtension, null, "light",
        `<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="${iconSize}" height="${iconSize}" viewBox="0 0 1000 1000" enable-background="new 0 0 1000 1000" xml:space="preserve"><g><g transform="translate(0.000000,512.000000) scale(0.100000,-0.100000)"><path fill="#b86c02" fill-rule="evenodd" d="M8155.1,4994c-166.8-55.6-247.4-122.7-747.8-625.1l-517.7-519.6l926.1-924.2L8740,1999l525.4,525.4c615.5,615.5,634.7,642.3,634.7,922.3c0,281.9-17.3,306.8-648.1,933.8c-462.1,458.3-556.1,542.6-646.2,581C8477.3,5015.1,8266.4,5030.4,8155.1,4994z"/><path d="M3647.3,605L1016.6-2027.6L562.1-3389c-249.3-747.8-458.3-1371-462.1-1384.4c-5.7-11.5,615.5,186,1380.5,441l1392.1,464l2630.7,2630.7L8136,1395l-920.4,920.4c-506.2,506.2-924.2,920.4-930,920.4C6279.9,3235.7,5093,2050.8,3647.3,605z" fill="#b86c02" fill-rule="evenodd"/></g></g></svg>`
        , "Edita la configuración de la práctica");

    editButton.onclick = function () {
        window.open('./form-practica.html?file=' + file, '_blank');
    };


    btnGroupDiv.appendChild(createButton);
    btnGroupDiv.appendChild(executeButton);
    btnGroupDiv.appendChild(stopButton);
    btnGroupDiv.appendChild(infoButton);
    btnGroupDiv.appendChild(destroyButton);
    btnGroupDiv.appendChild(resetButton);
    btnGroupDiv.appendChild(editButton);
    inputTextDiv.appendChild(inputText);

    item.appendChild(btnGroupDiv);
    item.appendChild(inputTextDiv);

    return item;
}