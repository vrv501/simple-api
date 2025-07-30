package apihandler

import (
	"context"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Add a new pet to the store.
// (POST /pet)
func (h *ApiHandler) AddPet(ctx context.Context, request genRouter.AddPetRequestObject) (genRouter.AddPetResponseObject, error) {

	panic("not implemented") // TODO: Implement
}

// Update an existing pet.
// (PUT /pet)
func (h *ApiHandler) UpdatePet(ctx context.Context, request genRouter.UpdatePetRequestObject) (genRouter.UpdatePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Finds Pets by status.
// (GET /pet/findByStatus)
func (h *ApiHandler) FindPetsByStatus(ctx context.Context, request genRouter.FindPetsByStatusRequestObject) (genRouter.FindPetsByStatusResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Finds Pets by tags.
// (GET /pet/findByTags)
func (h *ApiHandler) FindPetsByTags(ctx context.Context, request genRouter.FindPetsByTagsRequestObject) (genRouter.FindPetsByTagsResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Deletes a pet.
// (DELETE /pet/{petId})
func (h *ApiHandler) DeletePet(ctx context.Context, request genRouter.DeletePetRequestObject) (genRouter.DeletePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Find pet by ID.
// (GET /pet/{petId})
func (h *ApiHandler) GetPetById(ctx context.Context, request genRouter.GetPetByIdRequestObject) (genRouter.GetPetByIdResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Updates a pet in the store with form data.
// (POST /pet/{petId})
func (h *ApiHandler) UpdatePetWithForm(ctx context.Context, request genRouter.UpdatePetWithFormRequestObject) (genRouter.UpdatePetWithFormResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Uploads an image.
// (POST /pet/{petId}/uploadImage)
func (h *ApiHandler) UploadFile(ctx context.Context, request genRouter.UploadFileRequestObject) (genRouter.UploadFileResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
