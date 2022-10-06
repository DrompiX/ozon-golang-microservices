package payment

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/events"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/domain"
)

type succeededHandler struct {
	repo domain.Repository
}

func NewSucceededHandler(repo domain.Repository) *succeededHandler {
	return &succeededHandler{repo}
}

func (h *succeededHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *succeededHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *succeededHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")

		var event events.PaymentSucceededEvent
		err := events.DecodeEvent(msg.Value, &event)
		if err != nil {
			log.Printf("Error in income data %v: %v", string(msg.Value), err)
			continue
		}
		log.Printf("Order was successfully created and payed: %+v", event)

		// TODO: add transaction
		ctx := context.Background()
		order, err := h.repo.GetOrder(ctx, event.OrderId)
		if err != nil {
			log.Printf("Error in payment success processing: %s", err)
			continue
		}

		order.MarkPayed(event.PaymentId)
		h.repo.UpdateOrder(ctx, order)
	}
	return nil
}
