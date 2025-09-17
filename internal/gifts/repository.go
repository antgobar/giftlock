package gifts

import (
	"context"
	"giftlock/internal/model"
)

type Repository interface {
	Create(ctx context.Context, gift *model.Gift)
	Edit(ctx context.Context, gift *model.Gift)
	Delete(ctx context.Context, giftId *model.GiftId)
	Claim(ctx context.Context, userId *model.UserId, giftId *model.GiftId)
}
