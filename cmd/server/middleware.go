package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newLoggingMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogProtocol:  true,
		LogLatency:   true,
		LogHost:      true,
		LogURIPath:   true,
		LogRequestID: true,
		LogError:     true,
		HandleError:  true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			msg := fmt.Sprintf("%s - %s %s %s %d", v.RemoteIP, v.Method, v.URIPath, v.Protocol, v.Status)
			protocol := strings.Split(v.Protocol, "/")
			scheme, version := protocol[0], protocol[1]
			httpMap := map[string]string{
				"url":         fmt.Sprintf("%s://%s%s", strings.ToLower(scheme), v.Host, v.URI),
				"path":        v.URIPath,
				"method":      v.Method,
				"version":     version,
				"host":        v.Host,
				"status_code": fmt.Sprint(v.Status),
			}

			if v.Error == nil {
				slog.LogAttrs(context.Background(),
					slog.LevelInfo,
					msg,
					slog.String("level", "info"),
					slog.Time("timestamp", time.Now()),
					slog.Int64("duration", v.Latency.Nanoseconds()),
					slog.String("request_id", v.RequestID),
					slog.Any("http", httpMap),
				)
			} else {
				slog.LogAttrs(context.Background(),
					slog.LevelError,
					msg,
					slog.String("level", "error"),
					slog.Time("timestamp", time.Now()),
					slog.Int64("duration", v.Latency.Nanoseconds()),
					slog.String("request_id", v.RequestID),
					slog.Any("http", httpMap),
					slog.String("error", v.Error.Error()),
				)
			}

			return nil
		},
	},
	)
}
