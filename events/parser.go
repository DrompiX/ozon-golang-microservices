package events

import (
	"bytes"
	"encoding/json"
)

type Event interface {
	*OrderCreatedEvent | 
	*PaymentSucceededEvent | *PaymentFailedEvent |
	*ReservationCreatedEvent | *ReservationRejectedEvent
}

func DecodeEvent[T Event](data []byte, event T) error {
	jdec := json.NewDecoder(bytes.NewReader(data))
	jdec.DisallowUnknownFields()
	return jdec.Decode(event)
}

func EncodeEvent(event interface{}) ([]byte, error) {
	return json.Marshal(event)
}