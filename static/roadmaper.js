(() => {
    const handleError = msg => {
        console.error(msg);
    };

    const refreshTimeline = () => {
        const drawing = document.getElementById('drawing');

        if (!drawing || !drawing.dataset || !drawing.dataset.from || !drawing.dataset.to) {
            return;
        }

        const h = 75,
            projectFrom = new Date(drawing.dataset.from),
            projectTo = new Date(drawing.dataset.to),
            lineColor = window.getComputedStyle(drawing, null).borderColor,
            textColor = window.getComputedStyle(drawing, null).color;

        let line, text, draw;

        drawing.innerHTML = '';

        window.setTimeout(() => {
            const w = drawing.offsetWidth,
                today = new Date(),
                pink = '#f88',
                red = '#f33',
                left = (today.getTime() - projectFrom.getTime()) / (projectTo.getTime() - projectFrom.getTime()) * w;

            draw = SVG().addTo(drawing).size(w, h);

            // baseline
            line = draw.line(0, 0, w, 0).move(0, 40);
            line.stroke({color: lineColor, width: 2});

            // first day
            line = draw.line(0, 0, 0, 10).move(10, 35);
            line.stroke({color: textColor, width: 2});
            text = draw.text(projectFrom.toLocaleDateString());
            text.move(5, 10).font({});

            // last day
            line = draw.line(0, 0, 0, 10).move(w - 10, 35);
            line.stroke({color: textColor, width: 2});
            text = draw.text(projectTo.toLocaleDateString());
            text.move(w - 5, 10).font({anchor: 'end'});

            // today
            if (projectFrom.getTime() < today.getTime() && projectTo.getTime() > today.getTime()) {
                line = draw.line(0, 0, 0, h).move(left, 0);
                line.stroke({color: pink, width: 1});
                text = draw.text(today.toLocaleDateString());
                text.move(left, 45).font({fill: red, anchor: 'middle'});
            }
        }, 0);
    };

    const refreshRoadmap = () => {
        const table = document.getElementById('roadmap'),
            control = document.getElementById('control');

        table.innerHTML = '';

        const buildTable = (p, container) => {
            const
                thead = document.createElement('thead'),
                tbody = document.createElement('tbody'),
                projectFrom = new Date(p.From),
                projectTo = new Date(p.To),
                fullDiff = projectTo.getTime() - projectFrom.getTime();

            container.appendChild(thead);
            container.appendChild(tbody);

            displayHeader(p, thead, p.From, p.To);

            if (p.Children !== null) {
                p.Children.forEach(c => displayProject(c, tbody, 1, projectFrom, fullDiff));
            }
        };

        const displayHeader = (p, container, projectFrom, projectTo) => {
            const row = document.createElement('tr'),
                left = document.createElement('th'),
                right = document.createElement('th'),
                roadmapTitle = document.createElement('p'),
                drawing = document.createElement('div');

            roadmapTitle.innerHTML = '&nbsp;';

            row.appendChild(left);
            row.appendChild(right);

            drawing.id = 'drawing';
            drawing.dataset.from = projectFrom;
            drawing.dataset.to = projectTo;

            right.append(drawing);
            right.classList.add('timeline');

            left.appendChild(roadmapTitle);

            container.appendChild(row);

            refreshTimeline();
        };

        const displayProject = (p, container, level, projectFrom, fullDiff) => {
            const row = document.createElement('tr'),
                left = document.createElement('th'),
                right = document.createElement('td'),
                projectTitle = document.createElement('span'),
                projectContainer = document.createElement('div'),
                projectBar = document.createElement('div'),
                from = new Date(p.From),
                to = new Date(p.To),
                diff = to.getTime() - from.getTime(),
                durationDays = diff / 86400000,
                w = diff / fullDiff * 100,
                l = (from.getTime() - projectFrom.getTime()) / fullDiff * 100,
                tooltip = `${p.Percentage}%, ${from.toLocaleDateString()} - ${to.toLocaleDateString()}, ${durationDays} days`;

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
                p.Children.forEach(c => displayProject(c, container, level + 1, projectFrom, fullDiff));
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

    window.addEventListener('DOMContentLoaded', () => {
        const txt = document.getElementById('txt');

        if (!roadmap || !roadmap.Children) {
            $('#roadmap-dashboard, .roadmap-dashboard-link').remove();
        } else {
            refreshRoadmap();
        }

        txt.addEventListener('paste', e => {
            let lines = (e.clipboardData || window.clipboardData).getData('text').split("\n");

            // Remove empty lines from the top
            while (lines.length > 0 && lines[0].trim() === "") {
                lines.shift();
            }

            if (lines.length === 0) {
                return handleError('empty lines');
            }

            // Find out if all lines are indented
            const ws = lines[0].match(/^\s+/g),
                ind = (ws && ws.length > 0 ? ws[0] : ''),
                allIndented =
                    ind.length === 0 ||
                    lines.every(line => {
                        return line.indexOf(ind) === 0;
                    });

            if (!allIndented) {
                return handleError('not all indented');
            }

            // Remove common indentation
            if (ind.length > 0) {
                for (const x in lines) {
                    lines[x] = lines[x].substr(ind.length).trimRight();
                }
            }

            // Find the indentation of the first indented line
            let t;
            for (const line of lines) {
                const m = line.match(/^\s+/g);
                if (m && m.length) {
                    t = m[0];
                    break;
                }
            }

            // Turn indentation into tabs
            let val = lines.join("\n");
            if (t !== "\t") {
                val = val.replace(new RegExp(`${t}`, 'g'), "\t");
            }

            txt.value = val;

            e.preventDefault();
        });

        const handleTab = (e, field) => {
            const
                lineStart = (text, start) => {
                    const prevNL = text.substr(0, start).lastIndexOf("\n");

                    if (prevNL >= 0) {
                        return prevNL;
                    }

                    return 0;
                },
                lineEnd = (text, end) => {
                    const nextNL = text.substr(end).indexOf("\n");

                    if (nextNL >= 0) {
                        return end + nextNL;
                    }

                    return end;
                },
                applyShift = (text, hasShift) => {
                    if (hasShift) {
                        return text.replace(/\n\t/g, "\n");
                    }

                    return text.replace(/\n/g, "\n\t");
                };

            const text = field.value,
                s0 = field.selectionStart,
                e0 = field.selectionEnd,
                s1 = lineStart(text, s0),
                e1 = lineEnd(text, e0),
                v0 = text.substring(s1, e1),
                v1 = applyShift(v0, e.shiftKey);

            field.value = text.substring(0, s1) + v1 + text.substring(e1);

            field.selectionStart = s1 + 1;
            field.selectionEnd = field.selectionStart + v1.length - 1;

            e.preventDefault();
        };

        const handleSpace = (e, field) => {
            if (field.selectionStart !== field.selectionEnd) {
                return;
            }

            const
                text = field.value,
                s0 = field.selectionStart,
                e0 = field.selectionEnd,
                lineStart = (text, start) => {
                    const prevNL = text.substr(0, start).lastIndexOf("\n");

                    if (prevNL >= 0) {
                        return prevNL + 1;
                    }

                    return 0;
                },
                start = lineStart(field.value, s0),
                lineToCur = text.substr(start, s0 - start),
                m = lineToCur.match(/^\s+$/);

            // Do nothing if we're not somewhere inside the task definition
            if (m === null) {
                return;
            }

            field.value = text.substring(0, s0) + "\t" + text.substring(e0);

            field.selectionStart = s0 + 1;
            field.selectionEnd = e0 + 1;

            e.preventDefault();
        };

        txt.addEventListener('keydown', e => {
            switch (e.key) {
                case 'Tab':
                    return handleTab(e, txt);
                case ' ':
                    return handleSpace(e, txt);
            }
        });

        $('[data-toggle="tooltip"]').tooltip();
    });

    window.addEventListener('resize', () => {
        refreshTimeline();
    });
})();
