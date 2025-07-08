package cmd

import (
	"fmt"
	// "log"
	// "net/http"

	// "github.com/gin-gonic/gin"
	"github.com/aidahputri/go-transaction/api"
	"github.com/aidahputri/go-transaction/repo"
	// "github.com/aidahputri/go-transaction/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	serverCmd := cobra.Command{
		Use:   "server",
		Short: "Start API server",
		Run: func(cmd *cobra.Command, args []string) {
			// konek db
			dbConn := Connect()
			defer dbConn.Close()

			// init repo
			accountRepo := repo.NewAccount(dbConn)
			transactionRepo := repo.NewTransaction(dbConn)
			handler := api.NewHandler(accountRepo, transactionRepo)

			// init kafka
			// init kafka producer
			// api.InitKafka("localhost:9092", "transfer-topic")

			// run kafka consumer di goroutine
			// go kafka.StartKafkaConsumer("localhost:9092", "transfer-topic", accountRepo)

			// init router
			router := api.InitRouter(handler)
			addr := viper.GetString("server.listen_addr")
			fmt.Println("Server listen on:", addr)
			router.Run(addr)
		},
	}

	rootCmd.AddCommand(&serverCmd)
}
