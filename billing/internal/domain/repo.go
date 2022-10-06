package domain

import (
	"context"
)

type Repository interface {
	CreatePayment(ctx context.Context, o *Payment) error
}
