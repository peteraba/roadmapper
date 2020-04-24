package roadmap

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/peteraba/roadmapper/pkg/colors"
)

// VisualRoadmap represent a roadmap in a way that is prepared for visualization
type VisualRoadmap struct {
	Title      string
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

	visual.Title = r.Title
	visual.Dates = r.ToDates()
	visual.Projects = r.Projects
	visual.Milestones = r.Milestones
	visual.DateFormat = r.DateFormat

	visual.calculateProjectDates().calculateProjectColors().calculatePercentages().applyBaseURL(r.BaseURL)

	projectMilestones := visual.collectProjectMilestones()
	visual.applyProjectMilestone(projectMilestones)

	return visual
}

// calculateProjectDates tries to find reasonable dates for all projects
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
			c = colors.PickFgColor(epicCount, taskCount, int(p.Indentation))
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
