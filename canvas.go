package main

import (
	"errors"
	"image/color"
	"time"

	"github.com/tdewolff/canvas"
)

func createImg(roadmap Roadmap, fullWidth, headerHeight, lineHeight float64, dateFormat string) *canvas.Canvas {
	fullHeight := headerHeight

	c := canvas.New(fullWidth, headerHeight)

	ctx := canvas.NewContext(c)

	if roadmap.IsPlanned() {
		drawHeader(ctx, roadmap.Dates.Start, roadmap.Dates.End, fullHeight, fullWidth/3*2, headerHeight, fullWidth/3, 0, dateFormat)
	}

	c.Fit(1.0)

	return c
}

func drawHeader(ctx *canvas.Context, start, end time.Time, fullHeight, headerWidth, headerHeight, dx, dy float64, dateFormat string) {
	drawHeaderLines(ctx, headerWidth, headerHeight, dx, dy)
}

func drawHeaderLines(ctx *canvas.Context, width, height, dx, dy float64) {
	p := &canvas.Path{}
	p.MoveTo(width-10+dx, (height/2)-5+dy).LineTo(width-10+dx, (height/2)+5+dy)
	p.MoveTo(0+dx, height/2+dy).LineTo(width+dx, height/2+dy)
	p.MoveTo(10+dx, (height/2)-5+dy).LineTo(10+dx, (height/2)+5+dy)

	stroke := p.Stroke(2.0, canvas.SquareCap, canvas.MiterJoin)

	fillColor, err := colorFromHexa("#212529")
	if err != nil {
		return
	}

	ctx.SetFillColor(fillColor)
	ctx.DrawPath(0, 0, stroke)
}

func colorFromHexa(s string) (color.RGBA, error) {
	if len(s) != 4 && len(s) != 7 {
		return color.RGBA{}, errors.New("invalid hexa color length")
	}

	if s[0] != '#' {
		return color.RGBA{}, errors.New("invalid first character for hexa color")
	}

	us, err := charsToUint8(s[1:])
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{R: us[0], G: us[1], B: us[2], A: 255}, nil
}
