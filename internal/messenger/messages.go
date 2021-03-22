package messenger

import (
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

type Message struct {
	ID                 string `db:"id"`
	Message            string `db:"message"`
	Status             string `db:"status"`
	InsertedAt         string `db:"inserted_at"`
	ConversationID     string `db:"conversation_id"`
	ConversationStatus string `db:"conversation_status"`
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db}
}

func (s *Service) GetConversationMessages(conversationID int) ([]Message, error) {
	messages := []Message{}

	var query = `
		SELECT
			m.inserted_at,
			m.message,
			ms.status
		FROM messages m
		LEFT JOIN message_statuses ms ON m.id = ms.message_id
		WHERE m.conversation_id = $1;
	`

	err := s.db.Select(&messages, query, conversationID)

	return messages, err
}

func (s *Service) UpdateMessageStatus(smsXID, status string) error {
	message := Message{}

	var query = `
		SELECT
			m.id,
			m.inserted_at,
			m.message,
			ms.status,
			c.status AS conversation_status,
			c.id AS conversation_id
		FROM messages m
		JOIN conversations c ON c.id = m.conversation_id
		LEFT JOIN message_statuses ms ON m.id = ms.message_id
		WHERE m.sms_xid = $1;
	`

	err := s.db.Get(&message, query, smsXID)

	if err != nil {
		return err
	}

	if message.ConversationStatus == "blocked" {
		// Do nothing when a conversation has been blocked
		return nil
	}

	// Update message status
	var updateStatus = `
		UPDATE message_statuses
		SET status = $1
		WHERE message_id = $2
  `

	s.db.MustExec(updateStatus, status, message.ID)

	if status == "failed" && message.ConversationStatus == "closed" {
		// Re-open conversation to notify user of failed message
		var updateConversationStatus = `
			UPDATE conversations
			SET status = 'open'
			WHERE id = $1
		`

		s.db.MustExec(updateConversationStatus, message.ConversationID)
	}

	return nil
}
