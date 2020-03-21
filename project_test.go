package main

import (
	"image/color"
	"image/color/palette"
	"reflect"
	"testing"
	"time"
)

func Test_charsToUint8(t *testing.T) {
	type args struct {
		part string
	}
	tests := []struct {
		name    string
		args    args
		want    [3]uint8
		wantErr bool
	}{
		{
			name:    "000000",
			args:    args{part: "000000"},
			want:    [3]uint8{0, 0, 0},
			wantErr: false,
		},
		{
			name:    "ffffff",
			args:    args{part: "ffffff"},
			want:    [3]uint8{255, 255, 255},
			wantErr: false,
		},
		{
			name:    "f2f3aa",
			args:    args{part: "f2f3aa"},
			want:    [3]uint8{15*16 + 2, 15*16 + 3, 10*16 + 10},
			wantErr: false,
		},
		{
			name:    "ff",
			args:    args{part: "ff"},
			want:    [3]uint8{},
			wantErr: true,
		},
		{
			name:    "ffff",
			args:    args{part: "ffff"},
			want:    [3]uint8{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := charsToUint8(tt.args.part)
			if (err != nil) != tt.wantErr {
				t.Errorf("charsToUint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("charsToUint8() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createProject(t *testing.T) {
	var (
		start = time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC)
		end   = time.Date(2020, 3, 20, 0, 0, 0, 0, time.UTC)
	)
	type args struct {
		line            string
		previousProject *internalProject
		pi              int
		colorNum        *uint8
		dateFormat      string
		baseUrl         string
	}

	var a, b, c, d, e, f, g, h, i uint8

	_, _, _, _, _, _, _, _, _ = a, b, c, d, e, f, g, h, i

	tests := []struct {
		name    string
		args    args
		want1   *internalProject
		want2   int
		wantErr bool
	}{
		{
			name: "too few indentations",
			args: args{
				line:            "asd",
				previousProject: &internalProject{},
				pi:              4,
				colorNum:        &a,
				dateFormat:      "2006-01-02",
				baseUrl:         "",
			},
			want1:   nil,
			want2:   0,
			wantErr: true,
		},
		{
			name: "too many indentations",
			args: args{
				line:            "\t\t\tasd",
				previousProject: &internalProject{},
				pi:              0,
				colorNum:        &a,
				dateFormat:      "2006-01-02",
				baseUrl:         "",
			},
			want1:   nil,
			want2:   0,
			wantErr: true,
		},
		{
			name: "level1 initial",
			args: args{
				line:            "asd",
				previousProject: &internalProject{},
				pi:              -1,
				colorNum:        &b,
				dateFormat:      "2006-01-02",
				baseUrl:         "",
			},
			want1: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want2:   0,
			wantErr: false,
		},
		{
			name: "level1 with default dates",
			args: args{
				line:            "asd [2020-03-15, 2020-03-20]",
				previousProject: &internalProject{},
				pi:              -1,
				colorNum:        &c,
				dateFormat:      "2006-01-02",
				baseUrl:         "",
			},
			want1: &internalProject{
				Title:      "asd",
				start:      &start,
				end:        &end,
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want2:   0,
			wantErr: false,
		},
		{
			name: "level1 with German dates",
			args: args{
				line:            "asd [15.03.2020, 20.03.2020]",
				previousProject: &internalProject{},
				pi:              -1,
				colorNum:        &d,
				dateFormat:      "02.01.2006",
				baseUrl:         "",
			},
			want1: &internalProject{
				Title:      "asd",
				start:      &start,
				end:        &end,
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want2:   0,
			wantErr: false,
		},
		{
			name: "level1 with url",
			args: args{
				line:            "asd [https://gist.github.com/]",
				previousProject: &internalProject{},
				pi:              -1,
				colorNum:        &e,
				dateFormat:      "02.01.2006",
				baseUrl:         "",
			},
			want1: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
				url:        "https://gist.github.com/",
			},
			want2:   0,
			wantErr: false,
		},
		{
			name: "level1 with url and base url",
			args: args{
				line:            "asd [dsa]",
				previousProject: &internalProject{},
				pi:              -1,
				colorNum:        &f,
				dateFormat:      "02.01.2006",
				baseUrl:         "https://gist.github.com/",
			},
			want1: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
				url:        "https://gist.github.com/dsa",
			},
			want2:   0,
			wantErr: false,
		},
		{
			name: "level1 with dates, long color code and url with base url",
			args: args{
				line:            "asd [15.03.2020, 20.03.2020, #a3a3a3, dsa]",
				previousProject: &internalProject{},
				pi:              -1,
				colorNum:        &f,
				dateFormat:      "02.01.2006",
				baseUrl:         "https://gist.github.com/",
			},
			want1: &internalProject{
				Title:      "asd",
				start:      &start,
				end:        &end,
				color:      color.RGBA{163, 163, 163, 255},
				percentage: 100,
				url:        "https://gist.github.com/dsa",
			},
			want2:   0,
			wantErr: false,
		},
		{
			name: "level1 after level1",
			args: args{
				line: "asd",
				previousProject: &internalProject{
					Title: "prev-Level1",
					parent: &internalProject{
						Title: "root",
					},
				},
				pi:         0,
				colorNum:   &g,
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want1: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want2:   0,
			wantErr: false,
		},
		{
			name: "level1 after level3",
			args: args{
				line: "asd",
				previousProject: &internalProject{
					Title: "prev-Level3",
					parent: &internalProject{
						Title: "prevParent-Level2",
						parent: &internalProject{
							Title: "prevParentParent-Level1",
							parent: &internalProject{
								Title: "root",
							},
						},
					},
				},
				pi:         2,
				colorNum:   &h,
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want1: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want2:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2, err := createProject(tt.args.line, tt.args.previousProject, tt.args.pi, tt.args.colorNum, tt.args.dateFormat, tt.args.baseUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("createProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != nil && !reflect.DeepEqual(got1.String(), tt.want1.String()) {
				t.Errorf("createProject() got1 = %v, want1 %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("createProject() got2 = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}

func Test_getNextColor(t *testing.T) {
	type args struct {
		colorNum *uint8
		nextNum  uint8
	}

	var a1, b1, c1 uint8 = 0, 125, 200
	var a2, b2, c2 uint8 = 71, 196, 271 % 256

	tests := []struct {
		name string
		args args
		want color.Color
	}{
		{
			name: "a",
			args: args{colorNum: &a1, nextNum: a2},
			want: palette.WebSafe[a2],
		},
		{
			name: "b",
			args: args{colorNum: &b1, nextNum: b2},
			want: palette.WebSafe[b2],
		},
		{
			name: "c",
			args: args{colorNum: &c1, nextNum: c2},
			want: palette.WebSafe[c2],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNextColor(tt.args.colorNum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNextColor() = %v, want %v", got, tt.want)
			}
			if *tt.args.colorNum != tt.args.nextNum {
				t.Errorf("getNextColor() -> colorNum: %v, want %v", *tt.args.colorNum, tt.args.nextNum)
			}
		})
	}
}

func Test_internalProject_GetStart(t *testing.T) {
	type fields struct {
		start         *time.Time
		childrenStart *time.Time
	}

	var now = time.Now()

	tests := []struct {
		name   string
		fields fields
		want   *time.Time
	}{
		{
			name: "returns start attribute value by default",
			fields: fields{
				start:         &now,
				childrenStart: nil,
			},
			want: &now,
		},
		{
			name: "returns childrenStart attribute value in case start is nil",
			fields: fields{
				start:         nil,
				childrenStart: &now,
			},
			want: &now,
		},
		{
			name: "returns nil if start and childrenStart are nil",
			fields: fields{
				start:         nil,
				childrenStart: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := internalProject{
				start:         tt.fields.start,
				childrenStart: tt.fields.childrenStart,
			}
			if got := p.GetStart(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_internalProject_GetEnd(t *testing.T) {
	type fields struct {
		end        *time.Time
		childrenTo *time.Time
	}

	var now = time.Now()

	tests := []struct {
		name   string
		fields fields
		want   *time.Time
	}{
		{
			name: "returns end attribute value by default",
			fields: fields{
				end:        &now,
				childrenTo: nil,
			},
			want: &now,
		},
		{
			name: "returns childrenEnd attribute value in case end is nil",
			fields: fields{
				end:        nil,
				childrenTo: &now,
			},
			want: &now,
		},
		{
			name: "returns nil if end and childrenEnd are nil",
			fields: fields{
				end:        nil,
				childrenTo: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := internalProject{
				end:         tt.fields.end,
				childrenEnd: tt.fields.childrenTo,
			}
			if got := p.GetEnd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnd() = %v, want %v", got, tt.want)
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
		want    color.Color
		wantErr bool
	}{
		{
			name:    "#333",
			args:    args{part: "#333"},
			want:    color.RGBA{hex33, hex33, hex33, 255},
			wantErr: false,
		},
		{
			name:    "#332133",
			args:    args{part: "#332133"},
			want:    color.RGBA{hex33, hex21, hex33, 255},
			wantErr: false,
		},
		{
			name:    "#fa2133",
			args:    args{part: "#fa2133"},
			want:    color.RGBA{hexfa, hex21, hex33, 255},
			wantErr: false,
		},
		{
			name:    "332133",
			args:    args{part: "332133"},
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name:    "#33",
			args:    args{part: "#33"},
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name:    "#33213",
			args:    args{part: "#33213"},
			want:    color.RGBA{},
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
			name:    "10",
			args:    args{part: "10"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "0.20",
			args:    args{part: "0.20"},
			want:    20,
			wantErr: false,
		},
		{
			name:    "0.300005",
			args:    args{part: "0.300005"},
			want:    30,
			wantErr: false,
		},
		{
			name:    "40.0005",
			args:    args{part: "0.400005"},
			want:    40,
			wantErr: false,
		},
		{
			name:    "50%",
			args:    args{part: "50%"},
			want:    50,
			wantErr: false,
		},
		{
			name:    "as",
			args:    args{part: "as"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "200",
			args:    args{part: "200"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "-20",
			args:    args{part: "-20"},
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

func Test_parseProject(t *testing.T) {
	type args struct {
		trimmed    string
		colorNum   *uint8
		dateFormat string
		baseUrl    string
	}

	var (
		a, b, c, d, e uint8
		dec4          = time.Date(2020, 12, 4, 0, 0, 0, 0, time.UTC)
		dec31         = time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	)

	_, _, _, _, _ = a, b, c, d, e
	_, _ = dec4, dec31

	tests := []struct {
		name    string
		args    args
		want    *internalProject
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				trimmed:    "lorem ipsum",
				colorNum:   &a,
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want: &internalProject{
				Title:      "lorem ipsum",
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			wantErr: false,
		},
		{
			name: "dates only",
			args: args{
				trimmed:    "dates only [2020-12-04, 2020-12-31]",
				colorNum:   &b,
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want: &internalProject{
				Title:      "dates only",
				start:      &dec4,
				end:        &dec31,
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			wantErr: false,
		},
		{
			name: "all there 1",
			args: args{
				trimmed:    "all there [2020-12-04, 2020-12-31, #f949b9, https://example.com/ah?os=linux&browser=firefox, 53%]",
				colorNum:   &c,
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want: &internalProject{
				Title:      "all there",
				start:      &dec4,
				end:        &dec31,
				color:      color.RGBA{R: 15*16 + 9, G: 4*16 + 9, B: 11*16 + 9, A: 255},
				percentage: 53,
				url:        "https://example.com/ah?os=linux&browser=firefox",
			},
			wantErr: false,
		},
		{
			name: "all there 2",
			args: args{
				trimmed:    "all there 2 [#F4b, http://example.com/ah?os=linux&browser=firefox, 2020-12-04, 2020-12-31, 0.53]",
				colorNum:   &d,
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want: &internalProject{
				Title:      "all there 2",
				start:      &dec4,
				end:        &dec31,
				color:      color.RGBA{R: 15*16 + 15, G: 4*16 + 4, B: 11*16 + 11, A: 255},
				percentage: 53,
				url:        "http://example.com/ah?os=linux&browser=firefox",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProject(tt.args.trimmed, tt.args.colorNum, tt.args.dateFormat, tt.args.baseUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseProject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseRoadmap(t *testing.T) {
	type args struct {
		lines      []string
		dateFormat string
		baseUrl    string
	}

	var (
		d0                     = time.Date(2020, 10, 8, 0, 0, 0, 0, time.UTC)
		d1                     = time.Date(2020, 11, 28, 0, 0, 0, 0, time.UTC)
		initialStart           = time.Date(2020, 2, 12, 0, 0, 0, 0, time.UTC)
		initialEnd             = time.Date(2020, 2, 20, 0, 0, 0, 0, time.UTC)
		selectAndPurchaseStart = time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC)
		selectAndPurchaseEnd   = time.Date(2020, 2, 25, 0, 0, 0, 0, time.UTC)
		createServerStart      = time.Date(2020, 2, 25, 0, 0, 0, 0, time.UTC)
		createServerEnd        = time.Date(2020, 2, 28, 0, 0, 0, 0, time.UTC)
	)
	_, _ = d0, d1
	tests := []struct {
		name    string
		args    args
		want    *internalProject
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			want:    &internalProject{},
			wantErr: false,
		},
		{
			name: "project without subprojects",
			args: args{
				lines:      []string{"Simple project, no brackets"},
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want: &internalProject{
				children: []*internalProject{
					{
						Title:      "Simple project, no brackets",
						color:      color.RGBA{R: 102, G: 51, B: 204, A: 255},
						percentage: 100,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "single, simple project",
			args: args{
				lines:      []string{"Simple project, dates only [2020-10-08, 2020-11-28]"},
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want: &internalProject{
				children: []*internalProject{
					{
						Title:      "Simple project, dates only",
						start:      &d0,
						end:        &d1,
						percentage: 100,
						color:      color.RGBA{R: 102, G: 51, B: 204, A: 255},
					},
				},
				childrenStart: &d0,
				childrenEnd:   &d1,
			},
			wantErr: false,
		},
		{
			name: "simple project with one sub-project and no dates",
			args: args{
				lines: []string{
					"Rather simple project",
					"\tSimple sub-project, dates only [2020-10-08, 2020-11-28]",
				},
				dateFormat: "2006-01-02",
				baseUrl:    "",
			},
			want: &internalProject{
				children: []*internalProject{
					{
						Title:         "Rather simple project",
						childrenStart: &d0,
						childrenEnd:   &d1,
						start:         &d0,
						end:           &d1,
						percentage:    100,
						color:         color.RGBA{R: 102, G: 51, B: 204, A: 255},
						children: []*internalProject{
							{
								Title:      "Simple sub-project, dates only",
								start:      &d0,
								end:        &d1,
								percentage: 100,
								color:      color.RGBA{R: 204, G: 51, B: 153, A: 255},
							},
						},
					},
				},
				childrenStart: &d0,
				childrenEnd:   &d1,
			},
			wantErr: false,
		},
		{
			name: "complex example",
			args: args{
				lines: []string{
					"Initial development [2020-02-12, 2020-02-20]",
					"Bring website online",
					"\tSelect and purchase domain [2020-02-04, 2020-02-25, 100%, /issues/1, #434]",
					"\tCreate server infrastructure [#434, 2020-02-25, 2020-02-28, 100%, https://github.com/peteraba/roadmapper/issues/2]",
					// "Command line tool",
					// "\tCreate backend SVG generation [2020-03-03, 2020-03-10, 100%]",
					// "\tReplace frontend SVG generation with backend [2020-03-08, 2020-03-12, 100%]",
					// "\tCreate documentation page [2020-03-13, 2020-03-31, 20%]",
				},
				dateFormat: "2006-01-02",
				baseUrl:    "https://github.com/peteraba/roadmapper",
			},
			want: &internalProject{
				children: []*internalProject{
					{
						Title:      "Initial development",
						start:      &initialStart,
						end:        &initialEnd,
						percentage: 100,
						color:      color.RGBA{R: 102, G: 51, B: 204, A: 255},
					},
					{
						Title:         "Bring website online",
						childrenStart: &selectAndPurchaseStart,
						childrenEnd:   &createServerEnd,
						percentage:    100,
						color:         color.RGBA{R: 204, G: 51, B: 153, A: 255},
						children: []*internalProject{
							{
								Title:      "Select and purchase domain",
								start:      &selectAndPurchaseStart,
								end:        &selectAndPurchaseEnd,
								percentage: 100,
								color:      color.RGBA{R: 68, G: 51, B: 68, A: 255},
								url:        "https://github.com/peteraba/roadmapper/issues/1",
							},
							{
								Title:      "Create server infrastructure",
								start:      &createServerStart,
								end:        &createServerEnd,
								percentage: 100,
								color:      color.RGBA{R: 68, G: 51, B: 68, A: 255},
								url:        "https://github.com/peteraba/roadmapper/issues/2",
							},
						},
					},
				},
				childrenStart: &selectAndPurchaseStart,
				childrenEnd:   &createServerEnd,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProjects(tt.args.lines, tt.args.dateFormat, tt.args.baseUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("parseProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setChildrenDates(t *testing.T) {
	type args struct {
		p *internalProject
	}

	d0, _ := time.Parse("2006-01-02", "2030-08-01")
	d1, _ := time.Parse("2006-01-02", "2030-09-17")
	d2, _ := time.Parse("2006-01-02", "2030-10-22")
	d3, _ := time.Parse("2006-01-02", "2030-12-01")
	dm1, _ := time.Parse("2/1/2006", "10/10/2010")
	d4, _ := time.Parse("2/1/2006", "8/8/2088")

	tests := []struct {
		name    string
		args    args
		want    *time.Time
		want1   *time.Time
		wantErr bool
	}{
		{
			name: "empty project",
			args: args{
				p: &internalProject{},
			},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
		{
			name: "starting date must not be before ending date",
			args: args{
				p: &internalProject{start: &d1, end: &d0},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "empty project with starting and end dates",
			args: args{
				p: &internalProject{start: &d0, end: &d1},
			},
			want:    &d0,
			want1:   &d1,
			wantErr: false,
		},
		{
			name: "project with children does not need starting or end dates",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{start: &d0, end: &d1},
					},
				},
			},
			want:    &d0,
			want1:   &d1,
			wantErr: false,
		},
		{
			name: "project with children will have minimum of starts and maximum of ends set as childrenStart and childrenEnd",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{start: &d0, end: &d2},
						{start: &d1, end: &d3},
					},
				},
			},
			want:    &d0,
			want1:   &d3,
			wantErr: false,
		},
		{
			name: "children with missing start or end date are ignored for finding minimum or maximum",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{start: &d0, end: &d2},
						{start: &dm1},
						{end: &d4},
						{start: &d1, end: &d3},
					},
				},
			},
			want:    &d0,
			want1:   &d3,
			wantErr: false,
		},
		{
			name: "project with children must have start equal childrenStart if set",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{start: &d0, end: &d2},
						{start: &d1, end: &d3},
					},
					start: &d1,
				},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "project with children must have from equal childrenStart if set",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{start: &d0, end: &d2},
						{start: &d1, end: &d3},
					},
					end: &d1,
				},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := setChildrenDates(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("setChildrenDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setChildrenDates() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("setChildrenDates() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parseProjectExtra(t *testing.T) {
	const dateFormat = "2006-01-02"

	type args struct {
		part       string
		f          *time.Time
		t          *time.Time
		u          string
		p          uint8
		c          color.Color
		dateFormat string
		baseUrl    string
	}

	var (
		nyeRaw = "2020-12-31"
		now    = time.Now()
		url    = "https://example.com/hello?nope=nope&why=why"
	)

	nye, _ := time.Parse(dateFormat, nyeRaw)

	tests := []struct {
		name  string
		args  args
		want  *time.Time
		want1 *time.Time
		want2 string
		want3 uint8
		want4 color.Color
	}{
		{
			name:  "parse from",
			args:  args{part: "2020-12-31", dateFormat: "2006-01-02", baseUrl: ""},
			want:  &nye,
			want1: nil,
			want2: "",
			want3: 0,
			want4: nil,
		},
		{
			name:  "parse to",
			args:  args{part: "2020-12-31", f: &nye, dateFormat: "2006-01-02", baseUrl: ""},
			want:  &nye,
			want1: &nye,
			want2: "",
			want3: 0,
			want4: nil,
		},
		{
			name:  "parsing to overwrites existing to",
			args:  args{part: "2020-12-31", f: &nye, t: &now, dateFormat: "2006-01-02", baseUrl: ""},
			want:  &nye,
			want1: &nye,
			want2: "",
			want3: 0,
			want4: nil,
		},
		{
			name:  "parse url",
			args:  args{part: url},
			want:  nil,
			want1: nil,
			want2: url,
			want3: 0,
			want4: nil,
		},
		{
			name:  "parsing url overwrites existing url",
			args:  args{part: url, u: "asd", dateFormat: "2006-01-02", baseUrl: ""},
			want:  nil,
			want1: nil,
			want2: url,
			want3: 0,
			want4: nil,
		},
		{
			name:  "parse percentage",
			args:  args{part: "60", dateFormat: "2006-01-02", baseUrl: ""},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 60,
			want4: nil,
		},
		{
			name:  "parsing percentage overwrites existing percentage",
			args:  args{part: "60", p: 30, dateFormat: "2006-01-02", baseUrl: ""},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 60,
			want4: nil,
		},
		{
			name:  "parse color",
			args:  args{part: "#ffffff", dateFormat: "2006-01-02", baseUrl: ""},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 0,
			want4: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:  "parsing color overwrites existing color",
			args:  args{part: "#ffffff", c: color.RGBA{R: 30, G: 20, B: 40, A: 30}, dateFormat: "2006-01-02", baseUrl: ""},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 0,
			want4: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, got4 := parseProjectExtra(tt.args.part, tt.args.f, tt.args.t, tt.args.u, tt.args.p, tt.args.c, tt.args.dateFormat, tt.args.baseUrl)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseProjectExtra() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseProjectExtra() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("parseProjectExtra() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("parseProjectExtra() got3 = %v, want %v", got3, tt.want3)
			}
			if !reflect.DeepEqual(got4, tt.want4) {
				t.Errorf("parseProjectExtra() got4 = %v, want %v", got4, tt.want4)
			}
		})
	}
}

func Test_internalProject_ToPublic(t *testing.T) {
	type fields struct {
		Title         string
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
	type args struct {
		roadmapStart *time.Time
		roadmapEnd   *time.Time
	}

	var (
		d0 = time.Date(2020, 0, 0, 0, 0, 0, 0, time.UTC)
		d1 = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		d2 = time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC)
		d3 = time.Date(2020, 3, 3, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Project
	}{
		{
			name: "complex example",
			fields: fields{
				Title:         "Lorem Ipsum",
				percentage:    30,
				childrenStart: &d0,
				childrenEnd:   &d3,
				children: []*internalProject{
					{
						Title:      "Nullam vulputate",
						start:      &d0,
						end:        &d2,
						percentage: 34,
					},
					{
						Title:      "Curabitur ullamcorper condimentum",
						start:      &d1,
						end:        &d3,
						percentage: 18,
					},
				},
			},
			args: args{
				roadmapStart: &d0,
				roadmapEnd:   &d3,
			},
			want: Project{
				Title:      "Lorem Ipsum",
				Percentage: 30,
				Dates:      &Dates{Start: d0, End: d3},
				Children: []Project{
					{
						Title:      "Nullam vulputate",
						Dates:      &Dates{Start: d0, End: d2},
						Percentage: 34,
					},
					{
						Title:      "Curabitur ullamcorper condimentum",
						Dates:      &Dates{Start: d1, End: d3},
						Percentage: 18,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &internalProject{
				Title:         tt.fields.Title,
				start:         tt.fields.start,
				end:           tt.fields.end,
				parent:        tt.fields.parent,
				color:         tt.fields.color,
				percentage:    tt.fields.percentage,
				url:           tt.fields.url,
				children:      tt.fields.children,
				childrenStart: tt.fields.childrenStart,
				childrenEnd:   tt.fields.childrenEnd,
			}
			if got := p.ToPublic(tt.args.roadmapStart, tt.args.roadmapEnd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPublic() = %v, want %v", got, tt.want)
			}
		})
	}
}
