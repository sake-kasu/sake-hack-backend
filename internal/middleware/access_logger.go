package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
	"go.uber.org/zap"
)

// AccessLoggerConfig はアクセスログミドルウェア設定
type AccessLoggerConfig struct {
	// SkipPaths に含まれるパスはログ出力しない
	SkipPaths map[string]struct{}
	// Logger を指定した場合はこのロガーを利用
	Logger *zap.Logger
}

// AccessLogger はリクエスト完了時に構造化アクセスログを出力する
func AccessLogger(cfg AccessLoggerConfig) gin.HandlerFunc {
	baseLogger := cfg.Logger
	if baseLogger == nil {
		baseLogger = logger.Get()
	}

	skipPaths := cfg.SkipPaths
	if skipPaths == nil {
		skipPaths = map[string]struct{}{
			"/health": {},
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if _, ok := skipPaths[path]; ok {
			c.Next()
			return
		}

		reqLogger := baseLogger
		if requestID, ok := c.Request.Context().Value(logger.RequestIDKey).(string); ok && requestID != "" {
			reqLogger = reqLogger.With(zap.String("request_id", requestID))
		}

		reqLogger.Info("HTTPリクエスト開始",
			zap.String("method", c.Request.Method),
			zap.String("uri", c.Request.URL.RequestURI()),
		)

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		route := c.FullPath()
		if route == "" {
			route = path
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("route", route),
			zap.String("uri", c.Request.URL.RequestURI()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Int64("latency_ms", latency.Milliseconds()),
			zap.Int("bytes", c.Writer.Size()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case c.Writer.Status() >= http.StatusInternalServerError:
			reqLogger.Error("HTTPリクエスト完了", fields...)
		case c.Writer.Status() >= http.StatusBadRequest:
			reqLogger.Warn("HTTPリクエスト完了", fields...)
		default:
			reqLogger.Info("HTTPリクエスト完了", fields...)
		}
	}
}
