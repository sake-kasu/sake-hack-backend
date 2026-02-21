package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// テスト: レート制限が無効の場合、リクエストがバイパスされる
func TestRateLimiter_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RateLimiter(nil, RateLimitConfig{
		Enabled:     false,
		MaxRequests: 10,
		Window:      60 * time.Second,
	}))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// レート制限ヘッダーが付与されないことを確認
	assert.Empty(t, w.Header().Get("X-RateLimit-Limit"))
	assert.Empty(t, w.Header().Get("X-RateLimit-Remaining"))
	assert.Empty(t, w.Header().Get("X-RateLimit-Reset"))
}

// テスト: レート制限が無効の場合、複数リクエストがすべて通る
func TestRateLimiter_Disabled_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RateLimiter(nil, RateLimitConfig{
		Enabled:     false,
		MaxRequests: 1,
		Window:      60 * time.Second,
	}))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// MaxRequests=1でも無効なら全てバイパス
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

// テスト: RateLimitConfig構造体のゼロ値確認
func TestRateLimitConfig_ZeroValue(t *testing.T) {
	cfg := RateLimitConfig{}

	assert.False(t, cfg.Enabled)
	assert.Equal(t, 0, cfg.MaxRequests)
	assert.Equal(t, time.Duration(0), cfg.Window)
}
