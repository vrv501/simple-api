package apihandler

import (
	"context"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Create user.
// (POST /user)
func (h *apiHandler) CreateUser(ctx context.Context, request genRouter.CreateUserRequestObject) (genRouter.CreateUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Creates list of users with given input array.
// (POST /user/createWithList)
func (h *apiHandler) CreateUsersWithListInput(ctx context.Context, request genRouter.CreateUsersWithListInputRequestObject) (genRouter.CreateUsersWithListInputResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Logs user into the system.
// (GET /user/login)
func (h *apiHandler) LoginUser(ctx context.Context, request genRouter.LoginUserRequestObject) (genRouter.LoginUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Logs out current logged in user session.
// (GET /user/logout)
func (h *apiHandler) LogoutUser(ctx context.Context, request genRouter.LogoutUserRequestObject) (genRouter.LogoutUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete user resource.
// (DELETE /user/{username})
func (h *apiHandler) DeleteUser(ctx context.Context, request genRouter.DeleteUserRequestObject) (genRouter.DeleteUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get user by user name.
// (GET /user/{username})
func (h *apiHandler) GetUserByName(ctx context.Context, request genRouter.GetUserByNameRequestObject) (genRouter.GetUserByNameResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Update user resource.
// (PUT /user/{username})
func (h *apiHandler) UpdateUser(ctx context.Context, request genRouter.UpdateUserRequestObject) (genRouter.UpdateUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
