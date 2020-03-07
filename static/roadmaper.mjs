import { refreshTimeline, refreshProjects } from './roadmap-dashboard.mjs';
import { roadmapForm } from './roadmap-form.mjs';

const app = () => {
    roadmapForm();

    if (!roadmap || !roadmap.Children) {
        document.querySelectorAll('.roadmap-dashboard-link').forEach(element => element.classList.add('disabled'));
        document.getElementById('roadmap-dashboard').remove();
    } else {
        refreshProjects();
    }

    window.addEventListener('resize', () => {
        refreshTimeline();
    });

    const tt = $('[data-toggle="tooltip"]');
    if (tt && tt['tooltip'] && typeof tt['tooltip'] === 'function') {
        tt['tooltip']();
    }
};

document.addEventListener("DOMContentLoaded", app);
