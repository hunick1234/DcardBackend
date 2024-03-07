package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Storager[T any, Q any] interface {
	Store(*T) error
	FindByFilter(context.Context, *Q) (*[]T, error)
	Aggregate(context.Context, Filter, any) error
}

type Filter interface {
	Pipeline() mongo.Pipeline
}
