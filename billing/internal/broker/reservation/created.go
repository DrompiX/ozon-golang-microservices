package reservation

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/app"
	"gitlab.ozon.dev/DrompiX/homework-3/events"
)

type createdHandler struct {
	producer sarama.SyncProducer
	service  *app.BillingService
}

func NewCreatedHandler(s *app.BillingService, p sarama.SyncProducer) *createdHandler {
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

		var event events.ReservationCreatedEvent
		err := events.DecodeEvent(msg.Value, &event)
		if err != nil {
			log.Printf("Error in income data %v: %v", string(msg.Value), err)
			continue
		}
		log.Printf("Received event with reservation to pay: %+v", event)

		var topic string
		var e interface{}

		pId, err := h.service.ProcessPayment(event.OrderId, event.UserId, event.SellerId, event.TotalCost)
		if err != nil {
			topic = "payment_failed"
			e = &events.PaymentFailedEvent{
				OrderId: event.OrderId,
				ItemId:  event.ItemId,
				Reason:  err.Error(),
			}
			log.Printf("Payment failed, creating event: %+v", e)
		} else {
			topic = "payment_succeeded"
			e = &events.PaymentSucceededEvent{OrderId: event.OrderId, PaymentId: *pId}
			log.Printf("Payment succeeded, creating event: %+v", e)
		}
		resp, err := events.EncodeEvent(e)
		if err != nil {
			log.Printf("Could not marshal payment event %s", err)
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
