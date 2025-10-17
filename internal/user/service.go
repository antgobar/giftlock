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

func (s *Service) SearchUserNotInGroup(ctx context.Context, groupId model.GroupId, username string) ([]*model.User, error) {
	return s.repo.SearchUserNotInGroup(ctx, groupId, username)
}
