package apihandler

import (
	"context"

	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

// Create user.
// (POST /users)
func (a *ApiHandler) CreateUser(ctx context.Context, request genRouter.CreateUserRequestObject) (genRouter.CreateUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete user resource.
// (DELETE /users/{username})
func (a *ApiHandler) DeleteUser(ctx context.Context, request genRouter.DeleteUserRequestObject) (genRouter.DeleteUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get user by user name.
// (GET /users/{username})
func (a *ApiHandler) GetUserByName(ctx context.Context, request genRouter.GetUserByNameRequestObject) (genRouter.GetUserByNameResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Replace user resource.
// (PUT /users/{username})
func (a *ApiHandler) ReplaceUser(ctx context.Context, request genRouter.ReplaceUserRequestObject) (genRouter.ReplaceUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
