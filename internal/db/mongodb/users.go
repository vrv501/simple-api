package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func (m *mongoClient) AddUser(ctx context.Context,
	userReq *genRouter.CreateUserJSONRequestBody) (*genRouter.UserJSONResponse, error) {
	_, err := m.mongoDbHandler.Collection(usersCollection).InsertOne(ctx, user{
		Username:    userReq.Username,
		Password:    userReq.Password,
		Address:     userReq.Address,
		FullName:    userReq.FullName,
		PhoneNumber: userReq.PhoneNumber,
		CreatedOn:   time.Now().UTC(),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			se := mongo.ServerError(nil)
			_ = errors.As(err, &se)
			if se.HasErrorMessage(phoneNumberField) {
				return nil, &dbErr.HintError{Key: phoneNumberField, Err: dbErr.ErrConflict}
			}
			return nil, &dbErr.HintError{Key: usernameField, Err: dbErr.ErrConflict}
		}
		return nil, err
	}

	return &genRouter.UserJSONResponse{
		Username:    userReq.Username,
		FullName:    userReq.FullName,
		PhoneNumber: userReq.PhoneNumber,
		Address:     userReq.Address,
	}, nil
}

func (m *mongoClient) GetUser(ctx context.Context,
	userID string) (*genRouter.UserJSONResponse, error) {
	bsonID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, dbErr.ErrInvalidValue
	}

	res := m.mongoDbHandler.Collection(usersCollection).FindOne(
		ctx,
		bson.M{iDField: bsonID, deletedOnField: bson.Null{}},
		options.FindOne().SetProjection(bson.M{
			iDField:       0,
			passwordField: 0,
		}),
	)
	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, dbErr.ErrNotFound
		}
		return nil, err
	}
	var userInstance user
	err = res.Decode(&userInstance)
	if err != nil {
		return nil, err
	}
	return &genRouter.UserJSONResponse{
		Username:    userInstance.Username,
		FullName:    userInstance.FullName,
		PhoneNumber: userInstance.PhoneNumber,
		Address:     userInstance.Address,
	}, nil
}

func (m *mongoClient) DeleteUser(aInctx context.Context, userID string) error {
	bsonID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return dbErr.ErrInvalidValue
	}

	_, err = m.performAdvisoryLockDBOperation(aInctx, bsonID, func(aCtx context.Context) (any, error) {
		errS := m.mongoDbHandler.Collection(petsCollection).FindOne(
			aCtx,
			bson.M{userIDField: bsonID, statusField: genRouter.Available},
			options.FindOne().SetProjection(bson.M{iDField: 1}),
		).Err()
		if errS != nil && !errors.Is(errS, mongo.ErrNoDocuments) {
			return nil, errS
		}
		if errS == nil {
			return nil, &dbErr.HintError{Key: petsCollection, Err: dbErr.ErrForeignKeyViolation}
		}

		errS = m.mongoDbHandler.Collection(ordersCollection).FindOne(
			aCtx,
			bson.M{
				userIDField: bsonID,
				statusField: bson.M{
					notInOperator: []string{
						string(genRouter.Delivered),
						string(genRouter.Cancelled),
					},
				},
			},
			options.FindOne().SetProjection(bson.M{iDField: 1}),
		).Err()
		if errS != nil && !errors.Is(errS, mongo.ErrNoDocuments) {
			return nil, errS
		}
		if errS == nil {
			return nil, &dbErr.HintError{Key: ordersCollection, Err: dbErr.ErrForeignKeyViolation}
		}

		res, errS := m.mongoDbHandler.Collection(usersCollection).
			UpdateOne(aCtx, bson.M{iDField: bsonID, deletedOnField: bson.Null{}},
				bson.M{setOperator: bson.M{deletedOnField: time.Now().UTC()}})
		if errS == nil {
			if res.MatchedCount == 0 {
				return nil, dbErr.ErrNotFound
			}
			return nil, nil
		}
		return nil, errS
	})
	return err
}

func (m *mongoClient) PatchUser(ctx context.Context, userID string,
	userReq *genRouter.PatchUserApplicationMergePatchPlusJSONRequestBody) (*genRouter.UserJSONResponse, error) {
	bsonID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, dbErr.ErrInvalidValue
	}
	updateDoc := bson.M{}
	if userReq.FullName != nil {
		updateDoc[fullNameField] = *userReq.FullName
	}
	if userReq.Password != nil {
		updateDoc[passwordField] = *userReq.Password
	}
	if userReq.PhoneNumber != nil {
		updateDoc[phoneNumberField] = *userReq.PhoneNumber
	}
	if userReq.Address != nil {
		updateDoc[addressField] = *userReq.Address
	}
	updateDoc[updatedOnField] = time.Now().UTC()

	res := m.mongoDbHandler.Collection(usersCollection).
		FindOneAndUpdate(
			ctx,
			bson.M{iDField: bsonID, deletedOnField: bson.Null{}},
			bson.M{setOperator: updateDoc},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		)
	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, dbErr.ErrNotFound
		}
		if mongo.IsDuplicateKeyError(err) {
			return nil, dbErr.ErrConflict
		}
		return nil, err
	}
	var userInstance user
	err = res.Decode(&userInstance)
	if err != nil {
		return nil, err
	}

	return &genRouter.UserJSONResponse{
		Address:     userInstance.Address,
		FullName:    userInstance.FullName,
		PhoneNumber: userInstance.PhoneNumber,
		Username:    userInstance.Username,
	}, nil
}
