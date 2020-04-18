package roadmap

import "testing"

func Test_bootstrapRoadmap(t *testing.T) {
	type args struct {
		roadmap      *Roadmap
		appVersion   string
		matomoDomain string
		docBaseURL   string
		currentURL   string
		selfHosted   bool
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
			got, err := bootstrapRoadmap(tt.args.roadmap, tt.args.appVersion, tt.args.matomoDomain, tt.args.docBaseURL, tt.args.currentURL, tt.args.selfHosted)
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
