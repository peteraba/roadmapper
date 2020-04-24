package middleware

import (
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
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	if config.Logger == nil {
		config.Logger, _ = zap.NewProduction()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			l := stop.Sub(start)
			latency := strconv.FormatInt(int64(l), 10)

			requestSize := req.Header.Get(echo.HeaderContentLength)
			if requestSize == "" {
				requestSize = "0"
			}

			responseSize := strconv.FormatInt(res.Size, 10)

			config.Logger.Info("echo request",
				zap.String("time", time.Now().Format(time.RFC3339Nano)),
				zap.String("id", id),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("user_agent", req.UserAgent()),
				zap.Int("status", res.Status),
				zap.String("latency", latency),
				zap.String("latency_human", stop.Sub(start).String()),
				zap.String("bytes_in", requestSize),
				zap.String("bytes_out", responseSize),
			)

			err = config.Logger.Sync()

			return
		}
	}
}
