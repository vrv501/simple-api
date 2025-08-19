package db

//go:generate go tool mockgen -package db -destination db_mock.go . Handler

import (
	"context"
	"os"

	"github.com/vrv501/simple-api/internal/constants"
	"github.com/vrv501/simple-api/internal/db/mongodb"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

type Handler interface {
	animalCategoryHandler

	Close(ctx context.Context) error
}

type animalCategoryHandler interface {
	FindAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error)

	AddAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error)

	UpdateAnimalCategory(ctx context.Context, id, name string) (*genRouter.AnimalCategoryJSONResponse, error)

	DeleteAnimalCategory(ctx context.Context, id string) error
}

func NewDBHandler(ctx context.Context) Handler {
	switch dbEnv := os.Getenv("DB_TYPE"); dbEnv {
	case constants.MongoDB:
		return mongodb.NewInstance(ctx)
	default:
		return nil
	}
}
