export const roadmapForm = () => {
    const form = document.getElementById('roadmap-form'),
        titleField = document.getElementById('title'),
        txtField = document.getElementById('txt'),
        txtFieldValid = document.getElementById('txt-valid'),
        txtFieldInvalid = document.getElementById('txt-invalid'),
        dateFormatField = document.getElementById('date-format'),
        baseUrlField = document.getElementById('base-url'),
        saveBtn = document.getElementById('form-submit'),
        resetBtn = document.getElementById('reset-btn'),
        loadExampleBtn = document.getElementById('load-example-btn'),
        resetData = {'txt': txtField.value, 'dateFormat': dateFormatField.value, 'baseUrl': baseUrlField.value};

    let validationTimeout = false,
        txtFieldHistory = [txtField.value],
        txtFieldHistoryNum = 0;

    const handleInvalidTextarea = (txtField, txtFieldValid, txtFieldInvalid, msg, lines) => {
        txtField.classList.add('is-invalid');

        txtFieldValid.innerText = '';
        txtFieldValid.style.display = 'none';
        txtFieldInvalid.innerText = msg + " on lines: " + lines.join(", ");
        txtFieldInvalid.style.display = 'block';
    };

    const handleValidTextarea = (txtField, txtFieldValid, txtFieldInvalid, msg) => {
        if (txtFieldInvalid.innerText === '') {
            return;
        }

        txtField.classList.remove('is-invalid');

        txtFieldValid.innerText = msg;
        txtFieldValid.style.display = 'block';
        txtFieldInvalid.innerText = '';
        txtFieldInvalid.style.display = 'none';
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

    const handlePaste = (e, txtField, txtFieldValid, txtFieldInvalid) => {
        const origText = txtField.value,
            s0 =  txtField.selectionStart,
            e0 = txtField.selectionEnd,
            before = origText.substr(0, s0),
            after = origText.substr(e0);

        let lines = (e.clipboardData || window.clipboardData).getData('text').split("\n");

        if (lines.length === 0) {
            return handleInvalidTextarea(txtField, txtFieldValid, txtFieldInvalid, 'empty lines', [0]);
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
            return handleInvalidTextarea(txtField, txtFieldValid, txtFieldInvalid, 'some lines are not indented', []);
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

        txtField.value = before + val + after;

        txtField.selectionStart = s0 + val.length;
        txtField.selectionEnd = s0 + val.length;

        e.preventDefault();
    };

    txtField.addEventListener('paste', e => {
        saveHistory(txtField);
        handlePaste(e, txtField, txtFieldValid, txtFieldInvalid);
        saveHistory(txtField);
    });

    const handleTab = (e, txtField) => {
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

        const text = txtField.value,
            s0 = txtField.selectionStart,
            e0 = txtField.selectionEnd,
            s1 = lineStart(text, s0) - 1,
            e1 = lineEnd(text, e0),
            v0 = text.substring(s1, e1),
            v1 = applyShift(v0, e.shiftKey);

        txtField.value = text.substring(0, s1) + v1 + text.substring(e1);

        txtField.selectionStart = s0 + v1.length - v0.length;
        txtField.selectionEnd = e0 + v1.length - v0.length;

        e.preventDefault();
    };

    const handleSpace = (e, txtField) => {
        if (txtField.selectionStart !== txtField.selectionEnd) {
            return;
        }

        const
            text = txtField.value,
            s0 = txtField.selectionStart,
            e0 = txtField.selectionEnd,
            start = lineStart(txtField.value, s0),
            lineToCur = text.substr(start, s0 - start),
            m = lineToCur.match(/^\s+$/);

        // Do nothing if we're not somewhere inside the task definition
        if (m === null && lineToCur !== "") {
            return;
        }

        txtField.value = text.substring(0, s0) + "\t" + text.substring(e0);

        txtField.selectionStart = s0 + 1;
        txtField.selectionEnd = e0 + 1;

        e.preventDefault();
    };

    const handleEnter = (e, txtField) => {
        const
            text = txtField.value,
            s0 = txtField.selectionStart,
            e0 = txtField.selectionEnd,
            ls = lineStart(text, s0),
            tabs = getOpeningTabs(text.substr(ls, s0 - ls));

        txtField.value = text.substring(0, s0) + "\n" + tabs + text.substring(e0);

        txtField.selectionStart = s0 + 1 + tabs.length;
        txtField.selectionEnd = e0 + 1 + tabs.length;

        e.preventDefault();
    };

    const saveHistory = txtField => {
        txtFieldHistory = txtFieldHistory.slice(0, txtFieldHistoryNum);

        if (txtFieldHistory[txtFieldHistoryNum-1] !== txtField.value) {
            txtFieldHistory.push(txtField.value);
            txtFieldHistoryNum++;
        }
    };

    const applyHistory = (e, txtField) => {
        if (!e.ctrlKey) {
            return;
        }

        if (!e.shiftKey && txtFieldHistoryNum > 0) {
            txtFieldHistoryNum--;
            txtField.value = txtFieldHistory[txtFieldHistoryNum];
        } else if (e.shiftKey && txtFieldHistoryNum < txtFieldHistory.length - 1) {
            txtFieldHistoryNum++;
            txtField.value = txtFieldHistory[txtFieldHistoryNum];
        }
    };

    txtField.addEventListener('keydown', e => {
        switch (e.key) {
            case 'Tab':
                saveHistory(txtField);

                return handleTab(e, txtField);
            case ' ':
                return handleSpace(e, txtField);
            case 'Enter':
                saveHistory(txtField);

                return handleEnter(e, txtField);
            case 'z':
                return applyHistory(e, txtField);
        }
    });

    const findRoadmapErrors = txtField => {
        let prevIndCount = 0, errors = [];

        txtField.value.split("\n").forEach((val, idx) => {
            const m = val.match(/^\t*/),
                curIndCount = m && m[0] ? m[0].length : 0;

            if (idx === 0 && curIndCount > 0) {
                errors.push(idx);
            } else if (prevIndCount < curIndCount - 1) {
                errors.push(idx);
            }

            prevIndCount = curIndCount;
        });

        return errors;
    };

    const showWasValidated = (form, txtFieldValid) => {
        form.classList.add('was-validated');

        if (validationTimeout) {
            clearTimeout(validationTimeout);
        }

        window.setTimeout(() => {
            form.classList.remove('was-validated');
            txtFieldValid.innerText = '';
            txtFieldValid.style.display = 'none';
            validationTimeout = false;
        }, 5000);
    };

    const validateForm = (form, txtField, txtFieldValid, txtFieldInvalid, saveBtn) => {
        const roadmapErrors = findRoadmapErrors(txtField),
            hasError = roadmapErrors.length > 0;

        showWasValidated(form, txtFieldValid);

        if (hasError) {
            handleInvalidTextarea(txtField, txtFieldValid, txtFieldInvalid, 'invalid indentation', roadmapErrors);
            saveBtn.disabled = true;
        } else {
            handleValidTextarea(txtField, txtFieldValid, txtFieldInvalid, 'valid roadmap');
        }

        saveBtn.disabled = hasError;

        return hasError;
    };

    form.addEventListener('submit', e => {
        const hasError = validateForm(form, txtField, txtFieldValid, txtFieldInvalid, saveBtn);

        if (hasError) {
            e.preventDefault();
        }
    });

    txtField.addEventListener('keydown', e => {
        validateForm(form, txtField, txtFieldValid, txtFieldInvalid, saveBtn);
    });

    txtField.addEventListener('keydown', e => {
        validateForm(form, txtField, txtFieldValid, txtFieldInvalid, saveBtn);
    });

    txtField.addEventListener('keydown', e => {
        validateForm(form, txtField, txtFieldValid, txtFieldInvalid, saveBtn);
    });

    const setSelectedIndex = (selectField, value) => {
        const opts = selectField.options;

        for (let j = 0; j < opts.length; j++) {
            if (opts[j].value === value) {
                selectField.selectedIndex = j;
                break;
            }
        }
    };

    resetBtn.addEventListener('click', _ => {
        txtField.value = resetData.txt;
        setSelectedIndex(dateFormatField, resetData.dateFormat);
        baseUrlField.value = resetData.baseUrl;

        validateForm(form, txtField, txtFieldValid, txtFieldInvalid, saveBtn);
    });

    loadExampleBtn.addEventListener('click', _ => {
        saveHistory(titleField);
        titleField.value = 'Example Roadmap';
        txtField.value = `Monocle ipsum dolor sit amet
Ettinger punctual izakaya concierge [2020-02-02, 2020-02-20, 60%]
	Zürich Baggu bureaux [/issues/1]
		Toto Comme des Garçons liveable [2020-02-04, 2020-02-25, 100%, /issues/2]
		Winkreative boutique St Moritz [2020-02-06, 2020-02-22, 55%, /issues/3]
	Toto joy perfect Porter  [2020-02-25, 2020-03-01, 100%, |1]
Craftsmanship artisanal
	Marylebone exclusive [2020-03-03, 2020-03-10, 100%]
	Beams elegant destination [2020-03-08, 2020-03-12, 100%, |1]
	Winkreative ryokan hand-crafted [2020-03-13, 2020-03-31, 20%]
Nordic Toto first-class Singap
	Concierge cutting-edge Zürich global bureaux
		Sunspel sophisticated lovely uniforms [2020-03-17, 2020-03-31]
		Share blog post on social media [2020-03-17, 2020-03-31, 80%]
	Talk about the tool in relevant meetups [2020-04-01, 2020-06-15, 20%]
Melbourne handsome boutique
	Boutique magna iconic
		Carefully curated laborum destination [2020-03-28, 2020-05-01, 60%]
	Qui incididunt sleepy
		Scandinavian occaecat culpa [2020-03-26, 2020-04-01, 90%]
Hand-crafted K-pop boulevard
	Charming sed quality [2020-03-18, 2020-05-31, 20%]
	Sunspel alluring ut dolore [2020-04-15, 2020-04-30, 30%]
Business class Shinkansen [2020-04-01, 2020-05-31, 45%]
	Nisi excepteur hand-crafted hub
	Ettinger Airbus A380
Essential conversation bespoke
Muji enim

|Laboris ullamco
|Muji enim finest [2020-02-12, https://example.com/abc, bcdef]`;
        setSelectedIndex(dateFormatField, '2006-01-02');
        baseUrlField.value = 'https://example.com/foo';

        validateForm(form, txtField, txtFieldValid, txtFieldInvalid, saveBtn);
        saveHistory(titleField);
    });
};
