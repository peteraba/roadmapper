package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"image/color/palette"
	"math"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var roadmapRegexp *regexp.Regexp

func init() {
	roadmapRegexp = regexp.MustCompile(`^([^\]]*)(\s*\[(.*)\]\s*)?$`)
}

type Dates struct {
	Start time.Time
	End   time.Time
}

type Project struct {
	Title      string
	Dates      *Dates
	Children   []Project
	Color      color.Color
	Percentage uint8
	URL        string
}

func (p Project) IsPlanned() bool {
	return p.Dates != nil
}

type internalProject struct {
	title         string
	start         *time.Time
	end           *time.Time
	parent        *internalProject
	color         color.Color
	percentage    uint8
	url           string
	children      []*internalProject
	childrenStart *time.Time
	childrenEnd   *time.Time
}

func (p internalProject) GetStart() *time.Time {
	if p.start != nil {
		return p.start
	}

	if p.childrenStart != nil {
		return p.childrenStart
	}

	return nil
}

func (p internalProject) GetEnd() *time.Time {
	if p.end != nil {
		return p.end
	}

	if p.childrenEnd != nil {
		return p.childrenEnd
	}

	return nil
}

func (p internalProject) GetChildren() []*internalProject {
	return p.children
}

func (p internalProject) GetColor() color.Color {
	return p.color
}

func (p internalProject) GetPercentage() uint8 {
	return p.percentage
}

func (p internalProject) GetURL() string {
	return p.url
}

func (p *internalProject) ToPublic(roadmapStart, roadmapEnd *time.Time) Project {
	project := Project{Title: p.title, Color: p.GetColor(), Percentage: p.GetPercentage(), URL: p.GetURL()}
	project.Dates = p.GetDates()

	for _, c := range p.GetChildren() {
		project.Children = append(project.Children, c.ToPublic(roadmapStart, roadmapEnd))
	}

	return project
}

func (p *internalProject) GetDates() *Dates {
	if p.start == nil && p.childrenStart == nil {
		return nil
	}
	if p.end == nil && p.childrenEnd == nil {
		return nil
	}

	s, e := p.start, p.end
	if s == nil {
		s = p.childrenStart
	}
	if e == nil {
		e = p.childrenEnd
	}

	return &Dates{Start: *s, End: *e}
}

func (p internalProject) String() string {
	b, err := json.Marshal(p.ToPublic(p.GetStart(), p.GetEnd()))
	if err != nil {
		return ""
	}

	return string(b)
}

func parseRoadmap(lines []string, dateFormat, baseUrl string) (*internalProject, error) {
	var (
		err                error
		roadmap                  = internalProject{}
		currentProject           = &roadmap
		currentIndentation       = -1
		colorNum           uint8 = 11
	)

	for _, line := range lines {
		if line == "" {
			continue
		}

		currentProject, currentIndentation, err = createProject(line, currentProject, currentIndentation, &colorNum, dateFormat, baseUrl)
		if err != nil {
			return nil, err
		}
	}

	if len(roadmap.children) == 0 {
		return &roadmap, nil
	}

	_, _, err = setChildrenDates(&roadmap)
	if err != nil {
		return nil, err
	}

	return &roadmap, nil
}

func createProject(line string, previousProject *internalProject, parentIndentation int, colorNum *uint8, dateFormat, baseUrl string) (*internalProject, int, error) {
	trimmed := strings.TrimLeft(line, "\t")
	lineIndentation := len(line) - len(trimmed)

	newProject, err := parseProject(trimmed, colorNum, dateFormat, baseUrl)
	if err != nil {
		return nil, 0, err
	}

	pp := previousProject

	for currentIndentation := parentIndentation + 1; currentIndentation >= 0; currentIndentation -= 1 {
		if currentIndentation == lineIndentation {
			newProject.parent = pp
			break
		}

		if pp.parent == nil {
			return nil, 0, errors.New("invalid indentation")
		}

		pp = pp.parent
	}

	newProject.parent.children = append(newProject.parent.children, newProject)

	return newProject, lineIndentation, nil
}

func parseProject(trimmed string, colorNum *uint8, dateFormat, baseUrl string) (*internalProject, error) {
	res := roadmapRegexp.FindAllSubmatch([]byte(trimmed), -1)

	if res == nil {
		return nil, fmt.Errorf("failed to parse line: %s", trimmed)
	}

	if res[0][3] == nil {
		return &internalProject{title: string(res[0][1]), color: getNextColor(colorNum), percentage: 100}, nil
	}

	parts := strings.Split(string(res[0][3]), ", ")

	var (
		title      = strings.Trim(string(res[0][1]), " ")
		start, end *time.Time
		u          string
		p          uint8       = 100
		c          color.Color = color.RGBA{}
	)

	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			break
		}

		start, end, u, p, c = parseProjectExtra(parts[i], start, end, u, p, c, dateFormat, baseUrl)
	}

	if (reflect.DeepEqual(c, color.RGBA{})) {
		c = getNextColor(colorNum)
	}

	return &internalProject{title: title, start: start, end: end, color: c, percentage: p, url: u}, nil
}

func parseProjectExtra(part string, f, t *time.Time, u string, p uint8, c color.Color, dateFormat, baseUrl string) (*time.Time, *time.Time, string, uint8, color.Color) {
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
		return f, t, part, p, c
	}

	if baseUrl != "" {
		prefixedUrl := fmt.Sprintf("%s/%s", strings.TrimRight(baseUrl, "/"), strings.TrimLeft(part, "/"))
		_, err = url.ParseRequestURI(prefixedUrl)
		if err == nil {
			return f, t, prefixedUrl, p, c
		}
	}

	return f, t, u, p, c
}

func parsePercentage(part string) (uint8, error) {
	if len(part) < 1 {
		return 0, errors.New("invalid uint8 string")
	}

	percentage := false
	if part[len(part)-1] == '%' {
		part = part[:len(part)-1]
		percentage = true
	}

	n, err := strconv.ParseUint(part, 10, 8)
	if err == nil {
		if n > 100 {
			return 0, errors.New("invalid uint8 string")
		}
		return uint8(n), nil
	}

	n2, err := strconv.ParseFloat(part, 64)
	if err != nil {
		return 0, err
	}
	if n2 < 0 {
		return 0, errors.New("invalid uint8 string")
	}
	if n2 < 1 && !percentage {
		n2 = n2 * 100
	}

	if n2 > 0 && n2 < 100 {
		return uint8(math.Round(n2)), nil
	}

	return 0, errors.New("invalid uint8 string")
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

func setChildrenDates(p *internalProject) (*time.Time, *time.Time, error) {
	if len(p.children) == 0 {
		if p.end != nil && p.start != nil && p.end.Sub(*p.start) < 0 {
			return nil, nil, fmt.Errorf("starting must not be before ending date: %s", p.title)
		}

		return p.start, p.end, nil
	}

	var minStart, maxEnd *time.Time

	for _, c := range p.children {
		start, end, err := setChildrenDates(c)
		if err != nil {
			return nil, nil, err
		}

		if start != nil && end != nil {
			if minStart == nil || minStart.Sub(*start) > 0 {
				minStart = start
			}
			if maxEnd == nil || maxEnd.Sub(*end) < 0 {
				maxEnd = end
			}
		}
	}

	p.childrenStart = minStart
	p.childrenEnd = maxEnd

	if p.start != nil && p.start != p.childrenStart {
		return nil, nil, fmt.Errorf("project from date does not match calculated value: %s", p.title)
	}

	if p.end != nil && p.end != p.childrenEnd {
		return nil, nil, fmt.Errorf("project from date does not match calculated value: %s", p.title)
	}

	return minStart, maxEnd, nil
}

func getNextColor(colorNum *uint8) color.Color {
	*colorNum = (*colorNum + 71) % uint8(len(palette.WebSafe))

	return palette.WebSafe[*colorNum]
}
