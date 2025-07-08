package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/aidahputri/go-transaction/model"
	"github.com/aidahputri/go-transaction/repo"
	"github.com/segmentio/kafka-go"
)

var (
	accountRepo repo.Account
)

func InitConsumerDependency(ar repo.Account) {
	accountRepo = ar
}

func Write(topic, message string) error {
	w := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    topic,
		Balancer: &kafka.Hash{}, // strategi pengiriman pesan ke partition
	}

	// write 1 pesan ke kafka dengan value berupa byte string
	return w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
}

func Read(topic string) {
	// bisa ngeset either groupid atau partitions, kalau kita ga ngeset itu kita bisa ngebaca semua data yang ada di kafka --> penting untuk set group id (disarankan set group id)
	// kalau reader ke restart dan ke start lagi dia ga akan baca data yang udah pernah dibaca sebelumnya
	// untuk satu reader biasanya dia ngebaca dari satu partition aja (satu partition bisa punya dua reader)
	// Partitions: 0,
	// selama group id nya beda, dia bisa ngakses berbarengan
	// kalau gamasukkin group id --> itu akan dapet semua message dari awal
	// nentuin jumlah partisi itu tergantung volume datanya seberapa banyak
	// kalau group id sama dan partisi sama --> akan masuk ke salah satu doang
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		GroupID:  "reader",
		MaxBytes: 10e6,
	})

	// loop for read message
	for {
		m, error := r.ReadMessage(context.Background())
		if error != nil {
			break
		}

		log.Printf("message at topic/partitions/offset %s/%d/%d: %s\n",
			m.Topic, m.Partition, m.Offset, string(m.Value))

		var tx model.Transaction
		if err := json.Unmarshal(m.Value, &tx); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			continue
		}

		fmt.Printf("Received Transaction: %+v\n", tx)

		// Cek dan update status underwatch
		ctx := context.Background()

		// Periksa sender
		fromAcc, err := accountRepo.Get(ctx, tx.FromAccount)
		if err == nil && fromAcc.Blacklisted {
			fromAcc.Underwatch = true
			if _, err := accountRepo.Update(ctx, fromAcc); err != nil {
				log.Printf("failed to update sender: %v", err)
			}
		}

		// Periksa receiver
		toAcc, err := accountRepo.Get(ctx, tx.ToAccount)
		if err == nil && toAcc.Blacklisted {
			toAcc.Underwatch = true
			if _, err := accountRepo.Update(ctx, toAcc); err != nil {
				log.Printf("failed to update receiver: %v", err)
			}
		}
	}

	if err := r.Close(); r != nil {
		log.Fatal("failed to close reader:", err)
	}

	log.Println("No new message")
}

func CreateTopic(topic string) {
	controllerConn := LeaderConnection() // manggil koneksi ke broker controller/leader kafka yang bertanggung jawab membuat topic
	defer controllerConn.Close()

	topicConfig := []kafka.TopicConfig{
		{
			Topic: topic,
			NumPartitions: 2,
			ReplicationFactor: 1,
		},
	}

	err := controllerConn.CreateTopics(topicConfig...)
	if err != nil {
		panic(err.Error())
	}
} 
// menghubungkan ke broker kafka yang menjadi controller (pemimpin untuk cluster kafka)
func LeaderConnection() *kafka.Conn {
	conn, err := kafka.Dial("tcp", "localhost:9092") // membuat koneksi ke kafka broker
	if err != nil {
		// Biasanya ketika kita dial kita belum tentu bisa nulis ke api yang kita connect
		log.Fatal("Dial error:", err)
	}

	// controller ngebalikin host dan port untuk kafka yang bertanggung jawab menerima request dan disebar prosesnya
	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("Dial error:", err.Error())
	}

	defer conn.Close()

	controller, err = conn.Controller()
	if err != nil {
		log.Fatal("Kafka connection, get controller:", err.Error())
	}

	log.Print("Leader:", controller.Host, controller.Port)

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatal("Kafka connection, dial leader:", err.Error())
	}

	return controllerConn
}
