package herr

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
