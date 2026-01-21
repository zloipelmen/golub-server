package storage

import "context"

type Device struct {
	ID        string
	UserID    string
	DeviceKey string
}

func (s *Store) CreateUser(ctx context.Context) (string, error) {
	var id string
	err := s.db.QueryRow(ctx, `INSERT INTO users DEFAULT VALUES RETURNING id`).Scan(&id)
	return id, err
}

func (s *Store) GetDeviceByKey(ctx context.Context, deviceKey string) (*Device, error) {
	row := s.db.QueryRow(ctx, `SELECT id, user_id, device_key FROM devices WHERE device_key=$1`, deviceKey)
	var d Device
	if err := row.Scan(&d.ID, &d.UserID, &d.DeviceKey); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Store) CreateDevice(ctx context.Context, userID, deviceKey, deviceName string) (string, error) {
	var id string
	err := s.db.QueryRow(ctx, `
INSERT INTO devices (user_id, device_key, device_name)
VALUES ($1,$2,$3)
RETURNING id
`, userID, deviceKey, deviceName).Scan(&id)
	return id, err
}

