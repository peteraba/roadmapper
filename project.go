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

type Project struct {
	Title      string
	From       time.Time
	To         time.Time
	Children   []Project
	Color      color.Color
	Percentage uint8
	URL        string
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

	if p.childrenFrom != nil {
		return *p.childrenFrom
	}

	return time.Time{}
}

func (p internalProject) GetTo() time.Time {
	if p.to != nil {
		return *p.to
	}

	if p.childrenTo != nil {
		return *p.childrenTo
	}

	return time.Time{}
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

func (p *internalProject) ToPublic(roadmapFrom, roadmapTo time.Time) Project {
	res := Project{Title: p.Title, From: p.GetFrom(), To: p.GetTo(), Color: p.GetColor(), Percentage: p.GetPercentage(), URL: p.GetURL()}

	for _, c := range p.GetChildren() {
		res.Children = append(res.Children, c.ToPublic(roadmapFrom, roadmapTo))
	}

	return res
}

func (p internalProject) String() string {
	b, err := json.Marshal(p.ToPublic(p.GetFrom(), p.GetTo()))
	if err != nil {
		return ""
	}

	return string(b)
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

	if len(roadmap.children) == 0 {
		return &roadmap, nil
	}

	_, _, err = setChildrenDates(&roadmap)
	if err != nil {
		return nil, err
	}

	return &roadmap, nil
}

func createProject(line string, previousProject *internalProject, pi int, colorNum *uint8) (*internalProject, int, error) {
	trimmed := strings.TrimLeft(line, " ")
	ni := len(line) - len(trimmed)

	newProject, err := parseProject(trimmed, colorNum)
	if err != nil {
		return nil, 0, err
	}

	pp := previousProject

	for ci := pi + 2; ci >= 0; ci -= 2 {
		if ci == ni {
			newProject.parent = pp
			break
		}

		if pp.parent == nil {
			return nil, 0, errors.New("invalid indentation")
		}

		pp = pp.parent
	}

	newProject.parent.children = append(newProject.parent.children, newProject)

	return newProject, ni, nil
}

func parseProject(trimmed string, colorNum *uint8) (*internalProject, error) {
	res := roadmapRegexp.FindAllSubmatch([]byte(trimmed), -1)

	if res == nil {
		return nil, fmt.Errorf("failed to parse line: %s", trimmed)
	}

	if res[0][3] == nil {
		return &internalProject{Title: string(res[0][1]), color: getNextColor(colorNum), percentage: 100}, nil
	}

	parts := strings.Split(string(res[0][3]), ", ")

	var (
		title = strings.Trim(string(res[0][1]), " ")
		f, t  *time.Time
		u     string
		p     uint8       = 100
		c     color.Color = color.RGBA{}
	)

	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			break
		}

		f, t, u, p, c = parseProjectExtra(parts[i], f, t, u, p, c)
	}

	if (reflect.DeepEqual(c, color.RGBA{})) {
		c = getNextColor(colorNum)
	}

	return &internalProject{Title: title, from: f, to: t, color: c, percentage: p, url: u}, nil
}

func parseProjectExtra(part string, f, t *time.Time, u string, p uint8, c color.Color) (*time.Time, *time.Time, string, uint8, color.Color) {
	t2, err := time.Parse(dateFormat, part)
	if err == nil {
		if f == nil {
			return &t2, t, u, p, c
		}

		return f, &t2, u, p, c
	}

	_, err = url.ParseRequestURI(part)
	if err == nil {
		return f, t, part, p, c
	}

	n, err := parsePercentage(part)
	if err == nil {
		return f, t, u, n, c
	}

	c2, err := parseColor(part)
	if err == nil {
		return f, t, u, p, c2
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
		if n < 0 || n > 100 {
			return 0, errors.New("invalid uint8 string")
		}
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
		if p.from == nil || p.to == nil {
			return nil, nil, fmt.Errorf("project needs starting and ending dates: %s", p.Title)
		}
		if p.to.Sub(*p.from) < 0 {
			return nil, nil, fmt.Errorf("starting must not be before ending date: %s", p.Title)
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
