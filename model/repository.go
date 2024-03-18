package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T any, Q any] interface {
	Store(context.Context, *T) error
	FindByFilter(context.Context, *Q) (*[]T, error)
	Aggregate(context.Context, Filter, any) error
}

type Filter interface {
	Pipeline() mongo.Pipeline
}
