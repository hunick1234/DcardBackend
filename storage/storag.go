package storage

import "go.mongodb.org/mongo-driver/mongo"

type Storager interface {
	Disconnect() error
	GetCollection(string) (*mongo.Collection, error)
}
