package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aidahputri/go-transaction/kafka"
)

func main() {
	topic := "test-topic"

	fmt.Println("Creating topic:", topic)
	kafka.CreateTopic(topic)

	fmt.Println("Kafka Producer started!")
	fmt.Println("Type messages to send (press Enter to send, type 'quit' to exit):")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Message: ")
		scanner.Scan()
		message := scanner.Text()

		if strings.ToLower(message) == "quit" {
			fmt.Println("Producer shutting down...")
			break
		}

		if message == "" {
			continue
		}

		err := kafka.Write(topic, message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		} else {
			fmt.Printf("Message sent: %s\n", message)
		}
	}
}
