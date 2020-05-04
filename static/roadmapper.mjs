import { formElements, roadmapForm } from './roadmap-form.mjs';
import { refreshSvg } from './roadmap-svg.mjs';
import { privacyPolicy } from './privacy-policy.mjs';

const app = () => {
    roadmapForm(formElements);
    privacyPolicy();

    if (!hasRoadmap) {
        document.querySelectorAll('.roadmap-dashboard-link').forEach(element => element.classList.add('disabled'));
        document.getElementById('roadmap-dashboard').remove();

        return;
    }

    const imgWidth = document.getElementById('img-width'),
        imgWidthEnabled = document.getElementById('img-width-enabled');

    window.addEventListener('resize', () => {
        if (imgWidth.disabled) {
            refreshSvg();
        }
    });

    imgWidth.addEventListener('change', () => {
        refreshSvg();
    });

    imgWidthEnabled.addEventListener('change', () => {
        imgWidth.disabled = !imgWidth.disabled;
    });

    refreshSvg();
};

document.addEventListener("DOMContentLoaded", app);
