package config

import "time"

// CacheConfig はValkey設定
type CacheConfig struct {
	Host         string        `env:"CACHE_HOST" envDefault:"localhost"`
	Port         int           `env:"CACHE_PORT" envDefault:"6379"`
	Password     string        `env:"CACHE_PASSWORD"`
	Database     int           `env:"CACHE_DATABASE" envDefault:"0"`
	PoolSize     int           `env:"CACHE_POOL_SIZE" envDefault:"10"`
	MinIdleConns int           `env:"CACHE_MIN_IDLE_CONNS" envDefault:"5"`
	MaxRetries   int           `env:"CACHE_MAX_RETRIES" envDefault:"3"`
	DialTimeout  time.Duration `env:"CACHE_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"CACHE_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"CACHE_WRITE_TIMEOUT" envDefault:"3s"`
}
