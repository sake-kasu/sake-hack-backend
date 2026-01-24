package config

import "time"

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port                    int           `env:"SERVER_PORT" envDefault:"8080"`
	Mode                    string        `env:"SERVER_MODE" envDefault:"debug"`
	GracefulShutdownTimeout time.Duration `env:"SERVER_GRACEFUL_SHUTDOWN_TIMEOUT" envDefault:"30s"`
}
