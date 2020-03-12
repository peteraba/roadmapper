package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	svg "github.com/peteraba/go-svg"
)

func createSvg(roadmap Project, fullWidth, headerHeight, lineHeight float64, dateFormat string) svg.SVG {
	var (
		fullHeight float64
		elements   []interface{}
	)

	if !roadmap.IsPlanned() {
		headerHeight = 0
	}

	elements = append(elements, createStripesPattern(), createStyle())

	titles, fullHeight := createSvgTitles(roadmap, fullWidth/3, headerHeight, lineHeight)

	elements = append(elements, titles)

	if roadmap.IsPlanned() {
		elements = append(elements, createSvgVisuals(roadmap, fullWidth, headerHeight, lineHeight, dateFormat))
		elements = append(elements, createSvgHeader(roadmap.Dates.Start, roadmap.Dates.End, fullHeight, fullWidth/3*2, headerHeight, fullWidth/3, 0, dateFormat))
	}

	elements = append(elements, createSvgTableLines(fullWidth, fullHeight, fullWidth/3, headerHeight, lineHeight))

	return svg.NewSVG(fullWidth, fullHeight, elements...)
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

func createSvgHeader(start, end time.Time, fullHeight, headerWidth, headerHeight, dx, dy float64, dateFormat string) svg.Group {
	var elements []interface{}

	elements = append(elements, createSvgHeaderLines(headerWidth, headerHeight, dx, dy))
	elements = append(elements, createSvgHeaderDates(start, end, headerWidth, headerHeight, dx, dy, dateFormat))
	elements = append(elements, createSvgHeaderToday(start, end, fullHeight, headerWidth, headerHeight, dx, dy, dateFormat))

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

func createSvgHeaderToday(start, end time.Time, fullHeight, width, height, dx, dy float64, dateFormat string) svg.Group {
	var elements []interface{}

	if time.Until(end) < 0 {
		return svg.NewGroup()
	}

	now := time.Now()
	untilToday := now.Sub(start).Hours()
	startToEnd := end.Sub(start).Hours()

	pos := untilToday / startToEnd * width

	strokeColor, _ := svg.ColorFromHexaString("#ff3333")
	elements = append(elements, svg.L(pos+dx, 0, pos+dx, fullHeight).SetStrokeWidth(2).SetStroke(strokeColor).AddAttr("stroke-dasharray", "10,10"))

	fillColor, _ := svg.ColorFromHexaString("#ff3333")
	tspan3 := svg.TS(now.Format(dateFormat))
	elements = append(elements, svg.T(pos+dx, (height/2)+20+dy, tspan3).SetFill(fillColor).SetTextAnchor(svg.Middle))

	return svg.NewGroup(elements)
}

func createSvgTitles(roadmap Project, titlesWidth, dy, lineHeight float64) (svg.Element, float64) {
	var (
		project  svg.Group
		projects []svg.Group
	)

	for _, p := range roadmap.Children {
		project, dy = createProjectTitle(roadmap, p, dy, 10, lineHeight, 30)
		projects = append(projects, project)
	}

	return svg.E("svg", "", "", map[string]string{"width": fs(titlesWidth), "height": fs(dy), "x": "0", "y": "0"}, projects), dy
}

func createProjectTitle(roadmap, project Project, dy, dx, lineHeight, indentWidth float64) (svg.Group, float64) {
	var (
		subProject  svg.Group
		subProjects []interface{}
	)

	title := createProjectTitleText(project, dx, dy, lineHeight)

	subProjects = append(subProjects, title)

	dy += lineHeight

	for _, c := range project.Children {
		subProject, dy = createProjectTitle(roadmap, c, dy, dx+indentWidth, lineHeight, indentWidth)
		subProjects = append(subProjects, subProject)
	}

	return svg.NewGroup(subProjects), dy
}

func createProjectTitleText(project Project, dx, dy, lineHeight float64) svg.Text {
	title := svg.TS(project.Title).AddAttr("class", "strong")
	if project.URL != "" {
		a := svg.NewA(project.URL, svg.TS(" "), title)
		return svg.T(dx, dy+lineHeight/2+5, a)
	}

	return svg.T(dx, dy+lineHeight/2+5, title)
}

func fs(num float64) string {
	return fmt.Sprintf("%v", num)
}

func createSvgVisuals(roadmap Project, fullWidth, dy, lineHeight float64, dateFormat string) svg.Group {
	var (
		project  svg.Group
		projects []svg.Group
	)

	for _, p := range roadmap.Children {
		project, dy = createSvgProjectVisuals(roadmap, p, fullWidth, dy, 10, lineHeight, 30, dateFormat)
		projects = append(projects, project)
	}

	return svg.NewGroup(projects)
}

func createSvgProjectVisuals(roadmap, project Project, fullWidth, dy, dx, lineHeight, indentWidth float64, dateFormat string) (svg.Group, float64) {
	var (
		subProject  svg.Group
		subProjects []interface{}
	)

	if project.IsPlanned() {
		subProjects = append(subProjects, createProjectVisual(*roadmap.Dates, project, fullWidth, dy, lineHeight, dateFormat)...)
	}

	dy += lineHeight

	for _, c := range project.Children {
		subProject, dy = createSvgProjectVisuals(roadmap, c, fullWidth, dy, dx+indentWidth, lineHeight, indentWidth, dateFormat)
		subProjects = append(subProjects, subProject)
	}

	return svg.NewGroup(subProjects), dy
}

func createProjectVisual(roadmapDates Dates, project Project, fullWidth, dy, lineHeight float64, dateFormat string) []interface{} {
	wl := lineHeight * 0.6
	rd, pd := roadmapDates, project.Dates
	rs, rw := fullWidth/3+12, fullWidth/3*2-24
	ps := rs + (rw * pd.Start.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())
	pe := rs + (rw * pd.End.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())

	r, g, b, _ := project.Color.RGBA()

	baseColor, _ := svg.ColorFromHexaString("#dedede")
	base := svg.R(ps, dy+lineHeight/2-wl/2, pe-ps, wl).
		SetFill(baseColor).
		AddAttr("rx", "5").
		AddAttr("ry", "5")

	stripesColor := svg.Color{RGBA: color.RGBA{uint8(r), uint8(g), uint8(b), 255}}
	stripesBase := svg.R(0, 0, pe-ps, wl).
		SetFill(stripesColor).
		AddAttr("rx", "5").
		AddAttr("ry", "5")

	start := project.Dates.Start
	end := project.Dates.End
	tooltip := fmt.Sprintf("%d%%, %s - %s, %d days", project.Percentage, start.Format(dateFormat), end.Format(dateFormat), int64(end.Sub(start).Hours()/24))
	title := svg.E("title", "", tooltip, nil)
	stripes := svg.R(0, 0, pe-ps, wl, title).
		AddAttr("rx", "5").
		AddAttr("ry", "5").
		AddAttr("fill", `url(#stripes)`).
		SetFillOpacity(svg.Opacity{Number: .2})

	stripesWidth := fs((pe - ps) * float64(project.Percentage) / 100)
	stripesY := fs(dy + lineHeight/2 - wl/2)
	stripesContainer := svg.E("svg", "", "", map[string]string{"width": stripesWidth, "height": fs(wl), "x": fs(ps), "y": stripesY}, stripesBase, stripes)

	if project.Percentage >= 100 {
		stripesBase = stripesBase.SetFillOpacity(svg.Opacity{Number: .4})
	}

	return []interface{}{base, stripesContainer}
}

func createSvgTableLines(fullWidth, fullHeight, headerX, headerHeight, lineHeight float64) []interface{} {
	var result []interface{}

	light, _ := svg.ColorFromHexaString("#eee")
	for y := headerHeight + lineHeight; y < fullHeight; y += lineHeight {
		result = append(result, svg.L(0, y, fullWidth, y).SetStrokeWidth(2).SetStroke(light))
	}

	dark, _ := svg.ColorFromHexaString("#999")
	result = append(result, svg.L(headerX, 0, headerX, fullHeight).SetStrokeWidth(1).SetStroke(dark))
	result = append(result, svg.L(0, headerHeight, fullWidth, headerHeight).SetStrokeWidth(1).SetStroke(dark))

	return result
}
