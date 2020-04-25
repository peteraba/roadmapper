package roadmap

import (
	"reflect"
	"testing"
	"time"
)

func TestRoadmap_viewHtml(t *testing.T) {
	type fields struct {
		ID         uint64
		PrevID     *uint64
		Title      string
		DateFormat string
		BaseURL    string
		Projects   []Project
		Milestones []Milestone
		CreatedAt  time.Time
		UpdatedAt  time.Time
		AccessedAt time.Time
	}
	type args struct {
		appVersion   string
		matomoDomain string
		docBaseURL   string
		currentURL   string
		selfHosted   bool
		err          error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Roadmap{
				ID:         tt.fields.ID,
				PrevID:     tt.fields.PrevID,
				Title:      tt.fields.Title,
				DateFormat: tt.fields.DateFormat,
				BaseURL:    tt.fields.BaseURL,
				Projects:   tt.fields.Projects,
				Milestones: tt.fields.Milestones,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
				AccessedAt: tt.fields.AccessedAt,
			}
			got, err := r.viewHtml(tt.args.appVersion, tt.args.matomoDomain, tt.args.docBaseURL, tt.args.currentURL, tt.args.selfHosted, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("viewHtml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("viewHtml() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoadmap_getProjectURLs(t *testing.T) {
	type fields struct {
		r *Roadmap
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]string
	}{
		{
			"empty for nil roadmap",
			fields{},
			map[string][]string{},
		},
		{
			"empty for nil projects",
			fields{
				r: &Roadmap{Projects: nil},
			},
			map[string][]string{},
		},
		{
			"urls are collected by title",
			fields{
				r: &Roadmap{Projects: []Project{
					{Title: "foo", URLs: []string{"foo-1", "foo-2"}},
					{Title: "bar", URLs: []string{"bar-1"}},
					{Title: "baz", URLs: []string{"baz-1", "baz-2"}},
				}},
			},
			map[string][]string{
				"foo": {"foo-1", "foo-2"},
				"bar": {"bar-1"},
				"baz": {"baz-1", "baz-2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.r.getProjectURLs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getProjectURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}
