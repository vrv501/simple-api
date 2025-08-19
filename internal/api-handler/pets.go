package apihandler

import (
	"context"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Find Pets using name, status, tags.
// (GET /pets)
func (a *APIHandler) FindPets(ctx context.Context, request genRouter.FindPetsRequestObject) (genRouter.FindPetsResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Add new pet to the store.
// (POST /pets)
func (a *APIHandler) AddPet(ctx context.Context, request genRouter.AddPetRequestObject) (genRouter.AddPetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete a pet.
// (DELETE /pets/{petId})
func (a *APIHandler) DeletePet(ctx context.Context, request genRouter.DeletePetRequestObject) (genRouter.DeletePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Find pet by ID.
// (GET /pets/{petId})
func (a *APIHandler) GetPetById(ctx context.Context, request genRouter.GetPetByIdRequestObject) (genRouter.GetPetByIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Replace existing pet data using Id.
// (PUT /pets/{petId})
func (a *APIHandler) ReplacePet(ctx context.Context, request genRouter.ReplacePetRequestObject) (genRouter.ReplacePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Upload a new image for a pet.
// (POST /pets/{petId}/images)
func (a *APIHandler) UploadPetImage(ctx context.Context, request genRouter.UploadPetImageRequestObject) (genRouter.UploadPetImageResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete a pet image.
// (DELETE /pets/{petId}/images/{imageId})
func (a *APIHandler) DeletePetImage(ctx context.Context, request genRouter.DeletePetImageRequestObject) (genRouter.DeletePetImageResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get a pet image using ID.
// (GET /pets/{petId}/images/{imageId})
func (a *APIHandler) GetImageByPetId(ctx context.Context, request genRouter.GetImageByPetIdRequestObject) (genRouter.GetImageByPetIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
