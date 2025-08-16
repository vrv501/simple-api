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

// Find Pets using name, status, tags.
// (GET /pets)
func (a *ApiHandler) FindPets(ctx context.Context, request genRouter.FindPetsRequestObject) (genRouter.FindPetsResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Add new pet to the store.
// (POST /pets)
func (a *ApiHandler) AddPet(ctx context.Context, request genRouter.AddPetRequestObject) (genRouter.AddPetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete a pet.
// (DELETE /pets/{petId})
func (a *ApiHandler) DeletePet(ctx context.Context, request genRouter.DeletePetRequestObject) (genRouter.DeletePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Find pet by ID.
// (GET /pets/{petId})
func (a *ApiHandler) GetPetById(ctx context.Context, request genRouter.GetPetByIdRequestObject) (genRouter.GetPetByIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Replace existing pet data using Id.
// (PUT /pets/{petId})
func (a *ApiHandler) ReplacePet(ctx context.Context, request genRouter.ReplacePetRequestObject) (genRouter.ReplacePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Upload a new image for a pet.
// (POST /pets/{petId}/images)
func (a *ApiHandler) UploadPetImage(ctx context.Context, request genRouter.UploadPetImageRequestObject) (genRouter.UploadPetImageResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete a pet image.
// (DELETE /pets/{petId}/images/{imageId})
func (a *ApiHandler) DeletePetImage(ctx context.Context, request genRouter.DeletePetImageRequestObject) (genRouter.DeletePetImageResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get a pet image using ID.
// (GET /pets/{petId}/images/{imageId})
func (a *ApiHandler) GetImageByPetId(ctx context.Context, request genRouter.GetImageByPetIdRequestObject) (genRouter.GetImageByPetIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
