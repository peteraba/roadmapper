package main

import (
	"errors"
	"fmt"
	"image/color"
	"image/color/palette"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Dates struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

type Roadmap struct {
	ID         uint64
	PrevID     *uint64
	DateFormat string
	BaseURL    string
	Projects   []Project
	Milestones []Milestone
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AccessedAt time.Time
}

type Project struct {
	Indentation uint8       `json:"indentation"`
	Title       string      `json:"title"`
	Dates       *Dates      `json:"dates,omitempty"`
	Color       *color.RGBA `json:"color,omitempty"`
	Percentage  uint8       `json:"percentage"`
	URLs        []string    `json:"urls,omitempty"` // nolint
	Milestone   uint8       `json:"milestone,omitempty"`
}

type Milestone struct {
	Title      string      `json:"title"`
	DeadlineAt *time.Time  `json:"deadline_at,omitempty"`
	Color      *color.RGBA `json:"color,omitempty"`
	URLs       []string    `json:"urls,omitempty"` // nolint
}

type Content string

func (c Content) ToLines() []string {
	if len(c) == 0 {
		return nil
	}

	return strings.Split(string(c), "\n")
}

func (c Content) ToRoadmap(id uint64, prevID *uint64, dateFormat, baseUrl string, now time.Time) Roadmap {
	r := Roadmap{
		ID:         id,
		PrevID:     prevID,
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

func (r Roadmap) ToContent() Content {
	return Content(r.ToString())
}

func (r Roadmap) ToString() string {
	var lines []string

	for _, p := range r.Projects {
		lines = append(lines, p.ToString(r.DateFormat))
	}

	if len(r.Projects) > 0 && len(r.Milestones) > 0 {
		lines = append(lines, "")
	}

	for _, m := range r.Milestones {
		lines = append(lines, m.ToString(r.DateFormat))
	}

	return strings.Join(lines, "\n")
}

func (p Project) ToString(dateFormat string) string {
	indentation := strings.Repeat("\t", int(p.Indentation))

	var extra []string

	if p.Dates != nil {
		extra = append(extra, p.Dates.StartAt.Format(dateFormat))
		extra = append(extra, p.Dates.EndAt.Format(dateFormat))
	}

	if p.Percentage < 100 {
		extra = append(extra, fmt.Sprintf("%d%%", p.Percentage))
	}

	if p.Color != nil {
		extra = append(extra, colorToHexa(p.Color))
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

func (m Milestone) ToString(dateFormat string) string {
	var extra []string

	if m.DeadlineAt != nil {
		extra = append(extra, m.DeadlineAt.Format(dateFormat))
	}

	if m.Color != nil {
		extra = append(extra, colorToHexa(m.Color))
	}

	extra = append(extra, m.URLs...)

	if len(extra) == 0 {
		return fmt.Sprintf("|%s", m.Title)
	}

	return fmt.Sprintf("|%s [%s]", m.Title, strings.Join(extra, ", "))
}

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

func isLineProject(line string) bool {
	line = strings.TrimLeft(line, "\t ")

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

func isLineMilestone(line string) bool {
	line = strings.TrimLeft(line, "\t ")

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

func parseExtra(extra, dateFormat, baseUrl string) (*time.Time, *time.Time, *color.RGBA, []string, uint8, uint8, error) {
	parts := strings.Split(extra, ", ")

	var (
		startAt, endAt *time.Time
		urls           []string
		c              *color.RGBA
		percent        uint8 = 100
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
			return f, t, append(u, prefixedUrl), c, p, m
		}
	}

	return f, t, u, c, p, m
}

var errCannotParsePercentage = errors.New("can not parse string as percentage")

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

func parseColor(part string) (*color.RGBA, error) {
	if len(part) != 4 && len(part) != 7 {
		return nil, errors.New("invalid hexa color length")
	}

	if part[0] != '#' {
		return nil, errors.New("invalid first character for hexa color")
	}

	s, err := charsToUint8(part[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse to uint8s: %w", err)
	}

	return &color.RGBA{R: s[0], G: s[1], B: s[2], A: 255}, nil
}

func charsToUint8(part string) ([3]uint8, error) {
	if len(part) != 3 && len(part) != 6 {
		return [3]uint8{}, errors.New("invalid hexadecimal color string")
	}

	part = strings.ToLower(part)

	tmp := []int{}
	for _, runeValue := range part {
		if idx := strings.IndexRune("0123456789abcdef", runeValue); idx > -1 {
			tmp = append(tmp, idx)
			if len(part) == 3 {
				tmp = append(tmp, idx)
			}
		}
	}

	res := [3]uint8{}
	res[0] = uint8(tmp[0]*16 + tmp[1])
	res[1] = uint8(tmp[2]*16 + tmp[3])
	res[2] = uint8(tmp[4]*16 + tmp[5])

	return res, nil
}

func colorToHexa(c color.Color) string {
	if c == nil {
		return ""
	}

	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%s%s%s", twoDigitHexa(r), twoDigitHexa(g), twoDigitHexa(b))
}

func twoDigitHexa(i uint32) string {
	if i > 0xf {
		return fmt.Sprintf("%x", uint8(i))
	}

	return fmt.Sprintf("0%x", uint8(i))
}

func getNextColor(colorNum *uint8) color.Color {
	*colorNum = (*colorNum + 71) % uint8(len(palette.WebSafe))

	return palette.WebSafe[*colorNum]
}
