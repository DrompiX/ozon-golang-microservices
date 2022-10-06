package app

import (
	"context"

	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/domain"
)

type BillingService struct {
	repo          domain.Repository
	client        domain.PaymentClient
}

func New(repo domain.Repository, c domain.PaymentClient) *BillingService {
	return &BillingService{repo: repo, client: c}
}

func (s *BillingService) ProcessPayment(orderId, fromUser, toUser, amount int64) (*int64, error) {
	ctx := context.Background()
	err := s.client.Transfer(ctx, fromUser, toUser, amount)
	if err != nil {
		return nil, err
	}

	p := domain.NewPayment(orderId, fromUser, toUser, amount)

	// TODO: think about how to make sure that this transaction will execute
	err = s.repo.CreatePayment(ctx, p)
	if err != nil {
		return nil, err
	}

	return &p.ID, nil
}
