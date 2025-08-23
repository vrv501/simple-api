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
			case se.HasErrorMessage("email"):
				return nil, &dbErr.ConflictError{Key: "email", Err: err}
			case se.HasErrorMessage("phone_number"):
				return nil, &dbErr.ConflictError{Key: "phone_number", Err: err}
			}
			return nil, &dbErr.ConflictError{Key: "username", Err: err}
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

func (m *mongoClient) DeleteUser(ctx context.Context, username string) error {

	res := m.mongoDbHandler.Collection(usersCollection).FindOne(
		ctx,
		bson.M{usernameField: username, deletedOnField: bson.Null{}},
		options.FindOne().SetProjection(bson.M{iDField: 1}),
	)
	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dbErr.ErrNotFound
		}
		return err
	}

	var user user
	err = res.Decode(&user)
	if err != nil {
		return err
	}

	err = m.mongoDbHandler.Collection(petsCollection).FindOne(
		ctx,
		bson.M{userIDField: user.ID, statusField: genRouter.Available},
		options.FindOne().SetProjection(bson.M{iDField: 1}),
	).Err()
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return &dbErr.ForeignKeyError{Key: "pets", Err: err}
	}

	err = m.mongoDbHandler.Collection(ordersCollection).FindOne(
		ctx,
		bson.M{
			userIDField: user.ID,
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
		return err
	}
	if err == nil {
		return &dbErr.ForeignKeyError{Key: "orders", Err: err}
	}

	_, err = m.mongoDbHandler.Collection(usersCollection).
		UpdateOne(ctx, bson.M{iDField: user.ID},
			bson.M{setOperator: bson.M{deletedOnField: time.Now().UTC()}})
	return err
}
