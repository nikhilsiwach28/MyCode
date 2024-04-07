package config

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     GetEnvWithDefault("POSTGRES_HOST", "localhost"),
		Port:     GetIntEnvWithDefault("POSTGRES_PORT", 5432),
		User:     GetEnvWithDefault("POSTGRES_USER", "username"),
		Password: GetEnvWithDefault("POSTGRES_PASSWORD", "password"),
		DBName:   GetEnvWithDefault("POSTGRES_DBNAME", "database_name"),
		SSLMode:  GetEnvWithDefault("POSTGRES_SSLMODE", "disable"),
	}
}
