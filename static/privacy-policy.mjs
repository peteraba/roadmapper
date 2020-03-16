export const privacyPolicy = () => {
    const pp = document.getElementById('privacy-policy');

    if (window.localStorage.getItem('show-privacy-policy') !== 'false') {
        $(pp).modal('show');
    }

    document.getElementById('privacy-policy-ok').addEventListener('click', event => {
        $(pp).modal('hide');

        event.preventDefault();
    });

    document.getElementById('privacy-policy-save').addEventListener('click', event => {
        window.localStorage.setItem('show-privacy-policy', 'false');
        $(pp).modal('hide');

        event.preventDefault();
    });
};
