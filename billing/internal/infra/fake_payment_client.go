package infra

import (
	"context"
	"fmt"
	"math/rand"

	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/domain"
)

// Results in successful payment with probability 0.67
type FakePaymentClient struct{}

func (c *FakePaymentClient) Transfer(ctx context.Context, from, to, amount int64) error {
	if rand.Float64() < 0.67 {
		return nil
	}

	return fmt.Errorf("payment failed from %d to %d for amount %d", from, to, amount)
}

var _ domain.PaymentClient = (*FakePaymentClient)(nil)
