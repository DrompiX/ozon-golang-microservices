package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/app"
	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/broker"
	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/broker/reservation"
	"gitlab.ozon.dev/DrompiX/homework-3/billing/internal/infra"
)

func dbConnPool(ctx context.Context) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "billing",
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

	repo := infra.NewPostgresRepository(dbConnPool(ctx))
	producer := broker.InitKafkaProducer()
	warehouseService := app.New(repo, &infra.FakePaymentClient{})

	handlers := map[string]sarama.ConsumerGroupHandler{
		"reservation_created": reservation.NewCreatedHandler(warehouseService, producer),
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Starting consumers with handlers: %v", handlers)
	broker.RunConsumers(ctx, handlers, "billing_service")
	<-ctx.Done()
}
