package herr

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpError(t *testing.T) {
	type args struct {
		err    error
		status int
	}
	tests := []struct {
		name string
		args args
		want HttpError
	}{
		{
			"foo",
			args{assert.AnError, 300},
			HttpError{assert.AnError, 300},
		},
		{
			"HttpError",
			args{HttpError{assert.AnError, 300}, 333},
			HttpError{assert.AnError, 300},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromError(tt.args.err, tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromError() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("panic on nil error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("the code did not panic")
			}
		}()

		var err error
		_ = NewFromError(err, 302)
	})
}

func TestToHttpCode(t *testing.T) {
	type args struct {
		err           error
		defaultStatus int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"nil",
			args{err: nil, defaultStatus: http.StatusTeapot},
			http.StatusTeapot,
		},
		{
			"http code error",
			args{err: HttpError{error: assert.AnError, status: http.StatusTeapot}, defaultStatus: http.StatusBadGateway},
			http.StatusTeapot,
		},
		{
			"not http code error",
			args{err: assert.AnError, defaultStatus: http.StatusBadGateway},
			http.StatusBadGateway,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToHttpCode(tt.args.err, tt.args.defaultStatus); got != tt.want {
				t.Errorf("ToHttpCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromString(t *testing.T) {
	type args struct {
		msg    string
		status int
	}
	tests := []struct {
		name string
		args args
		want HttpError
	}{
		{
			"default",
			args{"foo", 322},
			HttpError{fmt.Errorf("foo"), 322},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromString(tt.args.msg, tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
