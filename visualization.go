package main

import (
	"fmt"
	"image/color"
	"net/url"
	"strings"
	"time"

	"github.com/tdewolff/canvas"
)

type VisualRoadmap struct {
	Projects   []Project
	Milestones []Milestone
	Dates      *Dates
	DateFormat string
}

// ToVisual converts a roadmap to a visual roadmap
// main difference between the two is that the visual roadmap
// will contain calculated values where possible
func (r Roadmap) ToVisual() *VisualRoadmap {
	visual := &VisualRoadmap{}

	visual.Dates = r.ToDates()
	visual.Projects = r.Projects
	visual.Milestones = r.Milestones
	visual.DateFormat = r.DateFormat

	visual.calculateProjectDates().calculateProjectColors().calculatePercentages().applyBaseURL(r.BaseURL)

	projectMilestones := visual.collectProjectMilestones()
	visual.applyProjectMilestone(projectMilestones)

	return visual
}

// caluclateProjectDates tries to find reasonable dates for all projects
// first it tries to find dates bottom up, meaning that based on the sub-projects
// then it tries to find dates top down, meaning that it will copy over dates from parents
func (vr *VisualRoadmap) calculateProjectDates() *VisualRoadmap {
	for i := range vr.Projects {
		p := &vr.Projects[i]

		if p.Dates != nil {
			continue
		}

		p.Dates = vr.findDatesBottomUp(i)
	}

	for i := range vr.Projects {
		p := &vr.Projects[i]

		if p.Dates != nil {
			continue
		}

		p.Dates = vr.findDatesTopDown(i)
	}

	return vr
}

// findDatesBottomUp will look for the minimum start date and maximum end date of sub projects
func (vr *VisualRoadmap) findDatesBottomUp(start int) *Dates {
	if vr.Projects == nil || len(vr.Projects) < start {
		panic(fmt.Errorf("illegal start %d for finding visual dates", start))
	}

	if vr.Projects[start].Dates != nil {
		return vr.Projects[start].Dates
	}

	minIndentation := vr.Projects[start].Indentation + 1

	var dates *Dates
	for i := start + 1; i < len(vr.Projects); i++ {
		p := vr.Projects[i]
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

// findDatesTopDown will return the dates of the first ancestor it finds
func (vr *VisualRoadmap) findDatesTopDown(start int) *Dates {
	if vr.Projects == nil || len(vr.Projects) < start {
		panic(fmt.Errorf("illegal start %d for finding visual dates", start))
	}

	if vr.Projects[start].Dates != nil {
		return vr.Projects[start].Dates
	}

	currentIndentation := vr.Projects[start].Indentation

	var dates *Dates
	for i := start - 1; i >= 0; i-- {
		p := vr.Projects[i]
		if p.Indentation >= currentIndentation {
			continue
		}

		if p.Dates != nil {
			return p.Dates
		}

		currentIndentation = p.Indentation
	}

	return dates
}

// calculateProjectColors will set a color for each projects without one
func (vr *VisualRoadmap) calculateProjectColors() *VisualRoadmap {
	epicCount := -1
	taskCount := -1
	for i := range vr.Projects {
		p := &vr.Projects[i]

		if p.Indentation == 0 {
			epicCount++
			taskCount = -1
		}
		taskCount++

		c := p.Color
		if c == nil {
			c = pickFgColor(epicCount, taskCount, int(p.Indentation))
		}

		p.Color = c
	}

	return vr
}

// calculatePercentages will try to calculate the percentage of all projects without a percentage set bottom up,
// meaning looking at their subprojects
func (vr *VisualRoadmap) calculatePercentages() *VisualRoadmap {
	for i := range vr.Projects {
		p := &vr.Projects[i]

		if p.Percentage != 0 {
			continue
		}

		p.Percentage = vr.findPercentageBottomUp(i)
	}

	return vr
}

// findPercentageBottomUp will calculate the average percentage of sub-projects
func (vr *VisualRoadmap) findPercentageBottomUp(start int) uint8 {
	if vr.Projects == nil || len(vr.Projects) < start {
		panic(fmt.Errorf("illegal start %d for finding visual dates", start))
	}

	if vr.Projects[start].Percentage != 0 {
		return vr.Projects[start].Percentage
	}

	matchIndentation := vr.Projects[start].Indentation + 1

	var sum, count uint8
	for i := start + 1; i < len(vr.Projects); i++ {
		p := &vr.Projects[i]
		if p.Indentation < matchIndentation {
			break
		}

		if p.Indentation > matchIndentation {
			continue
		}

		if p.Percentage == 0 {
			p.Percentage = vr.findPercentageBottomUp(i)
		}

		sum += p.Percentage
		count++
	}

	if count == 0 {
		return 0
	}

	return sum / count
}

func (vr *VisualRoadmap) applyBaseURL(baseUrl string) *VisualRoadmap {
	if baseUrl == "" {
		return vr
	}

	for i := range vr.Projects {
		p := &vr.Projects[i]

		for j := range p.URLs {
			u := &p.URLs[j]

			parsedUrl, err := url.ParseRequestURI(*u)
			if err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "" {
				continue
			}

			*u = fmt.Sprintf("%s/%s", strings.TrimRight(baseUrl, "/"), strings.TrimLeft(*u, "/"))
		}
	}

	for i := range vr.Milestones {
		m := &vr.Milestones[i]

		for j := range m.URLs {
			u := &m.URLs[j]

			parsedUrl, err := url.ParseRequestURI(*u)
			if err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "" {
				continue
			}

			*u = fmt.Sprintf("%s/%s", strings.TrimRight(baseUrl, "/"), strings.TrimLeft(*u, "/"))
		}
	}

	return vr
}

// collectProjectMilestones creates temporary milestones based on project information
// these will then be used in applyProjectMilestones as default values from milestones
// at the moment only colors and deadlines are collected
func (vr *VisualRoadmap) collectProjectMilestones() map[int]*Milestone {
	foundMilestones := map[int]*Milestone{}
	for i := range vr.Projects {
		p := &vr.Projects[i]

		if p.Milestone == 0 {
			continue
		}

		if p.Color == nil && p.Dates == nil {
			continue
		}

		mk := int(p.Milestone) - 1

		var endAt *time.Time
		if p.Dates != nil {
			endAt = &p.Dates.EndAt
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

	return foundMilestones
}

// applyProjectMilestone will apply temporary milestones created in collectProjectMilestones
// as default values for milestones
// this means that if milestones don't have deadlines or colors set, they can be set using what was
// found or generated for the projects linked to a milestone
func (vr *VisualRoadmap) applyProjectMilestone(projectMilestones map[int]*Milestone) *VisualRoadmap {
	for i, m := range projectMilestones {
		if len(vr.Milestones) < i {
			panic("original milestone not found")
		}

		om := &vr.Milestones[i]

		if om.Color == nil {
			om.Color = m.Color
		}

		if om.DeadlineAt != nil {
			continue
		}

		if m.DeadlineAt == nil {
			continue
		}

		om.DeadlineAt = m.DeadlineAt
	}

	for i := range vr.Milestones {
		if vr.Milestones[i].Color == nil {
			vr.Milestones[i].Color = defaultMilestoneColor
		}
	}

	return vr
}

var fontFamily *canvas.FontFamily
var myLightGrey = color.RGBA{R: 220, G: 220, B: 220, A: 255}
var myDarkGray = color.RGBA{R: 95, G: 95, B: 95, A: 255}
var defaultMilestoneColor = &canvas.Darkgray

// Draw will draw a roadmap on a canvas.Canvas
func (vr *VisualRoadmap) Draw(fullW, lineH float64) *canvas.Canvas {
	headerH := 0.0
	if vr.Dates != nil {
		headerH = lineH * 3
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

	vr.drawBackground(ctx, fullW, fullH, headerH)
	vr.drawHeader(ctx, fullW, fullH, headerH, lineH, strokeW)

	vr.drawProjectBackgrounds(ctx, fullW, fullH, headerH, lineH)
	vr.writeProjects(ctx, fullW, fullH, headerH, lineH)
	vr.drawProjects(ctx, fullW, fullH, headerH, lineH, strokeW)

	vr.drawMilestones(ctx, fullW, fullH, headerH, lineH)

	vr.drawLines(ctx, fullW, fullH, headerH, lineH)

	vr.drawToday(ctx, fullW, fullH, lineH)

	vr.drawLogo(ctx, fullW, fullH, lineH)

	return c
}

func (vr *VisualRoadmap) drawBackground(ctx *canvas.Context, fullW, fullH, headerH float64) {
	p := &canvas.Path{}
	p.MoveTo(0, 0)
	p.LineTo(0, fullH+headerH)
	p.LineTo(fullW, fullH+headerH)
	p.LineTo(fullW, 0)
	p.LineTo(0, 0)
	p.Close()

	ctx.SetFillColor(canvas.White)
	ctx.DrawPath(0, 0, p)
}

func (vr *VisualRoadmap) drawHeader(ctx *canvas.Context, fullW, fullH, headerH, lineH, strokeW float64) {
	if vr.Dates == nil {
		return
	}

	vr.drawHeaderBaseline(ctx, fullW, fullH, headerH)

	vr.writeHeaderDates(ctx, fullW, fullH, lineH)

	vr.markHeaderDates(ctx, fullW, fullH, headerH, strokeW)
}

func (vr *VisualRoadmap) drawHeaderBaseline(ctx *canvas.Context, fullW, fullH, headerH float64) {
	p := &canvas.Path{}
	p.MoveTo(0, 0)
	p.LineTo(fullW*2/3, 0)

	x := fullW / 3
	y := fullH - headerH/3
	ctx.SetStrokeColor(canvas.Black)
	ctx.DrawPath(x, y, p)
}

func (vr *VisualRoadmap) writeHeaderDates(ctx *canvas.Context, fullW, fullH, lineH float64) {
	x := fullW / 3
	y := fullH
	face := fontFamily.Face(lineH*1.5, canvas.Black, canvas.FontBold, canvas.FontNormal)
	date := vr.Dates.StartAt.Format(vr.DateFormat)
	ctx.DrawText(x, y, canvas.NewTextBox(face, date, 0.0, lineH, canvas.Left, canvas.Center, 0.0, 0.0))

	x = fullW
	date = vr.Dates.EndAt.Format(vr.DateFormat)
	ctx.DrawText(x, y, canvas.NewTextBox(face, date, 0.0, lineH, canvas.Right, canvas.Center, 0.0, 0.0))
}

func (vr *VisualRoadmap) markHeaderDates(ctx *canvas.Context, fullW, fullH, headerH, strokeW float64) {
	markH := headerH / 10.0

	p1 := &canvas.Path{}
	p1.MoveTo(strokeW, markH/-2)
	p1.LineTo(strokeW, markH/2)

	p2 := &canvas.Path{}
	p2.MoveTo(fullW*2/3-strokeW, markH/-2)
	p2.LineTo(fullW*2/3-strokeW, markH/2)

	x := fullW / 3
	y := fullH - headerH/3
	ctx.SetStrokeColor(canvas.Black)
	ctx.DrawPath(x, y, p1, p2)
}

func (vr *VisualRoadmap) writeProjects(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	textW := fullW / 3
	ctx.SetFillColor(canvas.Black)

	for i, p := range vr.Projects {
		y := fullH - float64(i)*lineH - headerH
		indentW := float64(p.Indentation)*textW/20 + 2
		font := vr.createFont(p.Indentation, lineH)
		ctx.DrawText(0, y, canvas.NewTextBox(font, p.Title, textW, lineH, canvas.Left, canvas.Center, indentW, 0.0))
	}
}

func (vr *VisualRoadmap) createFont(indentation uint8, lineH float64) canvas.FontFace {
	switch indentation {
	case 0:
		fontSize := lineH * 1.5
		return fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal)
	case 1:
		fontSize := lineH * 1.5
		return fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	case 2:
		fontSize := lineH * 1.35
		return fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	case 3:
		fontSize := lineH * 1.2
		return fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	case 4:
		fontSize := lineH * 1.05
		return fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	default:
		fontSize := lineH * 0.9
		return fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	}
}

func (vr *VisualRoadmap) drawProjects(ctx *canvas.Context, fullW, fullH, headerH, lineH, strokeW float64) {
	if vr.Dates == nil {
		return
	}

	h := lineH / 2
	maxW := fullW * 2 / 3
	roadmapInterval := vr.Dates.EndAt.Sub(vr.Dates.StartAt).Hours()
	r := lineH / 5

	ctx.SetStrokeWidth(1.0)
	ctx.SetStrokeColor(canvas.Darkgray)

	for i, p := range vr.Projects {
		if p.Dates == nil {
			continue
		}

		startAtLeft := p.Dates.StartAt.Sub(vr.Dates.StartAt).Hours()
		endAtLeft := p.Dates.EndAt.Sub(vr.Dates.StartAt).Hours()
		x0 := startAtLeft / roadmapInterval * maxW
		x1 := endAtLeft / roadmapInterval * maxW
		w := x1 - x0
		y := fullH - float64(i)*lineH - headerH - lineH/4*3

		ctx.SetFillColor(myLightGrey)
		ctx.DrawPath(x0+fullW/3, y, canvas.RoundedRectangle(w, h, r))

		if p.Percentage > 0 {
			w *= float64(p.Percentage) / 100
			ctx.SetFillColor(p.Color)
			ctx.DrawPath(x0+fullW/3, y, canvas.RoundedRectangle(w, h, r))
		}
	}

	ctx.SetStrokeWidth(strokeW)
}

func (vr *VisualRoadmap) drawMilestones(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	if vr.Dates == nil {
		return
	}

	maxW := fullW * 2 / 3
	y := fullH - headerH/3
	roadmapInterval := vr.Dates.EndAt.Sub(vr.Dates.StartAt).Hours()

	ctx.SetDashes(0.0, 3.0, 3.0)
	for _, m := range vr.Milestones {
		if m.DeadlineAt == nil {
			continue
		}

		c := canvas.Darkgray
		if m.Color != nil {
			c = *m.Color
		}

		deadlineFromStart := m.DeadlineAt.Sub(vr.Dates.StartAt).Hours()
		w := deadlineFromStart / roadmapInterval * maxW

		p := &canvas.Path{}
		p.MoveTo(w, 0)
		p.LineTo(w, fullH)

		ctx.SetStrokeColor(c)
		ctx.DrawPath(fullW/3, 0, p)

		x := w + fullW/3
		face := fontFamily.Face(lineH*1.5, c, canvas.FontRegular, canvas.FontNormal)
		date := fmt.Sprintf("%s\n%s", m.DeadlineAt.Format(vr.DateFormat), m.Title)
		ctx.DrawText(x, y, canvas.NewTextBox(face, date, 0.0, lineH*2, canvas.Center, canvas.Center, 0.0, 0.0))
	}
	ctx.SetDashes(0.0)
}

func (vr *VisualRoadmap) drawToday(ctx *canvas.Context, fullW, fullH, lineH float64) {
	if vr.Dates == nil {
		return
	}

	now := time.Now()

	if now.Before(vr.Dates.StartAt) || now.After(vr.Dates.EndAt) {
		return
	}

	if vr.Dates == nil {
		return
	}

	maxW := fullW * 2 / 3
	roadmapInterval := vr.Dates.EndAt.Sub(vr.Dates.StartAt).Hours()

	todayFromStart := now.Sub(vr.Dates.StartAt).Hours()
	w := todayFromStart / roadmapInterval * maxW

	p := &canvas.Path{}
	p.MoveTo(w, 0)
	p.LineTo(w, fullH)

	ctx.SetStrokeColor(myDarkGray)
	ctx.SetDashes(0.0, 8.0, 12.0)
	ctx.DrawPath(fullW/3, 0, p)

	x := w + fullW/3
	y := fullH
	face := fontFamily.Face(lineH*1.5, myDarkGray, canvas.FontRegular, canvas.FontNormal)
	date := now.Format(vr.DateFormat)
	ctx.DrawText(x, y, canvas.NewTextBox(face, date, 0.0, lineH, canvas.Center, canvas.Center, 0.0, 0.0))
}

func (vr *VisualRoadmap) drawLines(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	vr.drawHeaderLine(ctx, fullW, fullH, headerH)
	vr.drawProjectLines(ctx, fullW, fullH, headerH, lineH)
}

func (vr *VisualRoadmap) drawHeaderLine(ctx *canvas.Context, fullW, fullH, headerH float64) {
	if vr.Dates == nil {
		return
	}

	p := &canvas.Path{}
	p.MoveTo(0, fullH-headerH)
	p.LineTo(fullW, fullH-headerH)

	ctx.SetStrokeColor(canvas.Lightgray)
	ctx.DrawPath(0, 0, p)
}

func (vr *VisualRoadmap) drawProjectLines(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
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

	ctx.SetStrokeColor(canvas.Darkgray)
	ctx.DrawPath(0, 0, paths...)
}

func (vr *VisualRoadmap) drawProjectBackgrounds(ctx *canvas.Context, fullW, fullH, headerH, lineH float64) {
	var c1 color.RGBA
	var epicCount = -1

	for i := range vr.Projects {
		h0 := fullH - headerH - (float64(i+1) * (lineH))
		h1 := h0 + lineH

		if vr.Projects[i].Indentation == 0 {
			epicCount++
			c1 = pickBgColor(epicCount)
		}

		p := &canvas.Path{}
		p.MoveTo(0, h0)
		p.LineTo(fullW, h0)
		p.LineTo(fullW, h1)
		p.LineTo(0, h1)
		p.Close()

		ctx.SetFillColor(c1)
		ctx.SetStrokeColor(canvas.Lightgray)
		ctx.DrawPath(0, 0, p)
	}
}

func (vr *VisualRoadmap) drawLogo(ctx *canvas.Context, fullW, fullH, lineH float64) {
	if vr.Dates == nil {
		return
	}

	xName := 2.0
	yName := fullH
	faceName := fontFamily.Face(lineH*4, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	ctx.DrawText(xName, yName, canvas.NewTextBox(faceName, "rdmp.app", fullW/3, lineH*2, canvas.Left, canvas.Top, 0.0, 0.0))

	yLink := fullH - lineH*1.8
	faceLink := fontFamily.Face(lineH*1.5, canvas.Blue, canvas.FontRegular, canvas.FontNormal)
	ctx.DrawText(xName, yLink, canvas.NewTextBox(faceLink, "https://rdmp.app/", fullW/3, lineH, canvas.Left, canvas.Top, 0.0, 0.0))
}
