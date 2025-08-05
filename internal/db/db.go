package db

import (
	"context"
	"os"

	"github.com/vrv501/simple-api/internal/constants"
	"github.com/vrv501/simple-api/internal/db/mongodb"
)

type DBHandler interface {
	AddPet(ctx context.Context)

	Close(ctx context.Context) error
}

func NewDBHandler(ctx context.Context) DBHandler {
	switch dbEnv := os.Getenv("DB_TYPE"); dbEnv {
	case constants.MongoDB:
		return mongodb.NewInstance(ctx)
	default:
		panic("Unsupported DB_TYPE: " + dbEnv)
	}
}
