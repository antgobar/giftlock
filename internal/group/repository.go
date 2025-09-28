package group

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	Delete(ctx context.Context, groupID model.GroupId) error
	GetByID(ctx context.Context, groupID model.GroupId) (*model.Group, error)
	Join(ctx context.Context, userID model.UserId, groupID model.GroupId) error
}
