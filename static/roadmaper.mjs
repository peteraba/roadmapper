import { refreshTimeline, refreshProjects } from './roadmap-dashboard.mjs';
import { roadmapForm } from './roadmap-form.mjs';

const app = () => {
    const $ = window.jQuery || document.querySelectorAll;

    roadmapForm();

    if (!roadmap || !roadmap.Children) {
        $('#roadmap-dashboard, .roadmap-dashboard-link').remove();
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
