package main

import (
	"time"

	"github.com/peteraba/go-svg"
)

func createSvg(roadmap Project, width, height float64, dateFormat string) svg.SVG {
	var children []interface{}

	children = append(children, createSvgHeader(roadmap.From, roadmap.To, width, height, 0, 0, dateFormat)...)

	return svg.NewSVG(width, height, children...)
}

func createSvgHeader(start, end time.Time, width, height, dx, dy float64, dateFormat string) []interface{} {
	var children []interface{}

	children = append(children, createSvgHeaderLines(width, height, dx, dy))
	children = append(children, createSvgHeaderDates(start, end, width, height, dx, dy, dateFormat))
	children = append(children, createSvgHeaderToday(start, end, width, height, dx, dy, dateFormat))

	return children
}

func createSvgHeaderLines(width, height, dx, dy float64) []interface{} {
	var children []interface{}

	strokeColor1, _ := svg.ParseHexaColor("#212529")
	children = append(children, svg.NewLine(width-10+dx, (height/2)-5+dy, width-10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
	children = append(children, svg.NewLine(0+dx, height/2+dy, width+dx, height/2+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
	children = append(children, svg.NewLine(10+dx, (height/2)-5+dy, 10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))

	return children
}

func createSvgHeaderDates(start, end time.Time, width, height, dx, dy float64, dateFormat string) []interface{} {
	var children []interface{}

	tspan1 := svg.NewTSpan(start.Format(dateFormat))
	children = append(children, svg.NewText(12+dx, (height/2)-10+dy, tspan1))

	tspan2 := svg.NewTSpan(end.Format(dateFormat))
	children = append(children, svg.NewText(width-12+dx, (height/2)-10+dy, tspan2).SetTextAnchor(svg.End))

	return children
}

func createSvgHeaderToday(start, end time.Time, width, height, dx, dy float64, dateFormat string) []interface{} {
	var children []interface{}

	if time.Until(end) < 0 {
		return children
	}

	now := time.Now()
	untilToday := now.Sub(start).Hours()
	startToEnd := end.Sub(start).Hours()

	pos := untilToday / startToEnd * width

	strokeColor2, _ := svg.ParseHexaColor("#ff8888")
	children = append(children, svg.NewLine(pos+dx, 0+dy, pos+dx, height+dy).SetStrokeWidth(2).SetStroke(strokeColor2))

	fillColor, _ := svg.ParseHexaColor("#ff3333")
	tspan3 := svg.NewTSpan(now.Format(dateFormat))
	children = append(children, svg.NewText(pos+dx, (height/2)+20+dy, tspan3).SetFill(fillColor).SetTextAnchor(svg.Middle))

	return children
}
