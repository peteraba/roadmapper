package main

import (
	"time"

	"github.com/peteraba/go-svg"
)

func createSvgHeader(start, end time.Time, width, height float64, dateFormat string) svg.SVG {
	var children []interface{}

	children = append(children, createSvgHeaderLines(width, height))
	children = append(children, createSvgHeaderDates(start, end, width, height, dateFormat))
	children = append(children, createSvgHeaderToday(start, end, width, height, dateFormat))

	return svg.NewSVG(width, height, children...)
}

func createSvgHeaderLines(width, height float64) []interface{} {
	var children []interface{}

	strokeColor1, _ := svg.ParseHexaColor("#212529")
	children = append(children, svg.NewLine(width-10, (height/2)-5, width-10, (height/2)+5).SetStrokeWidth(2).SetStroke(strokeColor1))
	children = append(children, svg.NewLine(0, height/2, width, height/2).SetStrokeWidth(2).SetStroke(strokeColor1))
	children = append(children, svg.NewLine(10, (height/2)-5, 10, (height/2)+5).SetStrokeWidth(2).SetStroke(strokeColor1))

	return children
}

func createSvgHeaderDates(start, end time.Time, width, height float64, dateFormat string) []interface{} {
	var children []interface{}

	tspan1 := svg.NewTSpan(start.Format(dateFormat)).SetX(6)
	children = append(children, svg.NewText(6, (height/2)-10, tspan1))

	tspan2 := svg.NewTSpan(end.Format(dateFormat)).SetX(911)
	children = append(children, svg.NewText(6, (height/2)-10, tspan2).SetTextAnchor(svg.End))

	return children
}

func createSvgHeaderToday(start, end time.Time, width, height float64, dateFormat string) []interface{} {
	var children []interface{}

	if time.Until(end) < 0 {
		return children
	}

	now := time.Now()
	untilToday := now.Sub(start).Hours()
	startToEnd := end.Sub(start).Hours()

	pos := untilToday / startToEnd * width

	strokeColor2, _ := svg.ParseHexaColor("#ff8888")
	children = append(children, svg.NewLine(pos, 0, pos, height).SetStrokeWidth(2).SetStroke(strokeColor2))

	fillColor, _ := svg.ParseHexaColor("#ff3333")
	tspan3 := svg.NewTSpan(now.Format(dateFormat)).SetX(pos)
	children = append(children, svg.NewText(6, (height/2)+20, tspan3).SetFill(fillColor).SetTextAnchor(svg.Middle))

	return children
}
