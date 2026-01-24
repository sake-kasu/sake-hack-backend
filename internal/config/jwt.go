package config

// JWTConfig はJWT設定
type JWTConfig struct {
	Secret       string `env:"JWT_SECRET,required,notEmpty"`
	Expiration   int    `env:"JWT_EXPIRATION" envDefault:"86400"`
	CookieSecure bool   `env:"JWT_COOKIE_SECURE" envDefault:"false"`
	CookieName   string `env:"JWT_COOKIE_NAME" envDefault:"sake_hack_token"`
	CookiePath   string `env:"JWT_COOKIE_PATH" envDefault:"/"`
	CookieDomain string `env:"JWT_COOKIE_DOMAIN"`
}
