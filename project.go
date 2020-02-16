package main

import (
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

type Project struct {
	Title      string
	From       time.Time
	To         time.Time
	Children   []Project
	Color      color.Color
	Percentage uint8
	URL        string
}

func ProjectToPublic(project *internalProject, roadmapFrom, roadmapTo time.Time) Project {
	p := Project{Title: project.Title, From: project.GetFrom(), To: project.GetTo(), Color: project.GetColor(), Percentage: project.GetPercentage(), URL: project.GetURL()}

	for _, c := range project.GetChildren() {
		p.Children = append(p.Children, ProjectToPublic(c, roadmapFrom, roadmapTo))
	}

	return p
}

const dateFormat = "2006-01-02"

type internalProject struct {
	Title        string
	from         *time.Time
	to           *time.Time
	parent       *internalProject
	color        color.Color
	percentage   uint8
	url          string
	children     []*internalProject
	childrenFrom *time.Time
	childrenTo   *time.Time
}

func (p internalProject) GetFrom() time.Time {
	if p.from != nil {
		return *p.from
	}

	return *p.childrenFrom
}

func (p internalProject) GetTo() time.Time {
	if p.to != nil {
		return *p.to
	}

	return *p.childrenTo
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

func parseRoadmap(lines []string) (*internalProject, error) {
	var (
		err                error
		roadmap                  = internalProject{}
		currentProject           = &roadmap
		currentIndentation       = -2
		colorNum           uint8 = 11
	)

	for _, line := range lines {
		if line == "" {
			continue
		}

		currentProject, currentIndentation, err = createProject(line, currentProject, currentIndentation, &colorNum)
		if err != nil {
			return nil, err
		}
	}

	_, _, err = setChildrenDates(&roadmap)
	if err != nil {
		return nil, err
	}

	return &roadmap, nil
}

func createProject(line string, previousProject *internalProject, pi int, colorNum *uint8) (*internalProject, int, error) {
	trimmed := strings.TrimLeft(line, " ")
	ci := len(line) - len(trimmed)

	newProject, err := parseProject(trimmed, colorNum)
	if err != nil {
		return nil, 0, err
	}

	switch ci {
	case pi + 2:
		newProject.parent = previousProject
		break
	case pi:
		newProject.parent = previousProject.parent
		break
	case pi - 2:
		newProject.parent = previousProject.parent.parent
		break
	default:
		return nil, 0, errors.New("invalid indentation")
	}

	newProject.parent.children = append(newProject.parent.children, newProject)

	return newProject, ci, nil
}

func parseProject(trimmed string, colorNum *uint8) (*internalProject, error) {
	res := roadmapRegexp.FindAllSubmatch([]byte(trimmed), -1)

	if res == nil {
		return nil, fmt.Errorf("failed to parse line: %s", trimmed)
	}

	if res[0][3] == nil {
		return &internalProject{Title: string(res[0][1]), color: getNextColor(colorNum), percentage: 100}, nil
	}

	var f, t *time.Time
	parts := strings.Split(string(res[0][3]), ", ")

	if len(parts) > 0 && parts[0] != "" {
		fv, err := time.Parse(dateFormat, parts[0])
		if err != nil {
			return nil, err
		}
		f = &fv
	}

	if len(parts) > 1 && parts[1] != "" {
		tv, err := time.Parse(dateFormat, parts[1])
		if err != nil {
			return nil, err
		}
		t = &tv
	}

	var (
		u string
		p uint8       = 100
		c color.Color = color.RGBA{}
	)

	for i := 2; i < len(parts); i++ {
		if parts[i] == "" {
			break
		}

		u, p, c = parseProjectExtra(parts[i], u, p, c)
	}

	if (reflect.DeepEqual(c, color.RGBA{})) {
		c = getNextColor(colorNum)
	}

	return &internalProject{Title: string(res[0][1]), from: f, to: t, color: c, percentage: p, url: u}, nil
}

func parseProjectExtra(part, u string, p uint8, c color.Color) (string, uint8, color.Color) {
	_, err := url.ParseRequestURI(part)
	if err == nil {
		return part, p, c
	}

	n, err := parsePercentage(part)
	if err == nil {
		return u, n, c
	}

	c2, err := parseColor(part)
	if err == nil {
		return u, p, c2
	}

	return u, p, c
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
		return uint8(n), nil
	}

	n2, err := strconv.ParseFloat(part, 64)
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

func parseColor(part string) (color.Color, error) {
	if len(part) != 4 && len(part) != 7 {
		return nil, errors.New("invalid hexa color length")
	}

	if part[0] != '#' {
		return nil, errors.New("invalid first character for hexa color")
	}

	s, err := charsToUint8(part[1:])
	if err != nil {
		return nil, err
	}

	return color.RGBA{R: s[0], G: s[1], B: s[2], A: 100}, nil
}

func charsToUint8(part string) ([3]uint8, error) {
	tmp := []int{}
	for i := 0; i < len(part); i++ {
		if idx := strings.Index("0123456789abcdef", string(part[i])); idx > -1 {
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
		if p.from == nil || p.to == nil {
			return nil, nil, fmt.Errorf("project needs starting and ending dates: %s", p.Title)
		}

		return p.from, p.to, nil
	}

	var minF, maxT *time.Time

	tmpF := time.Date(3000, 0, 0, 0, 0, 0, 0, time.UTC)
	minF = &tmpF

	tmpT := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	maxT = &tmpT

	for _, c := range p.children {
		f, t, err := setChildrenDates(c)
		if err != nil {
			return nil, nil, err
		}

		if minF.Sub(*f) > 0 {
			minF = f
		}
		if maxT.Sub(*t) < 0 {
			maxT = t
		}
	}

	p.childrenFrom = minF
	p.childrenTo = maxT

	if p.from != nil && p.from != p.childrenFrom {
		return nil, nil, fmt.Errorf("project from date does not match calculated value: %s", p.Title)
	}

	if p.to != nil && p.to != p.childrenTo {
		return nil, nil, fmt.Errorf("project from date does not match calculated value: %s", p.Title)
	}

	return minF, maxT, nil
}

func getNextColor(colorNum *uint8) color.Color {
	*colorNum = (*colorNum + 71) % uint8(len(palette.WebSafe))

	return palette.WebSafe[*colorNum]
}
