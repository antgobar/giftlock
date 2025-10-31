package user

import (
	"context"
	"giftlock/internal/model"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Register(ctx context.Context, username, password string) error {
	_, err := s.repo.Create(ctx, username, password)
	return err
}

func (s *Service) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	user, err := s.repo.SearchByUsername(ctx, username)
	if err == ErrUserNotExists {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (s *Service) SearchUserNotInGroup(ctx context.Context, groupId model.GroupId, username string) ([]*model.User, error) {
	limit := 1
	return s.repo.SearchUserNotInGroup(ctx, groupId, username, limit)
}
