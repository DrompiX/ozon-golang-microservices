package order

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/events"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/app"
)

type createdHandler struct {
	producer sarama.SyncProducer
	service  *app.WarehouseService
}

func NewCreatedHandler(s *app.WarehouseService, p sarama.SyncProducer) *createdHandler {
	return &createdHandler{producer: p, service: s}
}

func (h *createdHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *createdHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *createdHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")

		var event events.OrderCreatedEvent
		err := events.DecodeEvent(msg.Value, &event)
		if err != nil {
			log.Printf("Error in income data %v: %v", string(msg.Value), err)
			continue
		}
		log.Printf("Received event to reserve an item: %+v", event)

		var topic string
		var e interface{}

		cost, err := h.service.ReserveItemQuantity(event.ItemId, 1)
		if err != nil {
			topic = "reservation_rejected"
			e = events.ReservationRejectedEvent{OrderId: event.OrderId, Reason: err.Error()}
			log.Printf("Reservation failed for %+v, creating event: %+v", event, e)
		} else {
			topic = "reservation_created"
			e = events.ReservationCreatedEvent{
				UserId: event.UserId,
				OrderId: event.OrderId,
				ItemId: event.ItemId,
				TotalCost: cost,
			}
			log.Printf("Reservation created for %+v: %+v", event, e)
		}
		resp, err := events.EncodeEvent(e)
		if err != nil {
			log.Printf("Could not marshal reservation event %+v: %s", e, err)
			continue
		}

		producerMsg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("%d", event.OrderId)),
			Value: sarama.StringEncoder(resp),
		}
		_, _, err = h.producer.SendMessage(producerMsg)
		if err != nil {
			log.Printf("Failed to send message %+v: %s", producerMsg, err)
		}
	}
	return nil
}
