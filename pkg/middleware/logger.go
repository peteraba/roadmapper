package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		Logger *zap.Logger
	}
)

var (
	// DefaultLoggerConfig is the default Logger middleware config.
	DefaultLoggerConfig = LoggerConfig{
		Skipper: DefaultSkipper,
	}
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(&DefaultLoggerConfig)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(lc *LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if lc.Skipper == nil {
		lc.Skipper = DefaultLoggerConfig.Skipper
	}
	if lc.Logger == nil {
		lc.Logger, _ = zap.NewProduction()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error

			if lc.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			ip := c.RealIP()

			err = lc.log(req, res, start, stop, ip)

			return err
		}
	}
}

func (lc *LoggerConfig) log(req *http.Request, res *echo.Response, start, stop time.Time, ip string) error {
	lc.Logger.Info("echo request",
		zap.String("time", getTime(stop)),
		zap.String("id", getID(req, res)),
		zap.String("remote_ip", ip),
		zap.String("host", req.Host),
		zap.String("method", req.Method),
		zap.String("uri", req.RequestURI),
		zap.String("user_agent", req.UserAgent()),
		zap.Int("status", res.Status),
		zap.String("latency", getLatency(start, stop)),
		zap.String("latency_human", getLatencyHuman(start, stop)),
		zap.String("bytes_in", getRequestSize(req)),
		zap.String("bytes_out", getResponseSize(res)),
	)

	err := lc.Logger.Sync()

	return err
}

func getTime(stop time.Time) string {
	return stop.Format(time.RFC3339Nano)
}

func getID(req *http.Request, res *echo.Response) string {
	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}

	return id
}

func getLatency(start, stop time.Time) string {
	l := stop.Sub(start)
	return strconv.FormatInt(int64(l), 10)
}

func getLatencyHuman(start, stop time.Time) string {
	return stop.Sub(start).String()
}

func getRequestSize(req *http.Request) string {
	requestSize := req.Header.Get(echo.HeaderContentLength)
	if requestSize == "" {
		requestSize = "0"
	}

	return requestSize
}

func getResponseSize(res *echo.Response) string {
	return strconv.FormatInt(res.Size, 10)
}
