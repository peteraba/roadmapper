export const privacyPolicy = () => {
    const pp = document.getElementById('privacy-policy');

    if (!pp) {
        return;
    }

    if (window.localStorage.getItem('show-privacy-policy') === 'false') {
        $(pp).modal('dispose');
        pp.remove();

        return;
    }

    $(pp).modal('show');

    document.getElementById('privacy-policy-ok').addEventListener('click', event => {
        $(pp).modal('dispose');

        event.preventDefault();
    });

    document.getElementById('privacy-policy-save').addEventListener('click', event => {
        window.localStorage.setItem('show-privacy-policy', 'false');
        $(pp).modal('dispose');

        event.preventDefault();
    });
};
