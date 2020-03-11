import { roadmapForm } from './roadmap-form.mjs';
import { refreshSvg } from './roadmap-svg.mjs';

const app = () => {
    roadmapForm();

    if (!roadmap || !roadmap.Children) {
        document.querySelectorAll('.roadmap-dashboard-link').forEach(element => element.classList.add('disabled'));
        document.getElementById('roadmap-dashboard').remove();

        return;
    }
    
    window.addEventListener('resize', () => {
        refreshSvg();
    });

    refreshSvg();

    const tt = $('[data-toggle="tooltip"]');
    if (tt && tt['tooltip'] && typeof tt['tooltip'] === 'function') {
        tt['tooltip']();
    }
};

document.addEventListener("DOMContentLoaded", app);
