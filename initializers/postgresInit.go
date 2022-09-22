package initializers

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

func PostgresCreateConnection() {
	log.Println("Attempting to establish connection to the PostgreSQL database...")
	conn, conErr := pgx.Connect(context.Background(), os.Getenv("POSTGRES_DB_URL"))
	if conErr != nil {
		log.Printf("Unable to connect to PostgreSQL database: %v\n", conErr)
		os.Exit(1)
	}
	if conErr == nil {
		log.Printf("Connected to PostgreSQL!\n")
	}
	defer conn.Close(context.Background())
}
