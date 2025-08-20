package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func (m *mongoClient) AddUser(ctx context.Context,
	userReq *genRouter.CreateUserJSONRequestBody) (*genRouter.UserJSONResponse, error) {
	currTime := time.Now().UTC()
	_, err := m.mongoDbHandler.Collection(usersCollection).InsertOne(ctx, user{
		Username:    userReq.Username,
		Password:    userReq.Password,
		Email:       userReq.Email,
		FullName:    userReq.FullName,
		PhoneNumber: userReq.PhoneNumber,
		CreatedOn:   currTime,
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
		CreatedAt:   currTime,
	}, nil
}
