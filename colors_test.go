package main

import (
	"fmt"
	"image/color"
	"reflect"
	"testing"
)

func Test_pickFgColor(t *testing.T) {
	type args struct {
		epicCount   int
		taskCount   int
		indentation int
	}
	tests := []struct {
		name string
		args args
		want *color.RGBA
	}{
		{
			"first epic",
			args{
				epicCount:   0,
				taskCount:   0,
				indentation: 0,
			},
			&color.RGBA{R: 0, G: 166, B: 237, A: 255},
		},
		{
			"first sub-project of first epic",
			args{
				epicCount:   0,
				taskCount:   1,
				indentation: 1,
			},
			&color.RGBA{R: 0, G: 136, B: 194, A: 255},
		},
		{
			"second sub-project of first epic",
			args{
				epicCount:   0,
				taskCount:   2,
				indentation: 1,
			},
			&color.RGBA{R: 0, G: 106, B: 151, A: 255},
		},
		{
			"first sub-sub-project of first epic",
			args{
				epicCount:   0,
				taskCount:   2,
				indentation: 2,
			},
			&color.RGBA{R: 139, G: 214, B: 246, A: 255},
		},
		{
			"second epic",
			args{
				epicCount:   1,
				taskCount:   0,
				indentation: 0,
			},
			&color.RGBA{R: 255, G: 0, B: 114, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pickFgColor(tt.args.epicCount, tt.args.taskCount, tt.args.indentation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pickFgColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pickBgColor(t *testing.T) {
	type args struct {
		epic int
	}
	tests := []struct {
		name string
		args args
		want color.RGBA
	}{
		{
			"first background color",
			args{epic: 0},
			colors[0][len(colors[0])-1],
		},
		{
			"second background color",
			args{epic: 1},
			colors[1][len(colors[1])-1],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pickBgColor(tt.args.epic); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pickBgColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mustParseColor(t *testing.T) {
	type args struct {
		part string
	}
	tests := []struct {
		name string
		args args
		want color.RGBA
	}{
		{
			"#959595",
			args{part: "#959595"},
			color.RGBA{R: 149, G: 149, B: 149, A: 255},
		},
		{
			"#395273",
			args{part: "#395273"},
			color.RGBA{R: 57, G: 82, B: 115, A: 255},
		},
		{
			"#a4f",
			args{part: "#a4f"},
			color.RGBA{R: 170, G: 68, B: 255, A: 255},
		},
		{
			"#A4F",
			args{part: "#A4F"},
			color.RGBA{R: 170, G: 68, B: 255, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustParseColor(tt.args.part); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mustParseColor() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("panic on invalid length", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("the code did not panic")
			}
		}()

		mustParseColor("")
	})

	t.Run("panic on invalid first character", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("the code did not panic")
			}
		}()

		mustParseColor("A4FA")
	})

	t.Run("panic on invalid hexadecimal digit", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("the code did not panic")
			}
		}()

		mustParseColor("#00g")
	})
}

func Test_colorsAreUnique(t *testing.T) {
	var colorsFound = map[string]string{}

	for i := range colors {
		cg := colors[i]

		for j := range cg {
			c := fmt.Sprintf("%d, %d, %d", cg[j].R, cg[j].G, cg[j].B)
			w := fmt.Sprintf("%d / %d", i, j)

			if _, ok := colorsFound[c]; ok {
				t.Errorf("duplicate found: %s @ (%s)", c, w)
			}

			colorsFound[c] = w
		}
	}
}
