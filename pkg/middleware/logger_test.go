package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func setupLogger() (*zap.Logger, *zaptest.Buffer) {
	buf := &zaptest.Buffer{}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
		buf,
		zap.DebugLevel,
	))

	return logger, buf
}

func TestLoggerConfig_log(t *testing.T) {
	var (
		foo         = "foo"
		h0          = http.Header{}
		logger, buf = setupLogger()
		lc          = DefaultLoggerConfig
	)

	lc.Logger = logger

	h0.Add(echo.HeaderXRequestID, foo)
	h0.Add(echo.HeaderXRequestID, "")

	rec := httptest.NewRecorder()
	rec.HeaderMap.Add(echo.HeaderXRequestID, foo) // nolint
	res := &echo.Response{Writer: rec}

	type args struct {
		req   *http.Request
		res   *echo.Response
		start time.Time
		stop  time.Time
		ip    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"default",
			args{
				req: &http.Request{
					Header: h0,
				},
				res:   res,
				start: time.Date(2020, 5, 19, 9, 0, 0, 0, time.UTC),
				stop:  time.Date(2020, 5, 19, 9, 0, 0, 0, time.UTC),
				ip:    "127.0.0.1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := lc.log(tt.args.req, tt.args.res, tt.args.start, tt.args.stop, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("log() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, `{"time":"2020-05-19T09:00:00Z","id":"foo","remote_ip":"127.0.0.1","host":"","method":"","uri":"","user_agent":"","status":0,"latency":"0","latency_human":"0s","bytes_in":"0","bytes_out":"0"}`, strings.Trim(buf.String(), "\n"))
		})
	}
}

func Test_getTime(t *testing.T) {
	type args struct {
		stop time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"default",
			args{stop: time.Date(2020, 5, 19, 9, 0, 0, 0, time.UTC)},
			"2020-05-19T09:00:00Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTime(tt.args.stop); got != tt.want {
				t.Errorf("getTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getID(t *testing.T) {
	var (
		foo = "foo"
		h0  = http.Header{}
		h1  = http.Header{}
	)

	h0.Add(echo.HeaderXRequestID, foo)
	h0.Add(echo.HeaderXRequestID, "")

	rec := httptest.NewRecorder()
	rec.HeaderMap.Add(echo.HeaderXRequestID, foo) // nolint
	res := &echo.Response{Writer: rec}

	type args struct {
		req *http.Request
		res *echo.Response
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"request header first",
			args{
				req: &http.Request{
					Header: h0,
				},
				res: res,
			},
			foo,
		},
		{
			"request header first 2",
			args{
				req: &http.Request{
					Header: h1,
				},
				res: res,
			},
			foo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getID(tt.args.req, tt.args.res); got != tt.want {
				t.Errorf("getID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLatency(t *testing.T) {
	type args struct {
		stop  time.Time
		start time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"start first",
			args{
				start: time.Date(2020, 3, 19, 10, 0, 0, 0, time.UTC),
				stop:  time.Date(2020, 3, 19, 15, 0, 0, 0, time.UTC),
			},
			"-18000000000000",
		},
		{
			"stop first",
			args{
				start: time.Date(2020, 3, 19, 20, 0, 0, 0, time.UTC),
				stop:  time.Date(2020, 3, 19, 15, 0, 0, 0, time.UTC),
			},
			"18000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLatency(tt.args.stop, tt.args.start); got != tt.want {
				t.Errorf("getLatency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLatencyHuman(t *testing.T) {
	type args struct {
		stop  time.Time
		start time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"start first",
			args{
				start: time.Date(2020, 3, 19, 10, 0, 0, 0, time.UTC),
				stop:  time.Date(2020, 3, 19, 15, 0, 0, 0, time.UTC),
			},
			"-5h0m0s",
		},
		{
			"stop first",
			args{
				start: time.Date(2020, 3, 19, 20, 0, 0, 0, time.UTC),
				stop:  time.Date(2020, 3, 19, 15, 0, 0, 0, time.UTC),
			},
			"5h0m0s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLatencyHuman(tt.args.stop, tt.args.start); got != tt.want {
				t.Errorf("getLatencyHuman() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRequestSize(t *testing.T) {
	var foo = "foo"

	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"default",
			args{&http.Request{}},
			"0",
		},
		{
			"foo",
			args{&http.Request{
				Header: http.Header(map[string][]string{
					echo.HeaderContentLength: {foo},
				}),
			}},
			"foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRequestSize(tt.args.req); got != tt.want {
				t.Errorf("getRequestSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getResponseSize(t *testing.T) {
	res := &echo.Response{}

	type args struct {
		size int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"423",
			args{423},
			"423",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res.Size = tt.args.size
			if got := getResponseSize(res); got != tt.want {
				t.Errorf("getResponseSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
