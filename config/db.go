package config

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Username: GetEnv("db_username", "postgres"),
		Password: GetEnv("db_password", "1qaz2wsx"),
		Host:     GetEnv("db_host", "localhost"),
		Port:     GetEnv("db_port", "5432"),
		Name:     GetEnv("db_name", "bank"),
	}
}
