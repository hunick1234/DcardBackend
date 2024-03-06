package storage

import "go.mongodb.org/mongo-driver/mongo"

type Storager interface {
	GetDBName() string
	Ping() error
	Disconnect() error
	GetCollection(string) (*mongo.Collection, error)
}
