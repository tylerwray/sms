package ingest_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/tylerwray/sms/internal/ingest"
	"github.com/tylerwray/sms/internal/messenger"
)

func TestIngestCSV(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}

	mock.MatchExpectationsInOrder(false)

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	failedRows := sqlmock.NewRows([]string{
		"id",
		"inserted_at",
		"message",
		"status",
		"conversation_status",
		"conversation_id",
	}).AddRow(
		1,
		"2021-06-17T18:19:35Z",
		"Hello, Moto",
		"pending",
		"open",
		1,
	)

	sentRows := sqlmock.NewRows([]string{
		"id",
		"inserted_at",
		"message",
		"status",
		"conversation_status",
		"conversation_id",
	}).AddRow(
		1,
		"2021-06-17T18:19:35Z",
		"Hello, Moto",
		"pending",
		"open",
		1,
	)

	mock.ExpectQuery("^SELECT (.+) FROM messages m").WithArgs("761964A1-FD61-18E9-4BC0-8EC572413436").WillReturnRows(sentRows)
	mock.ExpectQuery("^SELECT (.+) FROM messages m").WithArgs("47B4906F-84EE-C509-7B5A-0F071B6DAD04").WillReturnRows(failedRows)
	mock.ExpectExec("UPDATE message_statuses").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE message_statuses").WillReturnResult(sqlmock.NewResult(1, 1))

	ms := messenger.NewService(sqlxDB)

	ingest.FromCSV(ms, "./testdata/message_status_updates.csv")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("Not all queries we run.")
	}
}

func TestIngestCSVReopenConversation(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}

	mock.MatchExpectationsInOrder(false)

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{
		"id",
		"inserted_at",
		"message",
		"status",
		"conversation_status",
		"conversation_id",
	}).AddRow(
		1,
		"2021-06-17T18:19:35Z",
		"Hello, Moto",
		"pending",
		"closed",
		1,
	)

	mock.ExpectQuery("^SELECT (.+) FROM messages m").WillReturnRows(rows)
	mock.ExpectExec("UPDATE message_statuses").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE conversations").WillReturnResult(sqlmock.NewResult(1, 1))

	ms := messenger.NewService(sqlxDB)

	ingest.FromCSV(ms, "./testdata/conversation_status_reopen.csv")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("Not all queries we run.")
	}
}
