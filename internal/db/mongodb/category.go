package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"

	dbErr "github.com/vrv501/simple-api/internal/db/errors"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func (m *mongoClient) FindAnimalCategory(ctx context.Context,
	name string) (*genRouter.AnimalCategoryJSONResponse, error) {
	// Try exact match
	res := m.mongoDbHandler.Collection(animalCategoryCollection).FindOne(ctx,
		bson.M{nameField: name},
	)
	err := res.Err()
	if err == nil {
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
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	// Try fuzzy search
	pipeline := []bson.M{
		// Stage1: Fuzzy match
		{
			"$search": bson.M{
				"index": "name_fuzzy_match",
				"text": bson.M{
					"query": name,
					"path":  nameField,
					"fuzzy": bson.M{
						"maxEdits":     2,
						"prefixLength": 1,
					},
				},
				// Sort in decreasing order based on relevance score
				"sort": bson.M{"score": bson.M{"$meta": "searchScore"}},
			},
		},
		// Stage2: Return the top most element
		{
			limitOperator: 1,
		},
	}

	// Atlas Search requires local read concern
	collection := m.mongoDbHandler.Collection(animalCategoryCollection,
		options.Collection().SetReadConcern(readconcern.Local()))
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var animalCategoryRes animalCategory
		err = cursor.Decode(&animalCategoryRes)
		if err != nil {
			return nil, err
		}
		return &genRouter.AnimalCategoryJSONResponse{
			Id:   animalCategoryRes.ID.Hex(),
			Name: animalCategoryRes.Name,
		}, nil
	}

	return nil, dbErr.ErrNotFound
}

func (m *mongoClient) AddAnimalCategory(ctx context.Context, name string) (*genRouter.AnimalCategoryJSONResponse, error) {
	res, err := m.mongoDbHandler.Collection(animalCategoryCollection).InsertOne(ctx, animalCategory{
		Name:      name,
		CreatedOn: time.Now().UTC(),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, dbErr.ErrConflict
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
		return nil, dbErr.ErrInvalidValue
	}

	err = m.mongoDbHandler.Collection(animalCategoryCollection).FindOneAndUpdate(ctx,
		bson.M{iDField: bsonID},
		bson.M{setOperator: bson.M{nameField: name, updatedOnField: time.Now().UTC()}},
	).Err()
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, dbErr.ErrConflict
		}
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, dbErr.ErrNotFound
		}
		return nil, err
	}

	return &genRouter.AnimalCategoryJSONResponse{
		Id:   id,
		Name: name,
	}, nil
}
