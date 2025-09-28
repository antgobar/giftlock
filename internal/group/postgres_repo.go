package group

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

func (s *PostgresRepo) Create(ctx context.Context, group *model.Group) (*model.Group, error) {
	sql := `
		INSERT INTO groups (
			name, description, created_by
		) VALUES (
			$1, $2, $3
		) RETURNING id, name, description, created_by, created_at;
	`
	row := s.db.QueryRow(ctx, sql,
		group.Name,
		group.Description,
		group.CreatedBy,
	)
	var created model.Group
	err := row.Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.CreatedBy,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (s *PostgresRepo) Delete(ctx context.Context, groupID model.GroupId) error {
	sql := `DELETE FROM groups WHERE id = $1;`
	_, err := s.db.Exec(ctx, sql, groupID)
	return err
}

func (s *PostgresRepo) GetByID(ctx context.Context, groupID model.GroupId) (*model.Group, error) {
	sql := `
		SELECT id, name, description, created_by, created_at
		FROM groups
		WHERE id = $1;
	`
	row := s.db.QueryRow(ctx, sql, groupID)
	var group model.Group
	err := row.Scan(
		&group.ID,
		&group.Name,
		&group.Description,
		&group.CreatedBy,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *PostgresRepo) Join(ctx context.Context, userID model.UserId, groupID model.GroupId) error {
	sql := `
		INSERT INTO group_members (user_id, group_id)
		VALUES ($1, $2);
	`
	_, err := s.db.Exec(ctx, sql, userID, groupID)
	return err
}
