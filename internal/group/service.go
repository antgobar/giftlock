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

func (s *Service) DeleteGroup(ctx context.Context, userId model.UserId, groupID model.GroupId) error {
	return s.repo.Delete(ctx, userId, groupID)
}

func (s *Service) GetCreatedGroups(ctx context.Context, userID model.UserId) ([]*model.Group, error) {
	return s.repo.ListCreated(ctx, userID)
}

func (s *Service) JoinGroup(ctx context.Context, userID model.UserId, groupID model.GroupId) error {
	return s.repo.Join(ctx, userID, groupID)
}

func (s *Service) ViewGroup(ctx context.Context, userId model.UserId, groupId model.GroupId) (*model.GroupDetails, error) {
	return s.repo.GetGroupDetails(ctx, userId, groupId)
}
