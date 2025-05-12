package config

import "os"

type JwtConfig struct {
	Secret string
	Exp    int
}

func LoadJwtConfig() JwtConfig {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "dev-secret-jwt"
	}

	return JwtConfig{
		Secret: secret,
		Exp:    24,
	}
}
