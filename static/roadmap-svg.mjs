export const refreshSvg = () => {
    const container = document.getElementById('roadmap-svg'),
        imgWidth = document.getElementById('img-width'),
        imgWidthEnabled = document.getElementById('img-width-enabled'),
        loc = window.location,
        xhttp = new XMLHttpRequest(),
        widthMin = 800,
        widthMax = 30000;

    if (!container) {
        return;
    }

    let width = Math.min(Math.max(container.clientWidth, widthMin), widthMax);

    if (imgWidthEnabled.checked) {
        width = Math.min(Math.max(imgWidth.value, widthMin), widthMax);
    } else {
        imgWidth.value = width;
    }

    const urls = {
            "svg": `${loc.origin}${loc.pathname}/svg?width=${width}`,
            "png": `${loc.origin}${loc.pathname}/png?width=${width}`,
            "jpg": `${loc.origin}${loc.pathname}/jpg?width=${width}`,
            "gif": `${loc.origin}${loc.pathname}/gif?width=${width}`,
            "pdf": `${loc.origin}${loc.pathname}/pdf?width=${width}`,
        },
        downloadButtons = document.querySelectorAll('.roadmap-download-buttons a')

    xhttp.onreadystatechange = function() {
        if (this.readyState === 4 && this.status === 200) {
            container.innerHTML = this.responseText.replace(/mm"/g, '"');

            updateSvg();
        }
    };
    xhttp.open("GET", urls.svg, true);
    xhttp.send();

    downloadButtons.forEach(elem => {
        if (elem.dataset.fileformat && urls[elem.dataset.fileformat]) {
            elem.href = urls[elem.dataset.fileformat];
        }
    });
};

const updateSvg = () => {
    if (!projectUrls) {
        return;
    }

    document.querySelectorAll('#roadmap-svg > svg > text > tspan').forEach(e => {
        if (typeof projectUrls[e.innerHTML] === "undefined") {
            return;
        }

        let html = e.innerHTML;
        projectUrls[e.innerHTML].forEach(u => {
            if (u.substr(0, 4) !== 'http') {
                u = `${baseUrl}${u}`;
            }

            html += `&nbsp;<a href="${u}" target="_blank" style="text-decoration: none;">&#x2B08;</a>`;
        });

        e.innerHTML = html;
    });
};
