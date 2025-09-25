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
			title, description, link, price, created_by
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, title, description, link, price, created_by, created_at;
	`
	row := s.db.QueryRow(ctx, sql,
		gift.Title,
		gift.Description,
		gift.Link,
		gift.Price,
		gift.CreatedBy,
	)
	var created model.Gift
	err := row.Scan(
		&created.ID,
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

func (s *PostgresRepo) GetAllUser(ctx context.Context, userId model.UserId) ([]*model.Gift, error) {
	sql := `
		SELECT id, title, description, price, link, claimed_at, created_by
		FROM gifts
		WHERE gifts.created_by = $1
	`
	rows, err := s.db.Query(ctx, sql, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve gifts for user %v: %w", userId, err)
	}
	defer rows.Close()

	var gifts = make([]*model.Gift, 0)
	for rows.Next() {
		var gift model.Gift
		if err := rows.Scan(
			&gift.ID, &gift.Title, &gift.Description, &gift.Link, &gift.ClaimedAt, &gift.CreatedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		gifts = append(gifts, &gift)
	}
	return gifts, nil
}

func (s *PostgresRepo) Edit(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	sql := `
		UPDATE gifts 
		SET title = $1, description = $2, link = $3, price = $4
		WHERE id = $5 AND created_by = $6
		RETURNING id, title, description, link, price, created_by, created_at, claimed_by, claimed_at;
	`
	row := s.db.QueryRow(ctx, sql,
		gift.Title,
		gift.Description,
		gift.Link,
		gift.Price,
		gift.ID,
		gift.CreatedBy,
	)
	var updated model.Gift
	err := row.Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.Link,
		&updated.Price,
		&updated.CreatedBy,
		&updated.CreatedAt,
		&updated.ClaimedBy,
		&updated.ClaimedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update gift: %w", err)
	}
	return &updated, nil
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

func (s *PostgresRepo) Claim(ctx context.Context, giftId model.GiftId, userId model.UserId) (*model.Gift, error) {
	sql := `
		UPDATE gifts 
		SET claimed_by = $1, claimed_at = NOW()
		WHERE id = $2 AND claimed_by IS NULL
		RETURNING id, title, description, link, price, created_by, created_at, claimed_by, claimed_at;
	`
	row := s.db.QueryRow(ctx, sql, userId, giftId)

	var claimed model.Gift
	err := row.Scan(
		&claimed.ID,
		&claimed.Title,
		&claimed.Description,
		&claimed.Link,
		&claimed.Price,
		&claimed.CreatedBy,
		&claimed.CreatedAt,
		&claimed.ClaimedBy,
		&claimed.ClaimedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to claim gift (may already be claimed): %w", err)
	}
	return &claimed, nil
}
