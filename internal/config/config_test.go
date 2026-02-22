package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// clearEnv は環境変数を全てクリアし、テスト終了時に元の状態に復元する
func clearEnv(t *testing.T) {
	t.Helper()
	original := os.Environ()
	os.Clearenv()
	t.Cleanup(func() {
		os.Clearenv()
		for _, kv := range original {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				_ = os.Setenv(parts[0], parts[1])
			}
		}
	})
}

// TestLoad_WithEnvironmentVariables は環境変数から設定を読み込むテスト
func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// 環境変数を設定
	t.Setenv("SERVER_PORT", "9999")
	t.Setenv("SERVER_MODE", "release")
	t.Setenv("SERVER_GRACEFUL_SHUTDOWN_TIMEOUT", "60s")

	t.Setenv("DB_HOST", "test-db")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "test_db")
	t.Setenv("DB_USER", "test_user")
	t.Setenv("DB_PASSWORD", "test_pass")
	t.Setenv("DB_SSL_MODE", "require")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")
	t.Setenv("DB_CONN_MAX_LIFETIME", "10m")

	t.Setenv("CACHE_HOST", "test-valkey")
	t.Setenv("CACHE_PORT", "6380")
	t.Setenv("CACHE_PASSWORD", "test_valkey_pass")
	t.Setenv("CACHE_DATABASE", "1")
	t.Setenv("CACHE_POOL_SIZE", "20")
	t.Setenv("CACHE_MIN_IDLE_CONNS", "10")
	t.Setenv("CACHE_MAX_RETRIES", "5")
	t.Setenv("CACHE_DIAL_TIMEOUT", "10s")
	t.Setenv("CACHE_READ_TIMEOUT", "5s")
	t.Setenv("CACHE_WRITE_TIMEOUT", "5s")

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRATION", "3600")
	t.Setenv("JWT_COOKIE_SECURE", "true")
	t.Setenv("JWT_COOKIE_NAME", "test_token")
	t.Setenv("JWT_COOKIE_PATH", "/test")
	t.Setenv("JWT_COOKIE_DOMAIN", "test.com")

	t.Setenv("CORS_ALLOWED_ORIGINS", "http://test.com")
	t.Setenv("CORS_ALLOWED_METHODS", "GET,POST")
	t.Setenv("CORS_ALLOWED_HEADERS", "Content-Type")
	t.Setenv("CORS_EXPOSED_HEADERS", "X-Total-Count")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "false")
	t.Setenv("CORS_MAX_AGE", "7200")

	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("LOG_FORMAT", "json")

	// テスト実行
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Server設定の検証
	assert.Equal(t, 9999, cfg.Server.Port)
	assert.Equal(t, "release", cfg.Server.Mode)
	assert.Equal(t, 60*time.Second, cfg.Server.GracefulShutdownTimeout)

	// Database設定の検証
	assert.Equal(t, "test-db", cfg.Database.Host)
	assert.Equal(t, 5433, cfg.Database.Port)
	assert.Equal(t, "test_db", cfg.Database.Database)
	assert.Equal(t, "test_user", cfg.Database.User)
	assert.Equal(t, "test_pass", cfg.Database.Password)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10, cfg.Database.MaxIdleConns)
	assert.Equal(t, 10*time.Minute, cfg.Database.ConnMaxLifetime)

	// Cache設定の検証
	assert.Equal(t, "test-valkey", cfg.Cache.Host)
	assert.Equal(t, 6380, cfg.Cache.Port)
	assert.Equal(t, "test_valkey_pass", cfg.Cache.Password)
	assert.Equal(t, 1, cfg.Cache.Database)
	assert.Equal(t, 20, cfg.Cache.PoolSize)
	assert.Equal(t, 10, cfg.Cache.MinIdleConns)
	assert.Equal(t, 5, cfg.Cache.MaxRetries)
	assert.Equal(t, 10*time.Second, cfg.Cache.DialTimeout)
	assert.Equal(t, 5*time.Second, cfg.Cache.ReadTimeout)
	assert.Equal(t, 5*time.Second, cfg.Cache.WriteTimeout)

	// JWT設定の検証
	assert.Equal(t, "test-secret", cfg.JWT.Secret)
	assert.Equal(t, 3600, cfg.JWT.Expiration)
	assert.True(t, cfg.JWT.CookieSecure)
	assert.Equal(t, "test_token", cfg.JWT.CookieName)
	assert.Equal(t, "/test", cfg.JWT.CookiePath)
	assert.Equal(t, "test.com", cfg.JWT.CookieDomain)

	// CORS設定の検証
	assert.Equal(t, []string{"http://test.com"}, cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST"}, cfg.CORS.AllowedMethods)
	assert.Equal(t, []string{"Content-Type"}, cfg.CORS.AllowedHeaders)
	assert.Equal(t, []string{"X-Total-Count"}, cfg.CORS.ExposedHeaders)
	assert.False(t, cfg.CORS.AllowCredentials)
	assert.Equal(t, 7200, cfg.CORS.MaxAge)

	// Logger設定の検証
	assert.Equal(t, "info", cfg.Logger.Level)
	assert.Equal(t, "json", cfg.Logger.Format)
}

// TestLoad_WithDefaultValues はデフォルト値を使用するテスト
func TestLoad_WithDefaultValues(t *testing.T) {
	// 環境変数を全てクリアし、テスト終了時に元の状態に復元する
	// (ローカルのシェル変数や.envの影響を排除するため)
	clearEnv(t)

	// .envファイルが読み込まれないように、一時ディレクトリに移動
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(originalWd)
	})
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// 必須環境変数のみ設定
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("JWT_SECRET", "test_secret")

	// テスト実行
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// デフォルト値が設定されていることを確認
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, 30*time.Second, cfg.Server.GracefulShutdownTimeout)

	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "sake_hack_app", cfg.Database.Database)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "test_password", cfg.Database.Password)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, 25, cfg.Database.MaxOpenConns)
	assert.Equal(t, 5, cfg.Database.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, cfg.Database.ConnMaxLifetime)

	assert.Equal(t, "localhost", cfg.Cache.Host)
	assert.Equal(t, 6379, cfg.Cache.Port)
	assert.Equal(t, "", cfg.Cache.Password)
	assert.Equal(t, 0, cfg.Cache.Database)
	assert.Equal(t, 10, cfg.Cache.PoolSize)
	assert.Equal(t, 5, cfg.Cache.MinIdleConns)
	assert.Equal(t, 3, cfg.Cache.MaxRetries)
	assert.Equal(t, 5*time.Second, cfg.Cache.DialTimeout)
	assert.Equal(t, 3*time.Second, cfg.Cache.ReadTimeout)
	assert.Equal(t, 3*time.Second, cfg.Cache.WriteTimeout)

	assert.Equal(t, "test_secret", cfg.JWT.Secret)
	assert.Equal(t, 86400, cfg.JWT.Expiration)
	assert.False(t, cfg.JWT.CookieSecure)
	assert.Equal(t, "sake_hack_token", cfg.JWT.CookieName)
	assert.Equal(t, "/", cfg.JWT.CookiePath)
	assert.Equal(t, "", cfg.JWT.CookieDomain)

	assert.Equal(t, []string{"http://localhost:3000", "http://localhost:8080"}, cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}, cfg.CORS.AllowedMethods)
	assert.Equal(t, []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Like-Token"}, cfg.CORS.AllowedHeaders)
	assert.Equal(t, []string{"Content-Length"}, cfg.CORS.ExposedHeaders)
	assert.True(t, cfg.CORS.AllowCredentials)
	assert.Equal(t, 43200, cfg.CORS.MaxAge)

	assert.Equal(t, "debug", cfg.Logger.Level)
	assert.Equal(t, "console", cfg.Logger.Format)
}

// TestLoad_MissingRequiredFields は必須フィールドが欠けている場合のテスト
func TestLoad_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T)
		wantErr bool
	}{
		{
			name: "DB_PASSWORD欠如",
			setup: func(t *testing.T) {
				t.Setenv("JWT_SECRET", "test_secret")
				// DB_PASSWORDを設定しない
			},
			wantErr: true,
		},
		{
			name: "JWT_SECRET欠如",
			setup: func(t *testing.T) {
				t.Setenv("DB_PASSWORD", "test_password")
				// JWT_SECRETを設定しない
			},
			wantErr: true,
		},
		{
			name: "両方設定",
			setup: func(t *testing.T) {
				t.Setenv("DB_PASSWORD", "test_password")
				t.Setenv("JWT_SECRET", "test_secret")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 環境変数をクリアし、テスト終了時に元の状態に復元する
			clearEnv(t)

			// テストごとの環境変数設定
			tt.setup(t)

			// テスト実行
			cfg, err := Load()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

// TestLoad_InvalidDurationFormat は不正なDurationフォーマットのテスト
func TestLoad_InvalidDurationFormat(t *testing.T) {
	// 必須環境変数を設定
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("JWT_SECRET", "test_secret")

	// 不正なDuration形式を設定
	t.Setenv("SERVER_GRACEFUL_SHUTDOWN_TIMEOUT", "invalid_duration")

	// テスト実行
	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

// TestLoad_CORSArrayParsing はCORS配列のパースをテスト
func TestLoad_CORSArrayParsing(t *testing.T) {
	// 必須環境変数を設定
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("JWT_SECRET", "test_secret")

	// カンマ区切りの配列を設定
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080,https://example.com")
	t.Setenv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE")

	// テスト実行
	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 配列が正しくパースされていることを確認
	assert.Equal(t, []string{"http://localhost:3000", "http://localhost:8080", "https://example.com"}, cfg.CORS.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE"}, cfg.CORS.AllowedMethods)
}
