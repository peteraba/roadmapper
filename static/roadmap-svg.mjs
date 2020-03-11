export const refreshSvg = () => {
    const container = document.getElementById('roadmap-svg'), loc = window.location, xhttp = new XMLHttpRequest();

    if (!container) {
        return;
    }

    const width = Math.max(container.clientWidth, 800), url = `${loc.origin}${loc.pathname}/svg?width=${width}`;

    xhttp.onreadystatechange = function() {
        if (this.readyState === 4 && this.status === 200) {
            container.innerHTML = this.responseText;
        }
    };
    xhttp.open("GET", url, true);
    xhttp.send();


};
