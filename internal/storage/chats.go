package storage

import "context"

func (s *Store) IsChatMember(ctx context.Context, chatID, userID string) (bool, error) {
	var ok bool
	err := s.db.QueryRow(ctx, `
SELECT EXISTS(
  SELECT 1 FROM chat_members WHERE chat_id=$1 AND user_id=$2
)`, chatID, userID).Scan(&ok)
	return ok, err
}

func (s *Store) ListChatMemberUserIDs(ctx context.Context, chatID string) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT user_id FROM chat_members WHERE chat_id=$1`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

