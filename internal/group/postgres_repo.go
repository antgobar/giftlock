package group

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

func (s *PostgresRepo) Delete(ctx context.Context, userID model.UserId, groupID model.GroupId) error {
	membersSql := `DELETE FROM group_members WHERE group_id = $1;`
	_, err := s.db.Exec(ctx, membersSql, groupID)
	if err != nil {
		return err
	}

	sql := `DELETE FROM groups WHERE id = $1 AND created_by = $2;`
	_, err = s.db.Exec(ctx, sql, groupID, userID)
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

func (s *PostgresRepo) AddMember(ctx context.Context, ownerId, memberId model.UserId, groupID model.GroupId) error {
	sql := `
		INSERT INTO group_members (user_id, group_id)
		SELECT $1, $2
		WHERE EXISTS (
			SELECT 1 FROM groups 
			WHERE id = $2 AND created_by = $3
		);
	`
	result, err := s.db.Exec(ctx, sql, memberId, groupID, ownerId)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("group not found or user is not the owner")
	}

	return nil
}

func (s *PostgresRepo) GroupMemberDetails(ctx context.Context, userID model.UserId, groupId model.GroupId) ([]*model.GroupMemberDetails, error) {
	sql := `
		SELECT groups.id, groups.created_by, groups.name, groups.description, users.id, users.username
		FROM groups
		JOIN group_members ON groups.id = group_members.group_id
		JOIN users ON group_members.user_id = users.id
		WHERE groups.id = $1;
	`
	rows, err := s.db.Query(ctx, sql, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupMembersDetails []*model.GroupMemberDetails

	for rows.Next() {
		var groupMemberDetails model.GroupMemberDetails
		if err := rows.Scan(
			&groupMemberDetails.GroupId,
			&groupMemberDetails.GroupCreatorId,
			&groupMemberDetails.GroupName,
			&groupMemberDetails.GroupDescription,
			&groupMemberDetails.MemberId,
			&groupMemberDetails.MemberUsername,
		); err != nil {
			return nil, err
		}
		groupMembersDetails = append(groupMembersDetails, &groupMemberDetails)
	}
	return groupMembersDetails, nil
}
