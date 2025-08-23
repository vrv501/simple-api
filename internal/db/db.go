package db

import (
	"context"
	"os"
	"testing"

	"github.com/vrv501/simple-api/internal/db/mongodb"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

type Handler interface {
	animalCategoryHandler
	userHandler
	Close(ctx context.Context) error
}

type animalCategoryHandler interface {
	FindAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error)
	AddAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error)
	UpdateAnimalCategory(ctx context.Context, id, name string) (*genRouter.AnimalCategoryJSONResponse, error)
}

type userHandler interface {
	AddUser(ctx context.Context,
		userReq *genRouter.CreateUserJSONRequestBody) (*genRouter.UserJSONResponse, error)
	GetUser(ctx context.Context,
		username string) (*genRouter.UserJSONResponse, error)
	DeleteUser(ctx context.Context, username string) error
}

func NewDBHandler(ctx context.Context) Handler {
	switch dbEnv := os.Getenv("DB_TYPE"); dbEnv {
	case "mongodb":
		return mongodb.NewInstance(ctx)
	case "postgres":
		return nil
	default:
		if testing.Testing() {
			return nil
		}
		return mongodb.NewInstance(ctx)
	}
}
