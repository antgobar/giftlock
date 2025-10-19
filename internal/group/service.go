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

func (s *Service) CreateAndJoinGroup(ctx context.Context, userID model.UserId, groupName, groupDescription string) (*model.Group, error) {
	group := &model.Group{
		Name:        groupName,
		Description: groupDescription,
		CreatedBy:   userID,
	}
	createdGroup, err := s.repo.Create(ctx, group)
	if err != nil {
		return nil, err
	}
	if err := s.repo.AddMember(ctx, userID, userID, createdGroup.ID); err != nil {
		return nil, err
	}
	return createdGroup, nil
}

func (s *Service) AddMember(ctx context.Context, ownerId, memberId model.UserId, groupId model.GroupId) error {
	return s.repo.AddMember(ctx, ownerId, memberId, groupId)
}

func (s *Service) DeleteGroup(ctx context.Context, userId model.UserId, groupID model.GroupId) error {
	return s.repo.Delete(ctx, userId, groupID)
}

func (s *Service) GetJoinedGroups(ctx context.Context, userID model.UserId) ([]*model.Group, error) {
	return s.repo.ListJoined(ctx, userID)
}

func (s *Service) ViewGroup(ctx context.Context, userId model.UserId, groupId model.GroupId) ([]*model.GroupMemberDetails, error) {
	groupMembersDetails, err := s.repo.GroupMemberDetails(ctx, userId, groupId)
	if err != nil {
		return nil, err
	}
	return groupMembersDetails, nil
}
