package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/tylerwray/sms/internal/messenger"
	"github.com/tylerwray/sms/internal/server"
)

func TestGetConversationMessages(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"inserted_at", "message", "status"}).
		AddRow("2021-06-17T18:19:35Z", "Hello, Moto", "pending")

	mock.ExpectQuery("^SELECT (.+) FROM messages m").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/conversations/1/messages", nil)
	rr := httptest.NewRecorder()

	ms := messenger.NewService(sqlxDB)
	h := server.NewHandler(ms)
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	expected := []byte(`{"data":[{"inserted_at":"2021-06-17T18:19:35Z","message":"Hello, Moto","status":"pending","object":"message"}],"object":"list"}`)

	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %s want %s",
			rr.Body.Bytes(), expected)
	}
}
