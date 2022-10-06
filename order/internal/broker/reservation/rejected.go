package reservation

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/events"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/domain"
)

type rejectedHandler struct{
	repo domain.Repository
}

func NewRejectedHandler(repo domain.Repository) *rejectedHandler {
	return &rejectedHandler{repo}
}

func (h *rejectedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *rejectedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *rejectedHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")

		var event events.ReservationRejectedEvent
		err := events.DecodeEvent(msg.Value, &event)
		if err != nil {
			log.Printf("Error in income data %v: %v", string(msg.Value), err)
			continue
		}
		log.Printf("Order was rejected: %+v", event)

		// TODO: add transaction
		ctx := context.Background()
		order, err := h.repo.GetOrder(ctx, event.OrderId)
		if err != nil {
			log.Printf("Error in reservation rejection: %s", err)
			continue
		}

		order.MarkFailed()
		h.repo.UpdateOrder(ctx, order)
	}
	return nil
}
