package config

import "time"

// RateLimitConfig はレート制限設定
type RateLimitConfig struct {
	Enabled     bool          `env:"RATE_LIMIT_ENABLED" envDefault:"true"`
	MaxRequests int           `env:"RATE_LIMIT_MAX_REQUESTS" envDefault:"60"`
	Window      time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"60s"`
}
