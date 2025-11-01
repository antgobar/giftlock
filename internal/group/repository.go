package group

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	Delete(ctx context.Context, userId model.UserId, groupID model.GroupId) error
	ListJoined(ctx context.Context, userID model.UserId) ([]*model.Group, error)
	AddMember(ctx context.Context, ownerId, memberId model.UserId, groupID model.GroupId) error
	GroupMemberDetails(ctx context.Context, userID model.UserId, groupId model.GroupId) ([]*model.GroupMemberDetails, error)
	Leave(ctx context.Context, userId model.UserId, groupID model.GroupId) error
}
