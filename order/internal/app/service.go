package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/DrompiX/homework-3/events"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/domain"
	"gitlab.ozon.dev/DrompiX/homework-3/order/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderService struct {
	repo          domain.Repository
	kafkaProducer sarama.SyncProducer
	pb.UnimplementedOrderServiceServer
}

func New(repo domain.Repository, p sarama.SyncProducer) *orderService {
	return &orderService{repo: repo, kafkaProducer: p}
}

func (s *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Create and store new order
	order := domain.NewOrder(req.ItemId, req.UserId)
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Prepare event about order creation
	e := events.OrderCreatedEvent{
		UserId: req.UserId,
		OrderId: order.ID,
		ItemId: order.ItemId,
	}
	eventb, err := events.EncodeEvent(&e)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Send event that new order was created
	producerMsg := &sarama.ProducerMessage{
		Topic: "order_created",
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", req.UserId)),
		Value: sarama.StringEncoder(eventb),
	}
	_, _, err = s.kafkaProducer.SendMessage(producerMsg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateOrderResponse{}, nil
}
