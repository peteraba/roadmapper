export const refreshSvg = () => {
    const container = document.getElementById('roadmap-svg'),
        loc = window.location,
        xhttp = new XMLHttpRequest();

    if (!container) {
        return;
    }

    const width = Math.max(container.clientWidth, 800),
        urls = {
            "svg": `${loc.origin}${loc.pathname}/svg?width=${width}`,
            "png": `${loc.origin}${loc.pathname}/png?width=${width}`,
            "jpg": `${loc.origin}${loc.pathname}/jpg?width=${width}`,
            "gif": `${loc.origin}${loc.pathname}/gif?width=${width}`,
            "pdf": `${loc.origin}${loc.pathname}/pdf?width=${width}`,
        },
        downloadButtons = document.querySelectorAll('.roadmap-download-buttons a');

    xhttp.onreadystatechange = function() {
        if (this.readyState === 4 && this.status === 200) {
            container.innerHTML = this.responseText.replace(/mm"/g, '"');
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
