package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Host        string
	Port        string
	DatabaseUrl string
}

func loadConfig() *config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	return &config{
		Host:        os.Getenv("HOST"),
		Port:        os.Getenv("PORT"),
		DatabaseUrl: os.Getenv("DATABASE_URL"),
	}
}
