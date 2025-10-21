package gift

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	GetAllUser(ctx context.Context, userId model.UserId) ([]*model.Gift, error)
	GetAllGroupUser(ctx context.Context, groupId model.GroupId, userId model.UserId) ([]*model.Gift, error)
	Create(ctx context.Context, gift *model.Gift) (*model.Gift, error)
	Delete(ctx context.Context, giftId model.GiftId, userId model.UserId) error
	Claim(ctx context.Context, giftId model.GiftId, userId model.UserId) (*model.Gift, error)
}
