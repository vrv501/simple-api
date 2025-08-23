package apihandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

var (
	errMsgInvalidAnimalCategoryID = "Invalid animal category ID"
	errMsgAnimalCategoryNotFound  = "Animal category not found for id"
	errMsgAnimalCategoryExists    = "Animal category %s already exists"
)

// Find animal-category using name
// (GET /animal-categories)
func (a *APIHandler) FindAnimalCategory(ctx context.Context,
	request genRouter.FindAnimalCategoryRequestObject) (genRouter.FindAnimalCategoryResponseObject, error) {
	logger := log.Ctx(ctx)
	categoryName := request.Params.Name

	res, err := a.dbClient.FindAnimalCategory(ctx, categoryName)
	if err != nil {
		if errors.Is(err, dbErr.ErrNotFound) {
			return genRouter.FindAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: fmt.Sprintf("Animal category %s not found", categoryName),
				},
				StatusCode: http.StatusNotFound,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to find animal category")
		return genRouter.FindAnimalCategorydefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	return genRouter.FindAnimalCategory200JSONResponse{AnimalCategoryJSONResponse: *res}, nil
}

// Add new animal-category to the store.
// (POST /animal-categories)
func (a *APIHandler) AddAnimalCategory(ctx context.Context,
	request genRouter.AddAnimalCategoryRequestObject) (genRouter.AddAnimalCategoryResponseObject, error) {
	logger := log.Ctx(ctx)
	categoryName := request.Body.Name
	res, err := a.dbClient.AddAnimalCategory(ctx, categoryName)
	if err != nil {
		var conflictErr *dbErr.ConflictError
		if errors.As(err, &conflictErr) {
			return genRouter.AddAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: fmt.Sprintf(errMsgAnimalCategoryExists, categoryName),
				},
				StatusCode: http.StatusConflict,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to add animal category")
		return genRouter.AddAnimalCategorydefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Added animal category %s", categoryName)
	return genRouter.AddAnimalCategory201JSONResponse{AnimalCategoryJSONResponse: *res}, nil
}

// Replace existing animal-category data using Id.
// (PUT /animal-categories/{animalCategoryId})
func (a *APIHandler) ReplaceAnimalCategory(ctx context.Context,
	request genRouter.ReplaceAnimalCategoryRequestObject) (genRouter.ReplaceAnimalCategoryResponseObject, error) {
	logger := log.Ctx(ctx)
	categoryName := request.Body.Name
	id := request.AnimalCategoryId

	res, err := a.dbClient.UpdateAnimalCategory(ctx, id, categoryName)
	if err != nil {
		var conflictErr *dbErr.ConflictError
		switch {
		case errors.Is(err, dbErr.ErrInvalidID):
			return genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidAnimalCategoryID,
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		case errors.Is(err, dbErr.ErrNotFound):
			return genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgAnimalCategoryNotFound + " " + id,
				},
				StatusCode: http.StatusNotFound,
			}, nil
		case errors.As(err, &conflictErr):
			return genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: fmt.Sprintf(errMsgAnimalCategoryExists, categoryName),
				},
				StatusCode: http.StatusUnprocessableEntity,
			}, nil
		}

		logger.Error().Err(err).Msg("Failed to replace animal category")
		return genRouter.ReplaceAnimalCategorydefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Replaced animal category with ID %s", id)
	return genRouter.ReplaceAnimalCategory200JSONResponse{AnimalCategoryJSONResponse: *res}, nil
}
