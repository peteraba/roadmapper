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
	type args struct {
		line            string
		previousProject *internalProject
		pi              int
		colorNum        *uint8
	}

	var a, b, c, d, e uint8

	_, _, _, _, _ = a, b, c, d, e

	tests := []struct {
		name    string
		args    args
		want    *internalProject
		want1   int
		wantErr bool
	}{
		{
			name: "too few indentations",
			args: args{
				line:            "asd",
				previousProject: &internalProject{},
				pi:              4,
				colorNum:        &a,
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "too many indentations",
			args: args{
				line:            "      asd",
				previousProject: &internalProject{},
				pi:              0,
				colorNum:        &a,
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "level1 initial",
			args: args{
				line:            "asd",
				previousProject: &internalProject{},
				pi:              -2,
				colorNum:        &b,
			},
			want: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want1:   0,
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
				pi:       0,
				colorNum: &c,
			},
			want: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want1:   0,
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
				pi:       4,
				colorNum: &d,
			},
			want: &internalProject{
				Title:      "asd",
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := createProject(tt.args.line, tt.args.previousProject, tt.args.pi, tt.args.colorNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("createProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("createProject() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("createProject() got1 = %v, want %v", got1, tt.want1)
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

func Test_internalProject_GetFrom(t *testing.T) {
	type fields struct {
		from         *time.Time
		childrenFrom *time.Time
	}

	var now = time.Now()

	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "returns from attribute value by default",
			fields: fields{
				from:         &now,
				childrenFrom: nil,
			},
			want: now,
		},
		{
			name: "returns childrenFrom attribute value in case from is nil",
			fields: fields{
				from:         nil,
				childrenFrom: &now,
			},
			want: now,
		},
		{
			name: "always returns a time.Time instance",
			fields: fields{
				from:         nil,
				childrenFrom: nil,
			},
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := internalProject{
				from:         tt.fields.from,
				childrenFrom: tt.fields.childrenFrom,
			}
			if got := p.GetFrom(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_internalProject_GetTo(t *testing.T) {
	type fields struct {
		to         *time.Time
		childrenTo *time.Time
	}

	var now = time.Now()

	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "returns to attribute value by default",
			fields: fields{
				to:         &now,
				childrenTo: nil,
			},
			want: now,
		},
		{
			name: "returns childrenTo attribute value in case from is nil",
			fields: fields{
				to:         nil,
				childrenTo: &now,
			},
			want: now,
		},
		{
			name: "always returns a time.Time instance",
			fields: fields{
				to:         nil,
				childrenTo: nil,
			},
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := internalProject{
				to:         tt.fields.to,
				childrenTo: tt.fields.childrenTo,
			}
			if got := p.GetTo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTo() = %v, want %v", got, tt.want)
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
		trimmed  string
		colorNum *uint8
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
				trimmed:  "lorem ipsum",
				colorNum: &a,
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
				trimmed:  "dates only [2020-12-04, 2020-12-31]",
				colorNum: &b,
			},
			want: &internalProject{
				Title:      "dates only",
				from:       &dec4,
				to:         &dec31,
				color:      palette.WebSafe[71],
				percentage: 100,
			},
			wantErr: false,
		},
		{
			name: "all there 1",
			args: args{
				trimmed:  "all there [2020-12-04, 2020-12-31, #f949b9, https://example.com/ah?os=linux&browser=firefox, 53%]",
				colorNum: &c,
			},
			want: &internalProject{
				Title:      "all there",
				from:       &dec4,
				to:         &dec31,
				color:      color.RGBA{R: 15*16 + 9, G: 4*16 + 9, B: 11*16 + 9, A: 255},
				percentage: 53,
				url:        "https://example.com/ah?os=linux&browser=firefox",
			},
			wantErr: false,
		},
		{
			name: "all there 2",
			args: args{
				trimmed:  "all there 2 [#F4b, http://example.com/ah?os=linux&browser=firefox, 2020-12-04, 2020-12-31, 0.53]",
				colorNum: &d,
			},
			want: &internalProject{
				Title:      "all there 2",
				from:       &dec4,
				to:         &dec31,
				color:      color.RGBA{R: 15*16 + 15, G: 4*16 + 4, B: 11*16 + 11, A: 255},
				percentage: 53,
				url:        "http://example.com/ah?os=linux&browser=firefox",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProject(tt.args.trimmed, tt.args.colorNum)
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
		lines []string
	}
	var (
		d0 = time.Date(2020, 10, 8, 0, 0, 0, 0, time.UTC)
		d1 = time.Date(2020, 11, 28, 0, 0, 0, 0, time.UTC)
	)
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
			name: "project without subprojects must have dates set",
			args: args{
				lines: []string{"Simple project, no brackets"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "single, simple project",
			args: args{
				lines: []string{"Simple project, dates only [2020-10-08, 2020-11-28]"},
			},
			want: &internalProject{
				children: []*internalProject{
					{
						Title:      "Simple project, dates only",
						from:       &d0,
						to:         &d1,
						percentage: 100,
						color:      color.RGBA{R: 102, G: 51, B: 204, A: 255},
					},
				},
				childrenFrom: &d0,
				childrenTo:   &d1,
			},
			wantErr: false,
		},
		{
			name: "simple project with one sub-project and no dates",
			args: args{
				lines: []string{
					"Rather simple project",
					"  Simple sub-project, dates only [2020-10-08, 2020-11-28]",
				},
			},
			want: &internalProject{
				children: []*internalProject{
					{
						Title:        "Rather simple project",
						childrenFrom: &d0,
						childrenTo:   &d1,
						from:         &d0,
						to:           &d1,
						percentage:   100,
						color:        color.RGBA{R: 102, G: 51, B: 204, A: 255},
						children: []*internalProject{
							{
								Title:      "Simple sub-project, dates only",
								from:       &d0,
								to:         &d1,
								percentage: 100,
								color:      color.RGBA{R: 204, G: 51, B: 153, A: 255},
							},
						},
					},
				},
				childrenFrom: &d0,
				childrenTo:   &d1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRoadmap(tt.args.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRoadmap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("parseRoadmap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setChildrenDates(t *testing.T) {
	type args struct {
		p *internalProject
	}

	d0, _ := time.Parse(dateFormat, "2030-08-01")
	d1, _ := time.Parse(dateFormat, "2030-09-17")
	d2, _ := time.Parse(dateFormat, "2030-10-22")
	d3, _ := time.Parse(dateFormat, "2030-12-01")

	tests := []struct {
		name    string
		args    args
		want    *time.Time
		want1   *time.Time
		wantErr bool
	}{
		{
			name: "empty project without children needs starting and end dates",
			args: args{
				p: &internalProject{},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "empty project without children needs end date",
			args: args{
				p: &internalProject{from: &d0},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "empty project without children needs starting date",
			args: args{
				p: &internalProject{to: &d0},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "starting date must not be before ending date",
			args: args{
				p: &internalProject{from: &d1, to: &d0},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "empty project with starting and end dates",
			args: args{
				p: &internalProject{from: &d0, to: &d1},
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
						{from: &d0, to: &d1},
					},
				},
			},
			want:    &d0,
			want1:   &d1,
			wantErr: false,
		},
		{
			name: "project with children will have minimum of froms and maximum of tos set as childrenFrom and childrenTo",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{from: &d0, to: &d2},
						{from: &d1, to: &d3},
					},
				},
			},
			want:    &d0,
			want1:   &d3,
			wantErr: false,
		},
		{
			name: "project with children must have from equal childrenFrom if set",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{from: &d0, to: &d2},
						{from: &d1, to: &d3},
					},
					from: &d1,
				},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "project with children must have from equal childrenFrom if set",
			args: args{
				p: &internalProject{
					children: []*internalProject{
						{from: &d0, to: &d2},
						{from: &d1, to: &d3},
					},
					to: &d1,
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
	type args struct {
		part string
		f    *time.Time
		t    *time.Time
		u    string
		p    uint8
		c    color.Color
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
			args:  args{part: "2020-12-31"},
			want:  &nye,
			want1: nil,
			want2: "",
			want3: 0,
			want4: nil,
		},
		{
			name:  "parse to",
			args:  args{part: "2020-12-31", f: &nye},
			want:  &nye,
			want1: &nye,
			want2: "",
			want3: 0,
			want4: nil,
		},
		{
			name:  "parsing to overwrites existing to",
			args:  args{part: "2020-12-31", f: &nye, t: &now},
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
			args:  args{part: url, u: "asd"},
			want:  nil,
			want1: nil,
			want2: url,
			want3: 0,
			want4: nil,
		},
		{
			name:  "parse percentage",
			args:  args{part: "60"},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 60,
			want4: nil,
		},
		{
			name:  "parsing percentage overwrites existing percentage",
			args:  args{part: "60", p: 30},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 60,
			want4: nil,
		},
		{
			name:  "parse color",
			args:  args{part: "#ffffff"},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 0,
			want4: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:  "parsing color overwrites existing color",
			args:  args{part: "#ffffff", c: color.RGBA{R: 30, G: 20, B: 40, A: 30}},
			want:  nil,
			want1: nil,
			want2: "",
			want3: 0,
			want4: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, got4 := parseProjectExtra(tt.args.part, tt.args.f, tt.args.t, tt.args.u, tt.args.p, tt.args.c)
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
	type args struct {
		roadmapFrom time.Time
		roadmapTo   time.Time
	}

	var (
		t0 = "Lorem Ipsum"
		t1 = "Nullam vulputate"
		t2 = "Curabitur ullamcorper condimentum"
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
				Title:        t0,
				percentage:   30,
				childrenFrom: &d0,
				childrenTo:   &d3,
				children: []*internalProject{
					{
						Title:      t1,
						from:       &d0,
						to:         &d2,
						percentage: 34,
					},
					{
						Title:      t2,
						from:       &d1,
						to:         &d3,
						percentage: 18,
					},
				},
			},
			args: args{
				roadmapFrom: d0,
				roadmapTo:   d3,
			},
			want: Project{
				Title:      t0,
				Percentage: 30,
				From:       d0,
				To:         d3,
				Children: []Project{
					{
						Title:      t1,
						From:       d0,
						To:         d2,
						Percentage: 34,
					},
					{
						Title:      t2,
						From:       d1,
						To:         d3,
						Percentage: 18,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &internalProject{
				Title:        tt.fields.Title,
				from:         tt.fields.from,
				to:           tt.fields.to,
				parent:       tt.fields.parent,
				color:        tt.fields.color,
				percentage:   tt.fields.percentage,
				url:          tt.fields.url,
				children:     tt.fields.children,
				childrenFrom: tt.fields.childrenFrom,
				childrenTo:   tt.fields.childrenTo,
			}
			if got := p.ToPublic(tt.args.roadmapFrom, tt.args.roadmapTo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPublic() = %v, want %v", got, tt.want)
			}
		})
	}
}
