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

func (m *mongoClient) AddAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error) {
	currTime := time.Now().UTC()
	res, err := m.mongoDbHandler.Collection(animalCategoryCollection).InsertOne(ctx, animalCategory{
		Name:      name,
		CreatedOn: currTime,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, dbErr.ErrConflict
		}
		return nil, err
	}

	return &genRouter.AnimalCategoryJSONResponse{
		Id:        res.InsertedID.(bson.ObjectID).Hex(),
		Name:      name,
		CreatedAt: currTime,
	}, nil
}

func (m *mongoClient) UpdateAnimalCategory(ctx context.Context, id, name string) (*genRouter.AnimalCategoryJSONResponse, error) {
	bsonID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, dbErr.ErrInvalidId
	}

	currTime := time.Now().UTC()
	res := m.mongoDbHandler.Collection(animalCategoryCollection).FindOneAndUpdate(ctx,
		bson.M{idField: bsonID, deletedOnField: bson.M{notEqOperator: bson.Null{}}},
		bson.M{nameField: name, updatedOnField: currTime},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	err = res.Err()
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, dbErr.ErrConflict
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
		Id:        animalCategoryRes.ID.Hex(),
		Name:      animalCategoryRes.Name,
		CreatedAt: animalCategoryRes.CreatedOn,
		UpdatedAt: animalCategoryRes.UpdatedOn,
	}, nil
}

func (m *mongoClient) DeleteAnimalCategory(ctx context.Context, id string) error {
	bsonID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return dbErr.ErrInvalidId
	}

	err = m.mongoDbHandler.Collection(petsCollection).
		FindOne(
			ctx,
			bson.M{categoryIdField: bsonID, deletedOnField: bson.Null{}},
			options.FindOne().SetProjection(bson.M{idField: 1}),
		).Err()
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil { // No error if one document is found
		return dbErr.ErrForeignKeyConstraint
	}

	res, err := m.mongoDbHandler.Collection(animalCategoryCollection).
		UpdateOne(ctx, bson.M{idField: bsonID, deletedOnField: bson.Null{}},
			bson.M{setOperator: bson.M{deletedOnField: time.Now().UTC()}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return dbErr.ErrNotFound
	}

	return nil
}
