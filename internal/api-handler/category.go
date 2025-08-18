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

// Find animal-category using name
// (GET /animal-categories)
func (a *ApiHandler) FindAnimalCategory(ctx context.Context, request genRouter.FindAnimalCategoryRequestObject) (genRouter.FindAnimalCategoryResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Add new animal-category to the store.
// (POST /animal-categories)
func (a *ApiHandler) AddAnimalCategory(ctx context.Context, request genRouter.AddAnimalCategoryRequestObject) (genRouter.AddAnimalCategoryResponseObject, error) {
	logger := log.Ctx(ctx)
	categoryName := string(request.Body.Name)
	id, createdTime, err := a.dbClient.AddAnimalCategory(ctx, categoryName)
	if err != nil {
		if errors.Is(err, dbErr.ErrConflict) {
			return genRouter.AddAnimalCategorydefaultJSONResponse{
				Body: genRouter.ApiResponse{
					Message: fmt.Sprintf("Animal Category: %s already exists", categoryName),
				},
				StatusCode: http.StatusConflict,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to add animal category")
		return genRouter.AddAnimalCategorydefaultJSONResponse{
			Body: genRouter.ApiResponse{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Added animal category %s", categoryName)
	return genRouter.AddAnimalCategory201JSONResponse{
		AnimalCategoryResponseJSONResponse: genRouter.AnimalCategoryResponseJSONResponse{
			Id:        id,
			Name:      categoryName,
			CreatedAt: createdTime,
		},
	}, nil
}

// Delete an animal-category.
// (DELETE /animal-categories/{animalCategoryId})
func (a *ApiHandler) DeleteAnimalCategory(ctx context.Context, request genRouter.DeleteAnimalCategoryRequestObject) (genRouter.DeleteAnimalCategoryResponseObject, error) {
	logger := log.Ctx(ctx)
	id := string(request.AnimalCategoryId)

	err := a.dbClient.DeleteAnimalCategory(ctx, id)
	if err != nil {
		if errors.Is(err, dbErr.ErrInvalidId) {
			return genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.ApiResponse{
					Message: "Invalid animal category ID",
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		if errors.Is(err, dbErr.ErrNotFound) {
			logger.Warn().Msgf("Animal category with ID %s not found", id)
			return genRouter.DeleteAnimalCategorydefaultJSONResponse{
				Body: genRouter.ApiResponse{
					Message: fmt.Sprintf("Animal category not found for id %s", id),
				},
				StatusCode: http.StatusNotFound,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to delete animal category")
		return genRouter.DeleteAnimalCategorydefaultJSONResponse{
			Body: genRouter.ApiResponse{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Deleted animal category with ID %s", id)
	return genRouter.DeleteAnimalCategory200Response{}, nil
}

// Replace existing animal-category data using Id.
// (PUT /animal-categories/{animalCategoryId})
func (a *ApiHandler) ReplaceAnimalCategory(ctx context.Context, request genRouter.ReplaceAnimalCategoryRequestObject) (genRouter.ReplaceAnimalCategoryResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
