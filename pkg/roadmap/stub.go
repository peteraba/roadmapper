package roadmap

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
)

func NewRoadmapExchangeStub(minProjects, minMilestones int, minDate, maxDate time.Time) RoadmapExchange {
	var (
		bu = newBaseURL()
		p  = gofakeit.Number(minProjects, 20)
		m  = gofakeit.Number(minMilestones, max(minMilestones, p))

		milestones []Milestone
		projects   []Project
		project    Project
		ind        = 0
		hasBU      = bu != ""
	)

	for i := 0; i < m; i++ {
		milestones = append(milestones, NewMilestoneStub(minDate, maxDate, hasBU))
	}

	for i := 0; i < p; i++ {
		project = NewProjectStub(m, ind, minDate, maxDate, hasBU)
		projects = append(projects, project)
		ind = nextIndentation(ind)
	}

	return RoadmapExchange{
		Title:      newWords(),
		DateFormat: "2006-01-02",
		BaseURL:    bu,
		Projects:   projects,
		Milestones: milestones,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func newBaseURL() string {
	if gofakeit.Bool() {
		return ""
	}

	return gofakeit.URL()
}

func NewProjectStub(milestoneCount, ind int, minDate, maxDate time.Time, hasBU bool) Project {
	m := gofakeit.Number(0, milestoneCount)
	d := newDates(minDate, maxDate)
	p := gofakeit.Number(0, 100)

	project := Project{
		Indentation: uint8(ind),
		Title:       newWords(),
		Milestone:   uint8(m),
		Dates:       d,
		Percentage:  uint8(p),
		URLs:        getURLs(hasBU),
	}

	return project
}

func NewMilestoneStub(minDate, maxDate time.Time, hasBU bool) Milestone {
	return Milestone{
		Title:      newWords(),
		DeadlineAt: newDateOptional(minDate, maxDate),
		URLs:       getURLs(hasBU),
	}
}

func nextIndentation(indentation int) int {
	return indentation - gofakeit.Number(-1, indentation)
}

func newWords() string {
	var w []string

	for i := 0; i < gofakeit.Number(1, 5); i++ {
		w = append(w, gofakeit.HipsterWord())
	}

	return strings.Join(w, " ")
}

func newDates(minDate, maxDate time.Time) *Dates {
	if gofakeit.Bool() {
		return nil
	}

	var (
		d0 = gofakeit.DateRange(minDate, maxDate)
		d1 = gofakeit.DateRange(minDate, maxDate)
	)

	if d0.Before(d1) {
		return &Dates{
			StartAt: d0,
			EndAt:   d1,
		}
	}

	return &Dates{
		StartAt: d1,
		EndAt:   d0,
	}
}

func getURLs(hasBU bool) []string {
	var (
		urls []string
	)

	for i := 0; i < gofakeit.Number(0, 2); i++ {
		urls = append(urls, gofakeit.URL())
	}

	if hasBU {
		for i := 0; i < gofakeit.Number(0, 2); i++ {
			urls = append(urls, gofakeit.Word())
		}
	}

	return urls
}

func newDateOptional(minDate, maxDate time.Time) *time.Time {
	var (
		optionalDate *time.Time
	)

	if gofakeit.Bool() {
		date := gofakeit.DateRange(minDate, maxDate)
		optionalDate = &date
	}

	return optionalDate
}
