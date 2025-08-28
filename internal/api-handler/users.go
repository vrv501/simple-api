package apihandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	contextKeys "github.com/vrv501/simple-api/internal/context-keys"
	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

const (
	errMsgAlreadyInUse  = "already in use"
	errMsgUserNotFound  = "user not found"
	errMsgInvalidUserID = "Invalid user ID"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func comparePasswords(hashedPswd, plainTextPswd string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPswd),
		[]byte(plainTextPswd),
	)
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
	_ genRouter.DeleteUserRequestObject) (genRouter.DeleteUserResponseObject, error) {
	logger := log.Ctx(ctx)
	userID, ok := contextKeys.UserIDFromContext(ctx)
	if !ok {
		logger.Error().Msg("userID not found in context")
		return genRouter.DeleteUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	err := a.dbClient.DeleteUser(ctx, userID)
	if err != nil {
		var fKeyErr *dbErr.ForeignKeyError
		switch {
		case errors.Is(err, dbErr.ErrInvalidID):
			return genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidUserID,
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		case errors.Is(err, dbErr.ErrNotFound):
			return genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgUserNotFound,
				},
				StatusCode: http.StatusNotFound,
			}, nil
		case errors.As(err, &fKeyErr):
			return genRouter.DeleteUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "User cannot be deleted as there are pending " + fKeyErr.Key,
				},
				StatusCode: http.StatusUnprocessableEntity,
			}, nil
		}

		logger.Error().Err(err).Msg("Failed to soft-delete user")
		return genRouter.DeleteUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	logger.Info().Msgf("UserID %s soft-deleted", userID)
	return genRouter.DeleteUser204Response{}, nil
}

// Get user by user name.
// (GET /users/{username})
func (a *APIHandler) GetUser(ctx context.Context,
	_ genRouter.GetUserRequestObject) (genRouter.GetUserResponseObject, error) {
	logger := log.Ctx(ctx)
	userID, ok := contextKeys.UserIDFromContext(ctx)
	if !ok {
		logger.Error().Msg("userID not foudn in context")
		return genRouter.GetUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	res, err := a.dbClient.GetUser(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, dbErr.ErrInvalidID):
			return genRouter.GetUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidUserID,
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		case errors.Is(err, dbErr.ErrNotFound):
			return genRouter.GetUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgUserNotFound,
				},
				StatusCode: http.StatusNotFound,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to get user")
		return genRouter.GetUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	return genRouter.GetUser200JSONResponse{UserJSONResponse: *res}, nil
}

// Replace user resource.
// (PUT /users/{username})
func (a *APIHandler) PatchUser(ctx context.Context,
	request genRouter.PatchUserRequestObject) (genRouter.PatchUserResponseObject, error) {
	logger := log.Ctx(ctx)
	userID, ok := contextKeys.UserIDFromContext(ctx)
	if !ok {
		logger.Error().Msg("userID not found in context")
		return genRouter.PatchUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	userReq := request.Body
	if userReq.Address == nil && userReq.Password == nil &&
		userReq.FullName == nil && userReq.PhoneNumber == nil {
		return genRouter.PatchUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: "Nothing to Update",
			},
			StatusCode: http.StatusBadRequest,
		}, nil
	}
	if userReq.Password != nil {
		hashedPswd, err := hashPassword(*userReq.Password)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to hash password")
			return genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: http.StatusText(http.StatusInternalServerError),
				},
				StatusCode: http.StatusInternalServerError,
			}, nil
		}
		userReq.Password = &hashedPswd
	}

	resp, err := a.dbClient.PatchUser(ctx, userID, userReq)
	if err != nil {
		var conflictErr *dbErr.ConflictError
		switch {
		case errors.Is(err, dbErr.ErrInvalidID):
			return genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidUserID,
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		case errors.Is(err, dbErr.ErrNotFound):
			return genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgUserNotFound,
				},
				StatusCode: http.StatusNotFound,
			}, nil
		case errors.As(err, &conflictErr):
			return genRouter.PatchUserdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: conflictErr.Key + " " + errMsgAlreadyInUse,
				},
				StatusCode: http.StatusConflict,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to patch user")
		return genRouter.PatchUserdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	logger.Info().Msgf("UserID %s patched", userID)
	return genRouter.PatchUser200JSONResponse{
		UserJSONResponse: *resp,
	}, nil
}
