package payment

import (
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/events"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/app"
)

type failedHandler struct{
	producer sarama.SyncProducer
	service  *app.WarehouseService
}

func NewFailedHandler(s *app.WarehouseService, p sarama.SyncProducer) *failedHandler {
	return &failedHandler{producer: p, service: s}
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

		var event events.PaymentFailedEvent
		err := events.DecodeEvent(msg.Value, &event)
		if err != nil {
			log.Printf("Error in income data %v: %v", string(msg.Value), err)
			continue
		}
		log.Printf("Received event for a failed payment: %+v", event)

		// TODO: hardcoded quantity for now
		err = h.service.AddItemQuantity(event.ItemId, 1)
		if err != nil {
			log.Printf("Failed to execute compensation for reservation: %s", err)
			continue
		}

		log.Printf("Compensated warehouse with for event: %+v", event)
		session.MarkMessage(msg, "")
	}
	return nil
}
