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
var strokeWidth = 1.0
var cr = canvas.ButtCap
var jr = canvas.RoundJoin

func (vr VisualRoadmap) Draw(fullW, headerH, lineH float64) *canvas.Canvas {
	if vr.Dates == nil {
		headerH = 0.0
	}

	fullH := lineH*float64(len(vr.Projects)) + headerH

	fontFamily = canvas.NewFontFamily("roboto")
	err := fontFamily.LoadFontFile("fonts/Roboto/Roboto-Regular.ttf", canvas.FontRegular)
	if err != nil {
		panic(fmt.Errorf("font not loaded: %w", err))
	}

	fontFamily.Use(canvas.CommonLigatures)

	c := canvas.New(fullW, fullH)

	ctx := canvas.NewContext(c)

	vr.drawLines(ctx, fullW, fullH, headerH, lineH)

	vr.drawHeader(ctx, fullW, fullH, headerH)

	vr.writeProjects(ctx, fullW, fullH, headerH, lineH)

	vr.drawMilestones(ctx, fullW, fullH)

	return c
}

func (vr VisualRoadmap) drawLines(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	var paths []*canvas.Path

	if vr.Dates != nil {
		p := &canvas.Path{}
		p.Stroke(strokeWidth, cr, jr)
		p.MoveTo(0, fullH-headerH)
		p.LineTo(fullW, fullH-headerH)
		paths = append(paths, p)
	}

	for i := range vr.Projects {
		h := fullH - headerH - (float64(i) * lineH)

		p := &canvas.Path{}
		p.Stroke(strokeWidth, cr, jr)
		p.MoveTo(0, h)
		p.LineTo(fullW, h)
		paths = append(paths, p)
	}

	p := &canvas.Path{}
	p.Stroke(strokeWidth, cr, jr)
	p.MoveTo(fullW/3, 0)
	p.LineTo(fullW/3, fullH)
	paths = append(paths, p)

	ctx.SetStrokeColor(canvas.Lightgray)
	ctx.DrawPath(0, 0, paths...)
}

func (vr VisualRoadmap) drawHeader(ctx *canvas.Context, fullW, fullH, headerH float64) {
	if vr.Dates == nil {
		return
	}

	p := &canvas.Path{}
	p.Stroke(strokeWidth, cr, jr)
	p.MoveTo(0, 0)
	p.LineTo(fullW*2/3, 0)

	ctx.SetStrokeColor(canvas.Black)
	ctx.DrawPath(fullW/3, fullH-headerH*2/3, p)
}

func (vr VisualRoadmap) writeProjects(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	face := fontFamily.Face(lineH*1.5, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	ctx.SetFillColor(canvas.Black)

	for i, p := range vr.Projects {
		y := fullH - float64(i)*lineH - headerH - lineH/4
		ctx.DrawText(0, y, canvas.NewTextBox(face, p.Title, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	}
}

func (vr VisualRoadmap) drawMilestones(ctx *canvas.Context, fullW, fullH float64) {
	if vr.Dates == nil {
		return
	}

	fullSub := vr.Dates.EndAt.Sub(vr.Dates.StartAt)

	for _, m := range vr.Milestones {
		if m.DeadlineAt == nil {
			continue
		}

		c := canvas.Red
		if m.Color != nil {
			c = *m.Color
		}

		maxW := fullW * 2 / 3
		msSub := vr.Dates.EndAt.Sub(*m.DeadlineAt)
		w := float64(fullSub/msSub) * maxW

		p := &canvas.Path{}
		p.Stroke(strokeWidth, cr, jr)
		p.MoveTo(w, 0)
		p.LineTo(w, fullH)

		ctx.SetStrokeColor(c)
		ctx.DrawPath(fullW/3, 0, p)
	}
}
