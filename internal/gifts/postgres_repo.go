package gifts

import (
	"context"
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
