package db

import (
	"context"
	"os"
)

type DBHandler interface {
	AddPet(ctx context.Context)

	Close()
}

func NewDBHandler() DBHandler {
	switch dbEnv := os.Getenv("DB_TYPE"); dbEnv {

	default:
		panic("Unsupported DB_TYPE: " + dbEnv)
	}
}
