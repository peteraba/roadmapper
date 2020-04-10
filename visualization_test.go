package main

import (
	"image/color"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/tdewolff/canvas"
)

func TestRoadmap_ToVisual(t *testing.T) {
	rand.Seed(0)

	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)
	dates0419 := time.Date(2020, 4, 19, 0, 0, 0, 0, time.UTC)
	dates0420 := time.Date(2020, 4, 20, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	_, _, _, _, _, _, _ = dates0402, dates0405, dates0408, dates0415, dates0418, dates0419, dates0420

	var percentage1 uint8 = 40

	urls1 := []string{"/foo", "https://example.com/foo"}
	urls2 := []string{"bar"}

	_, _, _ = percentage1, urls1, urls2

	color1 := &color.RGBA{255, 0, 0, 255}
	color2 := &color.RGBA{0, 255, 0, 255}
	color3 := pickFgColor(0, 0, 0)
	color4 := pickFgColor(1, 1, 1)
	color5 := pickFgColor(1, 2, 1)
	color6 := pickFgColor(2, 1, 1)
	color7 := pickFgColor(3, 0, 0)

	_, _, _, _, _, _, _ = color1, color2, color3, color4, color5, color6, color7

	rand.Seed(0)

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
		want   *VisualRoadmap
	}{
		{
			"empty",
			fields{0, nil, "", "", nil, nil, dates0402, dates0402, dates0402},
			&VisualRoadmap{nil, nil, nil},
		},
		{
			"complex",
			fields{
				123,
				nil,
				"02.01.2006",
				"https://example.com/",
				[]Project{
					{Title: "Initial development", Dates: &Dates{StartAt: dates0402, EndAt: dates0405}, URLs: urls1},
					{Title: "Bring website online", Milestone: 1, Color: color1},
					{Title: "Select and purchase domain", Dates: &Dates{StartAt: dates0402, EndAt: dates0415}, Indentation: 1},
					{Title: "Create server infrastructure", Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1},
					{Title: "Command line tool", Percentage: percentage1, Dates: &Dates{StartAt: dates0418, EndAt: dates0419}, Milestone: 1, Color: color2},
					{Title: "Create backend SVG generation", Indentation: 1},
					{Title: "Marketing"},
				},
				[]Milestone{
					{Title: "Milestone 0.1", URLs: urls2},
					{Title: "Milestone 0.2", DeadlineAt: &dates0420},
				},
				now,
				now,
				now,
			},
			&VisualRoadmap{
				Projects: []Project{
					{Title: "Initial development", Dates: &Dates{StartAt: dates0402, EndAt: dates0405}, URLs: urls1, Color: color3},
					{Title: "Bring website online", Dates: &Dates{StartAt: dates0402, EndAt: dates0418}, Color: color1, Milestone: 1},
					{Title: "Select and purchase domain", Dates: &Dates{StartAt: dates0402, EndAt: dates0415}, Indentation: 1, Color: color4},
					{Title: "Create server infrastructure", Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1, Color: color5},
					{Title: "Command line tool", Percentage: percentage1, Dates: &Dates{StartAt: dates0418, EndAt: dates0419}, Color: color2, Milestone: 1},
					{Title: "Create backend SVG generation", Dates: &Dates{StartAt: dates0418, EndAt: dates0419}, Indentation: 1, Color: color6},
					{Title: "Marketing", Color: color7},
				},
				Milestones: []Milestone{
					{Title: "Milestone 0.1", DeadlineAt: &dates0419, URLs: urls2, Color: color1},
					{Title: "Milestone 0.2", DeadlineAt: &dates0420, Color: &canvas.Darkgray},
				},
				Dates: &Dates{StartAt: dates0402, EndAt: dates0420},
			},
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
			if got := r.ToVisual(); !reflect.DeepEqual(got, tt.want) {
				for i := range got.Projects {
					if !reflect.DeepEqual(got.Projects[i], tt.want.Projects[i]) {
						t.Errorf("ToVisual().Projects[%d] = %v, want %v", i, got.Projects[i], tt.want.Projects[i])
					}
				}
				for i := range got.Milestones {
					if !reflect.DeepEqual(got.Milestones[i], tt.want.Milestones[i]) {
						t.Errorf("ToVisual().Milestones[%d] = %v, want %v", i, got.Milestones[i], tt.want.Milestones[i])
					}
				}
				t.Errorf("ToVisual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_calculateProjectDates(t *testing.T) {
	rand.Seed(0)

	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)
	dates0419 := time.Date(2020, 4, 19, 0, 0, 0, 0, time.UTC)

	_, _, _, _, _, _ = dates0402, dates0405, dates0408, dates0415, dates0418, dates0419

	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	tests := []struct {
		name   string
		fields fields
		want   *VisualRoadmap
	}{
		{
			"empty",
			fields{},
			&VisualRoadmap{},
		},
		{
			"find dates bottom up",
			fields{
				Projects: []Project{
					{Title: "Bring website online"},
					{Title: "Select and purchase domain", Dates: &Dates{StartAt: dates0402, EndAt: dates0415}, Indentation: 1},
					{Title: "Create server infrastructure", Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1},
				},
			},
			&VisualRoadmap{
				Projects: []Project{
					{Title: "Bring website online", Dates: &Dates{StartAt: dates0402, EndAt: dates0418}},
					{Title: "Select and purchase domain", Dates: &Dates{StartAt: dates0402, EndAt: dates0415}, Indentation: 1},
					{Title: "Create server infrastructure", Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1},
				},
			},
		},
		{
			"find dates top down",
			fields{
				Projects: []Project{
					{Title: "Command line tool", Dates: &Dates{StartAt: dates0418, EndAt: dates0419}},
					{Title: "Create backend SVG generation", Indentation: 1},
					{Title: "Create documentation page", Indentation: 1},
				},
			},
			&VisualRoadmap{
				Projects: []Project{
					{Title: "Command line tool", Dates: &Dates{StartAt: dates0418, EndAt: dates0419}},
					{Title: "Create backend SVG generation", Dates: &Dates{StartAt: dates0418, EndAt: dates0419}, Indentation: 1},
					{Title: "Create documentation page", Dates: &Dates{StartAt: dates0418, EndAt: dates0419}, Indentation: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.calculateProjectDates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateProjectDates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_findDatesBottomUp(t *testing.T) {
	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)
	dates0420 := time.Date(2020, 4, 20, 0, 0, 0, 0, time.UTC)

	_, _, _, _, _, _ = dates0402, dates0405, dates0408, dates0415, dates0418, dates0420

	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	type args struct {
		start int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Dates
	}{
		{
			"project has dates",
			fields{
				Projects: []Project{
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}},
				},
			},
			args{
				0,
			},
			&Dates{StartAt: dates0405, EndAt: dates0415},
		},
		{
			"project has no dates, but nothing can be found bottom-up",
			fields{
				Projects: []Project{
					{},
					{Indentation: 1},
					{Indentation: 2},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 1},
				},
			},
			args{
				1,
			},
			nil,
		},
		{
			"project does not have dates, but children do",
			fields{
				Projects: []Project{
					{Indentation: 1},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}, Indentation: 1},
				},
			},
			args{
				0,
			},
			&Dates{StartAt: dates0405, EndAt: dates0418},
		},
		{
			"project and children do not have dates, but grand-children do",
			fields{
				Projects: []Project{
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}},
					{Indentation: 1},
					{Indentation: 2},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 3},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 3},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}, Indentation: 1},
				},
			},
			args{
				1,
			},
			&Dates{StartAt: dates0405, EndAt: dates0418},
		},
		{
			"start 0 has dates, sub-projects are not checked",
			fields{
				Projects: []Project{
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0408}, Indentation: 1},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}, Indentation: 1},
				},
			},
			args{
				0,
			},
			&Dates{StartAt: dates0408, EndAt: dates0408},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.findDatesBottomUp(tt.args.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findDatesBottomUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_findDatesTopDown(t *testing.T) {
	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)

	_, _, _, _, _ = dates0402, dates0405, dates0408, dates0415, dates0418

	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	type args struct {
		start int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Dates
	}{
		{
			"project has a date",
			fields{
				Projects: []Project{
					{},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 1},
				},
			},
			args{
				1,
			},
			&Dates{StartAt: dates0405, EndAt: dates0415},
		},
		{
			"project does not have a date and none can be found",
			fields{
				Projects: []Project{
					{},
					{Indentation: 1},
					{Indentation: 2},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 2},
				},
			},
			args{
				1,
			},
			nil,
		},
		{
			"first parent is closest neighbor and has date",
			fields{
				Projects: []Project{
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}},
					{Indentation: 1},
				},
			},
			args{
				1,
			},
			&Dates{StartAt: dates0405, EndAt: dates0415},
		},
		{
			"first parent has date, siblings and their children are skipped",
			fields{
				Projects: []Project{
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0418}},
					{Indentation: 1},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1},
					{Dates: &Dates{StartAt: dates0415, EndAt: dates0418}, Indentation: 2},
					{Indentation: 1},
				},
			},
			args{
				4,
			},
			&Dates{StartAt: dates0405, EndAt: dates0418},
		},
		{
			"second parent has date, first parent is skipped",
			fields{
				Projects: []Project{
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0418}},
					{Indentation: 1},
					{Indentation: 2},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0415, EndAt: dates0418}, Indentation: 3},
					{Indentation: 2},
				},
			},
			args{
				5,
			},
			&Dates{StartAt: dates0405, EndAt: dates0418},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.findDatesTopDown(tt.args.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findDatesTopDown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_calculateProjectColors(t *testing.T) {
	rand.Seed(0)

	color1 := &color.RGBA{255, 0, 0, 255}
	color2 := &color.RGBA{0, 255, 0, 255}
	color3 := pickFgColor(0, 0, 0)
	color4 := pickFgColor(0, 1, 1)
	color5 := pickFgColor(0, 2, 1)

	_, _, _, _, _ = color1, color2, color3, color4, color5

	rand.Seed(0)

	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	tests := []struct {
		name   string
		fields fields
		want   *VisualRoadmap
	}{
		{
			"empty",
			fields{},
			&VisualRoadmap{},
		},
		{
			"does not overwrite existing colors",
			fields{
				Projects: []Project{
					{Title: "Bring website online", Color: color1},
					{Title: "Command line tool", Color: color2},
				},
			},
			&VisualRoadmap{
				Projects: []Project{
					{Title: "Bring website online", Color: color1},
					{Title: "Command line tool", Color: color2},
				},
			},
		},
		{
			"sets colors for all projects without a color",
			fields{
				Projects: []Project{
					{Title: "Initial development"},
					{Title: "Select and purchase domain", Indentation: 1},
					{Title: "Create server infrastructure", Indentation: 1},
				},
			},
			&VisualRoadmap{
				Projects: []Project{
					{Title: "Initial development", Color: color3},
					{Title: "Select and purchase domain", Indentation: 1, Color: color4},
					{Title: "Create server infrastructure", Indentation: 1, Color: color5},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.calculateProjectColors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateProjectColors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_calculatePercentage(t *testing.T) {
	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	tests := []struct {
		name   string
		fields fields
		want   *VisualRoadmap
	}{
		{
			"empty",
			fields{},
			&VisualRoadmap{},
		},
		{
			"does not overwrite existing percentages",
			fields{
				Projects: []Project{
					{Percentage: 35},
					{Percentage: 32, Indentation: 1},
				},
			},
			&VisualRoadmap{
				Projects: []Project{
					{Percentage: 35},
					{Percentage: 32, Indentation: 1},
				},
			},
		},
		{
			"calculates average percentage of sub-projects if missing",
			fields{
				Projects: []Project{
					{},
					{Percentage: 43, Indentation: 1},
					{Percentage: 45, Indentation: 1},
					{Percentage: 47, Indentation: 1},
					{Percentage: 49, Indentation: 1},
				},
			},
			&VisualRoadmap{
				Projects: []Project{
					{Percentage: 46},
					{Percentage: 43, Indentation: 1},
					{Percentage: 45, Indentation: 1},
					{Percentage: 47, Indentation: 1},
					{Percentage: 49, Indentation: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.calculatePercentage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculatePercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_findPercentageBottomUp(t *testing.T) {
	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	type args struct {
		start int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint8
	}{
		{
			"percentage is set",
			fields{
				Projects: []Project{
					{Percentage: 34},
				},
			},
			args{0},
			34,
		},
		{
			"percentage is not set and there are no children",
			fields{
				Projects: []Project{
					{},
					{Percentage: 34},
					{Percentage: 43, Indentation: 1},
				},
			},
			args{0},
			0,
		},
		{
			"percentage is not set average of children is used",
			fields{
				Projects: []Project{
					{},
					{Percentage: 32, Indentation: 1},
					{Percentage: 34},
					{Percentage: 43, Indentation: 1},
				},
			},
			args{0},
			32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.findPercentageBottomUp(tt.args.start); got != tt.want {
				t.Errorf("findPercentageBottomUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_collectProjectMilestones(t *testing.T) {
	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)

	_, _, _, _, _ = dates0402, dates0405, dates0408, dates0415, dates0418

	color1 := &color.RGBA{R: 34, G: 23, B: 73, A: 255}
	color2 := &color.RGBA{R: 53, G: 82, B: 19, A: 255}

	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	tests := []struct {
		name   string
		fields fields
		want   map[int]*Milestone
	}{
		{
			"empty",
			fields{},
			map[int]*Milestone{},
		},
		{
			"project without colors and dates are skipped",
			fields{
				Projects: []Project{
					{Milestone: 1},
				},
			},
			map[int]*Milestone{},
		},
		{
			"project with milestones are found",
			fields{
				Projects: []Project{
					{Milestone: 1},
					{Milestone: 1, Dates: &Dates{StartAt: dates0402, EndAt: dates0418}},
				},
			},
			map[int]*Milestone{
				0: {DeadlineAt: &dates0418},
			},
		},
		{
			"latest project is used for a given milestone",
			fields{
				Projects: []Project{
					{Milestone: 1},
					{Milestone: 1, Dates: &Dates{StartAt: dates0408, EndAt: dates0415}},
					{Milestone: 1, Dates: &Dates{StartAt: dates0415, EndAt: dates0418}},
					{Milestone: 1, Dates: &Dates{StartAt: dates0402, EndAt: dates0415}},
				},
			},
			map[int]*Milestone{
				0: {DeadlineAt: &dates0418},
			},
		},
		{
			"projects with colors are not skipped",
			fields{
				Projects: []Project{
					{Milestone: 1},
					{Milestone: 1, Dates: &Dates{StartAt: dates0402, EndAt: dates0415}},
					{Milestone: 2, Color: color1},
					{Milestone: 1, Dates: &Dates{StartAt: dates0402, EndAt: dates0418}},
				},
			},
			map[int]*Milestone{
				0: {DeadlineAt: &dates0418},
				1: {Color: color1},
			},
		},
		{
			"color of the first found project is used",
			fields{
				Projects: []Project{
					{Milestone: 1},
					{Milestone: 1, Dates: &Dates{StartAt: dates0402, EndAt: dates0415}, Color: color2, Indentation: 1},
					{Milestone: 2, Color: color1},
					{Milestone: 2, Color: color2, Indentation: 1},
				},
			},
			map[int]*Milestone{
				0: {DeadlineAt: &dates0415, Color: color2},
				1: {Color: color1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.collectProjectMilestones(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectProjectMilestones() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisualRoadmap_applyProjectMilestone(t *testing.T) {
	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)

	_, _, _, _, _ = dates0402, dates0405, dates0408, dates0415, dates0418

	rand.Seed(0)

	color1 := &color.RGBA{255, 0, 0, 255}
	color2 := &color.RGBA{0, 255, 0, 255}
	color3 := pickFgColor(0, 0, 0)
	color4 := pickFgColor(1, 1, 1)
	color5 := pickFgColor(1, 2, 1)
	color6 := pickFgColor(2, 1, 1)
	color7 := pickFgColor(3, 0, 0)

	_, _, _, _, _, _, _ = color1, color2, color3, color4, color5, color6, color7

	rand.Seed(0)

	type fields struct {
		Projects   []Project
		Milestones []Milestone
		Dates      *Dates
	}
	type args struct {
		projectMilestones map[int]*Milestone
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *VisualRoadmap
	}{
		{
			"empty",
			fields{},
			args{},
			&VisualRoadmap{},
		},
		{
			"default milestone color is used when nothing is found",
			fields{
				Milestones: []Milestone{
					{},
				},
			},
			args{},
			&VisualRoadmap{
				Milestones: []Milestone{
					{Color: defaultMilestoneColor},
				},
			},
		},
		{
			"milestone color can be set by project linked to milestone",
			fields{
				Milestones: []Milestone{
					{},
					{},
				},
			},
			args{
				projectMilestones: map[int]*Milestone{
					0: {Color: color1},
				},
			},
			&VisualRoadmap{
				Milestones: []Milestone{
					{Color: color1},
					{Color: defaultMilestoneColor},
				},
			},
		},
		{
			"project color will not override the milestone color",
			fields{
				Milestones: []Milestone{
					{Color: color1},
					{},
				},
			},
			args{
				projectMilestones: map[int]*Milestone{
					0: {Color: color2},
				},
			},
			&VisualRoadmap{
				Milestones: []Milestone{
					{Color: color1},
					{Color: defaultMilestoneColor},
				},
			},
		},
		{
			"milestone deadline can be set by project linked to milestone",
			fields{
				Milestones: []Milestone{
					{},
					{},
				},
			},
			args{
				projectMilestones: map[int]*Milestone{
					0: {DeadlineAt: &dates0415},
				},
			},
			&VisualRoadmap{
				Milestones: []Milestone{
					{DeadlineAt: &dates0415, Color: defaultMilestoneColor},
					{Color: defaultMilestoneColor},
				},
			},
		},
		{
			"project deadline will not override the milestone deadline",
			fields{
				Milestones: []Milestone{
					{DeadlineAt: &dates0415},
				},
			},
			args{
				projectMilestones: map[int]*Milestone{
					0: {DeadlineAt: &dates0418},
				},
			},
			&VisualRoadmap{
				Milestones: []Milestone{
					{DeadlineAt: &dates0415, Color: defaultMilestoneColor},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := &VisualRoadmap{
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				Dates:      tt.fields.Dates,
			}
			if got := vr.applyProjectMilestone(tt.args.projectMilestones); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyProjectMilestone() = %v, want %v", got, tt.want)
			}
		})
	}
}
