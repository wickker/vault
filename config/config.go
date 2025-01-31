package config

type EnvConfig struct {
	Env            string `env:"ENV" envDefault:"dev"`
	ClerkSecretKey string `env:"CLERK_SECRET_KEY"`
}
