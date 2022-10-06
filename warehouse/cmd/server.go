package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/app"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/broker"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/broker/order"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/broker/payment"
	"gitlab.ozon.dev/DrompiX/homework-3/warehouse/internal/infra"
)

func dbConnPool(ctx context.Context) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "warehouse",
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
	warehouseService := app.New(repo, producer)

	handlers := map[string]sarama.ConsumerGroupHandler{
		"order_created": order.NewCreatedHandler(warehouseService, producer),
		"payment_failed": payment.NewFailedHandler(warehouseService, producer),
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Printf("Starting consumers with handlers: %v", handlers)
	broker.RunConsumers(ctx, handlers, "warehouse-service")
	<-ctx.Done()
}
