package roadmap

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/peteraba/roadmapper/pkg/bindata"
	"github.com/peteraba/roadmapper/pkg/herr"
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

func (r *Roadmap) viewHtml(appVersion, matomoDomain, docBaseURL, currentURL string, selfHosted bool, origErr error) (string, error) {
	writer := bytes.NewBufferString("")

	layoutTemplate := bindata.MustAsset("res/templates/index.html")

	t, err := template.New("layout").Parse(string(layoutTemplate))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var (
		pageTitle    = "New Roadmap"
		roadmapTitle = ""
		dateFormat   string
		baseURL      string
		raw          string
		hasRoadmap   bool
		projectURLs  = r.getProjectURLs()
	)

	if r != nil {
		dateFormat = r.DateFormat
		baseURL = r.BaseURL
		raw = string(r.ToContent())
		hasRoadmap = true
		pageTitle = r.Title
		roadmapTitle = r.Title
	}

	data := struct {
		MatomoDomain  string
		DocBaseURL    string
		DateFormat    string
		BaseURL       string
		CurrentURL    string
		PageTitle     string
		RoadmapTitle  string
		SelfHosted    bool
		HasRoadmap    bool
		Raw           string
		DateFormats   []string
		DateFormatMap map[string]string
		Version       string
		ProjectURLs   map[string][]string
		Error         error
	}{
		MatomoDomain:  matomoDomain,
		DocBaseURL:    docBaseURL,
		DateFormat:    dateFormat,
		BaseURL:       baseURL,
		CurrentURL:    currentURL,
		PageTitle:     pageTitle,
		RoadmapTitle:  roadmapTitle,
		SelfHosted:    selfHosted,
		HasRoadmap:    hasRoadmap,
		Raw:           raw,
		DateFormats:   dateFormats,
		DateFormatMap: dateFormatMap,
		Version:       appVersion,
		ProjectURLs:   projectURLs,
		Error:         origErr,
	}

	err = t.Execute(writer, data)
	if err != nil {
		return "", herr.NewFromError(err, http.StatusInternalServerError)
	}

	return writer.String(), nil
}

func (r *Roadmap) getProjectURLs() map[string][]string {
	projectURLs := map[string][]string{}

	if r == nil || r.Projects == nil {
		return projectURLs
	}

	for _, p := range r.Projects {
		if len(p.URLs) < 1 {
			continue
		}
		projectURLs[p.Title] = p.URLs
	}

	return projectURLs
}

func (r *Roadmap) pushAssets(pusher http.Pusher, version string) {
	_ = pusher.Push(strings.Join([]string{"/static/roadmapper.css?", version}, ""), nil)
	_ = pusher.Push(strings.Join([]string{"/static/roadmapper.mjs?", version}, ""), nil)
	_ = pusher.Push("/static/roadmap-form.mjs", nil)
	_ = pusher.Push("/static/roadmap-svg.mjs", nil)
}
