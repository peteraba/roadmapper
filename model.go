package main

import (
	"errors"
	"fmt"
	"image/color"
	"image/color/palette"
	"math"
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
	Dates       *Dates      `json:"dates"`
	Color       color.Color `json:"color"`
	Percentage  uint8       `json:"percentage"`
	URLs        []string    `json:"urls"` // nolint
}

type Milestone struct {
	Title      string      `json:"title"`
	DeadlineAt *time.Time  `json:"deadline_at"`
	Color      color.Color `json:"color"`
	URLs       []string    `json:"urls"` // nolint
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

		startAt, endAt, c, urls, percent, err := parseExtra(extra, dateFormat, baseUrl)
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

		deadlineAt, endAt, c, urls, _, err := parseExtra(extra, dateFormat, baseUrl)
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

func parseExtra(extra, dateFormat, baseUrl string) (*time.Time, *time.Time, color.Color, []string, uint8, error) {
	parts := strings.Split(extra, ", ")

	var (
		startAt, endAt *time.Time
		urls           []string
		percent        uint8 = 100
		c              color.Color
	)

	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			break
		}

		startAt, endAt, urls, percent, c = parseExtraPart(parts[i], startAt, endAt, urls, percent, c, dateFormat, baseUrl)
	}

	return startAt, endAt, c, urls, percent, nil
}

func parseExtraPart(part string, f, t *time.Time, u []string, p uint8, c color.Color, dateFormat, baseUrl string) (*time.Time, *time.Time, []string, uint8, color.Color) {
	t2, err := time.Parse(dateFormat, part)
	if err == nil {
		if f == nil {
			return &t2, t, u, p, c
		}

		return f, &t2, u, p, c
	}

	n, err := parsePercentage(part)
	if err == nil {
		return f, t, u, n, c
	}

	c2, err := parseColor(part)
	if err == nil {
		return f, t, u, p, c2
	}

	parsedUrl, err := url.ParseRequestURI(part)
	if err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "" {
		return f, t, append(u, part), p, c
	}

	if baseUrl != "" {
		prefixedUrl := fmt.Sprintf("%s/%s", strings.TrimRight(baseUrl, "/"), strings.TrimLeft(part, "/"))
		_, err = url.ParseRequestURI(prefixedUrl)
		if err == nil {
			return f, t, append(u, prefixedUrl), p, c
		}
	}

	return f, t, u, p, c
}

var errCannotParsePercentage = errors.New("can not parse string as percentage")

func parsePercentage(part string) (uint8, error) {
	if len(part) < 1 {
		return 0, errCannotParsePercentage
	}

	percentage := false
	if part[len(part)-1] == '%' {
		part = part[:len(part)-1]
		percentage = true
	}

	n, err := strconv.ParseUint(part, 10, 8)
	if err == nil {
		if n > 100 {
			return 0, errCannotParsePercentage
		}
		return uint8(n), nil
	}

	n2, err := strconv.ParseFloat(part, 64)
	if err != nil {
		return 0, errCannotParsePercentage
	}
	if n2 < 0 {
		return 0, errCannotParsePercentage
	}
	if n2 < 1 && !percentage {
		n2 = n2 * 100
	}

	if n2 > 0 && n2 < 100 {
		return uint8(math.Round(n2)), nil
	}

	return 0, errCannotParsePercentage
}

func parseColor(part string) (color.RGBA, error) {
	if len(part) != 4 && len(part) != 7 {
		return color.RGBA{}, errors.New("invalid hexa color length")
	}

	if part[0] != '#' {
		return color.RGBA{}, errors.New("invalid first character for hexa color")
	}

	s, err := charsToUint8(part[1:])
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{R: s[0], G: s[1], B: s[2], A: 255}, nil
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
