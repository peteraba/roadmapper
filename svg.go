package main

import (
	"image/color"
	"time"

	"github.com/peteraba/go-svg"
)

func createSvg(roadmap Project, fullWidth, headerHeight float64, dateFormat string) svg.SVG {
	var (
		roadmapDy, lineHeight = 0.0, 30.0
	)

	var elements []interface{}

	if roadmap.IsPlanned() {
		elements = append(elements, createSvgHeader(roadmap.Dates.Start, roadmap.Dates.End, fullWidth/3*2, headerHeight, fullWidth/3, 0, dateFormat)...)
		roadmapDy += headerHeight
	}

	roadmapElements, roadmapDy := createSvgRoadmap(roadmap, fullWidth, roadmapDy, lineHeight, dateFormat)

	elements = append(elements, roadmapElements...)

	if roadmap.IsPlanned() {
		elements = append(elements, createSvgTableLines(fullWidth, roadmapDy, fullWidth/3, headerHeight)...)
	}

	return svg.NewSVG(fullWidth, roadmapDy, elements...)
}

func createSvgHeader(start, end time.Time, width, height, dx, dy float64, dateFormat string) []interface{} {
	var elements []interface{}

	elements = append(elements, createSvgHeaderLines(width, height, dx, dy))
	elements = append(elements, createSvgHeaderDates(start, end, width, height, dx, dy, dateFormat))
	elements = append(elements, createSvgHeaderToday(start, end, width, height, dx, dy, dateFormat))

	return elements
}

func createSvgHeaderLines(width, height, dx, dy float64) []interface{} {
	var elements []interface{}

	strokeColor1, _ := svg.ParseHexaColor("#212529")
	elements = append(elements, svg.NewLine(width-10+dx, (height/2)-5+dy, width-10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
	elements = append(elements, svg.NewLine(0+dx, height/2+dy, width+dx, height/2+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
	elements = append(elements, svg.NewLine(10+dx, (height/2)-5+dy, 10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))

	return elements
}

func createSvgHeaderDates(start, end time.Time, width, height, dx, dy float64, dateFormat string) []interface{} {
	var elements []interface{}

	tspan1 := svg.NewTSpan(start.Format(dateFormat))
	elements = append(elements, svg.NewText(12+dx, (height/2)-10+dy, tspan1))

	tspan2 := svg.NewTSpan(end.Format(dateFormat))
	elements = append(elements, svg.NewText(width-12+dx, (height/2)-10+dy, tspan2).SetTextAnchor(svg.End))

	return elements
}

func createSvgHeaderToday(start, end time.Time, width, height, dx, dy float64, dateFormat string) []interface{} {
	var elements []interface{}

	if time.Until(end) < 0 {
		return elements
	}

	now := time.Now()
	untilToday := now.Sub(start).Hours()
	startToEnd := end.Sub(start).Hours()

	pos := untilToday / startToEnd * width

	strokeColor, _ := svg.ParseHexaColor("#ff8888")
	elements = append(elements, svg.NewLine(pos+dx, 0+dy, pos+dx, height+dy).SetStrokeWidth(2).SetStroke(strokeColor))

	fillColor, _ := svg.ParseHexaColor("#ff3333")
	tspan3 := svg.NewTSpan(now.Format(dateFormat))
	elements = append(elements, svg.NewText(pos+dx, (height/2)+20+dy, tspan3).SetFill(fillColor).SetTextAnchor(svg.Middle))

	return elements
}

func createSvgRoadmap(roadmap Project, fullWidth, dy, lineHeight float64, dateFormat string) ([]interface{}, float64) {
	var (
		elements, result []interface{}
	)

	for _, p := range roadmap.Children {
		elements, dy = createSvgProject(roadmap, p, fullWidth, dy, 10, lineHeight, 30, dateFormat)
		result = append(result, elements)
	}

	return result, dy
}

func createSvgProject(roadmap, project Project, fullWidth, dy, dx, lineHeight, indentWidth float64, dateFormat string) ([]interface{}, float64) {
	var (
		elements, result []interface{}
	)

	strokeColor, _ := svg.ParseHexaColor("#eee")

	result = append(result, svg.NewText(dx, dy+lineHeight/2+5, svg.NewTSpan(project.Title)))

	if roadmap.IsPlanned() {
		rd, pd := roadmap.Dates, project.Dates
		rs, rw := fullWidth/3+12, fullWidth/3*2-24
		ps := rs + (rw * pd.Start.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())
		pe := rs + (rw * pd.End.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())

		r, g, b, a := project.Color.RGBA()
		strokeColor := svg.Color{RGBA: color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}}
		result = append(result, svg.NewLine(ps, dy+lineHeight/2, pe, dy+lineHeight/2).SetStrokeWidth(15).SetStroke(strokeColor))
	}

	dy += lineHeight

	for _, c := range project.Children {
		result = append(result, svg.NewLine(0, dy, fullWidth, dy).SetStrokeWidth(1).SetStroke(strokeColor))

		elements, dy = createSvgProject(roadmap, c, fullWidth, dy, dx+indentWidth, lineHeight, indentWidth, dateFormat)
		result = append(result, elements)
	}

	result = append(result, svg.NewLine(0, dy, fullWidth, dy).SetStrokeWidth(1).SetStroke(strokeColor))

	return result, dy
}

func createSvgTableLines(fullWidth, fullHeight, headerX, headerHeight float64) []interface{} {
	var result []interface{}

	strokeColor, _ := svg.ParseHexaColor("#999")
	result = append(result, svg.NewLine(headerX, 0, headerX, fullHeight).SetStrokeWidth(1).SetStroke(strokeColor))
	result = append(result, svg.NewLine(0, headerHeight, fullWidth, headerHeight).SetStrokeWidth(1).SetStroke(strokeColor))

	return result
}
