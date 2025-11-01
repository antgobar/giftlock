package gift

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

func (s *PostgresRepo) Create(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	sql := `
		INSERT INTO gifts (
			group_id, title, description, link, price, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id, group_id, title, description, link, price, created_by, created_at;
	`
	row := s.db.QueryRow(ctx, sql,
		gift.GroupId,
		gift.Title,
		gift.Description,
		gift.Link,
		gift.Price,
		gift.CreatedBy,
	)
	var created model.Gift
	err := row.Scan(
		&created.ID,
		&created.GroupId,
		&created.Title,
		&created.Description,
		&created.Link,
		&created.Price,
		&created.CreatedBy,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (s *PostgresRepo) GetAllUser(ctx context.Context, userId model.UserId) ([]*model.GroupGift, error) {
	sql := `
		SELECT gifts.id, gifts.group_id, gifts.title, gifts.description, gifts.link, groups.name
		FROM gifts
		INNER JOIN groups ON gifts.group_id = groups.id
		WHERE gifts.created_by = $1
	`
	rows, err := s.db.Query(ctx, sql, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve gifts for user %v: %w", userId, err)
	}
	defer rows.Close()

	var gifts = make([]*model.GroupGift, 0)
	for rows.Next() {
		var gift model.GroupGift
		if err := rows.Scan(
			&gift.ID, &gift.GroupId, &gift.Title, &gift.Description, &gift.Link, &gift.GroupName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		gifts = append(gifts, &gift)
	}
	return gifts, nil
}

func (s *PostgresRepo) GetAllGroupUser(ctx context.Context, groupId model.GroupId, userId model.UserId) ([]*model.Gift, error) {
	sql := `
		SELECT id, title, description, link, price, created_by, claimed_by
		FROM gifts
		WHERE gifts.created_by = $1 AND gifts.group_id = $2
	`
	rows, err := s.db.Query(ctx, sql, userId, groupId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve gifts for group: %v, user: %v: %w", groupId, userId, err)
	}
	defer rows.Close()

	var gifts = make([]*model.Gift, 0)
	for rows.Next() {
		var gift model.Gift
		if err := rows.Scan(
			&gift.ID, &gift.Title, &gift.Description, &gift.Link, &gift.Price, &gift.CreatedBy, &gift.ClaimedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		gifts = append(gifts, &gift)
	}
	return gifts, nil
}

func (s *PostgresRepo) Delete(ctx context.Context, giftId model.GiftId, userId model.UserId) error {
	sql := `
		DELETE FROM gifts 
		WHERE id = $1 AND created_by = $2
	`
	result, err := s.db.Exec(ctx, sql, giftId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete gift: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("gift not found or user not authorized to delete")
	}

	return nil
}

func (s *PostgresRepo) Claim(ctx context.Context, giftId model.GiftId, userId model.UserId) error {
	sql := `
		UPDATE gifts 
		SET claimed_by = $1, claimed_at = NOW()
		WHERE id = $2 
		  AND claimed_by IS NULL
		  AND EXISTS (
		    SELECT 1 FROM group_members 
		    WHERE group_members.group_id = gifts.group_id 
		      AND group_members.user_id = $1
		  )
	`
	result, err := s.db.Exec(ctx, sql, userId, giftId)
	if err != nil {
		return fmt.Errorf("failed to claim gift: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("gift not found, already claimed, or user not a member of the group")
	}

	return nil
}

func (s *PostgresRepo) Unclaim(ctx context.Context, giftId model.GiftId, userId model.UserId) error {
	sql := `
		UPDATE gifts 
		SET claimed_by = NULL, claimed_at = NULL
		WHERE id = $1 
		  AND claimed_by = $2
		  AND EXISTS (
		    SELECT 1 FROM group_members 
		    WHERE group_members.group_id = gifts.group_id 
		      AND group_members.user_id = $2
		  )
	`
	result, err := s.db.Exec(ctx, sql, giftId, userId)
	if err != nil {
		return fmt.Errorf("failed to unclaim gift: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("gift not found, not claimed by this user, or user not a member of the group")
	}

	return nil
}
