package user

import (
	"context"
	"errors"
	"fmt"
	"giftlock/internal/model"
	"giftlock/internal/security"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

var ErrUsernameTaken = errors.New("username taken")
var ErrUserNotExists = errors.New("user does not exist")

func isUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return true
	}
	return false
}

func isNoRowsFoundError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (s *PostgresRepo) Create(ctx context.Context, username string, password string) (*model.User, error) {
	sql := `
		INSERT INTO users (username, hashed_password)
		VALUES ($1, $2)
		Returning id, username, created_at
	`

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		log.Println("ERROR:", err.Error())
		return nil, security.ErrHashingError
	}

	user := model.User{
		Username:       username,
		HashedPassword: hashedPassword,
	}

	row := s.db.QueryRow(ctx, sql, user.Username, user.HashedPassword)
	err = row.Scan(&user.ID, &user.Username, &user.CreatedAt)

	if isUniqueViolationError(err) {
		return nil, ErrUsernameTaken
	}
	if err != nil {
		return nil, fmt.Errorf("failed to register user %v: %w", user, err)
	}

	return &user, nil
}

func (s *PostgresRepo) GetFromCreds(ctx context.Context, username string, password string) (*model.User, error) {
	sql := `
		SELECT id, username, hashed_password, created_at
		FROM users 
		WHERE username = $1
	`

	user := model.User{
		Username: username,
	}

	row := s.db.QueryRow(ctx, sql, user.Username)
	if err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt); err != nil {
		if isNoRowsFoundError(err) {
			return nil, ErrUserNotExists
		}
		return nil, fmt.Errorf("failed to find user %v: %w", user, err)
	}

	if !security.CheckPasswordHash(password, user.HashedPassword) {
		return nil, errors.New("incorrect username or password")
	}

	return &user, nil
}

func (s *PostgresRepo) SearchUserNotInGroup(ctx context.Context, groupId model.GroupId, username string, limit int) ([]*model.User, error) {
	sql := `
		SELECT id, username
		FROM users 
		WHERE username ILIKE $1
		AND id NOT IN (
			SELECT user_id 
			FROM group_members 
			WHERE group_id = $2
		)
		ORDER BY username
		LIMIT $3
	`
	searchPattern := "%" + username + "%"
	rows, err := s.db.Query(ctx, sql, searchPattern, groupId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *PostgresRepo) SearchByUsername(ctx context.Context, username string) (*model.User, error) {
	sql := `
		SELECT id, username
		FROM users 
		WHERE username = $1
	`

	user := model.User{}
	row := s.db.QueryRow(ctx, sql, username)
	if err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt); err != nil {
		if isNoRowsFoundError(err) {
			return nil, ErrUserNotExists
		}
		return nil, fmt.Errorf("failed to find user with username %s: %w", username, err)
	}

	return &user, nil
}

func (s *PostgresRepo) Delete(ctx context.Context, userId model.UserId) error {
	sql := `DELETE FROM users WHERE id = $1`

	result, err := s.db.Exec(ctx, sql, userId)
	if err != nil {
		return fmt.Errorf("failed to delete user with id %d: %w", userId, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotExists
	}

	return nil
}

func (s *PostgresRepo) List(ctx context.Context) ([]*model.User, error) {
	sql := `
		SELECT id, username, created_at, role
		FROM users 
		ORDER BY username
	`

	rows, err := s.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Role); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	return users, nil
}
