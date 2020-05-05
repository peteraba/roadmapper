export const formElements = {
    form: document.getElementById('roadmap-form'),
    titleField: document.getElementById('title'),
    txtField: document.getElementById('txt'),
    txtFieldValid: document.getElementById('txt-valid'),
    txtFieldInvalid: document.getElementById('txt-invalid'),
    dateFormatField: document.getElementById('base-url'),
    baseUrlField: document.getElementById('form-submit'),
    saveBtn: document.getElementById('form-submit'),
    resetBtn: document.getElementById('reset-btn'),
    loadExampleBtn: document.getElementById('load-example-btn'),
    areYouAHumanGroup: document.getElementById('are-you-a-human-group'),
    ts: document.getElementById('ts')
}

export const roadmapForm = elems => {
    const form = elems.form,
        titleField = elems.titleField,
        txtField = elems.txtField,
        txtFieldValid = elems.txtFieldValid,
        txtFieldInvalid = elems.txtFieldInvalid,
        dateFormatField = elems.dateFormatField,
        baseUrlField = elems.baseUrlField,
        saveBtn = elems.saveBtn,
        resetBtn = elems.resetBtn,
        loadExampleBtn = elems.loadExampleBtn,
        resetData = {
            'txt': txtField.value,
            'dateFormat': dateFormatField.value,
            'baseUrl': baseUrlField.value
        },
        areYouAHumanGroup = elems.areYouAHumanGroup,
        ts = elems.ts;

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
            s0 = txtField.selectionStart,
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

        if (txtFieldHistory[txtFieldHistoryNum - 1] !== txtField.value) {
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
        titleField.value = 'How To Start a Startup';
        txtField.value = `Find the idea [2019-07-20, 2020-01-20, 100%]
	Look for things missing in life
	Formalize your idea, run thought experiments [https://example.com/initial-plans]
	Survey friends, potential users or customers [https://example.com/survey-results]
	Go back to the drawing board [https://example.com/reworked-plans]
Validate the idea [2020-01-21, 2020-04-20]
	Make a prototype #1 [2020-01-21, 2020-04-10, 100%, TCK-1, https://github.com/peteraba/roadmapper, |1]
	Show the prototype to 100 people #1 [2020-04-11, 2020-04-20, 80%, TCK-123]
	Analyse results [2020-04-21, 2020-05-05]
	Improve prototype [2020-05-06, 2020-06-06]
	Show the prototype to 100 people #2 [2020-06-07, 2020-06-16]
	Analyse results [2020-06-16, 2020-06-30]
	Improve prototype [2020-07-01, 2020-07-16]
	Show the prototype to 100 people #2 [2020-07-17, 2020-07-25]
Start a business
	Learn about your options about various company types [2019-07-20, 2020-08-31]
	Learn about your options for managing equity [2019-07-20, 2020-08-01]
	Find a co-founder [2020-04-20, 2020-08-31]
	Register your business [2020-08-01, 2020-09-30, |2]
	Look for funding [2020-08-01, 2020-10-31]
	Build a team [2020-11-01, 2020-12-15]
Build version one [2021-01-01, 2021-04-15]
	Build version one [2021-01-01, 2021-03-31]
	Launch [2021-04-01, 2021-04-15, |3]
Grow [2021-04-16, 2021-12-31]
	Follow up with users
	Iterate / Pivot
	Launch again
	Get to 1,000 users
	Plan next steps

|Create the first prototype
|Start your business
|Lunch version one`;
        setSelectedIndex(dateFormatField, '2006-01-02');
        baseUrlField.value = 'https://example.com/foo';

        validateForm(form, txtField, txtFieldValid, txtFieldInvalid, saveBtn);
        saveHistory(titleField);
    });

    areYouAHumanGroup.remove();

    const disableSaveForTime = secs => {
        let t = Math.floor(ts.value);

        saveBtn.disabled = true;

        let si = setInterval(
            _ => {
                t++;
                ts.value = `${t}`;
                if (t === secs) {
                    saveBtn.disabled = false;
                    window.clearInterval(si);
                }
            },
            1000
        )
    }
    disableSaveForTime(5);
};
