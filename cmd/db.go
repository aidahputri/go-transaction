package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func init() {
	dbCmd := cobra.Command{
		Use:   "db",
		Short: "database example",
		Run: func(cmd *cobra.Command, args []string) {
			Connect()
		},
	}

	rootCmd.AddCommand(&dbCmd)
}

func Connect() *sql.DB {
	cs := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost:5432",
		Path:     "/go-transaction",
		RawQuery: "sslmode=disable",
	}
	fmt.Println(cs.String())

	db, err := sql.Open("postgres", cs.String())
	if err != nil {
		log.Fatal("unable to open db:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("unable to ping db:", err)
	}

	fmt.Println("connected to db")
	return db
}