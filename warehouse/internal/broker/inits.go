package broker

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

func InitKafkaProducer() sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	// TODO: move kafka address to configuration
	producer, err := sarama.NewSyncProducer([]string{"localhost:9095"}, cfg)
	if err != nil {
		log.Fatalf("Kafka error: %s", err)
	}
	return producer
}

func RunConsumers(ctx context.Context, handlers map[string]sarama.ConsumerGroupHandler, groupSuffix string) {
	consumerGrpoups := make(map[string]*sarama.ConsumerGroup)
	for topic := range handlers {
		consumerGrpoups[topic] = initGroup(topic + "-" + groupSuffix)
	}

	for topic, group := range consumerGrpoups {
		go func(topic string, group *sarama.ConsumerGroup) {
			defer func() {
				if r := recover(); r != nil {
					log.Panicf("%s consumer recovery failure: %s", topic, r)
				}
			}()

			for {
				err := (*group).Consume(ctx, []string{topic}, handlers[topic])
				if err != nil {
					log.Printf("%s consumer error: %s", topic, err)
				}
			}
		}(topic, group)
	}
}

func initGroup(groupName string) *sarama.ConsumerGroup {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Return.Errors = true

	// TODO: move kafka address to configuration
	group, err := sarama.NewConsumerGroup([]string{"localhost:9095"}, groupName, cfg)
	if err != nil {
		log.Printf("consumer group creation error: %s", err)
		return nil
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Panicf("group recovery failure: %s", r)
			}
		}()

		for err := range group.Errors() {
			log.Printf("consumer group error: %s", err)
		}
	}()

	return &group
}
