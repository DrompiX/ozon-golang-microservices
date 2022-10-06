package domain

import "time"

type OrderStatus string

var (
	OrderPending OrderStatus = "PENDING"
	OrderFailed  OrderStatus = "FAILED"
	OrderCreated OrderStatus = "CREATED"
)

type Order struct {
	ID        int64       `json:"id"`
	ItemId    int64       `json:"item_id"`
	UserId    int64       `json:"user_id"`
	PaymentId *int64      `json:"payment_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewOrder(itemId, userId int64) *Order {
	return &Order{
		ItemId: itemId,
		UserId: userId,
		Status: OrderPending,
	}
}

func (o *Order) MarkPayed(paymentId int64) {
	o.PaymentId = &paymentId
	o.Status = OrderCreated
}

func (o *Order) MarkFailed() {
	o.Status = OrderFailed
}
