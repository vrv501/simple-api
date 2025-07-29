package apiv1

import (
	"context"

	"github.com/vrv501/go-template/internal/generated"
)

type V1Api struct{}

// Returns all pets
// (GET /pets)
func (a *V1Api) FindPets(ctx context.Context, request generated.FindPetsRequestObject) (generated.FindPetsResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Creates a new pet
// (POST /pets)
func (a *V1Api) AddPet(ctx context.Context, request generated.AddPetRequestObject) (generated.AddPetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Deletes a pet by ID
// (DELETE /pets/{id})
func (a *V1Api) DeletePet(ctx context.Context, request generated.DeletePetRequestObject) (generated.DeletePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Returns a pet by ID
// (GET /pets/{id})
func (a *V1Api) FindPetByID(ctx context.Context, request generated.FindPetByIDRequestObject) (generated.FindPetByIDResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
