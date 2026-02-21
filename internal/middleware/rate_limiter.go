package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"

	"github.com/sake-kasu/sake-hack-backend/api/generated"
	"github.com/sake-kasu/sake-hack-backend/internal/apperror"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// RateLimitConfig はレート制限ミドルウェアの設定
type RateLimitConfig struct {
	Enabled     bool
	MaxRequests int
	Window      time.Duration
}

// RateLimiter はValkeyを使用したIPベースのレート制限ミドルウェア
func RateLimiter(client valkey.Client, cfg RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enabled {
			c.Next()
			return
		}

		ctx := c.Request.Context()
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		count, ttl, err := incrementAndGetTTL(ctx, client, key, cfg.Window)
		if err != nil {
			// Valkeyエラー時はレート制限をバイパス(サービス継続を優先)
			logger.Warn(ctx, "レート制限チェックに失敗しました(バイパスします)")
			c.Next()
			return
		}

		remaining := cfg.MaxRequests - int(count)
		if remaining < 0 {
			remaining = 0
		}

		// レート制限ヘッダーを付与
		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		if int(count) > cfg.MaxRequests {
			errors := []generated.APIError{
				{
					Code:    apperror.ErrCodeBadRequest,
					Message: "リクエスト数が上限を超えました。しばらく待ってからお試しください",
				},
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, generated.ErrorResponse{
				Data:   nil,
				Errors: &errors,
			})
			return
		}

		c.Next()
	}
}

// incrementAndGetTTL はレート制限カウンターをインクリメントしTTLを返す
func incrementAndGetTTL(ctx context.Context, client valkey.Client, key string, window time.Duration) (int64, time.Duration, error) {
	// INCR + EXPIRE パターン(Luaスクリプトで原子的に実行)
	script := valkey.NewLuaScript(`
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('PEXPIRE', KEYS[1], ARGV[1])
		end
		local ttl = redis.call('PTTL', KEYS[1])
		return {count, ttl}
	`)

	windowMs := strconv.FormatInt(window.Milliseconds(), 10)

	resp := script.Exec(ctx, client, []string{key}, []string{windowMs})
	if resp.Error() != nil {
		return 0, 0, resp.Error()
	}

	arr, err := resp.ToArray()
	if err != nil {
		return 0, 0, err
	}

	count, err := arr[0].ToInt64()
	if err != nil {
		return 0, 0, err
	}

	ttlMs, err := arr[1].ToInt64()
	if err != nil {
		return 0, 0, err
	}

	return count, time.Duration(ttlMs) * time.Millisecond, nil
}
