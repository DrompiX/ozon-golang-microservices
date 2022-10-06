package domain

type Payment struct {
	ID         int64 `json:"id"`
	OrderId    int64 `json:"order_id"`
	FromUserId int64 `json:"from_user_id"`
	ToUserId   int64 `json:"to_user_id"`
	Amount     int64 `json:"amount"`
}

func NewPayment(orderId, fromUserId, toUserId, amount int64) *Payment {
	return &Payment{
		OrderId: orderId,
		FromUserId: fromUserId,
		ToUserId: toUserId,
		Amount: amount,
	}
}
