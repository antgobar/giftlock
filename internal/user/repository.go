package user

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	Create(ctx context.Context, username string, password string) (*model.User, error)
	GetFromCreds(ctx context.Context, username string, password string) (*model.User, error)
	SearchUserNotInGroup(ctx context.Context, groupId model.GroupId, username string) ([]*model.User, error)
}
