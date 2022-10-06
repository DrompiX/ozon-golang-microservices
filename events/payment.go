package events

type PaymentSucceededEvent struct {
	OrderId   int64 `json:"order_id"`
	PaymentId int64 `json:"payment_id"`
}

type PaymentFailedEvent struct {
	OrderId int64  `json:"order_id"`
	ItemId  int64  `json:"item_id"`
	Reason  string `json:"reason"`
}
