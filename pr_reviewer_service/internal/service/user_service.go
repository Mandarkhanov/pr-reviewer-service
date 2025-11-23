package service

import (
	"context"
	"pr-reviewer/internal/domain"
)

func (s *Service) SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := s.repoUsers.SetIsActive(ctx, s.db, userID, isActive)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
