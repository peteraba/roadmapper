package main

import (
	"bytes"
	"html/template"
	"strings"
)

const layoutTemplate = `<!doctype html>
<html lang="en">
<head>
	<!-- Required meta tags -->
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	<!-- Bootstrap -->
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha256-L/W5Wfqfa0sdBNIKN9cG6QA5F2qx4qICmU2VgLruv9Y=" crossorigin="anonymous" />
	<!-- Font Awesome -->
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.12.1/css/all.min.css" integrity="sha256-mmgLkCYLUQbXn0B1SRqzHar6dCnv9oZFPEC1g1cwlkk=" crossorigin="anonymous" />
	<!-- Custom -->
	<link rel="stylesheet" href="/static/roadmaper.css">

	<title>{{.Roadmap.Title}}</title>

	<!-- Favicon from https://favicon.io/favicon-generator/ [r, rounded, archivo black, 150, #343a40, #fff] -->
	<link rel="apple-touch-icon" sizes="180x180" href="/static/apple-touch-icon.png">
	<link rel="icon" type="image/png" sizes="32x32" href="/static/favicon-32x32.png">
	<link rel="icon" type="image/png" sizes="16x16" href="/static/favicon-16x16.png">
	<link rel="manifest" href="/static/site.webmanifest">

	{{ if .MatomoDomain }}
	<!-- Matomo -->
	<script type="text/javascript">
	  var _paq = window._paq || [];
	  /* tracker methods like "setCustomDimension" should be called before "trackPageView" */
	  _paq.push(['trackPageView']);
	  _paq.push(['enableLinkTracking']);
	  (function() {
		var u="//{{ .MatomoDomain }}/";
		_paq.push(['setTrackerUrl', u+'matomo.php']);
		_paq.push(['setSiteId', '1']);
		var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
		g.type='text/javascript'; g.async=true; g.defer=true; g.src=u+'matomo.js'; s.parentNode.insertBefore(g,s);
	  })();
	</script>
	<noscript><p><img src="//{{ .MatomoDomain }}/matomo.php?idsite=1&amp;rec=1" style="border:0;" alt="" /></p></noscript>
	<!-- End Matomo Code -->
	{{ end }}

</head>
<body>
	<nav class="navbar sticky-top navbar-expand-lg navbar-dark bg-dark">
		<a class="navbar-brand" href="/">rdmp.app <small>(alpha)</small></a>
		<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarNavAltMarkup">
			<div class="navbar-nav mr-auto">
				<a class="nav-item nav-link roadmap-dashboard-link" href="#roadmap-dashboard">Dashboard <i class="fas fa-eye"></i></a>
				<a class="nav-item nav-link" href="#roadmap-edit">Edit <i class="fas fa-edit"></i></a>
				<a class="nav-item nav-link" href="https://docs.rdmp.app/about/rdmp.app/" target="_blank">About <i class="fas fa-question-circle"></i></a>
				<a class="nav-item nav-link" href="https://docs.rdmp.app/" target="_blank">Docs <i class="fas fa-book"></i></a>
				<a class="nav-item nav-link" href="https://docs.rdmp.app/privacy/privacy-policy/" target="_blank">Data Privacy <i class="fas fa-user-secret"></i></a>
				<a class="nav-item nav-link" href="https://github.com/peteraba/roadmapper" target="_blank">Source Code <i class="fas fa-code-branch"></i></a>
			</div>
			<div class="navbar-nav">
				<a class="btn btn-light" href="/">New <i class="fas fa-plus-circle"></i></a>
			</div>
		</div>
	</nav>

{{ if .SelfHosted }}
	<div class="modal fade" id="privacy-policy" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle" aria-hidden="false">
		<div class="modal-dialog modal-dialog-centered" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h4 class="modal-title" id="exampleModalLongTitle"><i class="fas fa-shield-alt"></i> Data Protection</h4>
					<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
					</button>
				</div>
				<div class="modal-body">
					<p>Please do not store sensitive data on <a href="https://rdmp.app/">https://rdmp.app/</a>.<br>
						At the moment it is mainly a <strong>tech demo</strong> <i class="fas fa-exclamation"></i><br><br>
						Please use <a href="https://github.com/peteraba/roadmapper">Roadmapper</a> instead,
						or read our <a href="https://docs.rdmp.app/privacy/privacy-policy/" target="_blank">Data Privacy <i class="fas fa-external-link-alt"></i>.</a>
					</p>
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-link" id="privacy-policy-save">I understand, I promise!</button>
					<button type="button" class="btn btn-primary" id="privacy-policy-ok">OK.</button>
				</div>
			</div>
		</div>
	</div>
{{ end }}

	<div class="container-fluid roadmap-dashboard section" id="roadmap-dashboard">
		<h1 class="h1">Dashboard</h1>
		<div id="roadmap-svg"></div>
		<table class="table table-bordered roadmap" id="roadmap">
		  <!-- Content here -->
		</table>
		<form id="control"></form>
		<hr class="hr">
	</div>

	<div class="container roadmap-edit section" id="roadmap-edit">
		<h1 class="h1">Form</h1>
		<form action="" method="POST" id="roadmap-form">
			<div class="form-group">
				<label for="txt">Raw roadmap</label>
				<textarea class="form-control" id="txt" name="txt" aria-describedby="txtHelp" rows="20">{{ .Raw }}</textarea>
				<div class="valid-feedback" id="txt-valid"></div>
				<div class="invalid-feedback" id="txt-invalid"></div>
				<small id="txtHelp" class="form-text text-muted"><a href="https://docs.rdmp.app/usage/format/">Format documentation</a></small>
			</div>
			<div class="form-group">
				<label for="dateFormat">Date format</label>
				<select id="dateFormat" name="dateFormat" class="form-control">
				{{range $val := .DateFormats }}
					 <option value="{{ $val }}"{{ if eq $val $.DateFormat }} selected{{ end }}>{{ index $.DateFormatMap $val }}</option>
				{{end}}
				</select>
			</div>
			<div class="form-group">
				<label for="baseUrl">Base URL</label>
				<input class="form-control" id="baseUrl" name="baseUrl" type="url" aria-describedby="baseUrlHelp" value="{{ .BaseUrl }}" />
				<small id="baseUrlHelp" class="form-text text-muted">URL to prepend for URLs in your roadmap</small>
			</div>
			<button type="submit" class="btn btn-primary" id="form-submit">Submit</button>
		</form>
		<hr class="hr">
	</div>

<script>
var roadmap = {{ .Roadmap }};
</script>
	
<!-- Optional JavaScript -->
<!-- jQuery first, then Popper.js, then Bootstrap JS -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.4.1/jquery.min.js" integrity="sha256-CSXorXvZcTkaix6Yvo6HppcZGetbYMGWSFlBw8HfCJo=" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.16.1/umd/popper.min.js" integrity="sha256-/ijcOLwFf26xEYAjW75FizKVo5tnTYiQddPZoLUHHZ8=" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha256-WqU1JavFxSAMcLP2WIOI+GB2zWmShMI82mTpLDcqFUg=" crossorigin="anonymous"></script>
<!-- Custom -->
<script type="module" src="/static/roadmaper.mjs"></script>
</body>
</html>
`

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

func bootstrapRoadmap(roadmap Project, lines []string, matomoDomain, dateFormat, baseUrl string, selfHosted bool) (string, error) {
	writer := bytes.NewBufferString("")

	t, err := template.New("layout").Parse(layoutTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		Roadmap       Project
		MatomoDomain  string
		DateFormat    string
		BaseUrl       string
		SelfHosted    bool
		Raw           string
		DateFormats   []string
		DateFormatMap map[string]string
	}{
		Roadmap:       roadmap,
		MatomoDomain:  matomoDomain,
		DateFormat:    dateFormat,
		BaseUrl:       baseUrl,
		SelfHosted:    selfHosted,
		Raw:           strings.Join(lines, "\n"),
		DateFormats:   dateFormats,
		DateFormatMap: dateFormatMap,
	}

	err = t.Execute(writer, data)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}
