export const refreshTimeline = () => {
    const drawing = document.getElementById('drawing');

    if (!drawing || !drawing.dataset || !drawing.dataset.start || !drawing.dataset.end) {
        return;
    }

    const h = 75,
        projectStart = new Date(drawing.dataset.start),
        projectEnd = new Date(drawing.dataset.end),
        lineColor = window.getComputedStyle(drawing, null).borderColor,
        textColor = window.getComputedStyle(drawing, null).color;

    let line, text, draw;

    drawing.innerHTML = '';

    window.setTimeout(() => {
        const w = drawing.offsetWidth,
            today = new Date(),
            pink = '#f88',
            red = '#f33',
            left = (today.getTime() - projectStart.getTime()) / (projectEnd.getTime() - projectStart.getTime()) * w;

        draw = SVG().addTo(drawing).size(w, h);

        // baseline
        line = draw.line(0, 0, w, 0).move(0, 40);
        line.stroke({color: lineColor, width: 2});

        // first day
        line = draw.line(0, 0, 0, 10).move(10, 35);
        line.stroke({color: textColor, width: 2});
        text = draw.text(projectStart.toLocaleDateString());
        text.move(5, 10).font({});

        // last day
        line = draw.line(0, 0, 0, 10).move(w - 10, 35);
        line.stroke({color: textColor, width: 2});
        text = draw.text(projectEnd.toLocaleDateString());
        text.move(w - 5, 10).font({anchor: 'end'});

        // today
        if (projectStart.getTime() < today.getTime() && projectEnd.getTime() > today.getTime()) {
            line = draw.line(0, 0, 0, h).move(left, 0);
            line.stroke({color: pink, width: 1});
            text = draw.text(today.toLocaleDateString());
            text.move(left, 45).font({fill: red, anchor: 'middle'});
        }
    }, 0);
};

export const refreshProjects = () => {
    const table = document.getElementById('roadmap'),
        control = document.getElementById('control');

    table.innerHTML = '';

    const buildTable = (p, container) => {
        const
            thead = document.createElement('thead'),
            tbody = document.createElement('tbody'),
            projectStart = new Date(p.Dates.Start),
            projectEnd = new Date(p.Dates.End),
            fullDiff = projectEnd.getTime() - projectStart.getTime();

        container.appendChild(thead);
        container.appendChild(tbody);

        displayHeader(p, thead, p.Dates.Start, p.Dates.End);

        if (p.Children !== null) {
            p.Children.forEach(c => displayProject(c, tbody, 1, projectStart, fullDiff));
        }
    };

    const displayHeader = (p, container, projectStart, projectEnd) => {
        const row = document.createElement('tr'),
            left = document.createElement('th'),
            right = document.createElement('th'),
            roadmapTitle = document.createElement('p'),
            drawing = document.createElement('div');

        roadmapTitle.innerHTML = '&nbsp;';

        row.appendChild(left);
        row.appendChild(right);

        drawing.id = 'drawing';
        drawing.dataset.start = projectStart;
        drawing.dataset.end = projectEnd;

        right.append(drawing);
        right.classList.add('timeline');

        left.appendChild(roadmapTitle);

        container.appendChild(row);

        refreshTimeline();
    };

    const displayProject = (p, container, level, projectStart, fullDiff) => {
        const row = document.createElement('tr'),
            left = document.createElement('th'),
            right = document.createElement('td'),
            projectTitle = document.createElement('span'),
            projectContainer = document.createElement('div'),
            projectBar = document.createElement('div'),
            start = new Date(p.Dates.Start),
            end = new Date(p.Dates.End),
            diff = end.getTime() - start.getTime(),
            durationDays = diff / 86400000,
            w = diff / fullDiff * 100,
            l = (start.getTime() - projectStart.getTime()) / fullDiff * 100,
            tooltip = `${p.Percentage}%, ${start.toLocaleDateString()} - ${end.toLocaleDateString()}, ${durationDays} days`;

        let nbsp, a, i;

        projectTitle.textContent = p.Title;
        projectTitle.setAttribute('title', tooltip);
        projectTitle.classList.add('project-title');
        projectTitle.onclick = (_ => toggleSubprojects(projectTitle));

        if (p.URL) {
            nbsp = document.createElement('a');
            nbsp.innerHTML = '&nbsp;';

            a = document.createElement('a');
            a.href = p.URL;
            a.setAttribute('target', '_blank');

            i = document.createElement('i');
            i.classList.add('fas');
            i.classList.add('fa-external-link-alt');

            a.appendChild(i);
            projectTitle.appendChild(nbsp);
            projectTitle.appendChild(a);
        }

        p.Color.A /= 255;
        if (p.Percentage === 100) {
            p.Color.A = 0.5;
        }

        projectBar.classList.add('progress-bar');
        projectBar.classList.add('progress-bar-striped');
        projectBar.role = 'progressbar';
        projectBar.style.width = `${p.Percentage}%`;
        projectBar.style.backgroundColor = `rgba(${p.Color.R}, ${p.Color.G}, ${p.Color.B}, ${p.Color.A})`;
        projectBar.setAttribute("aria-valuenow", "50");
        projectBar.setAttribute("aria-valuemin", "0");
        projectBar.setAttribute("aria-valuemax", "100");

        projectContainer.appendChild(projectBar);
        projectContainer.classList.add('progress');
        projectContainer.setAttribute('title', tooltip);
        projectContainer.style.width = `${w}%`;
        projectContainer.style.marginLeft = `${l}%`;

        left.appendChild(projectTitle);
        right.appendChild(projectContainer);

        row.dataset.level = level;
        row.classList.add(`level`);
        row.classList.add(`level${level}`);
        row.appendChild(left);
        row.appendChild(right);

        right.classList.add('timeline');

        container.appendChild(row);

        if (p.Children !== null) {
            p.Children.forEach(c => displayProject(c, container, level + 1, projectStart, fullDiff));
        }
    };

    const toggleSubprojects = (project) => {
        const tr = project.parentElement.parentElement,
            tbody = tr.parentElement;

        let cur = tr, hide = true;

        if (tbody.lastElementChild === cur) {
            return;
        } else {
            hide = cur.nextElementSibling.style.display === 'table-row';
        }

        while (tbody.lastElementChild !== cur) {
            cur = cur.nextElementSibling;

            if (parseInt(cur.dataset.level) <= parseInt(tr.dataset.level)) {
                break;
            }

            if (hide || parseInt(cur.dataset.level) === parseInt(tr.dataset.level) + 1) {
                cur.style.display = hide ? 'none' : 'table-row';
            }
        }
    };

    const buildControl = (roadmap, control) => {
        const toggleBtn = document.createElement('button'),
            clearIcon = document.createElement('i'),
            levels = document.querySelectorAll('.level'),
            level1s = document.querySelectorAll('.level1');

        clearIcon.classList.add('fas');
        clearIcon.classList.add('fa-eye-slash');

        toggleBtn.classList.add('btn');
        toggleBtn.classList.add('btn-primary');
        toggleBtn.innerHTML = 'Hide Sublevels&nbsp;';
        toggleBtn.type = 'button';
        toggleBtn.appendChild(clearIcon);

        control.appendChild(toggleBtn);

        toggleBtn.addEventListener('click', event => {
            event.preventDefault();

            let hide = true;

            levels.forEach(l => {
                if (l.style.display === 'none') {
                    hide = false;
                }
            });

            levels.forEach(elem => elem.style.display = hide ? 'none' : 'table-row');

            level1s.forEach(elem => elem.style.display = 'table-row');
        });
    };

    if (roadmap) {
        buildTable(roadmap, table);

        buildControl(roadmap, control);
    }
};
