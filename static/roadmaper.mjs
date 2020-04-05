import { roadmapForm } from './roadmap-form.mjs';
import { refreshSvg } from './roadmap-svg.mjs';
import { privacyPolicy } from './privacy-policy.mjs';

const app = () => {
    roadmapForm();
    privacyPolicy();

    if (!hasRoadmap) {
        document.querySelectorAll('.roadmap-dashboard-link').forEach(element => element.classList.add('disabled'));
        document.getElementById('roadmap-dashboard').remove();

        return;
    }

    window.addEventListener('resize', () => {
        refreshSvg();
    });

    refreshSvg();
};

document.addEventListener("DOMContentLoaded", app);
