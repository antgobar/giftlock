package claim

import (
	"context"
	"fmt"
	"giftlock/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (s *PostgresRepo) ViewUserClaims(ctx context.Context, userId model.UserId) ([]*model.GroupGift, error) {
	sql := `
		SELECT gifts.id, gifts.group_id, gifts.title, gifts.description, gifts.link, 
		       gifts.created_by, gifts.claimed_by, gifts.claimed_at, gifts.created_at, groups.name
		FROM gifts
		INNER JOIN groups ON gifts.group_id = groups.id
		WHERE gifts.claimed_by = $1
		ORDER BY gifts.claimed_at DESC
	`
	rows, err := s.db.Query(ctx, sql, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve claimed gifts for user %v: %w", userId, err)
	}
	defer rows.Close()

	var gifts = make([]*model.GroupGift, 0)
	for rows.Next() {
		var gift model.GroupGift
		if err := rows.Scan(
			&gift.ID, &gift.GroupId, &gift.Title, &gift.Description, &gift.Link,
			&gift.CreatedBy, &gift.ClaimedBy, &gift.ClaimedAt, &gift.CreatedAt, &gift.GroupName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		gifts = append(gifts, &gift)
	}
	return gifts, nil
}
