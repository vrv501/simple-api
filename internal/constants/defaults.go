package constants

import "time"

const (
	DefaultServerPort      = 8300
	DefaultShutdownTimeout = 3 * time.Minute

	StatusPath = "/status"

	MongoDB  = "mongodb"
	Postgres = "postgres"
)
