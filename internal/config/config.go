package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config はアプリケーション全体の設定
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Logger   LoggerConfig
}

// Load は設定を読み込む
// .envファイルが存在する場合は先に読み込み、環境変数から設定をパースする
func Load() (*Config, error) {
	// 開発環境: .envファイルがあれば読み込み (エラーは無視)
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("環境変数のパースに失敗しました: %w", err)
	}

	return cfg, nil
}
