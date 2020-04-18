package roadmap

import (
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/peteraba/roadmapper/pkg/colors"
)

// Dates represents a pair of start end end dates
type Dates struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

// Roadmap represents a roadmap, the main entity of Roadmapper
type Roadmap struct {
	ID         uint64
	PrevID     *uint64
	Title      string
	DateFormat string
	BaseURL    string
	Projects   []Project
	Milestones []Milestone
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AccessedAt time.Time
}

// Project represents a project that belongs to a Roadmap
type Project struct {
	Indentation uint8       `json:"indentation"`
	Title       string      `json:"title"`
	Dates       *Dates      `json:"dates,omitempty"`
	Color       *color.RGBA `json:"color,omitempty"`
	Percentage  uint8       `json:"percentage"`
	URLs        []string    `json:"urls,omitempty"` // nolint
	Milestone   uint8       `json:"milestone,omitempty"`
}

// Milestone represents a milestone set for the roadmap
type Milestone struct {
	Title      string      `json:"title"`
	DeadlineAt *time.Time  `json:"deadline_at,omitempty"`
	Color      *color.RGBA `json:"color,omitempty"`
	URLs       []string    `json:"urls,omitempty"` // nolint
}

// Content represents a raw string version of a roadmap
type Content string

// ToLines splits a Content into lines, a slice of strings
func (c Content) ToLines() []string {
	if len(c) == 0 {
		return nil
	}

	return strings.Split(string(c), "\n")
}

// ToRoadmap converts a Content to a Roadmap ready to be persisted or to be turned into a VisualRoadmap which can then be rendered
func (c Content) ToRoadmap(id uint64, prevID *uint64, title, dateFormat, baseUrl string, now time.Time) Roadmap {
	r := Roadmap{
		ID:         id,
		PrevID:     prevID,
		Title:      title,
		DateFormat: dateFormat,
		BaseURL:    baseUrl,
		CreatedAt:  now,
		UpdatedAt:  now,
		AccessedAt: now,
	}

	indentation := c.findIndentation()

	r.Projects = c.toProjects(indentation, r.DateFormat, r.BaseURL)

	r.Milestones = c.toMilestones(r.DateFormat, r.BaseURL)

	return r
}

// findIndentation looks at the beginning of the lines of Content and return the first string of spaces and tabs found
func (c Content) findIndentation() string {
	lines := c.ToLines()

	for _, l := range lines {
		if strings.Trim(l, "\t ") == "" {
			continue
		}
		if l[0] != ' ' && l[0] != '	' {
			continue
		}

		t := strings.TrimLeft(l, " \t")
		d := len(l) - len(t)

		return l[0:d]
	}

	return "\t"
}

// toProjects converts Content to a slice of projects
func (c Content) toProjects(indentation, dateFormat, baseUrl string) []Project {
	var projects []Project

	for _, line := range c.ToLines() {
		var dates *Dates

		if !isLineProject(line) {
			continue
		}

		indentation, title, extra := splitLine(line, indentation)

		startAt, endAt, c, urls, percent, milestone, err := parseExtra(extra, dateFormat, baseUrl)
		if err != nil {
			continue
		}

		if startAt != nil && endAt != nil && startAt.Before(*endAt) {
			dates = &Dates{StartAt: *startAt, EndAt: *endAt}
		}

		projects = append(
			projects,
			Project{
				Indentation: indentation,
				Title:       title,
				Dates:       dates,
				Color:       c,
				Percentage:  percent,
				URLs:        urls,
				Milestone:   milestone,
			},
		)
	}

	return projects
}

// ToContent converts a Roadmap into a Content
func (r Roadmap) ToContent() Content {
	return Content(r.String())
}

// String converts a Roadmap into a string
func (r Roadmap) String() string {
	var lines []string

	for _, p := range r.Projects {
		lines = append(lines, p.String(r.DateFormat))
	}

	if len(r.Projects) > 0 && len(r.Milestones) > 0 {
		lines = append(lines, "")
	}

	for _, m := range r.Milestones {
		lines = append(lines, m.String(r.DateFormat))
	}

	return strings.Join(lines, "\n")
}

// String converts a Project into a string
func (p Project) String(dateFormat string) string {
	indentation := strings.Repeat("\t", int(p.Indentation))

	var extra []string

	if p.Dates != nil {
		extra = append(extra, p.Dates.StartAt.Format(dateFormat))
		extra = append(extra, p.Dates.EndAt.Format(dateFormat))
	}

	if p.Percentage > 0 {
		extra = append(extra, fmt.Sprintf("%d%%", p.Percentage))
	}

	if p.Color != nil {
		extra = append(extra, colors.ToHexa(p.Color))
	}

	extra = append(extra, p.URLs...)

	if p.Milestone > 0 {
		extra = append(extra, fmt.Sprintf("|%d", p.Milestone))
	}

	if len(extra) == 0 {
		return fmt.Sprintf("%s%s", indentation, p.Title)
	}

	return fmt.Sprintf("%s%s [%s]", indentation, p.Title, strings.Join(extra, ", "))
}

// String converts a Milestone into a string
func (m Milestone) String(dateFormat string) string {
	var extra []string

	if m.DeadlineAt != nil {
		extra = append(extra, m.DeadlineAt.Format(dateFormat))
	}

	if m.Color != nil {
		extra = append(extra, colors.ToHexa(m.Color))
	}

	extra = append(extra, m.URLs...)

	if len(extra) == 0 {
		return fmt.Sprintf("|%s", m.Title)
	}

	return fmt.Sprintf("|%s [%s]", m.Title, strings.Join(extra, ", "))
}

// ToDates converts a Roadmap into a Dates pointer
func (r Roadmap) ToDates() *Dates {
	var d *Dates

	for _, p := range r.Projects {
		if p.Dates == nil {
			continue
		}

		if d == nil {
			d = &Dates{StartAt: p.Dates.StartAt, EndAt: p.Dates.EndAt}
		}

		if p.Dates.StartAt.Before(d.StartAt) {
			d.StartAt = p.Dates.StartAt
		}

		if p.Dates.EndAt.After(d.EndAt) {
			d.EndAt = p.Dates.EndAt
		}
	}

	for _, m := range r.Milestones {
		if m.DeadlineAt == nil {
			continue
		}

		if d == nil {
			break
		}

		if m.DeadlineAt.Before(d.StartAt) {
			d.StartAt = *m.DeadlineAt
		}

		if m.DeadlineAt.After(d.EndAt) {
			d.EndAt = *m.DeadlineAt
		}
	}

	return d
}

// toMilestones converts Content into a slice of Milestones
func (c Content) toMilestones(dateFormat, baseUrl string) []Milestone {
	var milestones []Milestone

	for _, line := range c.ToLines() {
		if !isLineMilestone(line) {
			continue
		}

		_, title, extra := splitLine(line, "")

		deadlineAt, endAt, c, urls, _, _, err := parseExtra(extra, dateFormat, baseUrl)
		if err != nil {
			continue
		}

		if endAt != nil {
			continue
		}

		milestones = append(
			milestones,
			Milestone{
				Title:      title,
				DeadlineAt: deadlineAt,
				Color:      c,
				URLs:       urls,
			},
		)
	}

	return milestones
}

// splitLine splits a Content line into a title and extra information, plus returns the indentation level found
func splitLine(line, indentation string) (uint8, string, string) {
	var (
		n uint8
	)

	for indentation != "" {
		if len(line) < len(indentation) {
			break
		}

		if line[:len(indentation)] != indentation {
			break
		}

		line = line[len(indentation):]
		n++
	}

	if len(line) > 0 && line[0] == '|' {
		line = line[1:]
	}

	lo := strings.LastIndex(line, "[")
	lc := strings.LastIndex(line, "]")

	if lc < 0 || lo < 0 || lc < lo {
		return n, line, ""
	}

	return n, strings.Trim(line[:lo], "\t\r "), strings.Trim(line[lo+1:lc], "\t\r ")
}

// isLineProject returns true if a given line appears to represent a project
func isLineProject(line string) bool {
	line = strings.Trim(line, "\t\r ")

	// empty line should be skipped
	if len(line) == 0 {
		return false
	}

	// milestones should be skipped
	if line[0] == '|' {
		return false
	}

	return true
}

// isLineMilestone returns true if a given line appears to represent a milestone
func isLineMilestone(line string) bool {
	line = strings.Trim(line, "\t\r ")

	// empty line should be skipped
	if len(line) == 0 {
		return false
	}

	// projects should be skipped
	if line[0] != '|' {
		return false
	}

	return true
}

// parseExtra returns data found in extra parts of lines representing projects and milestones
func parseExtra(extra, dateFormat, baseUrl string) (*time.Time, *time.Time, *color.RGBA, []string, uint8, uint8, error) {
	parts := strings.Split(extra, ", ")

	var (
		startAt, endAt *time.Time
		urls           []string
		c              *color.RGBA
		percent        uint8
		milestone      uint8
	)

	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			break
		}

		startAt, endAt, urls, c, percent, milestone = parseExtraPart(parts[i], startAt, endAt, urls, c, percent, milestone, dateFormat, baseUrl)
	}

	return startAt, endAt, c, urls, percent, milestone, nil
}

// parseExtraPart returns data found in one piece of extra information
func parseExtraPart(part string, f, t *time.Time, u []string, c *color.RGBA, p, m uint8, dateFormat, baseUrl string) (*time.Time, *time.Time, []string, *color.RGBA, uint8, uint8) {
	t2, err := time.Parse(dateFormat, part)
	if err == nil {
		if f == nil {
			return &t2, t, u, c, p, m
		}

		return f, &t2, u, c, p, m
	}

	p2, err := parsePercentage(part)
	if err == nil {
		return f, t, u, c, p2, m
	}

	m2, err := parseMilestone(part)
	if err == nil {
		return f, t, u, c, p, m2
	}

	c2, err := parseColor(part)
	if err == nil {
		return f, t, u, c2, p, m
	}

	parsedUrl, err := url.ParseRequestURI(part)
	if err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "" {
		return f, t, append(u, part), c, p, m
	}

	if baseUrl != "" {
		prefixedUrl := fmt.Sprintf("%s/%s", strings.TrimRight(baseUrl, "/"), strings.TrimLeft(part, "/"))
		_, err = url.ParseRequestURI(prefixedUrl)
		if err == nil {
			return f, t, append(u, part), c, p, m
		}
	}

	return f, t, u, c, p, m
}

var errCannotParsePercentage = errors.New("can not parse string as percentage")

// parsePercentage tries to parse a string as a percentage between 0 and 100
func parsePercentage(part string) (uint8, error) {
	if len(part) < 2 {
		return 0, errCannotParsePercentage
	}

	if part[len(part)-1] != '%' {
		return 0, errCannotParsePercentage
	}

	part = part[:len(part)-1]

	n, err := strconv.ParseUint(part, 10, 8)
	if err != nil {
		return 0, errCannotParsePercentage
	}
	if n > 100 {
		return 0, errCannotParsePercentage
	}

	return uint8(n), nil
}

var errCannotParseMilestone = errors.New("can not parse string as milestone")

// parseMilestone tries to parse a string as a milestone number
// project milestones start with a | character and continue with a positive integer
func parseMilestone(part string) (uint8, error) {
	if len(part) < 2 {
		return 0, errCannotParseMilestone
	}

	if part[0] != '|' {
		return 0, errCannotParseMilestone
	}

	part = part[1:]

	n, err := strconv.ParseUint(part, 10, 8)
	if err == nil {
		return uint8(n), nil
	}

	return 0, errCannotParseMilestone
}

// parseColor tries to parse a string as a color in a hexadecimal representation (e.g #fa3, #ffaa33)
func parseColor(part string) (*color.RGBA, error) {
	if len(part) != 4 && len(part) != 7 {
		return nil, errors.New("invalid hexa color length")
	}

	if part[0] != '#' {
		return nil, errors.New("invalid first character for hexa color")
	}

	s, err := colors.CharsToUint8(part[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse to uint8s: %w", err)
	}

	return &color.RGBA{R: s[0], G: s[1], B: s[2], A: 255}, nil
}
