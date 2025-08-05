package constants

import "time"

const (
	DefaultServerPort      = 8300
	StatusPath             = "/status"
	MongoDB                = "mongodb"
	Postgres               = "postgres"
	DefaultShutdownTimeout = 3 * time.Minute
)
