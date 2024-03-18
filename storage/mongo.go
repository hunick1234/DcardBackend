package storage

import (
	"context"
	"time"

	"github.com/hunick1234/DcardBackend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoAddress = "mongodb://localhost:27017"
)

type MongoDB struct {
	DB     *mongo.Database
	Client *mongo.Client
}

func init() {

	var _ Storager = (*MongoDB)(nil)
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		Client: nil,
		DB:     nil,
	}
}

func NewMongoConn(cfg *config.MongoCfg) (Storager, error) {
	opts := options.Client().ApplyURI(cfg.URI).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize)

	client, db, err := connectToMongo(opts, cfg.DB)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		Client: client,
		DB:     db,
	}, nil

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

// implement Storager interface
func (m *MongoDB) Disconnect() error {
	m.Client.Disconnect(context.Background())
	return nil
}

func (m *MongoDB) GetCollection(collection string) (*mongo.Collection, error) {
	coll := m.DB.Collection(collection)
	return coll, nil
}

func (m *MongoDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) GetDBName() string {
	return m.DB.Name()
}
