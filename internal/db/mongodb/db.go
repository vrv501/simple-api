package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"

	"github.com/vrv501/simple-api/internal/constants"
)

const (
	dbName string = "pet-store"
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
