package main

import (
	"image/color"
	"reflect"
	"testing"
	"time"

	"github.com/peteraba/go-svg"
)

func TestRoadmap_ToVisual(t *testing.T) {
	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)
	dates0419 := time.Date(2020, 4, 19, 0, 0, 0, 0, time.UTC)
	dates0420 := time.Date(2020, 4, 20, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	var percentage1 uint8 = 40
	urls1 := []string{"/foo", "https://example.com/foo"}
	urls2 := []string{"bar"}
	color1 := color.RGBA{255, 0, 0, 255}
	color2 := color.RGBA{0, 255, 0, 255}

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
		want   VisualRoadmap
	}{
		{
			"empty",
			fields{0, nil, "", "", nil, nil, dates0402, dates0402, dates0402},
			VisualRoadmap{nil, nil, nil},
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
			VisualRoadmap{
				Projects: []Project{
					{Title: "Initial development", Dates: &Dates{StartAt: dates0402, EndAt: dates0405}, URLs: urls1},
					{Title: "Bring website online", Dates: &Dates{StartAt: dates0402, EndAt: dates0418}, Color: color1},
					{Title: "Select and purchase domain", Dates: &Dates{StartAt: dates0402, EndAt: dates0415}, Indentation: 1},
					{Title: "Create server infrastructure", Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1},
					{Title: "Command line tool", Percentage: percentage1, Dates: &Dates{StartAt: dates0418, EndAt: dates0419}, Color: color2},
					{Title: "Marketing"},
				},
				Milestones: []Milestone{
					{Title: "Milestone 0.1", DeadlineAt: &dates0419, URLs: urls2, Color: color1},
					{Title: "Milestone 0.2", DeadlineAt: &dates0420},
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
				t.Errorf("ToVisual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findVisualDates(t *testing.T) {
	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)
	dates0420 := time.Date(2020, 4, 20, 0, 0, 0, 0, time.UTC)

	type args struct {
		projects []Project
		start    int
	}
	tests := []struct {
		name string
		args args
		want *Dates
	}{
		{
			"start 0 has dates",
			args{
				[]Project{
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}},
				},
				0,
			},
			&Dates{StartAt: dates0405, EndAt: dates0415},
		},
		{
			"start 0 does not have dates",
			args{
				[]Project{
					{Indentation: 1},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}, Indentation: 1},
				},
				0,
			},
			&Dates{StartAt: dates0405, EndAt: dates0418},
		},
		{
			"start 0 does dates, sub-projects are not checked",
			args{
				[]Project{
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0408}, Indentation: 1},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 2},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}, Indentation: 1},
				},
				0,
			},
			&Dates{StartAt: dates0408, EndAt: dates0408},
		},
		{
			"start 2 does not have dates",
			args{
				[]Project{
					{},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}, Indentation: 1},
					{},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 1},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}},
				},
				2,
			},
			&Dates{StartAt: dates0405, EndAt: dates0418},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findVisualDates(tt.args.projects, tt.args.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findVisualDates() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("panic on empty project list", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		_ = findVisualDates(nil, 1)
	})
}

func Test_createSvg(t *testing.T) {
	type args struct {
		roadmap      *Roadmap
		fullWidth    float64
		headerHeight float64
		lineHeight   float64
	}
	tests := []struct {
		name string
		args args
		want svg.SVG
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSvg(tt.args.roadmap, tt.args.fullWidth, tt.args.headerHeight, tt.args.lineHeight); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createSvg() = %v, want %v", got, tt.want)
			}
		})
	}
}
