package claim

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

func (s *Service) ViewUserClaims(ctx context.Context, userId model.UserId) ([]*model.GroupGift, error) {
	return s.repo.ViewUserClaims(ctx, userId)
}
