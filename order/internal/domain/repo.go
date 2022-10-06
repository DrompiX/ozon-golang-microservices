package domain

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("item not found")
)

type Repository interface {
	CreateOrder(ctx context.Context, o *Order) error
	UpdateOrder(ctx context.Context, e *Order) error
	GetOrder(ctx context.Context, orderId int64) (*Order, error)
}
