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

func (s *PostgresRepo) ListCreated(ctx context.Context, userID model.UserId) ([]*model.Group, error) {
	sql := `
		SELECT id, name, description, created_by, created_at
		FROM groups
		WHERE created_by = $1;
	`
	rows, err := s.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		var group model.Group
		if err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.CreatedBy,
			&group.CreatedAt,
		); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}

func (s *PostgresRepo) Join(ctx context.Context, userID model.UserId, groupID model.GroupId) error {
	sql := `
		INSERT INTO group_members (user_id, group_id)
		VALUES ($1, $2);
	`
	_, err := s.db.Exec(ctx, sql, userID, groupID)
	return err
}
