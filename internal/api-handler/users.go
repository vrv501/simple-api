package apihandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

const (
	errMsgAlreadyInUse = "already in use"
	errMsgUserNotFound = "user not found"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

// Create user.
// (POST /users)
func (a *APIHandler) CreateUser(ctx context.Context,
	request genRouter.CreateUserRequestObject) (genRouter.CreateUserResponseObject, error) {
	logger := log.Ctx(ctx)
	userReq := request.Body

	hashedPswd, err := hashPassword(userReq.Password)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to hash password")
		return genRouter.CreateUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	userReq.Password = hashedPswd

	res, err := a.dbClient.AddUser(ctx, userReq)
	if err != nil {
		var conflictErr *dbErr.ConflictError
		if errors.As(err, &conflictErr) {
			return genRouter.CreateUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: conflictErr.Key + " " + errMsgAlreadyInUse,
				},
				StatusCode: http.StatusConflict,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to add user")
		return genRouter.CreateUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Username %s created successfully", res.Username)
	return genRouter.CreateUser201JSONResponse{UserJSONResponse: *res}, nil
}

// Delete user resource.
// (DELETE /users/{username})
func (a *APIHandler) DeleteUser(ctx context.Context,
	request genRouter.DeleteUserRequestObject) (genRouter.DeleteUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get user by user name.
// (GET /users/{username})
func (a *APIHandler) GetUserByName(ctx context.Context,
	request genRouter.GetUserByNameRequestObject) (genRouter.GetUserByNameResponseObject, error) {
	logger := log.Ctx(ctx)
	res, err := a.dbClient.GetUser(ctx, request.Username)
	if err != nil {
		if errors.Is(err, dbErr.ErrNotFound) {
			return genRouter.GetUserByNamedefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgUserNotFound,
				},
				StatusCode: http.StatusNotFound,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to get user")
		return genRouter.GetUserByNamedefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	return genRouter.GetUserByName200JSONResponse{UserJSONResponse: *res}, nil
}

// Replace user resource.
// (PUT /users/{username})
func (a *APIHandler) ReplaceUser(ctx context.Context,
	request genRouter.ReplaceUserRequestObject) (genRouter.ReplaceUserResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
