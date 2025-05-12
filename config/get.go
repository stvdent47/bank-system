package config

import "os"

func GetEnv(key string, defaultValue string) string {
	secret := os.Getenv(key)

	if secret == "" {
		return defaultValue
	}

	return secret
}
