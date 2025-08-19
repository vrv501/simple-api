package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
)

func (m *mongoClient) AddAnimalCategory(ctx context.Context, name string) (string, time.Time, error) {
	currTime := time.Now().UTC()
	res, err := m.mongoDbHandler.Collection(animalCategoryCollection).InsertOne(ctx, animalCategory{
		Name:      name,
		CreatedOn: currTime,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", currTime, dbErr.ErrConflict
		}
		return "", currTime, err
	}
	return res.InsertedID.(bson.ObjectID).Hex(), currTime, nil
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
