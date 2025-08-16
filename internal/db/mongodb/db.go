package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"

	"github.com/vrv501/simple-api/internal/constants"
)

const (
	dbName string = "shop"
)

type MongoClient struct {
	client         *mongo.Client
	mongoDbHandler *mongo.Database
}

func NewInstance(ctx context.Context) *MongoClient {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	serverAPI.SetStrict(true)

	username := os.Getenv(constants.DBUsername)
	if username == "" {
		username = "apiUser"
	}
	pswd := os.Getenv(constants.DBPassword)
	if pswd == "" {
		pswd = "mongo"
	}

	// Default Connect timeout, Server Selection Timeout: 30s
	// Max PoolConnSize: 100
	// Retries are default activated for all sorts of operations
	c := options.Client()
	c.SetServerAPIOptions(serverAPI)
	c.SetAppName("pet-store-api-server")
	c.SetTimeout(5 * time.Minute) // Query timeout
	c.SetReadConcern(readconcern.Majority())
	c.SetReadPreference(readpref.PrimaryPreferred())
	c.SetWriteConcern(writeconcern.Majority())
	c.SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    dbName,
		Username:      username,
		Password:      pswd,
	})

	client, err := mongo.Connect(c.SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(fmt.Sprintf("Failed to create mongodb client %v", err))
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to ping db %v", err))
	}

	return &MongoClient{
		client:         client,
		mongoDbHandler: client.Database(dbName),
	}
}

func (m *MongoClient) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

type animalCategory struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`       // "bsonType": "string"
	CreatedOn time.Time     `bson:"created_on"` // "bsonType": "date"
	UpdatedOn *time.Time    `bson:"updated_on"` // "bsonType": ["date", "null"]
	DeletedOn *time.Time    `bson:"deleted_on"` // "bsonType": ["date", "null"]
}

type pet struct {
	ID         bson.ObjectID   `bson:"_id,omitempty"`
	Name       string          `bson:"name"`        // "bsonType": "string"
	CategoryID bson.ObjectID   `bson:"category_id"` // "bsonType": "objectId"
	PhotoURIs  []string        `bson:"photo_uris"`  // "bsonType": ["array", "null"]
	Price      bson.Decimal128 `bson:"price"`       // "bsonType": "decimal"
	Status     string          `bson:"status"`      // "bsonType": "string"
	Tags       []string        `bson:"tags"`        // "bsonType": ["array", "null"]
	CreatedOn  time.Time       `bson:"created_on"`  // "bsonType": "date"
	UpdatedOn  *time.Time      `bson:"updated_on"`  // "bsonType": ["date", "null"]
	DeletedOn  *time.Time      `bson:"deleted_on"`  // "bsonType": ["date", "null"]
}

type user struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Username    string        `bson:"username"`     // "bsonType": "string"
	FullName    string        `bson:"full_name"`    // "bsonType": "string"
	Email       string        `bson:"email"`        // "bsonType": "string"
	Password    string        `bson:"password"`     // "bsonType": "string"
	PhoneNumber string        `bson:"phone_number"` // "bsonType": "string"
	CreatedOn   time.Time     `bson:"created_on"`   // "bsonType": "date"
	UpdatedOn   *time.Time    `bson:"updated_on"`   // "bsonType": ["date", "null"]
	DeletedOn   *time.Time    `bson:"deleted_on"`   // "bsonType": ["date", "null"]
}

type order struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	UserID      bson.ObjectID `bson:"user_id"`      // "bsonType": "objectId"
	PetID       bson.ObjectID `bson:"pet_id"`       // "bsonType": "objectId"
	ShippedDate *time.Time    `bson:"shipped_date"` // "bsonType": ["date", "null"]
	Status      string        `bson:"status"`       // "bsonType": "string"
	CreatedOn   time.Time     `bson:"created_on"`   // "bsonType": "date"
	UpdatedOn   *time.Time    `bson:"updated_on"`   // "bsonType": ["date", "null"]
}

type fsFile struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	Length     int64         `bson:"length"`     // "bsonType": "long"
	ChunkSize  int32         `bson:"chunkSize"`  // "bsonType": "int"
	UploadDate time.Time     `bson:"uploadDate"` // "bsonType": "date"
	Filename   string        `bson:"filename"`   // "bsonType": "string"
	Metadata   bson.M        `bson:"metadata"`   // "bsonType": "object"
}

type fsChunk struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	FilesID bson.ObjectID `bson:"files_id"` // "bsonType": "objectId"
	Length  int32         `bson:"n"`        // "bsonType": "int"
	Data    []byte        `bson:"data"`     // "bsonType": "binData"
}
