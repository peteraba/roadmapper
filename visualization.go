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

func (vr VisualRoadmap) Draw(fullWidth, headerHeight, lineHeight float64) *canvas.Canvas {
	var (
		fullHeight float64
		// elements   []interface{}
		margin float64
		// d          = roadmap.ToDates()
	)

	c := canvas.New(fullWidth, fullHeight)

	ctx := canvas.NewContext(c)
	ctx.Push()

	// if d == nil {
	// 	headerHeight = 0
	// }
	//
	// elements = append(elements, createStripesPattern(), createStyle())
	//
	// titles, fullHeight := createSvgTitles(roadmap, fullWidth/3, headerHeight, lineHeight)
	//
	// elements = append(elements, titles)
	//
	// if roadmap.IsPlanned() {
	// 	elements = append(elements, createSvgVisuals(roadmap, fullWidth, headerHeight, lineHeight, dateFormat))
	// 	elements = append(elements, createSvgHeader(roadmap.Dates.Start, roadmap.Dates.End, fullHeight, fullWidth/3*2, headerHeight, fullWidth/3, 0, dateFormat))
	// }
	//
	// elements = append(elements, createSvgTableLines(fullWidth, fullHeight, fullWidth/3, headerHeight, lineHeight))

	c.Fit(margin)

	return c
}

// func createStripesPattern() svg.Element {
// 	polygons := []interface{}{
// 		svg.E("polygon", "", "", map[string]string{"points": "0,4 0,8 8,0 4,0", "fill": "white"}),
// 		svg.E("polygon", "", "", map[string]string{"points": "4,8 8,8 8,4", "fill": "white"}),
// 	}
// 	pattern := svg.E("pattern", "", "", map[string]string{"id": "stripes", "viewBox": "0,0,8,8", "width": "16", "height": "16", "patternUnits": "userSpaceOnUse"}, polygons...)
//
// 	def := svg.E("defs", "", "", nil, pattern)
//
// 	return def
// }
//
// func createStyle() svg.Element {
// 	rules := []string{
// 		`svg tspan {font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,"Noto Sans",sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol","Noto Color Emoji";}`,
// 		`svg .strong {font-weight: bold;}`,
// 		`svg a {fill: #06D; text-decoration: underline; cursor: pointer;}`,
// 		`svg a:hover, svg a:active {outline: dotted 1px blue; color: #06D;}`,
// 	}
// 	el := svg.E("style", "", strings.Join(rules, "\n"), nil)
//
// 	return el
// }
//
// func createSvgHeader(start, end time.Time, fullHeight, headerWidth, headerHeight, dx, dy float64, dateFormat string) svg.Group {
// 	var elements []interface{}
//
// 	elements = append(elements, createSvgHeaderLines(headerWidth, headerHeight, dx, dy))
// 	elements = append(elements, createSvgHeaderDates(start, end, headerWidth, headerHeight, dx, dy, dateFormat))
// 	elements = append(elements, createSvgHeaderToday(start, end, fullHeight, headerWidth, headerHeight, dx, dy, dateFormat))
//
// 	return svg.NewGroup(elements)
// }
//
// func createSvgHeaderLines(width, height, dx, dy float64) svg.Group {
// 	var elements []interface{}
//
// 	strokeColor1, _ := svg.ColorFromHexaString("#212529")
// 	elements = append(elements, svg.L(width-10+dx, (height/2)-5+dy, width-10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
// 	elements = append(elements, svg.L(0+dx, height/2+dy, width+dx, height/2+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
// 	elements = append(elements, svg.L(10+dx, (height/2)-5+dy, 10+dx, (height/2)+5+dy).SetStrokeWidth(2).SetStroke(strokeColor1))
//
// 	return svg.NewGroup(elements)
// }
//
// func createSvgHeaderDates(start, end time.Time, width, height, dx, dy float64, dateFormat string) svg.Group {
// 	var elements []interface{}
//
// 	tspan1 := svg.TS(start.Format(dateFormat))
// 	elements = append(elements, svg.T(12+dx, (height/2)-10+dy, tspan1))
//
// 	tspan2 := svg.TS(end.Format(dateFormat))
// 	elements = append(elements, svg.T(width-12+dx, (height/2)-10+dy, tspan2).SetTextAnchor(svg.End))
//
// 	return svg.NewGroup(elements)
// }
//
// func createSvgHeaderToday(start, end time.Time, fullHeight, width, height, dx, dy float64, dateFormat string) svg.Group {
// 	var elements []interface{}
//
// 	if time.Until(end) < 0 {
// 		return svg.NewGroup()
// 	}
//
// 	now := time.Now()
// 	untilToday := now.Sub(start).Hours()
// 	startToEnd := end.Sub(start).Hours()
//
// 	pos := untilToday / startToEnd * width
//
// 	strokeColor, _ := svg.ColorFromHexaString("#ff3333")
// 	elements = append(elements, svg.L(pos+dx, 0, pos+dx, fullHeight).SetStrokeWidth(2).SetStroke(strokeColor).AddAttr("stroke-dasharray", "10,10"))
//
// 	fillColor, _ := svg.ColorFromHexaString("#ff3333")
// 	tspan3 := svg.TS(now.Format(dateFormat))
// 	elements = append(elements, svg.T(pos+dx, (height/2)+20+dy, tspan3).SetFill(fillColor).SetTextAnchor(svg.Middle))
//
// 	return svg.NewGroup(elements)
// }
//
// func createSvgTitles(roadmap Project, titlesWidth, dy, lineHeight float64) (svg.Element, float64) {
// 	var (
// 		project  svg.Group
// 		projects []svg.Group
// 	)
//
// 	for _, p := range roadmap.Children {
// 		project, dy = createProjectTitle(roadmap, p, dy, 10, lineHeight, 30)
// 		projects = append(projects, project)
// 	}
//
// 	return svg.E("svg", "", "", map[string]string{"width": fs(titlesWidth), "height": fs(dy), "x": "0", "y": "0"}, projects), dy
// }
//
// func createProjectTitle(roadmap, project Project, dy, dx, lineHeight, indentWidth float64) (svg.Group, float64) {
// 	var (
// 		subProject  svg.Group
// 		subProjects []interface{}
// 	)
//
// 	title := createProjectTitleText(project, dx, dy, lineHeight)
//
// 	subProjects = append(subProjects, title)
//
// 	dy += lineHeight
//
// 	for _, c := range project.Children {
// 		subProject, dy = createProjectTitle(roadmap, c, dy, dx+indentWidth, lineHeight, indentWidth)
// 		subProjects = append(subProjects, subProject)
// 	}
//
// 	return svg.NewGroup(subProjects), dy
// }
//
// func createProjectTitleText(project Project, dx, dy, lineHeight float64) svg.Text {
// 	title := svg.TS(project.Title).AddAttr("class", "strong")
// 	if project.URL != "" {
// 		a := svg.NewA(project.URL, svg.TS(" "), title).AddAttr("target", "_blank")
// 		return svg.T(dx, dy+lineHeight/2+5, a)
// 	}
//
// 	return svg.T(dx, dy+lineHeight/2+5, title)
// }
//
// func fs(num float64) string {
// 	return fmt.Sprintf("%v", num)
// }
//
// func createSvgVisuals(roadmap Project, fullWidth, dy, lineHeight float64, dateFormat string) svg.Group {
// 	var (
// 		project  svg.Group
// 		projects []svg.Group
// 	)
//
// 	for _, p := range roadmap.Children {
// 		project, dy = createSvgProjectVisuals(roadmap, p, fullWidth, dy, 10, lineHeight, 30, dateFormat)
// 		projects = append(projects, project)
// 	}
//
// 	return svg.NewGroup(projects)
// }
//
// func createSvgProjectVisuals(roadmap, project Project, fullWidth, dy, dx, lineHeight, indentWidth float64, dateFormat string) (svg.Group, float64) {
// 	var (
// 		subProject  svg.Group
// 		subProjects []interface{}
// 	)
//
// 	if project.IsPlanned() {
// 		subProjects = append(subProjects, createProjectVisual(*roadmap.Dates, project, fullWidth, dy, lineHeight, dateFormat)...)
// 	}
//
// 	dy += lineHeight
//
// 	for _, c := range project.Children {
// 		subProject, dy = createSvgProjectVisuals(roadmap, c, fullWidth, dy, dx+indentWidth, lineHeight, indentWidth, dateFormat)
// 		subProjects = append(subProjects, subProject)
// 	}
//
// 	return svg.NewGroup(subProjects), dy
// }
//
// func createProjectVisual(roadmapDates Dates, project Project, fullWidth, dy, lineHeight float64, dateFormat string) []interface{} {
// 	wl := lineHeight * 0.6
// 	rd, pd := roadmapDates, project.Dates
// 	rs, rw := fullWidth/3+12, fullWidth/3*2-24
// 	ps := rs + (rw * pd.Start.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())
// 	pe := rs + (rw * pd.End.Sub(rd.Start).Hours() / rd.End.Sub(rd.Start).Hours())
//
// 	r, g, b, _ := project.Color.RGBA()
//
// 	baseColor, _ := svg.ColorFromHexaString("#dedede")
// 	base := svg.R(ps, dy+lineHeight/2-wl/2, pe-ps, wl).
// 		SetFill(baseColor).
// 		AddAttr("rx", "5").
// 		AddAttr("ry", "5")
//
// 	stripesColor := svg.Color{RGBA: color.RGBA{uint8(r), uint8(g), uint8(b), 255}}
// 	stripesBase := svg.R(0, 0, pe-ps, wl).
// 		SetFill(stripesColor).
// 		AddAttr("rx", "5").
// 		AddAttr("ry", "5")
//
// 	start := project.Dates.Start
// 	end := project.Dates.End
// 	tooltip := fmt.Sprintf("%d%%, %s - %s, %d days", project.Percentage, start.Format(dateFormat), end.Format(dateFormat), int64(end.Sub(start).Hours()/24))
// 	title := svg.E("title", "", tooltip, nil)
// 	stripes := svg.R(0, 0, pe-ps, wl, title).
// 		AddAttr("rx", "5").
// 		AddAttr("ry", "5").
// 		AddAttr("fill", `url(#stripes)`).
// 		SetFillOpacity(svg.Opacity{Number: .2})
//
// 	stripesWidth := fs((pe - ps) * float64(project.Percentage) / 100)
// 	stripesY := fs(dy + lineHeight/2 - wl/2)
// 	stripesContainer := svg.E("svg", "", "", map[string]string{"width": stripesWidth, "height": fs(wl), "x": fs(ps), "y": stripesY}, stripesBase, stripes)
//
// 	if project.Percentage >= 100 {
// 		stripesBase = stripesBase.SetFillOpacity(svg.Opacity{Number: .4})
// 	}
//
// 	return []interface{}{base, stripesContainer}
// }
//
// func createSvgTableLines(fullWidth, fullHeight, headerX, headerHeight, lineHeight float64) []interface{} {
// 	var result []interface{}
//
// 	light, _ := svg.ColorFromHexaString("#eee")
// 	for y := headerHeight + lineHeight; y < fullHeight; y += lineHeight {
// 		result = append(result, svg.L(0, y, fullWidth, y).SetStrokeWidth(2).SetStroke(light))
// 	}
//
// 	dark, _ := svg.ColorFromHexaString("#999")
// 	result = append(result, svg.L(headerX, 0, headerX, fullHeight).SetStrokeWidth(1).SetStroke(dark))
// 	result = append(result, svg.L(0, headerHeight, fullWidth, headerHeight).SetStrokeWidth(1).SetStroke(dark))
//
// 	return result
// }
