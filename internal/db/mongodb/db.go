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
	// ENV vars
	mongoURIEnvVar string = "MONGO_URI"       // Should be hostname:port
	mongoApplyURI  string = "MONGO_APPLY_URI" // Should be mongodb+srv://{{ rest of the URI }}

	defaultAPIUser       string = "apiUser"
	defaultMongoPassword string = "mongo"
	defaultMongoURI      string = "localhost:27017"

	dbName                   string = "shop"
	animalCategoryCollection string = "animal_categories"
	petsCollection           string = "pets"
	ordersCollection         string = "orders"
	usersCollection          string = "users"
	fsFilesCollection        string = "fs.files"
	fsChunksCollection       string = "fs.chunks"

	iDField         string = "_id"
	nameField       string = "name"
	usernameField   string = "username"
	passwordField   string = "password"
	categoryIDField string = "category_id"
	updatedOnField  string = "updated_on"
	deletedOnField  string = "deleted_on"

	setOperator   string = "$set"
	notEqOperator string = "$ne"
	limitOperator string = "$limit"
)

type mongoClient struct {
	client         *mongo.Client
	mongoDbHandler *mongo.Database
}

// Note: Mongo By default stores date in UTC timezone only
//
//revive:disable:unexported-return
func NewInstance(ctx context.Context) *mongoClient {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	serverAPI.SetStrict(false) // Atlas Search requires apiStrict: false

	// Default Connect timeout, Server Selection Timeout: 30s
	// Max PoolConnSize: 100
	// Retries are default activated for all sorts of operations
	c := options.Client()
	if applyURI := os.Getenv(mongoApplyURI); applyURI != "" {
		c = c.ApplyURI(applyURI)
	} else {
		mongoURI := os.Getenv(mongoURIEnvVar)
		if mongoURI == "" {
			mongoURI = defaultMongoURI
		}
		c.SetHosts([]string{mongoURI})

		username := os.Getenv(constants.DBUsername)
		if username == "" {
			username = defaultAPIUser
		}
		pswd := os.Getenv(constants.DBPassword)
		if pswd == "" {
			pswd = defaultMongoPassword
		}
		c.SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-256",
			AuthSource:    dbName,
			Username:      username,
			Password:      pswd,
		})
	}
	c.SetServerAPIOptions(serverAPI)
	c.SetAppName("pet-store-api-server")
	c.SetTimeout(5 * time.Minute) // Query timeout
	c.SetReadConcern(readconcern.Majority())
	c.SetReadPreference(readpref.PrimaryPreferred())
	c.SetWriteConcern(writeconcern.Majority())

	client, err := mongo.Connect(c.SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(fmt.Sprintf("Failed to create mongodb client %v", err))
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to ping db %v", err))
	}

	return &mongoClient{
		client:         client,
		mongoDbHandler: client.Database(dbName),
	}
}

func (m *mongoClient) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

type animalCategory struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`       // "bsonType": "string"
	CreatedOn time.Time     `bson:"created_on"` // "bsonType": "date"
	UpdatedOn *time.Time    `bson:"updated_on"` // "bsonType": ["date", "null"]
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
	Address     string        `bson:"address"`      // "bsonType": "string"
	PhoneNumber string        `bson:"phone_number"` // "bsonType": "string"
	CreatedOn   time.Time     `bson:"created_on"`   // "bsonType": "date"
	UpdatedOn   *time.Time    `bson:"updated_on"`   // "bsonType": ["date", "null"]
	DeletedOn   *time.Time    `bson:"deleted_on"`   // "bsonType": ["date", "null"]
}

type order struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	UserID        bson.ObjectID `bson:"user_id"`        // "bsonType": "objectId"
	PetID         bson.ObjectID `bson:"pet_id"`         // "bsonType": "objectId"
	ShippedDate   *time.Time    `bson:"shipped_date"`   // "bsonType": ["date", "null"]
	DeliveredDate *time.Time    `bson:"delivered_date"` // "bsonType": ["date", "null"]
	Status        string        `bson:"status"`         // "bsonType": "string"
	CreatedOn     time.Time     `bson:"created_on"`     // "bsonType": "date"
	UpdatedOn     *time.Time    `bson:"updated_on"`     // "bsonType": ["date", "null"]
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
