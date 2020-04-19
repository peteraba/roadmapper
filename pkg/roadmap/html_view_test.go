package roadmap

import (
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
			got, err := r.viewHtml(tt.args.appVersion, tt.args.matomoDomain, tt.args.docBaseURL, tt.args.currentURL, tt.args.selfHosted)
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
