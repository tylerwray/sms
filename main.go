package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	c := loadConfig()

	db, err := sqlx.Connect("postgres", c.DatabaseUrl)

	if err != nil {
		log.Println("Cannot connect to DB")
		log.Panic(err)
	}

	router := httprouter.New()
	router.GET("/conversations/:conversationID/messages", getMessagesHandler(&messageService{db}))

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	log.Printf("Server is running on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Println("Server could not be started, error:")
		log.Panic(err)
	}
}
