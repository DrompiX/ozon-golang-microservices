package payment

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/events"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/domain"
)

type failedHandler struct {
	repo domain.Repository
}

func NewFailedHandler(repo domain.Repository) *failedHandler {
	return &failedHandler{repo}
}

func (h *failedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *failedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *failedHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")

		var event events.PaymentFailedEvent
		err := events.DecodeEvent(msg.Value, &event)
		if err != nil {
			log.Printf("Error in income data %v: %v", string(msg.Value), err)
			continue
		}
		log.Printf("Order payment failed: %+v", event)

		// TODO: add transaction
		ctx := context.Background()
		order, err := h.repo.GetOrder(ctx, event.OrderId)
		if err != nil {
			log.Printf("Error in payment fail db processing: %s", err)
			continue
		}

		order.MarkFailed()
		h.repo.UpdateOrder(ctx, order)
	}
	return nil
}
