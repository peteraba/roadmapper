package main

import (
	"bytes"
	"fmt"
	"html/template"
)

var dateFormats = []string{
	"2006-01-02",
	"2006.01.02",
	"2006/01/02",
	"02.01.2006",
	"02/01/2006",
	"01/02/2020",
	"01.02.2020",
	"2006-1-2",
	"2006/1/2",
	"2.1.2006",
	"2/1/2006",
	"1/2/2020",
	"1.2.2020",
}
var dateFormatMap = map[string]string{
	"2006-01-02": "YYYY-MM-DD (2020-03-17)",
	"2006.01.02": "YYYY.MM.DD (2020.03.17)",
	"2006/01/02": "YYYY/MM/DD (2020/03/17)",
	"02.01.2006": "DD.MM.YYYY (17.03.2020)",
	"02/01/2006": "DD/MM/YYYY (17/03/2020)",
	"01/02/2020": "MM/DD/YYYY (03/17/2020)",
	"01.02.2020": "MM/DD/YYYY (03.17.2020)",
	"2006-1-2":   "YYYY-M-D (2020-3-7)",
	"2006/1/2":   "YYYY/M/D (2020/3/7)",
	"2.1.2006":   "D.M.YYYY (7.3.2020)",
	"2/1/2006":   "D/M/YYYY (7/3/2020)",
	"1/2/2020":   "M/D/YYYY (3/7/2020)",
	"1.2.2020":   "M/D/YYYY (3.7.2020)",
}

func bootstrapRoadmap(roadmap *Roadmap, matomoDomain, docBaseUrl, currentUrl string, selfHosted bool) (string, error) {
	writer := bytes.NewBufferString("")

	layoutTemplate, err := Asset("templates/index.html")
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	t, err := template.New("layout").Parse(string(layoutTemplate))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var dateFormat, baseUrl, title, raw string
	hasRoadmap := false

	if roadmap != nil {
		dateFormat = roadmap.DateFormat
		baseUrl = roadmap.BaseURL
		title = ""
		raw = string(roadmap.ToContent())
		hasRoadmap = true
	}

	data := struct {
		Roadmap       *Roadmap
		MatomoDomain  string
		DocBaseUrl    string
		DateFormat    string
		BaseUrl       string
		CurrentUrl    string
		Title         string
		SelfHosted    bool
		HasRoadmap    bool
		Raw           string
		DateFormats   []string
		DateFormatMap map[string]string
	}{
		Roadmap:       roadmap,
		MatomoDomain:  matomoDomain,
		DocBaseUrl:    docBaseUrl,
		DateFormat:    dateFormat,
		BaseUrl:       baseUrl,
		CurrentUrl:    currentUrl,
		Title:         title,
		SelfHosted:    selfHosted,
		HasRoadmap:    hasRoadmap,
		Raw:           raw,
		DateFormats:   dateFormats,
		DateFormatMap: dateFormatMap,
	}

	err = t.Execute(writer, data)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}
