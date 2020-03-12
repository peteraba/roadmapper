export const roadmapForm = () => {
    const form = document.getElementById('roadmap-form'),
        txt = document.getElementById('txt'),
        txtValid = document.getElementById('txt-valid'),
        txtInvalid = document.getElementById('txt-invalid');

    const handleInvalidTextarea = (field, msg, lines) => {
        field.classList.add('is-invalid');

        txtValid.innerText = '';
        txtValid.style.display = 'none';
        txtInvalid.innerText = msg + " on lines: " + lines.join(",");
        txtInvalid.style.display = 'block';
    };

    const handleValidTextarea = (field, msg) => {
        if (txtInvalid.innerText === '') {
            return;
        }

        form.classList.add('was-validated');
        field.classList.remove('is-invalid');

        txtValid.innerText = msg;
        txtValid.style.display = 'block';
        txtInvalid.innerText = '';
        txtInvalid.style.display = 'none';

        window.setTimeout(() => {
            form.classList.remove('was-validated');
            txtValid.innerText = '';
            txtValid.style.display = 'none';
        }, 2000);
    };

    const lineStart = (text, start) => {
        const prevNL = text.substr(0, start).lastIndexOf("\n");

        if (prevNL >= 0) {
            return prevNL + 1;
        }

        return 0;
    };

    const getOpeningTabs = (text) => {
        const m = text.match(/^\t/);

        return m && m.length > 0 ? m[0] : '';
    };

    const handlePaste = (e, field) => {
        const before = field.value.substr(0, field.selectionStart),
            after = field.value.substr(field.selectionEnd);

        let lines = (e.clipboardData || window.clipboardData).getData('text').split("\n");

        // Remove empty lines from the top
        while (lines.length > 0 && lines[0].trim() === "") {
            lines.shift();
        }

        if (lines.length === 0) {
            return handleInvalidTextarea(txt, 'empty lines', [0]);
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
            return handleInvalidTextarea(txt, 'some lines are not indented', []);
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

        txt.value = before + val + after;

        e.preventDefault();
    };

    txt.addEventListener('paste', e => {
        handlePaste(e, txt);
    });

    const handleTab = (e, field) => {
        const
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
            s1 = lineStart(text, s0) - 1,
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

    const handleEnter = (e, field) => {
        const
            text = field.value,
            s0 = field.selectionStart,
            e0 = field.selectionEnd,
            ls = lineStart(text, s0),
            tabs = getOpeningTabs(text.substr(ls, s0 - ls));

        field.value = text.substring(0, s0) + "\n" + tabs + text.substring(e0);

        field.selectionStart = s0 + 1 + tabs.length;
        field.selectionEnd = e0 + 1 + tabs.length;

        e.preventDefault();
    };

    txt.addEventListener('keydown', e => {
        switch (e.key) {
            case 'Tab':
                return handleTab(e, txt);
            case 'Enter':
                return handleEnter(e, txt);
            case ' ':
                return handleSpace(e, txt);
        }
    });

    const validateRoadmap = (e, field) => {
        let prevIndCount = 0, errors = [];

        field.value.split("\n").forEach((val, idx) => {
            const m = val.match(/^\t*/),
                curIndCount = m && m[0] ? m[0].length : 0;

            if (idx === 0 && curIndCount > 0) {
                errors.push(idx);
            } else if (prevIndCount < curIndCount - 1) {
                errors.push(idx);
            }

            prevIndCount = curIndCount;
        });

        if (errors.length > 0) {
            return handleInvalidTextarea(field, 'invalid indentation', errors);
        }

        return handleValidTextarea(field, 'valid roadmap');
    };

    txt.addEventListener('keydown', e => {
        validateRoadmap(e, txt);
    });
};
