package model

import (
	"context"
)

type Storage[T any, Q any] interface {
	Store(*T) error
	FindByFilter(context.Context, Q) ([]*T, error)
}
