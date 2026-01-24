package config

// LoggerConfig はロギング設定
type LoggerConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"debug"`
	Format string `env:"LOG_FORMAT" envDefault:"console"`
}
