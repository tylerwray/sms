package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tylerwray/sms/internal/messenger"
)

func NewHandler(ms *messenger.Service) *httprouter.Router {
	handler := httprouter.New()
	handler.GET("/conversations/:conversationID/messages", getMessagesHandler(ms))

	return handler
}

type response struct {
	Data   interface{} `json:"data"`
	Object string      `json:"object"`
}

type message struct {
	InsertedAt string `json:"inserted_at"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	Object     string `json:"object"`
}

func getMessagesHandler(ms *messenger.Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		conversationID, err := strconv.Atoi(ps.ByName("conversationID"))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Given conversationID is invalid"}`))
			return
		}

		messages, err := ms.GetConversationMessages(conversationID)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Internal Server Error"}`))
			return
		}

		var messagesRespose []message

		for _, m := range messages {
			newMessage := message{InsertedAt: m.InsertedAt,
				Message: m.Message,
				Status:  m.Status,
				Object:  "message",
			}

			messagesRespose = append(messagesRespose, newMessage)
		}

		json_resp, err := json.Marshal(response{Data: messagesRespose, Object: "list"})

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
