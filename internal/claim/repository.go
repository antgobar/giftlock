package claim

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	ViewUserClaims(ctx context.Context, userId model.UserId) ([]*model.GroupGift, error)
}
