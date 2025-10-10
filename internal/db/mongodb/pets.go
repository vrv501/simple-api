package mongodb

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func (m *mongoClient) AddPet(ctx context.Context, userID string,
	petReq *genRouter.AddPetMultipartBody) error {
	// Validate userID string
	userbsonID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return &dbErr.HintError{Key: userIDField, Err: dbErr.ErrInvalidValue}
	}
	// Validate price string
	price, err := bson.ParseDecimal128(petReq.Pet.Price)
	if err != nil {
		return &dbErr.HintError{Key: priceField, Err: dbErr.ErrInvalidValue}
	}

	// Validate animalCategory exists
	res := m.mongoDbHandler.Collection(animalCategoryCollection).
		FindOne(
			ctx,
			bson.M{nameField: petReq.Pet.Category},
			options.FindOne().SetProjection(bson.M{iDField: 1}),
		)
	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &dbErr.HintError{Key: animalCategoryCollection, Err: dbErr.ErrNotFound}
		}
		return err
	}
	var animalCategoryDetail animalCategory
	err = res.Decode(&animalCategoryDetail)
	if err != nil {
		return err
	}

	_, err = m.performAdvisoryLockDBOperation(ctx, userbsonID, func(aCtx context.Context) (any, error) {
		petID := bson.NewObjectID()
		session, errS := m.client.StartSession()
		if errS != nil {
			return nil, errS
		}
		defer session.EndSession(aCtx)

		// Validate userID exists
		errS = m.mongoDbHandler.Collection(usersCollection).
			FindOne(
				aCtx,
				bson.M{iDField: userbsonID, deletedOnField: bson.Null{}},
				options.FindOne().SetProjection(bson.M{iDField: 1}),
			).Err()
		if errS != nil {
			if errors.Is(errS, mongo.ErrNoDocuments) {
				return nil, &dbErr.HintError{Key: userIDField, Err: dbErr.ErrNotFound}
			}
			return nil, errS
		}

		_, errS = session.WithTransaction(
			aCtx,
			func(sessCtx context.Context) (any, error) {
				// Insert new pet record
				petInstance := pet{
					ID:         petID,
					Name:       petReq.Pet.Name,
					CategoryID: animalCategoryDetail.ID,
					UserID:     userbsonID,
					Price:      price,
					Status:     string(genRouter.Available),
					CreatedOn:  time.Now().UTC(),
				}
				if petReq.Pet.Tags != nil {
					petInstance.Tags = *petReq.Pet.Tags
				}
				_, errY := m.mongoDbHandler.Collection(petsCollection).
					InsertOne(sessCtx, petInstance)
				if errY != nil {
					if mongo.IsDuplicateKeyError(errY) {
						return nil, dbErr.ErrConflict
					}
					return nil, errY
				}

				if len(petReq.Photos) == 0 {
					return nil, nil
				}
				// Upload photos
				imageList := make([]image, len(petReq.Photos))
				for i := range petReq.Photos {
					imgBytes, _ := petReq.Photos[i].Bytes()
					imageList[i].PetID = petID
					imageList[i].Image = imgBytes
				}
				_, errY = m.mongoDbHandler.Collection(imagesCollection).
					InsertMany(sessCtx, imageList, options.InsertMany().SetOrdered(false))
				return nil, errY
			},
			// Transactions apparently require read preference to be primary
			options.Transaction().SetReadPreference(readpref.Primary()),
		)
		return nil, errS
	})
	return err
}

func (m *mongoClient) GetPetImage(ctx context.Context, imageID string) (io.Reader, int64, error) {
	bsonImageID, err := bson.ObjectIDFromHex(imageID)
	if err != nil {
		return nil, 0, dbErr.ErrInvalidValue
	}

	res := m.mongoDbHandler.Collection(imagesCollection).FindOne(ctx,
		bson.M{iDField: bsonImageID, deletedOnField: bson.Null{}},
		options.FindOne().SetProjection(bson.M{imageField: 1}))
	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, 0, dbErr.ErrNotFound
		}
		return nil, 0, err
	}
	var img image
	err = res.Decode(&img)
	if err != nil {
		return nil, 0, err
	}

	return bytes.NewReader(img.Image), int64(len(img.Image)), nil
}

func (m *mongoClient) UploadPhotos(ctx context.Context, userID, petID string,
	petReq *genRouter.UploadPetImageMultipartBody) error {
	// Validate userID string
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return &dbErr.HintError{Key: userIDField, Err: dbErr.ErrInvalidValue}
	}
	// Validate petID string
	bsonPetID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return &dbErr.HintError{Key: petIDField, Err: dbErr.ErrInvalidValue}
	}

	// Get advisory lock so that petID is not sold until images are uploaded
	_, err = m.performAdvisoryLockDBOperation(ctx, bsonPetID, func(aCtx context.Context) (any, error) {
		return m.performAdvisoryLockDBOperation(aCtx, bsonUserID, func(bCtx context.Context) (any, error) {
			// Validate userID exists
			errS := m.mongoDbHandler.Collection(usersCollection).
				FindOne(
					bCtx,
					bson.M{iDField: bsonUserID, deletedOnField: bson.Null{}},
					options.FindOne().SetProjection(bson.M{iDField: 1}),
				).Err()
			if errS != nil {
				if errors.Is(errS, mongo.ErrNoDocuments) {
					return nil, &dbErr.HintError{Key: userIDField, Err: dbErr.ErrNotFound}
				}
				return nil, errS
			}

			// Validate petID exists and belongs to userID
			res := m.mongoDbHandler.Collection(petsCollection).FindOne(bCtx,
				bson.M{iDField: bsonPetID, userIDField: bsonUserID},
				options.FindOne().SetProjection(bson.M{userIDField: 1, statusField: 1}),
			)
			errS = res.Err()
			if errS != nil {
				if errors.Is(errS, mongo.ErrNoDocuments) {
					return nil, dbErr.ErrUserIDMismatch
				}
				return nil, errS
			}
			var p pet
			errS = res.Decode(&p)
			if errS != nil {
				return nil, errS
			}
			// Only allow deletion of an image if pet is in "available" status
			if p.Status != string(genRouter.Available) {
				return nil, nil
			}

			session, errS := m.client.StartSession()
			if errS != nil {
				return nil, errS
			}
			defer session.EndSession(bCtx)

			_, errS = session.WithTransaction(
				bCtx,
				func(sessCtx context.Context) (any, error) {
					if len(petReq.Photos) == 0 {
						return nil, nil
					}
					// Upload photos
					imageList := make([]image, len(petReq.Photos))
					for i := range petReq.Photos {
						imgBytes, _ := petReq.Photos[i].Bytes()
						imageList[i].PetID = bsonPetID
						imageList[i].Image = imgBytes
					}
					_, errY := m.mongoDbHandler.Collection(imagesCollection).
						InsertMany(sessCtx, imageList, options.InsertMany().SetOrdered(false))
					return nil, errY
				})
			return nil, errS
		})
	})
	return err
}

func (m *mongoClient) GetPhotoCount(ctx context.Context, petID string) (int, error) {
	bsonPetID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return 0, &dbErr.HintError{Key: petIDField, Err: dbErr.ErrInvalidValue}
	}

	count, err := m.mongoDbHandler.Collection(imagesCollection).CountDocuments(
		ctx,
		bson.M{petIDField: bsonPetID, deletedOnField: bson.Null{}},
	)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (m *mongoClient) GetPetIDForImage(ctx context.Context, imageID string) (string, error) {
	bsonImageID, err := bson.ObjectIDFromHex(imageID)
	if err != nil {
		return "", dbErr.ErrInvalidValue
	}

	res := m.mongoDbHandler.Collection(imagesCollection).FindOne(ctx,
		bson.M{iDField: bsonImageID, deletedOnField: bson.Null{}},
		options.FindOne().SetProjection(bson.M{petIDField: 1}),
	)
	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", dbErr.ErrNotFound
		}
		return "", err
	}

	var img image
	err = res.Decode(&img)
	if err != nil {
		return "", err
	}

	return img.PetID.Hex(), nil
}

func (m *mongoClient) DeletePetImage(ctx context.Context, userID, petID, imageID string) error {
	bsonImageID, err := bson.ObjectIDFromHex(imageID)
	if err != nil {
		return &dbErr.HintError{Key: iDField, Err: dbErr.ErrInvalidValue}
	}
	bsonPetID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return &dbErr.HintError{Key: petIDField, Err: dbErr.ErrInvalidValue}
	}
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return &dbErr.HintError{Key: userIDField, Err: dbErr.ErrInvalidValue}
	}

	// Get advisory lock so that petID is not sold until image is deleted
	_, err = m.performAdvisoryLockDBOperation(ctx, bsonPetID, func(aCtx context.Context) (any, error) {
		// Get advisory lock so that userID is not deleted until image is deleted
		return m.performAdvisoryLockDBOperation(aCtx, bsonUserID, func(bCtx context.Context) (any, error) {
			// Validate userID exists
			errS := m.mongoDbHandler.Collection(usersCollection).
				FindOne(
					bCtx,
					bson.M{iDField: bsonUserID, deletedOnField: bson.Null{}},
					options.FindOne().SetProjection(bson.M{iDField: 1}),
				).Err()
			if errS != nil {
				if errors.Is(errS, mongo.ErrNoDocuments) {
					return nil, &dbErr.HintError{Key: userIDField, Err: dbErr.ErrNotFound}
				}
				return nil, errS
			}

			// Validate petID exists and belongs to userID
			res := m.mongoDbHandler.Collection(petsCollection).FindOne(bCtx,
				bson.M{iDField: bsonPetID, userIDField: bsonUserID},
				options.FindOne().SetProjection(bson.M{userIDField: 1, statusField: 1}),
			)
			errS = res.Err()
			if errS != nil {
				if errors.Is(errS, mongo.ErrNoDocuments) {
					return nil, dbErr.ErrUserIDMismatch
				}
				return nil, errS
			}
			var p pet
			errS = res.Decode(&p)
			if errS != nil {
				return nil, errS
			}
			// Only allow deletion of an image if pet is in "available" status
			if p.Status != string(genRouter.Available) {
				return nil, nil
			}

			// Mark image for soft-delete
			_, errS = m.mongoDbHandler.Collection(imagesCollection).UpdateOne(
				bCtx,
				bson.M{iDField: bsonImageID, deletedOnField: bson.Null{}},
				bson.M{setOperator: bson.M{deletedOnField: time.Now().UTC()}},
			)
			if errS != nil {
				if errors.Is(errS, mongo.ErrNoDocuments) {
					return nil, dbErr.ErrNotFound
				}
			}
			return nil, errS
		})
	})
	return err
}
