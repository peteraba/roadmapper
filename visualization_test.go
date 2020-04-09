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

	var percentage1 uint8 = 40

	urls1 := []string{"/foo", "https://example.com/foo"}
	urls2 := []string{"bar"}

	color1 := &color.RGBA{255, 0, 0, 255}
	color2 := &color.RGBA{0, 255, 0, 255}
	color3 := &color.RGBA{34, 9, 1, 255}
	color4 := &color.RGBA{148, 27, 12, 255}
	color5 := &color.RGBA{148, 27, 12, 255}
	color6 := &color.RGBA{224, 155, 26, 255}

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
					{Title: "Marketing", Color: color6},
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

func TestVisualRoadmap_findDatesBottomUp(t *testing.T) {
	dates0402 := time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC)
	dates0405 := time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC)
	dates0408 := time.Date(2020, 4, 8, 0, 0, 0, 0, time.UTC)
	dates0415 := time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC)
	dates0418 := time.Date(2020, 4, 18, 0, 0, 0, 0, time.UTC)
	dates0420 := time.Date(2020, 4, 20, 0, 0, 0, 0, time.UTC)

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
			"start 0 has dates",
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
			"start 0 does not have dates",
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
			"start 0 does dates, sub-projects are not checked",
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
		{
			"start 2 does not have dates",
			fields{
				Projects: []Project{
					{},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}, Indentation: 1},
					{},
					{Dates: &Dates{StartAt: dates0405, EndAt: dates0415}, Indentation: 1},
					{Dates: &Dates{StartAt: dates0408, EndAt: dates0418}, Indentation: 1},
					{Dates: &Dates{StartAt: dates0402, EndAt: dates0420}},
				},
			},
			args{
				2,
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
			if got := vr.findDatesBottomUp(tt.args.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findDatesBottomUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
