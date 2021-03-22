package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tylerwray/sms/internal/ingest"
	"github.com/tylerwray/sms/internal/messenger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var databaseUrl = os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", databaseUrl)

	if err != nil {
		log.Println("Cannot connect to DB")
		log.Panic(err)
	}

	ms := messenger.NewService(db)

	ingest.FromCSV(ms, os.Args[1])
}
