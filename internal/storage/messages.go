package storage

import (
	"context"
	"time"
)

type Message struct {
	ID           string
	ChatID       string
	SenderUserID string
	Text         string
	CreatedAt    time.Time
}

type CreateMessageParams struct {
	ChatID       string
	SenderUserID string
	ClientMsgID  string
	Text         string
}

func (s *Store) CreateMessage(ctx context.Context, p CreateMessageParams) (*Message, error) {
	var m Message
	err := s.db.QueryRow(ctx, `
INSERT INTO messages (chat_id, sender_user_id, client_msg_id, text)
VALUES ($1,$2,$3,$4)
ON CONFLICT (sender_user_id, client_msg_id)
DO UPDATE SET text = EXCLUDED.text
RETURNING id, chat_id, sender_user_id, text, created_at
`, p.ChatID, p.SenderUserID, p.ClientMsgID, p.Text).Scan(
		&m.ID, &m.ChatID, &m.SenderUserID, &m.Text, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Store) ListMessages(ctx context.Context, chatID string, limit int) ([]Message, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.Query(ctx, `
SELECT id, chat_id, sender_user_id, text, created_at
FROM messages
WHERE chat_id=$1
ORDER BY created_at DESC
LIMIT $2
`, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Message, 0, limit)
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.SenderUserID, &m.Text, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
