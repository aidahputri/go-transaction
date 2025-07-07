package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aidahputri/go-transaction/model"
	"github.com/aidahputri/go-transaction/repo"
	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumer(broker, topic string, accountRepo *repo.Account) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		GroupID:   "transfer-watchdog",
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})
	defer r.Close()

	log.Println("Kafka consumer started...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("error reading kafka message:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var tx model.Transaction
		if err := json.Unmarshal(m.Value, &tx); err != nil {
			log.Println("failed to unmarshal transaction:", err)
			continue
		}

		go handleTransactionMonitoring(tx, accountRepo)
	}
}

func handleTransactionMonitoring(tx model.Transaction, accountRepo *repo.Account) {
	ctx := context.Background()

	fromAcc, _ := accountRepo.Get(ctx, tx.FromAccount)
	toAcc, _ := accountRepo.Get(ctx, tx.ToAccount)

	if fromAcc.Blacklisted && !fromAcc.Underwatch {
		fromAcc.Underwatch = true
		accountRepo.Update(ctx, fromAcc)
		log.Printf("Account %s flagged under observation\n", fromAcc.AccountNumber)
	}
	if toAcc.Blacklisted && !toAcc.Underwatch {
		toAcc.Underwatch = true
		accountRepo.Update(ctx, toAcc)
		log.Printf("Account %s flagged under observation\n", toAcc.AccountNumber)
	}
}