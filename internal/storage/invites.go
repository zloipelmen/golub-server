package storage

import (
	"context"
	"errors"
	"time"
)

var ErrInviteInvalid = errors.New("invalid invite")

func (s *Store) ConsumeInvite(ctx context.Context, code string) error {
	const q = `
UPDATE invites
SET uses = uses + 1
WHERE code = $1
  AND disabled = false
  AND (expires_at IS NULL OR expires_at > now())
  AND uses < max_uses
`
	ct, err := s.db.Exec(ctx, q, code)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrInviteInvalid
	}
	return nil
}

func (s *Store) TouchDeviceLastSeen(ctx context.Context, deviceID string) error {
	_, err := s.db.Exec(ctx, `UPDATE devices SET last_seen_at=$1 WHERE id=$2`, time.Now().UTC(), deviceID)
	return err
}
