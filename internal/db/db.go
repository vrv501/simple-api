package db

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/vrv501/simple-api/internal/db/mongodb"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

type Handler interface {
	animalCategoryHandler
	userHandler
	petsHandler
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
		userID string) (*genRouter.UserJSONResponse, error)
	PatchUser(ctx context.Context, userID string,
		userReq *genRouter.PatchUserApplicationMergePatchPlusJSONRequestBody) (*genRouter.UserJSONResponse, error)
	DeleteUser(ctx context.Context, userID string) error
}

type petsHandler interface {
	AddPet(ctx context.Context, userID string,
		petReq *genRouter.AddPetMultipartBody) error
	GetPetImage(ctx context.Context, imageID string) (io.Reader, int64, error)
	DeletePetImage(ctx context.Context, userID, imageID string) error
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
