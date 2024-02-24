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
	Options        *options.ClientOptions
	CollectionName string
	DBNmae         string
}

var dbConnections map[string]Storager
var mu sync.Mutex

func init() {
	dbConnections = make(map[string]Storager, 10)
	var _ Storager = (*MongoDB)(nil)
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		Client:         nil,
		DB:             nil,
		CollectionName: "",
		DBNmae:         "",
		Options:        options.Client().ApplyURI(MongoAddress),
	}
}

func connect(opts *options.ClientOptions, dbName string) (Storager, error) {
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
		Options:        opts,
		DBNmae:         dbName,
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

func (m *MongoDB) Connect() (Storager, error) {
	storage, err := connect(m.Options, m.DBNmae)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

// implement Storager interface
func (m *MongoDB) Disconnect() error {
	m.Client.Disconnect(context.Background())
	return nil
}

func (m *MongoDB) GetCollection(collection string) (*mongo.Collection, error) {
	coll := m.DB.Collection(collection)
	return coll, nil
}
