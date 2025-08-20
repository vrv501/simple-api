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
	errMsgAnimalCategoryExists    = "Animal category already exists with name"
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
		if errors.Is(err, dbErr.ErrConflict) {
			return genRouter.AddAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgAnimalCategoryExists + " " + categoryName,
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

// Delete an animal-category.
// (DELETE /animal-categories/{animalCategoryId})
func (a *APIHandler) DeleteAnimalCategory(ctx context.Context,
	request genRouter.DeleteAnimalCategoryRequestObject) (genRouter.DeleteAnimalCategoryResponseObject, error) {
	logger := log.Ctx(ctx)
	id := request.AnimalCategoryId

	err := a.dbClient.DeleteAnimalCategory(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, dbErr.ErrInvalidID):
			return genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgInvalidAnimalCategoryID,
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		case errors.Is(err, dbErr.ErrNotFound):
			return genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgAnimalCategoryNotFound + " " + id,
				},
				StatusCode: http.StatusNotFound,
			}, nil
		case errors.Is(err, dbErr.ErrForeignKeyConstraint):
			return genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "Pets found for animal category " + id,
				},
				StatusCode: http.StatusUnprocessableEntity,
			}, nil
		}

		logger.Error().Err(err).Msg("Failed to delete animal category")
		return genRouter.DeleteAnimalCategorydefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Deleted animal category with ID %s", id)
	return genRouter.DeleteAnimalCategory204Response{}, nil
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
		case errors.Is(err, dbErr.ErrConflict):
			return genRouter.ReplaceAnimalCategorydefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgAnimalCategoryExists + " " + categoryName,
				},
				StatusCode: http.StatusUnprocessableEntity,
			}, nil
		}

		logger.Error().Err(err).Msg("Failed to update animal category")
		return genRouter.ReplaceAnimalCategorydefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Updated animal category with ID %s", id)
	return genRouter.ReplaceAnimalCategory200JSONResponse{AnimalCategoryJSONResponse: *res}, nil
}
