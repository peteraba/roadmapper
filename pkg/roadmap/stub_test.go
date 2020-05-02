package roadmap

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
)

func TestNewRoadmapExchangeStub(t *testing.T) {
	type args struct {
		minProjects   int
		minMilestones int
		minDate       time.Time
		maxDate       time.Time
	}
	tests := []struct {
		name               string
		seed               int64
		before             int
		args               args
		want               RoadmapExchange
		wantProjectCount   int
		wantMilestoneCount int
	}{
		{
			"first after seeding 1",
			1,
			0,
			args{
				minProjects:   2,
				minMilestones: 2,
				minDate:       time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate:       time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			RoadmapExchange{
				Title:      "readymade pabst listicle",
				DateFormat: "2006-01-02",
			},
			16,
			4,
		},
		{
			"second after seeding 1",
			1,
			1,
			args{
				minProjects:   2,
				minMilestones: 2,
				minDate:       time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate:       time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			RoadmapExchange{
				Title:      "mixtape",
				DateFormat: "2006-01-02",
			},
			19,
			19,
		},
		{
			"first after seeding 67",
			67,
			0,
			args{
				minProjects:   2,
				minMilestones: 2,
				minDate:       time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate:       time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			RoadmapExchange{
				Title:      "roof",
				DateFormat: "2006-01-02",
				BaseURL:    "http://www.leadrelationships.com/best-of-breed/collaborative/best-of-breed",
			},
			20,
			7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i := 0; i < tt.before; i++ {
				NewRoadmapExchangeStub(tt.args.minProjects, tt.args.minMilestones, tt.args.minDate, tt.args.maxDate)
			}

			got := NewRoadmapExchangeStub(tt.args.minProjects, tt.args.minMilestones, tt.args.minDate, tt.args.maxDate)

			assert.Len(t, got.Projects, tt.wantProjectCount, "wrong project count")
			assert.Len(t, got.Milestones, tt.wantMilestoneCount, "wrong milestone count")

			got.Projects = nil
			got.Milestones = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoadmapExchangeStub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_max(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"a smaller",
			args{1, 3},
			3,
		},
		{
			"a larger",
			args{3, 1},
			3,
		},
		{
			"equals",
			args{2, 2},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newBaseURL(t *testing.T) {
	tests := []struct {
		name string
		seed int64
		want []string
	}{
		{
			"after seeding 1",
			1,
			[]string{"", "", "", "", "",
				"http://www.dynamicdrive.biz/evolve/infomediaries/24/7",
				"https://www.dynamiccross-media.com/cross-media/end-to-end/global",
				"http://www.humandistributed.io/expedite/scalable/best-of-breed"},
		},
		{
			"after seeding 15646",
			15646,
			[]string{"", "",
				"http://www.productproductize.org/mesh/whiteboard/bandwidth",
				"", "",
				"http://www.globalharness.info/robust/enterprise",
				"",
				"https://www.leadscalable.name/e-services/e-business"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i, w := range tt.want {
				if got := newBaseURL(); got != w {
					t.Errorf("@ #%d newBaseURL() = %v, want %v", i, got, w)
				}
			}
		})
	}
}

func TestNewMilestoneStub(t *testing.T) {
	type args struct {
		minDate time.Time
		maxDate time.Time
		hasBU   bool
	}
	tests := []struct {
		name   string
		seed   int64
		before int
		args   args
		want   Milestone
	}{
		{
			"first after seeding 1",
			1,
			0,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:   false,
			},
			Milestone{
				Title: "schlitz wolf",
				URLs:  []string{"https://www.dynamicinnovate.net/infomediaries/24/7/cutting-edge/e-enable"},
			},
		},
		{
			"first after seeding 1 with base url",
			1,
			0,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:   true,
			},
			Milestone{
				Title: "schlitz wolf",
				URLs:  []string{"https://www.dynamicinnovate.net/infomediaries/24/7/cutting-edge/e-enable", "eum", "totam"},
			},
		},
		{
			"second after seeding 1",
			1,
			1,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:   false,
			},
			Milestone{
				Title: "humblebrag stumptown",
				URLs:  []string{"http://www.humandistributed.io/expedite/scalable/best-of-breed", "http://www.productfrictionless.com/engineer/initiatives/applications/portals"},
			},
		},
		{
			"first after seeding 76",
			76,
			0,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:   false,
			},
			Milestone{
				Title: "meditation messenger bag",
				URLs:  []string{"http://www.seniordot-com.net/bleeding-edge/transform/platforms/one-to-one"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i := 0; i < tt.before; i++ {
				NewMilestoneStub(tt.args.minDate, tt.args.maxDate, tt.args.hasBU)
			}
			if got := NewMilestoneStub(tt.args.minDate, tt.args.maxDate, tt.args.hasBU); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMilestoneStub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProjectStub(t *testing.T) {
	type args struct {
		milestoneCount int
		ind            int
		minDate        time.Time
		maxDate        time.Time
		hasBU          bool
	}
	tests := []struct {
		name   string
		seed   int64
		before int
		args   args
		want   Project
	}{
		{
			"first after seeding 1",
			1,
			0,
			args{
				milestoneCount: 0,
				ind:            2,
				minDate:        time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate:        time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:          false,
			},
			Project{
				Indentation: 2,
				Title:       "wolf taxidermy",
				Percentage:  82,
				URLs:        []string{"https://www.dynamicfunctionalities.io/24/7/cutting-edge/e-enable"},
			},
		},
		{
			"first after seeding 1 with base url",
			1,
			0,
			args{
				milestoneCount: 0,
				ind:            2,
				minDate:        time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate:        time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:          true,
			},
			Project{
				Indentation: 2,
				Title:       "wolf taxidermy",
				Percentage:  82,
				URLs:        []string{"https://www.dynamicfunctionalities.io/24/7/cutting-edge/e-enable", "eum", "totam"},
			},
		},
		{
			"second after seeding 1 with base url",
			1,
			1,
			args{
				milestoneCount: 0,
				ind:            2,
				minDate:        time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate:        time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:          true,
			},
			Project{
				Indentation: 2,
				Title:       "neutra",
				Dates: &Dates{
					StartAt: time.Date(2020, 1, 1, 12, 42, 6, 821791281, time.UTC),
					EndAt:   time.Date(2020, 1, 1, 14, 16, 46, 183294663, time.UTC),
				},
				Percentage: 14,
				URLs:       []string{"accusantium", "consectetur"},
			},
		},
		{
			"first after seeding 76",
			76,
			0,
			args{
				milestoneCount: 0,
				ind:            2,
				minDate:        time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate:        time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
				hasBU:          false,
			},
			Project{
				Indentation: 2,
				Title:       "messenger bag truffaut",
				Percentage:  43,
				URLs:        []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i := 0; i < tt.before; i++ {
				NewProjectStub(tt.args.milestoneCount, tt.args.ind, tt.args.minDate, tt.args.maxDate, tt.args.hasBU)
			}
			if got := NewProjectStub(tt.args.milestoneCount, tt.args.ind, tt.args.minDate, tt.args.maxDate, tt.args.hasBU); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProjectStub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextIndentation(t *testing.T) {
	type args struct {
		indentation int
	}
	tests := []struct {
		name   string
		seed   int64
		before int
		args   args
		want   int
	}{
		{
			"first after 0 when seeding 1",
			1,
			0,
			args{0},
			0,
		},
		{
			"sixth after 0 when seeding 1",
			1,
			5,
			args{0},
			1,
		},
		{
			"sixth after 2 when seeding 1",
			1,
			5,
			args{2},
			1,
		},
	}
	for _, tt := range tests {
		gofakeit.Seed(tt.seed)
		for i := 0; i < tt.before; i++ {
			nextIndentation(tt.args.indentation)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := nextIndentation(tt.args.indentation); got != tt.want {
				t.Errorf("nextIndentation() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("0 <= p <= p0", func(t *testing.T) {
		gofakeit.Seed(0)
		for i := 0; i < 20; i++ {
			ind := rand.Int()
			got := nextIndentation(ind)

			assert.LessOrEqual(t, got, ind+1)
			assert.GreaterOrEqual(t, got, 0)
		}
	})
}

func Test_newWords(t *testing.T) {
	tests := []struct {
		name   string
		seed   int64
		before int
		want   string
	}{
		{
			"first after seeding 1",
			1,
			0,
			"schlitz wolf",
		},
		{
			"second after seeding 1",
			1,
			1,
			"neutra",
		},
		{
			"first after seeding 67",
			67,
			1,
			"gluten-free kombucha",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i := 0; i < tt.before; i++ {
				newWords()
			}
			if got := newWords(); got != tt.want {
				t.Errorf("newWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newDates(t *testing.T) {
	var (
		d0   = time.Date(2020, 1, 1, 15, 10, 47, 632952942, time.UTC)
		d1   = time.Date(2020, 1, 1, 19, 48, 37, 331736048, time.UTC)
		res0 = Dates{StartAt: d0, EndAt: d1}
	)
	type args struct {
		minDate time.Time
		maxDate time.Time
	}
	tests := []struct {
		name   string
		seed   int64
		before int
		args   args
		want   *Dates
	}{
		{
			"first after seeding 1",
			1,
			0,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"sixth after seeding 1",
			1,
			5,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			&res0,
		},
		{
			"first after seeding 87",
			87,
			0,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i := 0; i < tt.before; i++ {
				newDates(tt.args.minDate, tt.args.maxDate)
			}
			if got := newDates(tt.args.minDate, tt.args.maxDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newURLs(t *testing.T) {
	type args struct {
		hasBU bool
	}
	tests := []struct {
		name   string
		before int
		seed   int64
		args   args
		want   []string
	}{
		{
			"first after seeding 1",
			0,
			1,
			args{false},
			[]string{"http://www.futureinfrastructures.io/infrastructures/channels/drive"},
		},
		{
			"first after seeding 1 with base url",
			0,
			1,
			args{true},
			[]string{"http://www.futureinfrastructures.io/infrastructures/channels/drive", "omnis"},
		},
		{
			"second after seeding 1",
			1,
			1,
			args{false},
			[]string{"http://www.global24/7.biz/disintermediate/cross-media/one-to-one", "http://www.centralglobal.io/scale/distributed/turn-key/transition"},
		},
		{
			"second after seeding 1 with base url",
			1,
			1,
			args{true},
			[]string{"https://www.dynamicdisintermediate.io/end-to-end/cross-media", "doloremque"},
		},
		{
			"second after seeding 666",
			1,
			666,
			args{false},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i := 0; i < tt.before; i++ {
				newURLs(tt.args.hasBU)
			}
			if got := newURLs(tt.args.hasBU); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newURLs(%t) = %v, want %v", tt.args.hasBU, got, tt.want)
			}
		})
	}
}

func Test_newDateOptional(t *testing.T) {
	var (
		res0 = time.Date(2020, 1, 1, 19, 13, 43, 81912589, time.UTC)
		res1 = time.Date(2020, 1, 1, 18, 40, 13, 980046881, time.UTC)
	)
	type args struct {
		minDate time.Time
		maxDate time.Time
	}
	tests := []struct {
		name   string
		seed   int64
		before int
		args   args
		want   *time.Time
	}{
		{
			"first after seeding 1",
			1,
			0,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			&res0,
		},
		{
			"sixth after seeding 1 is null",
			1,
			5,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"first after seeding 76",
			76,
			0,
			args{
				minDate: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
				maxDate: time.Date(2020, 1, 1, 20, 0, 0, 0, time.UTC),
			},
			&res1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gofakeit.Seed(tt.seed)
			for i := 0; i < tt.before; i++ {
				newDateOptional(tt.args.minDate, tt.args.maxDate)
			}
			if got := newDateOptional(tt.args.minDate, tt.args.maxDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDateOptional() = %v, want %v", got, tt.want)
			}
		})
	}
}
