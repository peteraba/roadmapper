package colors

import (
	"fmt"
	"image/color"
	"reflect"
	"testing"
)

func Test_PickFgColor(t *testing.T) {
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
			&colors[0][10],
		},
		{
			"first sub-project of first epic",
			args{
				epicCount:   0,
				taskCount:   1,
				indentation: 1,
			},
			&colors[0][8],
		},
		{
			"second sub-project of first epic",
			args{
				epicCount:   0,
				taskCount:   2,
				indentation: 1,
			},
			&colors[0][6],
		},
		{
			"6th sub-project of first epic",
			args{
				epicCount:   0,
				taskCount:   6,
				indentation: 1,
			},
			&colors[0][18],
		},
		{
			"first sub-sub-project of second epic",
			args{
				epicCount:   0,
				taskCount:   2,
				indentation: 2,
			},
			&colors[0][11],
		},
		{
			"8th sub-sub-project of second epic",
			args{
				epicCount:   0,
				taskCount:   9,
				indentation: 2,
			},
			&colors[0][17],
		},
		{
			"second epic",
			args{
				epicCount:   1,
				taskCount:   0,
				indentation: 0,
			},
			&colors[1][10],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PickFgColor(tt.args.epicCount, tt.args.taskCount, tt.args.indentation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pickFgColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_PickBgColor(t *testing.T) {
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
			if got := PickBgColor(tt.args.epic); !reflect.DeepEqual(got, tt.want) {
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

func Test_CharsToUint8(t *testing.T) {
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
			got, err := CharsToUint8(tt.args.part)
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

func TestToHexa(t *testing.T) {
	type args struct {
		c color.Color
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"white",
			args{color.RGBA{R: 255, G: 255, B: 255, A: 255}},
			"#ffffff",
		},
		{
			"red",
			args{color.RGBA{R: 255, G: 0, B: 0, A: 255}},
			"#ff0000",
		},
		{
			"black",
			args{color.RGBA{R: 0, G: 0, B: 0, A: 255}},
			"#000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToHexa(tt.args.c); got != tt.want {
				t.Errorf("ToHexa() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_twoDigitHexa(t *testing.T) {
	type args struct {
		i uint32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"zero",
			args{0},
			"00",
		},
		{
			"fifteen",
			args{15},
			"0f",
		},
		{
			"sixteen",
			args{16},
			"10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := twoDigitHexa(tt.args.i); got != tt.want {
				t.Errorf("twoDigitHexa() = %v, want %v", got, tt.want)
			}
		})
	}
}
