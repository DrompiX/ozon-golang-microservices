package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/app"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/broker"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/broker/payment"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/broker/reservation"
	"gitlab.ozon.dev/DrompiX/homework-3/order/internal/infra"
	"gitlab.ozon.dev/DrompiX/homework-3/order/pb"
	"google.golang.org/grpc"
)

func runOrderGrpcServer(service pb.OrderServiceServer, listener net.Listener) error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterOrderServiceServer(grpcServer, service)

	log.Printf("Starting GRPC OrderService at %s", listener.Addr())
	return grpcServer.Serve(listener)
}

func dbConnPool(ctx context.Context) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "orders",
	)
	// TODO: add environment parsing
	// dsn := fmt.Sprintf(
	// 	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	conf.Db.Host, conf.Db.Port, conf.Db.User, conf.Db.Password, conf.Db.Dbname,
	// )
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}

func main() {
	ctx := context.Background()
	pool := dbConnPool(ctx)
	defer pool.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8888))
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

	repo := infra.NewPostgresRepository(dbConnPool(ctx))
	producer := broker.InitKafkaProducer()
	service := app.New(repo, producer)

	// TODO: Refactor handlers to accept service instead of repo
	// And add service method for each unique case
	handlers := map[string]sarama.ConsumerGroupHandler{
		"reservation_rejected": reservation.NewRejectedHandler(repo),
		"payment_succeeded": payment.NewSucceededHandler(repo),
		"payment_failed": payment.NewFailedHandler(repo),
	}

	log.Printf("Starting consumers with handlers: %v", handlers)
	broker.RunConsumers(ctx, handlers, "order_service")
	runOrderGrpcServer(service, listener)
}
