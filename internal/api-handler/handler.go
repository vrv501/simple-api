package apihandler

import (
	"net/http"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

// (POST /pet)
func (h *handler) AddPet(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Update an existing pet.
// (PUT /pet)
func (h *handler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Finds Pets by status.
// (GET /pet/findByStatus)
func (h *handler) FindPetsByStatus(w http.ResponseWriter, r *http.Request,
	params genRouter.FindPetsByStatusParams) {
	panic("not implemented") // TODO: Implement
}

// Finds Pets by tags.
// (GET /pet/findByTags)
func (h *handler) FindPetsByTags(w http.ResponseWriter, r *http.Request,
	params genRouter.FindPetsByTagsParams) {
	panic("not implemented") // TODO: Implement
}

// Deletes a pet.
// (DELETE /pet/{petId})
func (h *handler) DeletePet(w http.ResponseWriter, r *http.Request, petId int64,
	params genRouter.DeletePetParams) {
	panic("not implemented") // TODO: Implement
}

// Find pet by ID.
// (GET /pet/{petId})
func (h *handler) GetPetById(w http.ResponseWriter, r *http.Request,
	petId int64) {
	panic("not implemented") // TODO: Implement
}

// Updates a pet in the store with form data.
// (POST /pet/{petId})
func (h *handler) UpdatePetWithForm(w http.ResponseWriter, r *http.Request, petId int64,
	params genRouter.UpdatePetWithFormParams) {
	panic("not implemented") // TODO: Implement
}

// Uploads an image.
// (POST /pet/{petId}/uploadImage)
func (h *handler) UploadFile(w http.ResponseWriter, r *http.Request, petId int64,
	params genRouter.UploadFileParams) {
	panic("not implemented") // TODO: Implement
}

// Returns pet inventories by status.
// (GET /store/inventory)
func (h *handler) GetInventory(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Place an order for a pet.
// (POST /store/order)
func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Delete purchase order by identifier.
// (DELETE /store/order/{orderId})
func (h *handler) DeleteOrder(w http.ResponseWriter, r *http.Request, orderId int64) {
	panic("not implemented") // TODO: Implement
}

// Find purchase order by ID.
// (GET /store/order/{orderId})
func (h *handler) GetOrderById(w http.ResponseWriter, r *http.Request, orderId int64) {
	panic("not implemented") // TODO: Implement
}

// Create user.
// (POST /user)
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Creates list of users with given input array.
// (POST /user/createWithList)
func (h *handler) CreateUsersWithListInput(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Logs user into the system.
// (GET /user/login)
func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request,
	params genRouter.LoginUserParams) {
	panic("not implemented") // TODO: Implement
}

// Logs out current logged in user session.
// (GET /user/logout)
func (h *handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Delete user resource.
// (DELETE /user/{username})
func (h *handler) DeleteUser(w http.ResponseWriter,
	r *http.Request, username string) {
	panic("not implemented") // TODO: Implement
}

// Get user by user name.
// (GET /user/{username})
func (h *handler) GetUserByName(w http.ResponseWriter, r *http.Request, username string) {
	panic("not implemented") // TODO: Implement
}

// Update user resource.
// (PUT /user/{username})
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, username string) {
	panic("not implemented") // TODO: Implement
}
