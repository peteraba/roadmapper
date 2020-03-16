package main

import "testing"

func Test_bootstrapRoadmap(t *testing.T) {
	type args struct {
		roadmap      Project
		lines        []string
		matomoDomain string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bootstrapRoadmap(tt.args.roadmap, tt.args.lines, tt.args.matomoDomain)
			if (err != nil) != tt.wantErr {
				t.Errorf("bootstrapRoadmap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bootstrapRoadmap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
