package events

type OrderCreatedEvent struct {
	UserId  int64 `json:"user_id"`
	OrderId int64 `json:"order_id"`
	ItemId  int64 `json:"item_id"`
}
