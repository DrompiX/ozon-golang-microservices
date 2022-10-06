package events

type ReservationCreatedEvent struct {
	UserId    int64 `json:"user_id"`
	SellerId  int64 `json:"seller_id"`
	OrderId   int64 `json:"order_id"`
	ItemId    int64 `json:"item_id"`
	TotalCost int64 `json:"total_cost"`
}

type ReservationRejectedEvent struct {
	OrderId int64  `json:"order_id"`
	Reason  string `json:"reason"`
}
