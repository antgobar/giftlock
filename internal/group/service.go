package group

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

func (s *Service) CreateAndJoinGroup(ctx context.Context, userID model.UserId, name, description string) (*model.Group, error) {
	group := &model.Group{
		Name:        name,
		Description: description,
		CreatedBy:   userID,
	}
	createdGroup, err := s.repo.Create(ctx, group)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Join(ctx, userID, createdGroup.ID); err != nil {
		return nil, err
	}
	return createdGroup, nil
}

func (s *Service) DeleteGroup(ctx context.Context, groupID model.GroupId) error {
	return s.repo.Delete(ctx, groupID)
}

func (s *Service) GetGroupByID(ctx context.Context, groupID model.GroupId) (*model.Group, error) {
	return s.repo.GetByID(ctx, groupID)
}

func (s *Service) JoinGroup(ctx context.Context, userID model.UserId, groupID model.GroupId) error {
	return s.repo.Join(ctx, userID, groupID)
}
