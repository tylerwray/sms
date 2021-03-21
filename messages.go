package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type messageService struct {
	DB *sqlx.DB
}

type message struct {
	InsertedAt string `db:"inserted_at" json:"inserted_at"`
	Message    string `db:"message" json:"message"`
	Status     string `db:"status" json:"status"`
	Object     string `db:"object" json:"object"`
}

func (ms *messageService) getConversationMessages(conversationID int) ([]message, error) {
	messages := []message{}

	var query = `
		SELECT
			m.inserted_at,
			m.message,
			ms.status,
			'message' AS object
		FROM messages m
		LEFT JOIN message_statuses ms ON m.id = ms.message_id
		WHERE m.conversation_id = $1;
	`

	err := ms.DB.Select(&messages, query, conversationID)

	return messages, err
}

type resp struct {
	Data   interface{} `json:"data"`
	Object string      `json:"object"`
}

func getMessagesHandler(ms *messageService) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		conversationID, err := strconv.Atoi(ps.ByName("conversationID"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Given conversationID is invalid"}`))
			return
		}

		messages, err := ms.getConversationMessages(conversationID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal Server Error"}`))
			return
		}

		json_resp, err := json.Marshal(resp{Data: messages, Object: "list"})

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal Server Error"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json_resp)
	}
}
