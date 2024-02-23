package storage

import "go.mongodb.org/mongo-driver/mongo/options"

type Storager interface {
	Connect(opts *options.ClientOptions, dbName string) (*MongoDB, error)
}
