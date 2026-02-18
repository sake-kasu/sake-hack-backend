package config

// JWTConfig はJWT設定
type JWTConfig struct {
	Secret       string `yaml:"secret" env:"JWT_SECRET"`
	Expiration   int    `yaml:"expiration" env:"JWT_EXPIRATION"`
	CookieSecure bool   `yaml:"cookieSecure" env:"JWT_COOKIE_SECURE"`
	CookieName   string `yaml:"cookieName" env:"JWT_COOKIE_NAME"`
	CookiePath   string `yaml:"cookiePath" env:"JWT_COOKIE_PATH"`
	CookieDomain string `yaml:"cookieDomain" env:"JWT_COOKIE_DOMAIN"`
}
