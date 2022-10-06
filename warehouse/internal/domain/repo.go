package domain

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("item not found")
)

type Repository interface {
	CreateItem(ctx context.Context, o *Item) error
	UpdateItem(ctx context.Context, e *Item) error
	GetItem(ctx context.Context, itemId int64) (*Item, error)
}
