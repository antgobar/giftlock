package user

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	Create(ctx context.Context, username string, password string) (*model.User, error)
	GetFromCreds(ctx context.Context, username string, password string) (*model.User, error)
	SearchUserNotInGroup(ctx context.Context, groupId model.GroupId, username string, limit int) ([]*model.User, error)
	SearchByUsername(ctx context.Context, username string) (*model.User, error)
	Delete(ctx context.Context, userId model.UserId) error
	List(ctx context.Context) ([]*model.User, error)
}
