package storage

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoAddress = "mongodb://localhost:27017"
)

type MongoDB struct {
	DB             *mongo.Database
	Client         *mongo.Client
	CollectionName string
}

var dbConnections map[string]*MongoDB
var mu sync.Mutex

func init() {
	dbConnections = make(map[string]*MongoDB, 10)
	//var _ Storage = (*MongoDB)(nil)
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		Client:         nil,
		DB:             nil,
		CollectionName: "",
	}
}

func Connect(opts *options.ClientOptions, dbName string) (*MongoDB, error) {
	var err error
	mu.Lock()
	defer mu.Unlock()
	//connect to db at first time
	if condition, ok := dbConnections[dbName]; ok {
		return condition, nil
	}

	client, db, err := connectToMongo(opts, dbName)
	if err != nil {
		return nil, err
	}
	dbConnections[dbName] = &MongoDB{
		Client:         client,
		DB:             db,
		CollectionName: "",
	}

	return dbConnections[dbName], nil
}

func connectToMongo(opts *options.ClientOptions, dbName string) (*mongo.Client, *mongo.Database, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}
	db := client.Database(dbName)

	return client, db, nil
}

func (m *MongoDB) Collection() error {
	return nil
}
