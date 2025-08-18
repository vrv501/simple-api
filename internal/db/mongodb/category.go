package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func (m *mongoClient) AddAnimalCategory(ctx context.Context, name string) (string, error) {
	res, err := m.mongoDbHandler.Collection(animalCategoryCollection).InsertOne(ctx, animalCategory{
		Name:      name,
		CreatedOn: time.Now().UTC(),
	})
	if err != nil {
		return "", err
	}
	return res.InsertedID.(bson.ObjectID).String(), nil
}
