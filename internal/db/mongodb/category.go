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

func (m *mongoClient) FindAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error) {
	res := m.mongoDbHandler.Collection(animalCategoryCollection).FindOne(ctx,
		bson.M{nameField: name},
	)
	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, dbErr.ErrNotFound
		}
		return nil, err
	}

	var animalCategoryRes animalCategory
	err = res.Decode(&animalCategoryRes)
	if err != nil {
		return nil, err
	}
	return &genRouter.AnimalCategoryJSONResponse{
		Id:   animalCategoryRes.ID.Hex(),
		Name: animalCategoryRes.Name,
	}, nil
}

func (m *mongoClient) AddAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error) {
	currTime := time.Now().UTC()
	res, err := m.mongoDbHandler.Collection(animalCategoryCollection).InsertOne(ctx, animalCategory{
		Name:      name,
		CreatedOn: currTime,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, &dbErr.ConflictError{Key: "name", Err: err}
		}
		return nil, err
	}

	return &genRouter.AnimalCategoryJSONResponse{
		Id:   res.InsertedID.(bson.ObjectID).Hex(),
		Name: name,
	}, nil
}

func (m *mongoClient) UpdateAnimalCategory(ctx context.Context, id,
	name string) (*genRouter.AnimalCategoryJSONResponse, error) {
	bsonID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, dbErr.ErrInvalidID
	}

	currTime := time.Now().UTC()
	res := m.mongoDbHandler.Collection(animalCategoryCollection).FindOneAndUpdate(ctx,
		bson.M{iDField: bsonID},
		bson.M{setOperator: bson.M{nameField: name, updatedOnField: currTime}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	err = res.Err()
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, &dbErr.ConflictError{Key: "name", Err: err}
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, dbErr.ErrNotFound
		}
		return nil, err
	}

	var animalCategoryRes animalCategory
	err = res.Decode(&animalCategoryRes)
	if err != nil {
		return nil, err
	}

	return &genRouter.AnimalCategoryJSONResponse{
		Id:   animalCategoryRes.ID.Hex(),
		Name: animalCategoryRes.Name,
	}, nil
}
