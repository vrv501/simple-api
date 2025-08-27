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
		Email:       userReq.Email,
		FullName:    userReq.FullName,
		PhoneNumber: userReq.PhoneNumber,
		CreatedOn:   time.Now().UTC(),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			se := mongo.ServerError(nil)
			_ = errors.As(err, &se)
			switch {
			case se.HasErrorMessage(emailField):
				return nil, &dbErr.ConflictError{Key: emailField, Err: err}
			case se.HasErrorMessage(phoneNumberField):
				return nil, &dbErr.ConflictError{Key: phoneNumberField, Err: err}
			}
			return nil, &dbErr.ConflictError{Key: usernameField, Err: err}
		}
		return nil, err
	}

	return &genRouter.UserJSONResponse{
		Username:    userReq.Username,
		Email:       userReq.Email,
		FullName:    userReq.FullName,
		PhoneNumber: userReq.PhoneNumber,
		Address:     userReq.Address,
	}, nil
}

func (m *mongoClient) GetUser(ctx context.Context,
	username string) (*genRouter.UserJSONResponse, error) {
	res := m.mongoDbHandler.Collection(usersCollection).FindOne(
		ctx,
		bson.M{usernameField: username, deletedOnField: bson.Null{}},
		options.FindOne().SetProjection(bson.M{
			iDField:       0,
			passwordField: 0,
		}),
	)
	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, dbErr.ErrNotFound
		}
		return nil, err
	}
	var user user
	err = res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &genRouter.UserJSONResponse{
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
	}, nil
}

func (m *mongoClient) DeleteUser(aInctx context.Context, username string) error {
	_, err := m.performAdvisoryLockDBOperation(aInctx, username, func(aCtx context.Context) (any, error) {
		res := m.mongoDbHandler.Collection(usersCollection).FindOne(
			aCtx,
			bson.M{usernameField: username, deletedOnField: bson.Null{}},
			options.FindOne().SetProjection(bson.M{iDField: 1}),
		)
		err := res.Err()
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, dbErr.ErrNotFound
			}
			return nil, err
		}
		var userInst user
		err = res.Decode(&userInst)
		if err != nil {
			return nil, err
		}

		err = m.mongoDbHandler.Collection(petsCollection).FindOne(
			aCtx,
			bson.M{userIDField: userInst.ID, statusField: genRouter.Available},
			options.FindOne().SetProjection(bson.M{iDField: 1}),
		).Err()
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		if err == nil {
			return nil, &dbErr.ForeignKeyError{Key: petsCollection, Err: err}
		}

		err = m.mongoDbHandler.Collection(ordersCollection).FindOne(
			aCtx,
			bson.M{
				userIDField: userInst.ID,
				statusField: bson.M{
					notInOperator: []string{
						string(genRouter.Delivered),
						string(genRouter.Cancelled),
					},
				},
			},
			options.FindOne().SetProjection(bson.M{iDField: 1}),
		).Err()
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		if err == nil {
			return nil, &dbErr.ForeignKeyError{Key: ordersCollection, Err: err}
		}

		_, err = m.mongoDbHandler.Collection(usersCollection).
			UpdateOne(aCtx, bson.M{iDField: userInst.ID, deletedOnField: bson.Null{}},
				bson.M{setOperator: bson.M{deletedOnField: time.Now().UTC()}})
		return nil, err
	})
	return err
}
