package main

import (
	"fmt"
	"time"

	"github.com/tdewolff/canvas"
)

type VisualRoadmap struct {
	Projects   []Project
	Milestones []Milestone
	Dates      *Dates
}

func (r Roadmap) ToVisual() VisualRoadmap {
	visual := VisualRoadmap{}

	visual.Dates = r.ToDates()

	foundMilestones := map[int]*Milestone{}
	for i, p := range r.Projects {
		dates := findVisualDates(r.Projects, i)
		visual.Projects = append(
			visual.Projects,
			Project{
				Indentation: p.Indentation,
				Title:       p.Title,
				Dates:       dates,
				Color:       p.Color,
				Percentage:  p.Percentage,
				URLs:        p.URLs,
			},
		)
		if p.Milestone > 0 {
			mk := int(p.Milestone) - 1

			var endAt *time.Time
			if dates != nil {
				endAt = &dates.EndAt
			}

			milestone, ok := foundMilestones[mk]
			if !ok {
				foundMilestones[mk] = &Milestone{
					DeadlineAt: endAt,
					Color:      p.Color,
				}
				continue
			}

			if milestone.Color == nil && p.Color != nil {
				milestone.Color = p.Color
			}

			if endAt == nil {
				continue
			}

			if milestone.DeadlineAt == nil {
				milestone.DeadlineAt = endAt
			}

			if milestone.DeadlineAt.Before(*endAt) {
				milestone.DeadlineAt = endAt
			}
		}
	}

	visual.Milestones = r.Milestones

	for i, m := range foundMilestones {
		if len(visual.Milestones) < i {
			panic("original milestone not found")
		}

		om := &visual.Milestones[i]

		if om.Color == nil {
			om.Color = m.Color
		}

		if m.DeadlineAt == nil {
			continue
		}

		if om.DeadlineAt == nil {
			om.DeadlineAt = m.DeadlineAt
		}

		if om.DeadlineAt.Before(*m.DeadlineAt) {
			om.DeadlineAt = m.DeadlineAt
		}
	}

	return visual
}

func findVisualDates(projects []Project, start int) *Dates {
	if projects == nil || len(projects) < start {
		panic(fmt.Errorf("illegal start %d for finding visual dates", start))
	}

	if projects[start].Dates != nil {
		return &Dates{StartAt: projects[start].Dates.StartAt, EndAt: projects[start].Dates.EndAt}
	}

	minIndentation := projects[start].Indentation + 1

	var dates *Dates
	for i := start + 1; i < len(projects); i++ {
		p := projects[i]
		if p.Indentation < minIndentation {
			break
		}

		if p.Dates == nil {
			continue
		}

		if dates == nil {
			dates = &Dates{StartAt: p.Dates.StartAt, EndAt: p.Dates.EndAt}
			continue
		}

		if dates.StartAt.After(p.Dates.StartAt) {
			dates.StartAt = p.Dates.StartAt
		}

		if dates.EndAt.Before(p.Dates.EndAt) {
			dates.EndAt = p.Dates.EndAt
		}
	}

	return dates
}

var fontFamily *canvas.FontFamily

func (vr VisualRoadmap) Draw(fullW, headerH, lineH float64, dateFormat string) *canvas.Canvas {
	if vr.Dates == nil {
		headerH = 0.0
	}

	strokeW := 2.0
	fullH := lineH*float64(len(vr.Projects)) + headerH

	fontFamily = canvas.NewFontFamily("roboto")
	err := fontFamily.LoadFontFile("fonts/Roboto/Roboto-Regular.ttf", canvas.FontRegular)
	if err != nil {
		panic(fmt.Errorf("font not loaded: %w", err))
	}

	fontFamily.Use(canvas.CommonLigatures)

	c := canvas.New(fullW, fullH)

	ctx := canvas.NewContext(c)
	ctx.SetStrokeWidth(strokeW)

	vr.drawHeader(ctx, fullW, fullH, headerH, lineH, strokeW, dateFormat)

	vr.writeProjects(ctx, fullW, fullH, headerH, lineH)

	vr.drawMilestones(ctx, fullW, fullH, headerH, lineH, dateFormat)

	vr.drawLines(ctx, fullW, fullH, headerH, lineH)

	return c
}

func (vr VisualRoadmap) drawHeader(ctx *canvas.Context, fullW, fullH, headerH, lineH, strokeW float64, dateFormat string) {
	if vr.Dates == nil {
		return
	}

	vr.drawHeaderBaseline(ctx, fullW, fullH, headerH)

	vr.writeHeaderDates(ctx, fullW, fullH, lineH, dateFormat)

	vr.markHeaderDates(ctx, fullW, fullH, headerH, strokeW)
}

func (vr VisualRoadmap) drawHeaderBaseline(ctx *canvas.Context, fullW, fullH, headerH float64) {
	p := &canvas.Path{}
	p.MoveTo(0, 0)
	p.LineTo(fullW*2/3, 0)

	x := fullW / 3
	y := fullH - headerH/2
	ctx.SetStrokeColor(canvas.Black)
	ctx.DrawPath(x, y, p)
}

func (vr VisualRoadmap) writeHeaderDates(ctx *canvas.Context, fullW, fullH, lineH float64, dateFormat string) {
	x := fullW / 3
	y := fullH
	face := fontFamily.Face(lineH*1.5, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	date := vr.Dates.StartAt.Format(dateFormat)
	ctx.DrawText(x, y, canvas.NewTextBox(face, date, 0.0, lineH, canvas.Left, canvas.Center, 0.0, 0.0))

	x = fullW
	date = vr.Dates.EndAt.Format(dateFormat)
	ctx.DrawText(x, y, canvas.NewTextBox(face, date, 0.0, lineH, canvas.Right, canvas.Center, 0.0, 0.0))
}

func (vr VisualRoadmap) markHeaderDates(ctx *canvas.Context, fullW, fullH, headerH, strokeW float64) {
	markH := headerH / 10.0

	p1 := &canvas.Path{}
	p1.MoveTo(strokeW, markH/-2)
	p1.LineTo(strokeW, markH/2)

	p2 := &canvas.Path{}
	p2.MoveTo(fullW*2/3-strokeW, markH/-2)
	p2.LineTo(fullW*2/3-strokeW, markH/2)

	x := fullW / 3
	y := fullH - headerH/2
	ctx.SetStrokeColor(canvas.Black)
	ctx.DrawPath(x, y, p1, p2)
}

func (vr VisualRoadmap) writeProjects(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	textW := fullW / 3
	indentationW := textW / 20
	face := fontFamily.Face(lineH*1.5, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	ctx.SetFillColor(canvas.Black)

	for i, p := range vr.Projects {
		y := fullH - float64(i)*lineH - headerH - lineH/4
		indent := float64(p.Indentation) * indentationW
		ctx.DrawText(0, y, canvas.NewTextBox(face, p.Title, textW, lineH, canvas.Left, canvas.Top, indent, 0.0))
	}
}

func (vr VisualRoadmap) drawMilestones(ctx *canvas.Context, fullW, fullH, headerH, lineH float64, dateFormat string) {
	if vr.Dates == nil {
		return
	}

	y := fullH - headerH/2
	roadmapInterval := vr.Dates.EndAt.Sub(vr.Dates.StartAt).Hours()

	for _, m := range vr.Milestones {
		if m.DeadlineAt == nil {
			continue
		}

		c := canvas.Red
		if m.Color != nil {
			c = *m.Color
		}

		maxW := fullW * 2 / 3
		deadlineFromStart := vr.Dates.EndAt.Sub(*m.DeadlineAt).Hours()
		w := deadlineFromStart / roadmapInterval * maxW

		p := &canvas.Path{}
		p.MoveTo(w, 0)
		p.LineTo(w, fullH)

		ctx.SetStrokeColor(c)
		ctx.DrawPath(fullW/3, 0, p)

		x := w + fullW/3
		face := fontFamily.Face(lineH*1.5, c, canvas.FontRegular, canvas.FontNormal)
		date := m.DeadlineAt.Format(dateFormat)
		ctx.DrawText(x, y, canvas.NewTextBox(face, date, 0.0, lineH, canvas.Center, canvas.Center, 0.0, 0.0))
	}
}

func (vr VisualRoadmap) drawLines(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	vr.drawHeaderLine(ctx, fullW, fullH, headerH)
	vr.drawProjectLines(ctx, fullW, fullH, headerH, lineH)

}

func (vr VisualRoadmap) drawHeaderLine(ctx *canvas.Context, fullW, fullH, headerH float64) {
	if vr.Dates == nil {
		return
	}

	p := &canvas.Path{}
	p.MoveTo(0, fullH-headerH)
	p.LineTo(fullW, fullH-headerH)

	ctx.SetStrokeColor(canvas.Lightgray)
	ctx.DrawPath(0, 0, p)
}

func (vr VisualRoadmap) drawProjectLines(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	var paths []*canvas.Path

	for i := range vr.Projects {
		h := fullH - headerH - (float64(i) * lineH)

		p := &canvas.Path{}
		p.MoveTo(0, h)
		p.LineTo(fullW, h)
		paths = append(paths, p)
	}

	p := &canvas.Path{}
	p.MoveTo(fullW/3, 0)
	p.LineTo(fullW/3, fullH)
	paths = append(paths, p)

	ctx.SetStrokeColor(canvas.Lightgray)
	ctx.DrawPath(0, 0, paths...)
}
