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

    <!-- Bootstrap CSS -->
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha256-L/W5Wfqfa0sdBNIKN9cG6QA5F2qx4qICmU2VgLruv9Y=" crossorigin="anonymous" />
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.12.1/css/all.min.css" integrity="sha256-mmgLkCYLUQbXn0B1SRqzHar6dCnv9oZFPEC1g1cwlkk=" crossorigin="anonymous" />
	<link rel="stylesheet" href="/static/roadmaper.css">

    <title>{{.Roadmap.Title}}</title>
</head>
<body>
	<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
		<a class="navbar-brand" href="#">Roadmapper</a>
		<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarNavAltMarkup">
			<div class="navbar-nav">
				<a class="nav-item nav-link active" href="#roadmap-dashboard">Dashboard <span class="sr-only">(current)</span></a>
				<a class="nav-item nav-link" href="#roadmap-edit">Edit</a>
			</div>
		</div>
	</nav>
	<div class="container-fluid roadmap-dashboard section" id="roadmap-dashboard">
		<table class="table table-bordered roadmap" id="roadmap">
		  <!-- Content here -->
		</table>
		<form id="control"></form>
	</div>
	<div class="container-fluid roadmap-edit section" id="roadmap-edit">
		<form action="" method="POST">
			<input type="hidden" name="_method" value="PUT" />
			<div class="form-group">
				<label for="roadmapRaw">Raw roadmap</label>
				<textarea class="form-control" id="roadmapRaw" name="roadmap" aria-describedby="roadmapRaw" rows="20">{{ .Raw }}</textarea>
				<small id="roadmapRaw" class="form-text text-muted">We'll never share your email with anyone else.</small>
			</div>
			<button type="submit" class="btn btn-primary">Submit</button>
		</form>
	</div>

	<script>
	var roadmap = {{ .Roadmap }};
	</script>
	
	<!-- Optional JavaScript -->
	<!-- jQuery first, then Popper.js, then Bootstrap JS -->
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.4.1/jquery.min.js" integrity="sha256-CSXorXvZcTkaix6Yvo6HppcZGetbYMGWSFlBw8HfCJo=" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.16.1/umd/popper.min.js" integrity="sha256-/ijcOLwFf26xEYAjW75FizKVo5tnTYiQddPZoLUHHZ8=" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha256-WqU1JavFxSAMcLP2WIOI+GB2zWmShMI82mTpLDcqFUg=" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/svg.js/3.0.16/svg.min.js" integrity="sha256-MCvBrhCuX8GNt0gmv06kZ4jGIi1R2QNaSkadjRzinFs=" crossorigin="anonymous"></script>
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
