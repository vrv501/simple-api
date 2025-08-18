package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
)

const (
	animalCategoryCollection string = "animal_categories"
)

type animalCategory struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`       // "bsonType": "string"
	CreatedOn time.Time     `bson:"created_on"` // "bsonType": "date"
	UpdatedOn *time.Time    `bson:"updated_on"` // "bsonType": ["date", "null"]
	DeletedOn *time.Time    `bson:"deleted_on"` // "bsonType": ["date", "null"]
}

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

	res, err := m.mongoDbHandler.Collection(animalCategoryCollection).
		UpdateOne(ctx, bson.M{idField: bsonID, deletedOnField: bson.Null{}}, bson.M{
			setOperator: bson.M{deletedOnField: time.Now().UTC()},
		})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return dbErr.ErrNotFound
	}

	return nil
}
