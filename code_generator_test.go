package main

import (
	"testing"
)

func TestCode64_String(t *testing.T) {
	tests := []struct {
		name string
		c    Code64
		want string
	}{
		{
			name: "zero",
			c:    Code64(0),
			want: "",
		},
		{
			name: "one",
			c:    Code64(1),
			want: "1",
		},
		{
			name: "two (0b10)",
			c:    Code64(2),
			want: "2",
		},
		{
			name: "0xe",
			c:    Code64(0xe),
			want: "e",
		},
		{
			name: "0xff",
			c:    Code64(0xff),
			want: "3~",
		},
		{
			name: "0xffff",
			c:    Code64(0xffff),
			want: "f~~",
		},
		{
			name: "0xffffff",
			c:    Code64(0xffffff),
			want: "~~~~",
		},
		{
			name: "0xffffffff",
			c:    Code64(0xffffffff),
			want: "3~~~~~",
		},
		{
			name: "0xffffffffff",
			c:    Code64(0xffffffffff),
			want: "f~~~~~~",
		},
		{
			name: "maxCode64",
			c:    Code64(maxCode64),
			want: "3~~~~~~~",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("panic on out of bound (lower)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		_ = Code64(-1).String()
	})

	t.Run("panic on out of bound (upper)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		_ = Code64(maxCode64 + 1).String()
	})
}

func TestNewCode64(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		if got := NewCode64(); int64(got) < 0 || int64(got) > maxCode64 {
			t.Errorf("NewCode64() is invalid: %v", got)
		}
	})
}

func Test_toCode64(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "zero",
			args: args{n: 0},
			want: "",
		},
		{
			name: "one",
			args: args{n: 1},
			want: "1",
		},
		{
			name: "two (0b10)",
			args: args{n: 2},
			want: "2",
		},
		{
			name: "0xe",
			args: args{n: 0xe},
			want: "e",
		},
		{
			name: "0xff",
			args: args{n: 0xff},
			want: "3~",
		},
		{
			name: "0xffff",
			args: args{n: 0xffff},
			want: "f~~",
		},
		{
			name: "0xffffff",
			args: args{n: 0xffffff},
			want: "~~~~",
		},
		{
			name: "0xffffffff",
			args: args{n: 0xffffffff},
			want: "3~~~~~",
		},
		{
			name: "0xffffffffff",
			args: args{n: 0xffffffffff},
			want: "f~~~~~~",
		},
		{
			name: "0xffffffffffff",
			args: args{n: 0xffffffffffff},
			want: "3~~~~~~~",
		},
		{
			name: "0x1ffffffffffff (overflow)",
			args: args{n: 0x1ffffffffffff},
			want: "3~~~~~~~",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toCode64(tt.args.n); got != tt.want {
				t.Errorf("toCode64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCode64FromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "zero",
			args: args{s: ""},
			want: 0,
		},
		{
			name: "nine",
			args: args{s: "9"},
			want: 9,
		},
		{
			name: "A",
			args: args{s: "A"},
			want: 36,
		},
		{
			name: "~",
			args: args{s: "~"},
			want: 63,
		},
		{
			name: "~A",
			args: args{s: "~A"},
			want: 63*64 + 36,
		},
		{
			name: "~~A",
			args: args{s: "~~A"},
			want: 63*64*64 + 63*64 + 36,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCode64FromString(tt.args.s)
			if err != nil {
				t.Errorf("NewCode64FromString() error = %v, wantErr %v", err, false)
				return
			}
			if int64(got) != tt.want {
				t.Errorf("NewCode64FromString() got = %v, want %v", int64(got), tt.want)
			}
		})
	}

	t.Run("error on out of bound", func(t *testing.T) {
		got, err := NewCode64FromString("abcdefghijklmnopq")
		if err == nil {
			t.Errorf("NewCode64FromString() error = %v, wantErr %v", err, true)
			return
		}
		if got > 0 {
			t.Errorf("NewCode64FromString() got = %v, want %v", int64(got), 0)
		}
	})

	t.Run("error on invalid character", func(t *testing.T) {
		got, err := NewCode64FromString("世界")
		if err == nil {
			t.Errorf("NewCode64FromString() error = %v, wantErr %v", err, true)
			return
		}
		if got > 0 {
			t.Errorf("NewCode64FromString() got = %v, want %v", int64(got), 0)
		}
	})
}
