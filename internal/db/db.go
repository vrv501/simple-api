package db

import (
	"context"
	"os"
	"time"

	"github.com/vrv501/simple-api/internal/constants"
	"github.com/vrv501/simple-api/internal/db/mongodb"
)

type DBHandler interface {
	animalCategoryHandler

	Close(ctx context.Context) error
}

type animalCategoryHandler interface {
	AddAnimalCategory(ctx context.Context, name string) (string, time.Time, error)

	DeleteAnimalCategory(ctx context.Context, id string) error
}

func NewDBHandler(ctx context.Context) DBHandler {
	switch dbEnv := os.Getenv("DB_TYPE"); dbEnv {
	case constants.MongoDB:
		return mongodb.NewInstance(ctx)
	default:
		return mongodb.NewInstance(ctx)
	}
}
