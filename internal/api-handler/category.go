package apihandler

import (
	"context"

	"github.com/rs/zerolog/log"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Find animal-category using name
// (GET /animal-categories)
func (a *ApiHandler) FindAnimalCategory(ctx context.Context, request genRouter.FindAnimalCategoryRequestObject) (genRouter.FindAnimalCategoryResponseObject, error) {
	logger := log.Ctx(ctx)
	panic("not implemented") // TODO: Implement
}

// Add new animal-category to the store.
// (POST /animal-categories)
func (a *ApiHandler) AddAnimalCategory(ctx context.Context, request genRouter.AddAnimalCategoryRequestObject) (genRouter.AddAnimalCategoryResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete an animal-category.
// (DELETE /animal-categories/{animalCategoryId})
func (a *ApiHandler) DeleteAnimalCategory(ctx context.Context, request genRouter.DeleteAnimalCategoryRequestObject) (genRouter.DeleteAnimalCategoryResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Replace existing animal-category data using Id.
// (PUT /animal-categories/{animalCategoryId})
func (a *ApiHandler) ReplaceAnimalCategory(ctx context.Context, request genRouter.ReplaceAnimalCategoryRequestObject) (genRouter.ReplaceAnimalCategoryResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
