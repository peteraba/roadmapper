package main

import (
	"time"
)

func createSvgHeader(start, end time.Time, width, height float64, dateFormat string) SVG {
	var children []interface{}

	children = append(children, createSvgHeaderLines(width, height))
	children = append(children, createSvgHeaderDates(start, end, width, height, dateFormat))
	children = append(children, createSvgHeaderToday(start, end, width, height, dateFormat))

	return NewSVG(width, height, children...)
}

func createSvgHeaderLines(width, height float64) []interface{} {
	var children []interface{}

	strokeColor1 := NewColorFromHexa("#212529")
	children = append(children, NewLine(width-10, (height/2)-5, width-10, (height/2)+5).SetStrokeWidth(2).SetStroke(strokeColor1))
	children = append(children, NewLine(0, height/2, width, height/2).SetStrokeWidth(2).SetStroke(strokeColor1))
	children = append(children, NewLine(10, (height/2)-5, 10, (height/2)+5).SetStrokeWidth(2).SetStroke(strokeColor1))

	return children
}

func createSvgHeaderDates(start, end time.Time, width, height float64, dateFormat string) []interface{} {
	var children []interface{}

	tspan1 := NewTSpan(start.Format(dateFormat)).SetX(6)
	children = append(children, NewText(6, (height/2)-10, tspan1))

	tspan2 := NewTSpan(end.Format(dateFormat)).SetX(911)
	children = append(children, NewText(6, (height/2)-10, tspan2).SetTextAnchor(End))

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

	strokeColor2 := NewColorFromHexa("#ff8888")
	children = append(children, NewLine(pos, 0, pos, height).SetStrokeWidth(2).SetStroke(strokeColor2))

	fillColor := NewColorFromHexa("#ff3333")
	tspan3 := NewTSpan(now.Format(dateFormat)).SetX(pos)
	children = append(children, NewText(6, (height/2)+20, tspan3).SetFill(fillColor).SetTextAnchor(Middle))

	return children
}
