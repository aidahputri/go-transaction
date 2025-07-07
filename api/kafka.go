package api

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func InitKafka(brokerAddress string, topic string) {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka writer initialized")
}

func PublishTransferMessage(message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := kafkaWriter.WriteMessages(ctx, kafka.Message{
		Value: message,
	})
	return err
}