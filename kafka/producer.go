package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type TransferEvent struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

var topic = "transfer-topic"

func PublishTransfer(event TransferEvent) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	bytes, _ := json.Marshal(event)
	err := writer.WriteMessages(context.Background(),
		kafka.Message{Value: bytes},
	)
	if err != nil {
		log.Println("Failed to publish transfer:", err)
	}
	writer.Close()
}
