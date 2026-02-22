package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestAccessLogger_InfoFor2xx(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, observed := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)

	router := gin.New()
	router.Use(RequestID())
	router.Use(AccessLogger(AccessLoggerConfig{
		Logger:    zapLogger,
		SkipPaths: map[string]struct{}{},
	}))
	router.GET("/test/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test/123?foo=bar", nil)
	req.Header.Set("User-Agent", "middleware-test")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, observed.FilterMessage("HTTPリクエスト開始").Len())
	assert.Equal(t, 1, observed.FilterMessage("HTTPリクエスト完了").Len())

	entry := observed.FilterMessage("HTTPリクエスト完了").All()[0]
	fields := entry.ContextMap()
	assert.Equal(t, zapcore.InfoLevel, entry.Level)
	assert.Equal(t, "GET", fields["method"])
	assert.Equal(t, "/test/:id", fields["route"])
	assert.Equal(t, "/test/123?foo=bar", fields["uri"])
	assert.Equal(t, int64(http.StatusOK), fields["status"])
	assert.Equal(t, "middleware-test", fields["user_agent"])

	// RequestIDミドルウェアが付与したIDが含まれていること
	_, hasRequestID := fields["request_id"]
	assert.True(t, hasRequestID)
}

func TestAccessLogger_WarnFor4xx(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, observed := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)

	router := gin.New()
	router.Use(AccessLogger(AccessLoggerConfig{
		Logger:    zapLogger,
		SkipPaths: map[string]struct{}{},
	}))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	entry := observed.FilterMessage("HTTPリクエスト完了").All()[0]
	assert.Equal(t, zapcore.WarnLevel, entry.Level)
}

func TestAccessLogger_ErrorFor5xx(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, observed := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)

	router := gin.New()
	router.Use(AccessLogger(AccessLoggerConfig{
		Logger:    zapLogger,
		SkipPaths: map[string]struct{}{},
	}))
	router.GET("/test", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	entry := observed.FilterMessage("HTTPリクエスト完了").All()[0]
	assert.Equal(t, zapcore.ErrorLevel, entry.Level)
}

func TestAccessLogger_SkipHealthPathByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, observed := observer.New(zapcore.DebugLevel)
	zapLogger := zap.New(core)

	router := gin.New()
	router.Use(AccessLogger(AccessLoggerConfig{
		Logger: zapLogger,
	}))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, observed.All(), 0)
}
