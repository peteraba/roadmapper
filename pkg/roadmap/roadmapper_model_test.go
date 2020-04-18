package roadmap

import (
	"image/color"
	"reflect"
	"testing"
	"time"
)

func TestContent_ToLines(t *testing.T) {
	tests := []struct {
		name string
		c    Content
		want []string
	}{
		{
			"empty",
			"",
			nil,
		},
		{
			"2 lines",
			"one\ntwo",
			[]string{"one", "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ToLines(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContent_ToRoadmap(t *testing.T) {
	prevID := uint64(43)
	now := time.Date(2020, 3, 4, 0, 0, 0, 0, time.UTC)

	type args struct {
		id         uint64
		prevID     *uint64
		title      string
		dateFormat string
		baseUrl    string
		now        time.Time
	}
	tests := []struct {
		name string
		c    Content
		args args
		want Roadmap
	}{
		{
			"empty roadmap",
			"",
			args{},
			Roadmap{},
		},
		{
			"example",
			`Initial development
Bring website online
	Select and purchase domain
	Create server infrastructure
Command line tool
	Create backend SVG generation
	Replace frontend SVG generation with backend
	Create documentation page
Marketing
	Create Facebook page
	Write blog posts
	Share blog post on social media
	Talk about the tool in relevant meetups

|Milestone 0.2`,
			args{
				id:         123,
				prevID:     &prevID,
				title:      "example",
				baseUrl:    "https://example.com",
				dateFormat: "2006-01-02",
				now:        now,
			},
			Roadmap{
				ID:         123,
				PrevID:     &prevID,
				Title:      "example",
				BaseURL:    "https://example.com",
				DateFormat: "2006-01-02",
				Projects: []Project{
					{Title: "Initial development", Percentage: 0},
					{Title: "Bring website online", Percentage: 0},
					{Title: "Select and purchase domain", Percentage: 0, Indentation: 1},
					{Title: "Create server infrastructure", Percentage: 0, Indentation: 1},
					{Title: "Command line tool", Percentage: 0},
					{Title: "Create backend SVG generation", Percentage: 0, Indentation: 1},
					{Title: "Replace frontend SVG generation with backend", Percentage: 0, Indentation: 1},
					{Title: "Create documentation page", Percentage: 0, Indentation: 1},
					{Title: "Marketing", Percentage: 0},
					{Title: "Create Facebook page", Percentage: 0, Indentation: 1},
					{Title: "Write blog posts", Percentage: 0, Indentation: 1},
					{Title: "Share blog post on social media", Percentage: 0, Indentation: 1},
					{Title: "Talk about the tool in relevant meetups", Percentage: 0, Indentation: 1},
				},
				Milestones: []Milestone{
					{Title: "Milestone 0.2"},
				},
				CreatedAt:  now,
				UpdatedAt:  now,
				AccessedAt: now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ToRoadmap(tt.args.id, tt.args.prevID, tt.args.title, tt.args.dateFormat, tt.args.baseUrl, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRoadmap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContent_findIndentation(t *testing.T) {
	tests := []struct {
		name string
		c    Content
		want string
	}{
		{
			"tab by default",
			"",
			"\t",
		},
		{
			"space missed in 'empty' line",
			" ",
			"\t",
		},
		{
			"space found in first line",
			" asd",
			" ",
		},
		{
			"double space found in second line",
			"asd\n  abd",
			"  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.findIndentation(); got != tt.want {
				t.Errorf("findIndentation() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestContent_toProjects(t *testing.T) {
	var (
		startAt = time.Date(2020, 2, 12, 0, 0, 0, 0, time.UTC)
		endAt   = time.Date(2020, 2, 20, 0, 0, 0, 0, time.UTC)
	)

	_, _ = startAt, endAt

	type args struct {
		indentation string
		dateFormat  string
		baseUrl     string
	}
	tests := []struct {
		name string
		c    Content
		args args
		want []Project
	}{
		{
			"empty",
			"",
			args{"", "2006-01-02", "http://example.com/"},
			nil,
		},
		{
			"1 milestone line",
			"|Initial milestone [2020-02-12]",
			args{"", "2006-01-02", "http://example.com/"},
			nil,
		},
		{
			"1 project line w/o dates",
			"Initial development",
			args{"  ", "2006-01-02", "http://example.com/"},
			[]Project{
				{0, "Initial development", nil, nil, 0, nil, 0},
			},
		},
		{
			"1 project line w dates",
			"Initial development [2020-02-12, 2020-02-20]",
			args{"  ", "2006-01-02", "http://example.com/"},
			[]Project{
				{0, "Initial development", &Dates{startAt, endAt}, nil, 0, nil, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.toProjects(tt.args.indentation, tt.args.dateFormat, tt.args.baseUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toProjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoadmap_ToContent(t *testing.T) {
	var (
		startAt1 = time.Date(2020, 2, 12, 0, 0, 0, 0, time.UTC)
		endAt1   = time.Date(2020, 2, 20, 0, 0, 0, 0, time.UTC)
		dates1   = Dates{StartAt: startAt1, EndAt: endAt1}
		startAt2 = time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC)
		endAt2   = time.Date(2020, 2, 25, 0, 0, 0, 0, time.UTC)
		dates2   = Dates{StartAt: startAt2, EndAt: endAt2}
		startAt3 = time.Date(2020, 2, 25, 0, 0, 0, 0, time.UTC)
		endAt3   = time.Date(2020, 2, 28, 0, 0, 0, 0, time.UTC)
		dates3   = Dates{StartAt: startAt3, EndAt: endAt3}
	)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   Content
	}{
		{
			"empty",
			fields{},
			"",
		},
		{
			"1 simple project",
			fields{
				DateFormat: "2006-01-02",
				Projects: []Project{
					{
						Title:      "Select and purchase domain",
						Percentage: 0,
					},
				},
			},
			"Select and purchase domain",
		},
		{
			"1 project with dates",
			fields{
				DateFormat: "2006-01-02",
				Projects: []Project{
					{
						Title:      "Select and purchase domain",
						Dates:      &dates1,
						Percentage: 0,
					},
				},
			},
			"Select and purchase domain [2020-02-12, 2020-02-20]",
		},
		{
			"1 project with color provided",
			fields{
				Projects: []Project{
					{
						Title:      "Select and purchase domain",
						Color:      &color.RGBA{R: 255, G: 0, B: 0, A: 255},
						Percentage: 0,
					},
				},
			},
			"Select and purchase domain [#ff0000]",
		},
		{
			"1 project with percentage",
			fields{
				Projects: []Project{
					{
						Title:      "Select and purchase domain",
						Percentage: 12,
					},
				},
			},
			"Select and purchase domain [12%]",
		},
		{
			"1 project with 1 url",
			fields{
				Projects: []Project{
					{
						Title:      "Select and purchase domain",
						URLs:       []string{"https://example.com/abc"},
						Percentage: 0,
					},
				},
			},
			"Select and purchase domain [https://example.com/abc]",
		},
		{
			"1 simple milestone",
			fields{
				DateFormat: "2006-01-02",
				Milestones: []Milestone{
					{Title: "Milestone 0.2"},
				},
			},
			"|Milestone 0.2",
		},
		{
			"1 milestone with deadline",
			fields{
				DateFormat: "2006-01-02",
				Milestones: []Milestone{
					{
						Title:      "Milestone 0.2",
						DeadlineAt: &startAt1,
					},
				},
			},
			"|Milestone 0.2 [2020-02-12]",
		},
		{
			"1 milestone with color",
			fields{
				DateFormat: "2006-01-02",
				Milestones: []Milestone{
					{
						Title: "Milestone 0.2",
						Color: &color.RGBA{R: 255, G: 0, B: 0, A: 255},
					},
				},
			},
			"|Milestone 0.2 [#ff0000]",
		},
		{
			"1 milestone with 1 url",
			fields{
				DateFormat: "2006-01-02",
				Milestones: []Milestone{
					{
						Title: "Milestone 0.2",
						URLs:  []string{"https://example.com/abc"},
					},
				},
			},
			"|Milestone 0.2 [https://example.com/abc]",
		},
		{
			"3 projects, 2 milestones",
			fields{
				DateFormat: "2006-01-02",
				Projects: []Project{
					{
						Title:      "Bring website online",
						Percentage: 0,
					},
					{
						Title:       "Select and purchase domain",
						Dates:       &dates2,
						Percentage:  85,
						Color:       &color.RGBA{R: 255, G: 0, B: 0, A: 255},
						URLs:        []string{"https://example.com/abc", "bcdef"},
						Indentation: 1,
						Milestone:   1,
					},
					{
						Title:       "Create server infrastructure",
						Dates:       &dates3,
						Percentage:  47,
						Indentation: 1,
						Milestone:   1,
					},
				},
				Milestones: []Milestone{
					{
						Title: "Milestone 0.1",
					},
					{
						Title:      "Milestone 0.2",
						DeadlineAt: &startAt1,
						Color:      &color.RGBA{R: 0, G: 255, B: 0, A: 255},
						URLs:       []string{"https://example.com/abc", "bcdef"},
					},
				},
			},
			`Bring website online
	Select and purchase domain [2020-02-04, 2020-02-25, 85%, #ff0000, https://example.com/abc, bcdef, |1]
	Create server infrastructure [2020-02-25, 2020-02-28, 47%, |1]

|Milestone 0.1
|Milestone 0.2 [2020-02-12, #00ff00, https://example.com/abc, bcdef]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Roadmap{
				ID:         tt.fields.ID,
				PrevID:     tt.fields.PrevID,
				DateFormat: tt.fields.DateFormat,
				BaseURL:    tt.fields.BaseURL,
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				AccessedAt: tt.fields.AccessedAt,
			}
			if got := r.ToContent(); got != tt.want {
				t.Errorf("ToContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoadmap_ToDates(t *testing.T) {
	var (
		startAt1   = time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
		endAt1     = time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC)
		dates1     = Dates{StartAt: startAt1, EndAt: endAt1}
		startAt2   = time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC)
		endAt2     = time.Date(2020, 2, 6, 0, 0, 0, 0, time.UTC)
		dates2     = Dates{StartAt: startAt2, EndAt: endAt2}
		dates1mix2 = Dates{StartAt: startAt1, EndAt: endAt2}
		startAt3   = time.Date(2020, 2, 3, 0, 0, 0, 0, time.UTC)
		endAt3     = time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC)
		dates3     = Dates{StartAt: startAt3, EndAt: endAt3}
	)

	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   *Dates
	}{
		{
			"empty",
			fields{},
			nil,
		},
		{
			"milestones don't count alone",
			fields{
				Milestones: []Milestone{
					{DeadlineAt: &startAt1},
				},
			},
			nil,
		},
		{
			"1 project",
			fields{
				Projects: []Project{
					{Dates: &dates1},
				},
			},
			&dates1,
		},
		{
			"2 overlapping projects",
			fields{
				Projects: []Project{
					{Dates: &dates1},
					{Dates: &dates2},
				},
			},
			&dates1mix2,
		},
		{
			"1 project with 2 milestones",
			fields{
				Projects: []Project{
					{Dates: &dates3},
				},
				Milestones: []Milestone{
					{DeadlineAt: &startAt1},
					{DeadlineAt: &endAt2},
				},
			},
			&dates1mix2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Roadmap{
				ID:         tt.fields.ID,
				PrevID:     tt.fields.PrevID,
				DateFormat: tt.fields.DateFormat,
				BaseURL:    tt.fields.BaseURL,
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				AccessedAt: tt.fields.AccessedAt,
			}
			if got := r.ToDates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContent_toMilestones(t *testing.T) {
	var (
		deadlineAt = time.Date(2020, 2, 12, 0, 0, 0, 0, time.UTC)
	)

	_ = deadlineAt

	type args struct {
		dateFormat string
		baseUrl    string
	}
	tests := []struct {
		name string
		c    Content
		args args
		want []Milestone
	}{
		{
			"empty",
			"",
			args{"2006-01-02", "http://example.com/"},
			nil,
		},
		{
			"1 project line",
			"Initial development [2020-02-12, 2020-02-20]",
			args{"2006-01-02", "http://example.com/"},
			nil,
		},
		{
			"1 milestone line",
			"|Initial milestone [2020-02-12]",
			args{"2006-01-02", "http://example.com/"},
			[]Milestone{
				{"Initial milestone", &deadlineAt, nil, nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.toMilestones(tt.args.dateFormat, tt.args.baseUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toMilestones() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitLine(t *testing.T) {
	type args struct {
		line        string
		indentation string
	}
	tests := []struct {
		name  string
		args  args
		want  uint8
		want1 string
		want2 string
	}{
		{
			"Bring website online",
			args{"Bring website online", "\t"},
			0,
			"Bring website online",
			"",
		},
		{
			"Initial development",
			args{"Initial development [2020-02-12, 2020-02-20]", "\t"},
			0,
			"Initial development",
			"2020-02-12, 2020-02-20",
		},
		{
			"Select and purchase domain",
			args{"\t\tSelect and purchase domain [2020-02-04, 2020-02-25, 100%, /issues/1]", "\t"},
			2,
			"Select and purchase domain",
			"2020-02-04, 2020-02-25, 100%, /issues/1",
		},
		{
			"Weird line",
			args{"\t\tSelect and purchase domain [2020-02-12, 2020-02-20] [2020-02-04, 2020-02-25, 100%, #f00, /issues/1, |2]", "\t"},
			2,
			"Select and purchase domain [2020-02-12, 2020-02-20]",
			"2020-02-04, 2020-02-25, 100%, #f00, /issues/1, |2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := splitLine(tt.args.line, tt.args.indentation)
			if got != tt.want {
				t.Errorf("splitLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitLine() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("splitLine() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_isLineProject(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"empty",
			args{line: ""},
			false,
		},
		{
			"simple project",
			args{line: "a|"},
			true,
		},
		{
			"simple project indented",
			args{line: "  a|"},
			true,
		},
		{
			"complex project",
			args{line: "  Select and purchase domain [2020-02-04, 2020-02-25, 100%, #f00, /issues/1, |2]"},
			true,
		},
		{
			"simple milestone",
			args{line: "|a"},
			false,
		},
		{
			"simple milestone indented",
			args{line: "  |a"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLineProject(tt.args.line); got != tt.want {
				t.Errorf("isLineProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLineMilestone(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"empty",
			args{line: ""},
			false,
		},
		{
			"simple project",
			args{line: "a|"},
			false,
		},
		{
			"simple project indented",
			args{line: "  a|"},
			false,
		},
		{
			"simple milestone",
			args{line: "|a"},
			true,
		},
		{
			"simple milestone indented",
			args{line: "  |a"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLineMilestone(tt.args.line); got != tt.want {
				t.Errorf("isLineMilestone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseExtraPart(t *testing.T) {
	const dateFormat = "2006-01-02"

	type args struct {
		part       string
		f          *time.Time
		t          *time.Time
		u          []string
		c          *color.RGBA
		p          uint8
		m          uint8
		dateFormat string
		baseUrl    string
	}

	var (
		nyeRaw = "2020-12-31"
		now    = time.Now()
		url    = "https://example.com/hello?nope=nope&why=why"
	)

	nye, _ := time.Parse(dateFormat, nyeRaw)

	_, _ = nye, now

	tests := []struct {
		name      string
		args      args
		startAt   *time.Time
		endAt     *time.Time
		urls      []string
		percent   uint8
		milestone uint8
		wantColor *color.RGBA
	}{
		{
			name:      "parse from",
			args:      args{part: "2020-12-31", dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   &nye,
			endAt:     nil,
			urls:      nil,
			percent:   0,
			milestone: 0,
			wantColor: nil,
		},
		{
			name:      "parse to",
			args:      args{part: "2020-12-31", f: &nye, dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   &nye,
			endAt:     &nye,
			urls:      nil,
			percent:   0,
			milestone: 0,
			wantColor: nil,
		},
		{
			name:      "parsing to overwrites existing to",
			args:      args{part: "2020-12-31", f: &nye, t: &now, dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   &nye,
			endAt:     &nye,
			urls:      nil,
			percent:   0,
			milestone: 0,
			wantColor: nil,
		},
		{
			name:      "parse url",
			args:      args{part: url},
			startAt:   nil,
			endAt:     nil,
			urls:      []string{url},
			percent:   0,
			milestone: 0,
			wantColor: nil,
		},
		{
			name:      "parsing url overwrites existing url",
			args:      args{part: "asd", u: []string{url}, dateFormat: "2006-01-02", baseUrl: "http://example.com/"},
			startAt:   nil,
			endAt:     nil,
			urls:      []string{url, "asd"},
			percent:   0,
			milestone: 0,
			wantColor: nil,
		},
		{
			name:      "parse color",
			args:      args{part: "#ffffff", dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   nil,
			endAt:     nil,
			urls:      nil,
			percent:   0,
			milestone: 0,
			wantColor: &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:      "parsing color overwrites existing color",
			args:      args{part: "#ffffff", c: &color.RGBA{R: 30, G: 20, B: 40, A: 30}, dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   nil,
			endAt:     nil,
			urls:      nil,
			percent:   0,
			milestone: 0,
			wantColor: &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:      "parse percentage",
			args:      args{part: "60%", dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   nil,
			endAt:     nil,
			urls:      nil,
			percent:   60,
			milestone: 0,
			wantColor: nil,
		},
		{
			name:      "parsing percentage overwrites existing percentage",
			args:      args{part: "60%", p: 30, dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   nil,
			endAt:     nil,
			urls:      nil,
			percent:   60,
			milestone: 0,
			wantColor: nil,
		},
		{
			name:      "parse milestone",
			args:      args{part: "|3", dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   nil,
			endAt:     nil,
			urls:      nil,
			percent:   0,
			milestone: 3,
			wantColor: nil,
		},
		{
			name:      "parsing milestone overwrites existing milestone",
			args:      args{part: "|3", m: 2, dateFormat: "2006-01-02", baseUrl: ""},
			startAt:   nil,
			endAt:     nil,
			urls:      nil,
			percent:   0,
			milestone: 3,
			wantColor: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startAt, endAt, urls, gotColor, percent, milestone := parseExtraPart(tt.args.part, tt.args.f, tt.args.t, tt.args.u, tt.args.c, tt.args.p, tt.args.m, tt.args.dateFormat, tt.args.baseUrl)
			if !reflect.DeepEqual(startAt, tt.startAt) {
				t.Errorf("parseProjectExtra() startAt = %v, want %v", startAt, tt.startAt)
			}
			if !reflect.DeepEqual(endAt, tt.endAt) {
				t.Errorf("parseProjectExtra() endAt = %v, want %v", endAt, tt.endAt)
			}
			if !reflect.DeepEqual(urls, tt.urls) {
				t.Errorf("parseProjectExtra() urls = %v, want %v", urls, tt.urls)
			}
			if !reflect.DeepEqual(gotColor, tt.wantColor) {
				t.Errorf("parseProjectExtra() gotColor = %v, want %v", gotColor, tt.wantColor)
			}
			if milestone != tt.milestone {
				t.Errorf("parseProjectExtra() milestone = %v, want %v", milestone, tt.milestone)
			}
			if percent != tt.percent {
				t.Errorf("parseProjectExtra() percent = %v, want %v", percent, tt.percent)
			}
		})
	}
}

func Test_parsePercentage(t *testing.T) {
	type args struct {
		part string
	}
	tests := []struct {
		name    string
		args    args
		want    uint8
		wantErr bool
	}{
		{
			name:    "10%",
			args:    args{part: "10%"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "foo",
			args:    args{part: "foo"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "bar%",
			args:    args{part: "bar%"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "200%",
			args:    args{part: "200%"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "-20%",
			args:    args{part: "-20%"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePercentage(tt.args.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePercentage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePercentage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMilestone(t *testing.T) {
	type args struct {
		part string
	}
	tests := []struct {
		name    string
		args    args
		want    uint8
		wantErr bool
	}{
		{
			name:    "|123",
			args:    args{part: "|123"},
			want:    123,
			wantErr: false,
		},
		{
			name:    "foo",
			args:    args{part: "foo"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "|bar",
			args:    args{part: "|bar"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "|",
			args:    args{part: "|"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "|234|",
			args:    args{part: "|234|"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMilestone(tt.args.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMilestone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMilestone() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseColor(t *testing.T) {
	type args struct {
		part string
	}

	const (
		hex21 = 2*16 + 1
		hex33 = 3*16 + 3
		hexfa = 15*16 + 10
	)

	tests := []struct {
		name    string
		args    args
		want    *color.RGBA
		wantErr bool
	}{
		{
			name:    "#333",
			args:    args{part: "#333"},
			want:    &color.RGBA{hex33, hex33, hex33, 255},
			wantErr: false,
		},
		{
			name:    "#332133",
			args:    args{part: "#332133"},
			want:    &color.RGBA{hex33, hex21, hex33, 255},
			wantErr: false,
		},
		{
			name:    "#fa2133",
			args:    args{part: "#fa2133"},
			want:    &color.RGBA{hexfa, hex21, hex33, 255},
			wantErr: false,
		},
		{
			name:    "332133",
			args:    args{part: "332133"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "#33",
			args:    args{part: "#33"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "#33213",
			args:    args{part: "#33213"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseColor(tt.args.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseColor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
