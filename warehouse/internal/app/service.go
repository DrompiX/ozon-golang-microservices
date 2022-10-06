package app

import (
	"context"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/domain"
)

type WarehouseService struct {
	repo          domain.Repository
	kafkaProducer sarama.SyncProducer
}

func New(repo domain.Repository, p sarama.SyncProducer) *WarehouseService {
	return &WarehouseService{repo: repo, kafkaProducer: p}
}

// TODO: wrap the logic inside a transaction to avoid inconsistency
func (s *WarehouseService) ReserveItemQuantity(itemId, quantity int64) (cost int64, err error) {
	ctx := context.Background()
	item, err := s.repo.GetItem(ctx, itemId)
	if err != nil {
		return
	}

	err = item.SubtractQuantity(quantity)
	if err != nil {
		return
	}

	if err = s.repo.UpdateItem(ctx, item); err == nil {
		cost := quantity * item.Price
		return cost, nil
	}
	return
}

// TODO: wrap the logic inside a transaction to avoid inconsistency
func (s *WarehouseService) AddItemQuantity(itemId, quantity int64) error {
	ctx := context.Background()
	item, err := s.repo.GetItem(ctx, itemId)
	if err != nil {
		return err
	}

	item.AddQuantity(quantity)
	return s.repo.UpdateItem(ctx, item)
}
