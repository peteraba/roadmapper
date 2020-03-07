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
</head>
<body>
	<nav class="navbar sticky-top navbar-expand-lg navbar-dark bg-dark">
		<a class="navbar-brand" href="#">rdmp.app <small>(alpha)</small></a>
		<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarNavAltMarkup">
			<div class="navbar-nav mr-auto">
				<a class="nav-item nav-link roadmap-dashboard-link" href="#roadmap-dashboard">Dashboard</a>
				<a class="nav-item nav-link" href="#roadmap-edit">Edit</a>
				<a class="nav-item nav-link" href="#roadmap-docs">Docs</a>
				<a class="nav-item nav-link" href="#roadmap-about">About</a>
				<a class="nav-item nav-link" href="#roadmap-privacy">Data Privacy</a>
				<a class="nav-item nav-link" href="https://github.com/peteraba/roadmapper" target="_blank">Source Code <i class="fas fa-external-link-alt"></i></a>
			</div>
			<div class="navbar-nav">
				<a class="nav-item nav-link" href="/">New</a>
			</div>
		</div>
	</nav>

	<div class="container-fluid roadmap-dashboard section" id="roadmap-dashboard">
		<h1 class="h1">Dashboard</h1>
		<table class="table table-bordered roadmap" id="roadmap">
		  <!-- Content here -->
		</table>
		<form id="control"></form>
		<hr class="hr">
	</div>

	<div class="container roadmap-edit section" id="roadmap-edit">
		<h1 class="h1">Form</h1>
		<form action="" method="POST">
			<div class="form-group">
				<label for="txt">Raw roadmap</label>
				<textarea class="form-control" id="txt" name="txt" aria-describedby="txtHelp" rows="20">{{ .Raw }}</textarea>
				<small id="txtHelp" class="form-text text-muted"><a href="#documentation-format">Format documentation</a></small>
			</div>
			<button type="submit" class="btn btn-primary">Submit</button>
		</form>
		<hr class="hr">
	</div>

	<div class="container roadmap-docs section" id="roadmap-docs">
		<h1 class="h1">Documentation</h1>
		<p>Roadmapper aims to provide various tools to help teams to keep track of their progress.</p>

		<h2 class="h2" id="documentation-format">Format</h2>
		<p>Roadmapper aims to be easy to understand and extend.</p>
		<ol>
			<li>You can start with a list of tasks to accomplish.</li>
			<li>Indent sub-tasks by starting the line with two spaces.</li>
			<li>Add extra information inside a pair of brackets at the end of the line. Separate each piece by a comma.</li>
			<li>Add the (estimated) start and end dates in the following format: <strong>YYYY-MM-DD</strong>.</li>
			<li>Optionally add the percentage of accomplishment as decimal value between 0 and 1 or percentage between 0% and 100%.</li>
			<li>Optionally add a hexadecimal color code. (e.g. #e9f8a3)</li>
		</ol>

		<h3 class="h3">Example</h3>
		<pre>
beta.0 - tabs[2020-02-12, 2020-02-18]
beta.1
	12-factor-app [2020-02-12, 2020-02-18, 60%, #f00]
	Validation improvements [2020-02-20, 2020-03-30, 0, beta.1]
idea-pool [idea-pool]
	Locking feature
		Investigations [2020-03-10, 2020-03-15, 0.006]
		Layout logic [2020-03-15, 2020-03-18, 70%]
	Map wysiwg editor languages
		Investigations [2020-03-10, 2020-03-15, 0.006]
</pre>

		<h2 class="h2">Usage</h2>

		<h3 class="h3">Command Line</h3>
		<p>The command line features of Roadmapper are meant to help teams with easily trackable roadmaps. They can be:</p>
		<ol>
			<li>Manually embedded into web-based documentations</li>
			<li>Integrated into static site generator workflows.</li>
		</ol>
		
		<h4>Usage</h4>
		<p>Roadmapper does not update documents.</p>
		<pre>
> roadmapper c -i 34sgkhkA
		</pre>

		<h4>Embedding into web-based documentations (e.g. Google Docs)</h4>
		<p>Roadmapper does not update documents.</p>

		<h4>Integrating into static site generator workflows (Hugo)</h4>
<pre>
> roadmapper c -i 34sgkhkA
> hugo --theme=hugo-bootstrap --verbose
</pre>

		<h3 class="h3">Online Service</h3>
		<p>The online service exists mainly to provide a super simple way to test the features of Roadmapper. It can be used for teams freely, but it does not come with warranties and it is not possible to delete roadmaps!</p>

		<h2 class="h2">Source code</h2>
		<p>The source code is available on <a href="https://github.com/peteraba/roadmapper" target="_blank">Github <i class="fas fa-external-link-alt"></i></a>.</p>
		<hr class="hr">
	</div>

	<div class="container roadmap-about section" id="roadmap-about">
		<h1 class="h1">About Roadmapper</h1>
		<p>Roadmapper aims to provide various tools to help teams to keep track of their progress. There are two main supported use cases at the moment:</p>
		<ul>
			<li>Online service (both self-hosted and on <a href="https://rdmp.app/">https://rdmp.app/</a>).</li>
			<li>Command line tool.</li>
		</ul>
		<p><a type="button" class="btn btn-primary" href="https://github.com/peteraba/roadmapper" target="_blank">Github <i class="fab fa-github"></i></a></p>
		<hr class="hr">
	</div>

	<div class="container roadmap-privacy section" id="roadmap-privacy">
		<h1 class="h1">Data Privacy</h1>
		<div class="alert alert-danger" role="alert">
  			<h4 class="alert-heading">Important!</h4>
		  	<p>
				We do not recommend using <a href="https://rdmp.app/">rdmp.app</a> for business-critical use cases!<br>
				We store your Roadmaps in plain text and we technically make them publicly available!<br>
				We delete Roadmaps that have not been opened for over 1 year.
			</p>
			<hr>
		  	<p class="mb-0">Please host <a href="https://github.com/peteraba/roadmapper" target="_blank">Roadmapper <i class="fas fa-external-link-alt"></i></a> yourself, or use its command line features instead!</p>
		</div>
		<h2 class="h2">Usage</h2>
		<ul>
			<li><a href="https://rdmp.app/">rdmp.app</a> is an open source software run a <a href="https://en.wikipedia.org/wiki/Software_as_a_service">SaaS service <i class="fas fa-external-link-alt"></i></a>.</li>
			<li><a href="https://rdmp.app/">rdmp.app</a> itself does not track you and does not store personal information about you. This means that it will not and can not share your personal information with 3rd-party software.</li>
			<li><a href="https://rdmp.app/">rdmp.app</a> uses various Google Products which may track you. If this is a problem for you, we recommend you to run Roadmapper on your own, it comes without 3rd-party integrations by default.</li>
		</ul>
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
<!-- SVG -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/svg.js/3.0.16/svg.min.js" integrity="sha256-MCvBrhCuX8GNt0gmv06kZ4jGIi1R2QNaSkadjRzinFs=" crossorigin="anonymous"></script>
<!-- Custom -->
<script src="/static/roadmaper.js"></script>
</body>
</html>
`

func bootstrapRoadmap(roadmap Project, lines []string) (string, error) {
	writer := bytes.NewBufferString("")

	t, err := template.New("layout").Parse(layoutTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		Roadmap Project
		Raw     string
	}{
		Roadmap: roadmap,
		Raw:     strings.Join(lines, "\n"),
	}

	err = t.Execute(writer, data)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}
