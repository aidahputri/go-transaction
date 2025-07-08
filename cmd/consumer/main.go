package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aidahputri/go-transaction/cmd"
	"github.com/aidahputri/go-transaction/kafka"
	"github.com/aidahputri/go-transaction/repo"
)

func main() {
	topic := "test-topic"
	dbConn := cmd.Connect()

	accountRepo := repo.NewAccount(dbConn)
	kafka.InitConsumerDependency(*accountRepo)

	fmt.Println("Kafka Consumer started!")
	fmt.Printf("Listening for messages on topic: %s\n", topic)
	fmt.Println("Press Ctrl+C to stop...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Consumer recovered from panic: %v", r)
			}
		}()

		kafka.Read(topic)
	}()

	<-sigChan
	fmt.Println("\nConsumer shutting down...")
}