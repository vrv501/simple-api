package apihandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"

	"github.com/vrv501/simple-api/internal/constants"
	contextKeys "github.com/vrv501/simple-api/internal/context-keys"
	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

const (
	errMsgIncorrectReqEncoding = "request body is not properly encoded"
)

// Find Pets using name, status, tags.
// (GET /pets)
func (a *APIHandler) FindPets(_ context.Context,
	_ genRouter.FindPetsRequestObject) (genRouter.FindPetsResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

func validateImage(r io.Reader) ([]byte, error) {
	imgData := make([]byte, 1+constants.MaxImgSize)
	n, err := io.ReadFull(r, imgData)
	if err != nil &&
		!errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, errors.New("image data is corrupted")
	}
	if n == 0 || n > constants.MaxImgSize {
		return nil, errors.New("images should have min size 1B and max size 250KB")
	}
	imgData = imgData[:n]

	reader := bytes.NewReader(imgData)
	imgDetails, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, errors.New("jpeg image is corrupted")
	}
	if imgDetails.Width < 256 || imgDetails.Width > 1920 ||
		imgDetails.Height < 256 || imgDetails.Height > 1080 {
		return nil,
			errors.New("supported min resolution for images is 256x256px & max is 1920x1080px")
	}

	reader.Seek(0, io.SeekStart)
	_, err = jpeg.Decode(reader)
	if err != nil {
		return nil, errors.New("jpeg image is corrupted")
	}
	return imgData, nil
}

// Add new pet to the store.
// (POST /pets)
func (a *APIHandler) AddPet(ctx context.Context,
	request genRouter.AddPetRequestObject) (genRouter.AddPetResponseObject, error) {
	logger := log.Ctx(ctx)
	userID, ok := contextKeys.UserIDFromContext(ctx)
	if !ok {
		logger.Error().Msg(errMsgUserIDNotFound)
		return genRouter.AddPetdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	var (
		part     *multipart.Part
		err      error
		petData  genRouter.Pet
		oapifile openapi_types.File
		mpReq    = genRouter.AddPetMultipartBody{Photos: genRouter.PetPhotos{}}
	)

	// Read multipart form data
	for {
		part, err = request.Body.NextPart()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			logger.Error().Err(err).Msg("failed to read multipart data")
			return genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsgIncorrectReqEncoding,
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		}

		// Decode based on form Name
		switch part.FormName() {
		case "pet":
			if err = json.NewDecoder(part).Decode(&petData); err != nil {
				return genRouter.AddPetdefaultJSONResponse{
					Body: genRouter.Generic{
						Message: errMsgIncorrectReqEncoding,
					},
					StatusCode: http.StatusBadRequest,
				}, nil
			}
			mpReq.Pet = petData
		case "photos":
			imgData, errS := validateImage(part)
			if errS != nil {
				return genRouter.AddPetdefaultJSONResponse{
					Body: genRouter.Generic{
						Message: errS.Error(),
					},
					StatusCode: http.StatusBadRequest,
				}, nil
			}
			oapifile.InitFromBytes(imgData, part.FileName())
			mpReq.Photos = append(mpReq.Photos, oapifile)
		default:
			return genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "unknown multipart field " + part.FormName(),
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		}
	}

	err = a.dbClient.AddPet(ctx, userID, &mpReq)
	if err != nil {
		var dberr *dbErr.HintError
		if errors.As(err, &dberr) {
			err = dberr.Err
			errMsg := "bad request"
			statusCode := http.StatusBadRequest
			switch {
			case errors.Is(err, dbErr.ErrInvalidValue):
				errMsg = "invalid value for " + dberr.Key
				statusCode = http.StatusBadRequest
			case errors.Is(err, dbErr.ErrNotFound):
				errMsg = dberr.Key + " not found"
				statusCode = http.StatusNotFound
			}
			return genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: errMsg,
				},
				StatusCode: statusCode,
			}, nil
		} else if errors.Is(err, dbErr.ErrConflict) {
			return genRouter.AddPetdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "Similar pet already exists",
				},
				StatusCode: http.StatusConflict,
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to insert new pet")
		return genRouter.AddPetdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msg("Successfully inserted pet")
	return genRouter.AddPet202Response{}, nil
}

// Delete a pet.
// (DELETE /pets/{petId})
func (a *APIHandler) DeletePet(_ context.Context,
	_ genRouter.DeletePetRequestObject) (genRouter.DeletePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Find pet by ID.
// (GET /pets/{petId})
func (a *APIHandler) GetPetByID(_ context.Context,
	_ genRouter.GetPetByIDRequestObject) (genRouter.GetPetByIDResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Replace existing pet data using Id.
// (PUT /pets/{petId})
func (a *APIHandler) ReplacePet(_ context.Context,
	_ genRouter.ReplacePetRequestObject) (genRouter.ReplacePetResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Upload a new image for a pet.
// (POST /pets/{petId}/images)
func (a *APIHandler) UploadPetImage(_ context.Context,
	_ genRouter.UploadPetImageRequestObject) (genRouter.UploadPetImageResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Delete a pet image.
// (DELETE /images/{imageId})
func (a *APIHandler) DeletePetImage(ctx context.Context,
	request genRouter.DeletePetImageRequestObject) (genRouter.DeletePetImageResponseObject, error) {
	logger := log.Ctx(ctx)
	userID, ok := contextKeys.UserIDFromContext(ctx)
	if !ok {
		logger.Error().Msg(errMsgUserIDNotFound)
		return genRouter.DeletePetImagedefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	err := a.dbClient.DeletePetImage(ctx, userID, request.ImageId)
	if err != nil {
		var invalidErr *dbErr.HintError
		if errors.As(err, &invalidErr) {
			return genRouter.DeletePetImagedefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "invalid " + invalidErr.Key,
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		} else if errors.Is(err, dbErr.ErrNotFound) {
			return genRouter.DeletePetImagedefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "image not found",
				},
				StatusCode: http.StatusNotFound,
			}, nil
		}
		logger.Error().Err(err).Msgf("failed to soft-delete pet image %s", request.ImageId)
		return genRouter.DeletePetImagedefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	logger.Info().Msgf("Successfully soft-deleted pet image %s", request.ImageId)
	return genRouter.DeletePetImage204Response{}, nil
}

// Get a pet image using ID.
// (GET /images/{imageId})
func (a *APIHandler) GetImageByID(ctx context.Context,
	request genRouter.GetImageByIDRequestObject) (genRouter.GetImageByIDResponseObject, error) {
	logger := log.Ctx(ctx)
	reader, contentLength, err := a.dbClient.GetPetImage(ctx, request.ImageId)
	if err != nil {
		switch {
		case errors.Is(err, dbErr.ErrNotFound):
			return genRouter.GetImageByIDdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "image not found",
				},
				StatusCode: http.StatusNotFound,
			}, nil
		case errors.Is(err, dbErr.ErrInvalidValue):
			return genRouter.GetImageByIDdefaultJSONResponse{
				Body: genRouter.Generic{
					Message: "invalid image ID",
				},
				StatusCode: http.StatusBadRequest,
			}, nil
		}
		logger.Error().Err(err).Msg("failed to get pet image")
		return genRouter.GetImageByIDdefaultJSONResponse{
			Body: genRouter.Generic{
				Message: http.StatusText(http.StatusInternalServerError),
			},
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	return genRouter.GetImageByID200ImagejpegResponse{
		Body:          reader,
		ContentLength: contentLength,
	}, nil
}
