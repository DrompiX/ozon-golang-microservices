package domain

import (
	"fmt"
)

type Item struct {
	ID       int64 `json:"id"`
	SellerId int64 `json:"seller_id"`
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
}

func NewItem(sellerId, price int64) *Item {
	return &Item{
		SellerId: sellerId,
		Price:    price,
	}
}

func (i *Item) AddQuantity(moreQuantity int64) {
	i.Quantity += moreQuantity	
}

func (i *Item) SubtractQuantity(quantity int64) error {
	if i.Quantity - quantity < 0 {
		return fmt.Errorf("not enough quantity (requested %d, available %d)", quantity, i.Quantity)
	}
	i.Quantity -= quantity
	return nil
}
