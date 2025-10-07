package config

import "strings"

type EnvConfig struct {
	Env             string `env:"ENV" envDefault:"dev"`
	ClerkSecretKey  string `env:"CLERK_SECRET_KEY"`
	DatabaseURL     string `env:"DATABASE_URL"`
	FrontendOrigins string `env:"FRONTEND_ORIGINS" envDefault:"http://localhost:5173"`
	EncryptionKey   string `env:"ENCRYPTION_KEY"`
	TWDatakitAddr   string `env:"TW_DATAKIT_ADDR"`
	ServiceName     string `env:"SERVICE_NAME"`
}

func (c EnvConfig) IsDev() bool {
	return strings.EqualFold(c.Env, "dev")
}
