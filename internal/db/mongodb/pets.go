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
			return &dbErr.HintError{Key: nameField, Err: dbErr.ErrNotFound}
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
				bson.M{iDField: userbsonID},
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
