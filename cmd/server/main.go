package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tylerwray/sms/internal/messenger"
	"github.com/tylerwray/sms/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var host = os.Getenv("HOST")
	var port = os.Getenv("PORT")
	var databaseUrl = os.Getenv("DATABASE_URL")

	db, err := sqlx.Connect("postgres", databaseUrl)

	if err != nil {
		log.Println("Cannot connect to DB")
		log.Panic(err)
	}

	ms := messenger.NewService(db)

	handler := server.NewHandler(ms)
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server is running on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Println("Server could not be started, error:")
		log.Panic(err)
	}
}
