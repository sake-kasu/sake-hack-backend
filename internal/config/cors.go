package config

// CORSConfig はCORS設定
type CORSConfig struct {
	AllowedOrigins   []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:3000,http://localhost:8080"`
	AllowedMethods   []string `env:"CORS_ALLOWED_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,DELETE,PATCH,OPTIONS"`
	AllowedHeaders   []string `env:"CORS_ALLOWED_HEADERS" envSeparator:"," envDefault:"Origin,Content-Type,Accept,Authorization"`
	ExposedHeaders   []string `env:"CORS_EXPOSED_HEADERS" envSeparator:"," envDefault:"Content-Length"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
	MaxAge           int      `env:"CORS_MAX_AGE" envDefault:"43200"`
}
