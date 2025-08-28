package constants

import "time"

// ENV Vars
const (
	LogLevel       = "LOG_LEVEL"
	ServerPort     = "SERVER_PORT"
	DBUsername     = "DB_USERNAME"
	DBPassword     = "DB_PASSWORD"
	AllowedOrigins = "ALLOWED_ORIGINS"
)

// Default values for various configurations
const (
	DefaultTimeout = 3 * time.Minute

	MongoDB  = "mongodb"
	Postgres = "postgres"
)
