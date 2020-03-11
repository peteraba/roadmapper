package main

import (
	"image/color"
	"strings"
	"time"

	svg "github.com/peteraba/go-svg"
)

func createSvg(roadmap Project, fullWidth, headerHeight, lineHeight float64, dateFormat string) svg.SVG {
	var (
		roadmapDy = 0.0
	)

	var elements []interface{}

	elements = append(elements, createStripesPattern(), createStyle())

	if roadmap.IsPlanned() {
		elements = append(elements, createSvgHeader(roadmap.Dates.Start, roadmap.Dates.End, fullWidth/3*2, headerHeight, fullWidth/3, 0, dateFormat))
		roadmapDy += headerHeight
	}

	svgRoadmap, roadmapDy := createSvgRoadmap(roadmap, fullWidth, roadmapDy, lineHeight, dateFormat)

	elements = append(elements, svgRoadmap)

	if roadmap.IsPlanned() {
		elements = append(elements, createSvgTableLines(fullWidth, roadmapDy, fullWidth/3, headerHeight))
	}

	return svg.NewSVG(fullWidth, roadmapDy, elements...)
}

func createStripesPattern() svg.Element {
	polygons := []interface{}{
		svg.E("polygon", "", "", map[string]string{"points": "0,4 0,8 8,0 4,0", "fill": "white"}),
		svg.E("polygon", "", "", map[string]string{"points": "4,8 8,8 8,4", "fill": "white"}),
	}
	pattern := svg.E("pattern", "", "", map[string]string{"id": "stripes", "viewBox": "0,0,8,8", "width": "16", "height": "16", "patternUnits": "userSpaceOnUse"}, polygons...)

	def := svg.E("defs", "", "", nil, pattern)

	return def
}

func createStyle() svg.Element {
	rules := []string{
		`tspan {font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,"Noto Sans",sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol","Noto Color Emoji";}`,
		`.strong {font-weight: bold;}`,
		`a {fill: #06D; text-decoration: underline; cursor: pointer;}`,
		`a:hover, a:active {outline: dotted 1px blue; color: #06D;}`,
	}
	el := svg.E("style", "", strings.Join(rules, "\n"), nil)

	return el
}

func createSvgHeader(start, end time.Time, width, height, dx, dy float64, dateFormat string) svg.Group {
	var elements []interface{}

	elements = append(elements, createSvgHeaderLines(width, height, dx, dy))
	elements = append(elements, createSvgHeaderDates(start, end, width, height, dx, dy, dateFormat))
	elements = append(elements, createSvgHeaderToday(start, end, width, height, dx, dy, dateFormat))

	return svg.NewGroup(elements)
}

func createSvgHeaderLines(width, height, dx, dy float64) svg.Group {
	var elements []interface{}

	strokeColor1, _ := svg.ColorFromHexaString("#212529")
	elements = append(elements, svg.L(width-10+dx, (height/2)-5+dy, width-10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
	elements = append(elements, svg.L(0+dx, height/2+dy, width+dx, height/2+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
	elements = append(elements, svg.L(10+dx, (height/2)-5+dy, 10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))

	return svg.NewGroup(elements)
}

func createSvgHeaderDates(start, end time.Time, width, height, dx, dy float64, dateFormat string) svg.Group {
	var elements []interface{}

	tspan1 := svg.TS(start.Format(dateFormat))
	elements = append(elements, svg.T(12+dx, (height/2)-10+dy, tspan1))

	tspan2 := svg.TS(end.Format(dateFormat))
	elements = append(elements, svg.T(width-12+dx, (height/2)-10+dy, tspan2).SetTextAnchor(svg.End))

	return svg.NewGroup(elements)
}

func createSvgHeaderToday(start, end time.Time, width, height, dx, dy float64, dateFormat string) svg.Group {
	var elements []interface{}

	if time.Until(end) < 0 {
		return svg.NewGroup()
	}

	now := time.Now()
	untilToday := now.Sub(start).Hours()
	startToEnd := end.Sub(start).Hours()

	pos := untilToday / startToEnd * width

	strokeColor, _ := svg.ColorFromHexaString("#ff8888")
	elements = append(elements, svg.L(pos+dx, 0+dy, pos+dx, height+dy).SetStrokeWidth(2).SetStroke(strokeColor))

	fillColor, _ := svg.ColorFromHexaString("#ff3333")
	tspan3 := svg.TS(now.Format(dateFormat))
	elements = append(elements, svg.T(pos+dx, (height/2)+20+dy, tspan3).SetFill(fillColor).SetTextAnchor(svg.Middle))

	return svg.NewGroup(elements)
}

func createSvgRoadmap(roadmap Project, fullWidth, dy, lineHeight float64, dateFormat string) (svg.Group, float64) {
	var (
		project  svg.Group
		projects []svg.Group
	)

	for _, p := range roadmap.Children {
		project, dy = createSvgProject(roadmap, p, fullWidth, dy, 10, lineHeight, 30, dateFormat)
		projects = append(projects, project)
	}

	return svg.NewGroup(projects), dy
}

func createSvgProject(roadmap, project Project, fullWidth, dy, dx, lineHeight, indentWidth float64, dateFormat string) (svg.Group, float64) {
	var (
		subProject  svg.Group
		subProjects []interface{}
	)

	subProjects = append(subProjects, createProjectTitle(project, dx, dy, lineHeight))

	if project.IsPlanned() {
		subProjects = append(subProjects, createProjectVisual(*roadmap.Dates, project, fullWidth, dy, lineHeight)...)
	}

	dy += lineHeight

	strokeColor, _ := svg.ColorFromHexaString("#eee")
	for _, c := range project.Children {
		subProjects = append(subProjects, svg.L(0, dy, fullWidth, dy).SetStrokeWidth(1).SetStroke(strokeColor))

		subProject, dy = createSvgProject(roadmap, c, fullWidth, dy, dx+indentWidth, lineHeight, indentWidth, dateFormat)
		subProjects = append(subProjects, subProject)
	}

	subProjects = append(subProjects, svg.L(0, dy, fullWidth, dy).SetStrokeWidth(1).SetStroke(strokeColor))

	return svg.NewGroup(subProjects), dy
}

func createProjectTitle(project Project, dx, dy, lineHeight float64) svg.Text {
	title := svg.TS(project.Title).AddAttr("class", "strong")
	if project.URL != "" {
		a := svg.NewA(project.URL, svg.TS(" "), title)
		return svg.T(dx, dy+lineHeight/2+5, a)
	}

	return svg.T(dx, dy+lineHeight/2+5, title)
}

func createProjectVisual(roadmapDates Dates, project Project, fullWidth, dy, lineHeight float64) []interface{} {
	wl := lineHeight * 0.6
	rd, pd := roadmapDates, project.Dates
	rs, rw := fullWidth/3+12, fullWidth/3*2-24
	ps := rs + (rw * pd.Start.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())
	pe := rs + (rw * pd.End.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())

	r, g, b, _ := project.Color.RGBA()
	strokeColor := svg.Color{RGBA: color.RGBA{uint8(r), uint8(g), uint8(b), 255}}
	s := svg.R(ps, dy+lineHeight/2-wl/2, pe-ps, wl).
		SetFill(strokeColor).
		AddAttr("rx", "5").
		AddAttr("ry", "5")

	s2 := svg.R(ps, dy+lineHeight/2-wl/2, pe-ps, wl).
		AddAttr("rx", "5").
		AddAttr("ry", "5").
		AddAttr("fill", `url(#stripes)`).
		SetFillOpacity(svg.Opacity{Number: .2})

	if project.Percentage >= 100 {
		s = s.SetFillOpacity(svg.Opacity{Number: .6})
	}

	return []interface{}{s, s2}
}

func createSvgTableLines(fullWidth, fullHeight, headerX, headerHeight float64) []interface{} {
	var result []interface{}

	strokeColor, _ := svg.ColorFromHexaString("#999")
	result = append(result, svg.L(headerX, 0, headerX, fullHeight).SetStrokeWidth(1).SetStroke(strokeColor))
	result = append(result, svg.L(0, headerHeight, fullWidth, headerHeight).SetStrokeWidth(1).SetStroke(strokeColor))

	return result
}
