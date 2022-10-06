package domain

import "context"

type PaymentClient interface {
	Transfer(ctx context.Context, from, to, amount int64) error
}