package group

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	Delete(ctx context.Context, userId model.UserId, groupID model.GroupId) error
	ListCreated(ctx context.Context, userID model.UserId) ([]*model.Group, error)
	Join(ctx context.Context, userID model.UserId, username string, groupID model.GroupId) error
	GroupMembers(ctx context.Context, groupId model.GroupId) ([]*model.GroupMember, error)
	GroupDetails(ctx context.Context, groupId model.GroupId) (*model.Group, error)
}
