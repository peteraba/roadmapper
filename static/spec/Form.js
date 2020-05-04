import {roadmapForm} from '../roadmap-form.mjs';

function newForm() {
    const
        form = document.createElement('form'),
        titleField = document.createElement('input'),
        txtField = document.createElement('textarea'),
        txtFieldValid = document.createElement('div'),
        txtFieldInvalid = document.createElement('div'),
        dateFormatField = document.createElement('input'),
        baseUrlField = document.createElement('input'),
        saveBtn = document.createElement('button'),
        resetBtn = document.createElement('button'),
        loadExampleBtn = document.createElement('button'),
        areYouAHumanGroup = document.createElement('input'),
        ts = document.createElement('input');

    return {
        form: form,
        titleField: titleField,
        txtField: txtField,
        txtFieldValid: txtFieldValid,
        txtFieldInvalid: txtFieldInvalid,
        dateFormatField: dateFormatField,
        baseUrlField: baseUrlField,
        saveBtn: saveBtn,
        resetBtn: resetBtn,
        loadExampleBtn: loadExampleBtn,
        areYouAHumanGroup: areYouAHumanGroup,
        ts: ts,
    }
}

describe("Form", function () {
    it("ts should get updated every second", function (done) {
        const elems = newForm(), tsVal = elems.ts.value;

        roadmapForm(elems);

        setTimeout(
            _ => {
                expect(elems.ts.value).not.toEqual(tsVal);
                expect(elems.ts.value).toBeGreaterThanOrEqual(0);
                done();
            },
            1000
        )
    });
});
