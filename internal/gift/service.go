package gift

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

func (s *Service) CreateOwnGift(ctx context.Context, userId model.UserId, groupId model.GroupId, title, description, link string) (*model.Gift, error) {
	gift := model.Gift{
		Title: title, Description: description, Link: link, GroupId: groupId, CreatedBy: userId,
	}
	return s.repo.Create(ctx, &gift)
}

func (s *Service) ViewUserGifts(ctx context.Context, userId model.UserId) ([]*model.Gift, error) {
	return s.repo.GetAllUser(ctx, userId)
}

func (s *Service) ViewGroupUserGifts(ctx context.Context, groupId model.GroupId, userId model.UserId) ([]*model.Gift, error) {
	return s.repo.GetAllGroupUser(ctx, groupId, userId)
}

func (s *Service) EditOwnGift(ctx context.Context, userId model.UserId, giftId model.GiftId, title, description, link string) (*model.Gift, error) {
	gift := model.Gift{
		Title: title, Description: description, Link: link, CreatedBy: userId, ID: giftId,
	}
	return s.repo.Edit(ctx, &gift)
}

func (s *Service) DeleteGift(ctx context.Context, userId model.UserId, giftId model.GiftId) error {
	return s.repo.Delete(ctx, giftId, userId)
}

func (s *Service) ClaimGift(ctx context.Context, userId model.UserId, giftId model.GiftId) (*model.Gift, error) {
	return s.repo.Claim(ctx, giftId, userId)
}
