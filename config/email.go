package config

type SMTPConfig struct {
	Host     string
	Port     string
	From     string
	User     string
	Password string
}

func LoadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     GetEnv("smtp.host", "smtp.example.com"),
		Port:     GetEnv("smtp.port", "587"),
		From:     GetEnv("smtp.from", "example@gmail.com"),
		User:     GetEnv("smtp.user", "noreply@example.com"),
		Password: GetEnv("smtp.password", "strong_password"),
	}
}
