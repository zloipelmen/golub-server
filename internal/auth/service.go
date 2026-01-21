package auth

import (
	"context"
	"errors"

	"messenger/internal/storage"
)

type Service struct {
	store *storage.Store
}

func NewService(store *storage.Store) *Service {
	return &Service{store: store}
}

var ErrAuth = errors.New("auth failed")

func (s *Service) AuthByInvite(ctx context.Context, inviteCode, deviceKey, deviceName string) (userID string, deviceID string, err error) {
	if inviteCode == "" || deviceKey == "" {
		return "", "", ErrAuth
	}

	// 1) Если девайс уже существует — НЕ тратим инвайт
	d, derr := s.store.GetDeviceByKey(ctx, deviceKey)
	if derr == nil && d != nil {
		_ = s.store.TouchDeviceLastSeen(ctx, d.ID)
		return d.UserID, d.ID, nil
	}

	// 2) Если девайс новый — тогда тратим инвайт
	if err := s.store.ConsumeInvite(ctx, inviteCode); err != nil {
		return "", "", err
	}

	// 3) Создаём нового пользователя и девайс
	uid, err := s.store.CreateUser(ctx)
	if err != nil {
		return "", "", err
	}
	did, err := s.store.CreateDevice(ctx, uid, deviceKey, deviceName)
	if err != nil {
		return "", "", err
	}
	return uid, did, nil
}

